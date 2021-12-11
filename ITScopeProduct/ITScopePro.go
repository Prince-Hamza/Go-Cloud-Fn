package ITScopeProduct

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	ApiSet "Main.go/Core/Api"
	CorsSet "Main.go/Core/Cors"
	StructSet "Main.go/Core/Structures"
	SuperJsonSet "Main.go/Core/SuperJson"
	// WooApiSet "Main.go/Core/WoocommerceApi"
)

type ITScopePro struct{}

type Error struct {
	Error    string
	Reason   string
	Solution string
}

type Response struct {
	Resp string
}

var ConsumerKey string = "ck_42a75ce7a233bc1e341e33779723c304e6d820cc"
var ConsumerSecret string = "cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526"

var waitGroup sync.WaitGroup = sync.WaitGroup{}
var awaitUpdate sync.WaitGroup = sync.WaitGroup{}

var fieldsInResponse [3]string
var finalInfo []string

// config
var create string
var update string
var titles string
var categories string
var attributes string
var images string
var descrption string
var priceStock string = "true"

func (ITS ITScopePro) ParseItScopeProduct(res http.ResponseWriter, req *http.Request) {

	fmt.Println("parsing It Scope Product")

	cors := CorsSet.CorsAccess{}
	super := SuperJsonSet.SuperJson{}
	cors.Cors(res, req)
	setParams(req)
	finalInfo = nil

	JsonString := super.Stringify(req.Body)
	fmt.Println("stringify:Done")
	fmt.Println("json strng input: ", JsonString)

	var itScopeJson StructSet.ITScopeInfo
	err := json.Unmarshal([]byte(JsonString), &itScopeJson) // Json to struct
	if err != nil {
		fmt.Println("error in parsing")
		sendError("Failed to Convert Req Body to ITScope Simplified Structure", "Invalid Json", "Valid Json", res)
		return
	}

	fmt.Println("parsing is done")

	waitGroup.Add(1)
	go parallelUpdate(itScopeJson, res)
	waitGroup.Wait()

	// send response

	fmt.Println("Wait Over")
	fmt.Println("Sending Final Response :) ")

	if len(finalInfo) > 0 {
		json.NewEncoder(res).Encode(finalInfo)
	} else {
		json.NewEncoder(res).Encode(`{data : "void" , updated: 0}`)
	}

	return

}

func setParams(req *http.Request) {
	images = req.URL.Query().Get("images")
	categories = req.URL.Query().Get("categories")
	attributes = req.URL.Query().Get("attributes")
	fmt.Println("param:images :: ", images)
	fmt.Println("param:categories :: ", categories)
	fmt.Println("param:attribs :: ", attributes)

}

func parallelUpdate(itScopeJson StructSet.ITScopeInfo, res http.ResponseWriter) {

	valid := Validate(itScopeJson, res)
	if valid != true {
		fmt.Println("Error in condition")
		waitGroup.Done()
		return
	}

	fmt.Println("Attributes Validated Successfully")

	fmt.Println("products length : ", len(itScopeJson.Products))

	for i := 0; i <= len(itScopeJson.Products)-1; i++ {
		fmt.Println("starting go routine")
		awaitUpdate.Add(1)
		go updateRoutine(itScopeJson, i, res)
	}

	awaitUpdate.Wait()
	waitGroup.Done()

}

func updateRoutine(itScopeJson StructSet.ITScopeInfo, index int, res http.ResponseWriter) {

	fmt.Println("update Product routine")
	fmt.Println("Sku : " , itScopeJson.Products[index].Sku)

	Api := ApiSet.Api{}
	//pricePackage := PriceSet.Prices{}

	product := itScopeJson.Products[index]

	productId := Api.Post("https://firewallforce.se/wp-json/wc/v3/idbysku?consumer_key="+ConsumerKey+"&consumer_secret="+ConsumerSecret, `{"Sku" : "`+product.Sku+`"}`)
	fmt.Println("Id By Sku : ", productId)

	setResponseFields()

	if productId != "0" {
		//finalPrice := pricePackage.GetFinalPrice(*product.Price)

		fmt.Println("Product Stock: ", *product.Stock)

		WooProduct := StructSet.WoocommerceInfo{
			FieldsInResponse: fieldsInResponse,
			Price:            *product.Price,
			RegularPrice:     *product.Price,
			Type:             "simple",
			ManageStock:      true,
			StockQuantity:    *product.Stock,
			StockStatus:      *product.StockStatus,
		}

		// WooProduct.Extra = product.Extra
		if attributes == "true" {
			WooProduct.Attributes = product.Attributes
		}
		if categories == "true" {
			WooProduct.Categories = product.Categories
		}
		if images == "true" {
			fmt.Println("including images")
			WooProduct.Images = product.Images
		}

		wooCommerceJson, _ := json.Marshal(WooProduct)
		fmt.Println("parsed json : ", string(wooCommerceJson))

		fmt.Println("updating Product")
		wooResponse := Api.Post("https://firewallforce.se/wp-json/wc/v3/products/"+productId+"?"+"consumer_key="+ConsumerKey+"&consumer_secret="+ConsumerSecret, string(wooCommerceJson))

		fmt.Println("update response", wooResponse)

		finalInfo = append(finalInfo, wooResponse)
	} else {


		WooProduct := StructSet.WoocommerceInfo{
			Sku: product.Sku,
			FieldsInResponse: fieldsInResponse,
			Price:            *product.Price,
			RegularPrice:     *product.Price,
			Type:             "simple",
			ManageStock:      true,
			StockQuantity:    *product.Stock,
			StockStatus:      *product.StockStatus,
		}

		// WooProduct.Extra = product.Extra
		if attributes == "true" {
			WooProduct.Attributes = product.Attributes
		}
		if categories == "true" {
			WooProduct.Categories = product.Categories
		}
		if images == "true" {
			fmt.Println("including images")
			WooProduct.Images = product.Images
		}
		WooProduct.Name = product.Title
		WooProduct.Description = product.Description

		wooCommerceJson, _ := json.Marshal(WooProduct)
		fmt.Println("parsed json : ", string(wooCommerceJson))

		fmt.Println("creating Product")
		wooResponse := Api.Post("https://firewallforce.se/wp-json/wc/v3/products?" +"consumer_key="+ConsumerKey+"&consumer_secret="+ConsumerSecret, string(wooCommerceJson))

		fmt.Println("create response", wooResponse)

		finalInfo = append(finalInfo, wooResponse)

	}

	fmt.Println("create: done")
	awaitUpdate.Done()

}

func setResponseFields() {
	fieldsInResponse[0] = "id"
	fieldsInResponse[1] = "sku"
	fieldsInResponse[2] = "price"
}
func Validate(itScopeJson StructSet.ITScopeInfo, res http.ResponseWriter) bool {

	fmt.Println("products to validate : ", itScopeJson.Products)

	for _, product := range itScopeJson.Products {

		if product.Stock == nil {

			sendError("Missing Attribute", "Stock is Missing", "send Valid Json", res)
			return false
		}

		if product.StockStatus == nil {
			sendError("Missing Attribute", "Stock Status is Missing", "send Valid Json", res)
			return false
		} else {
			if strings.Contains(*product.StockStatus, "Not") {
				*product.StockStatus = "outofstock"
			} else {
				*product.StockStatus = "instock"
			}
		}

		if product.Price == nil {
			sendError("Missing Attribute", "Price is Missing", "send Valid Json", res)
			return false
		}

		// if product.Category == nil {
		// 	sendError("Missing Attribute", "productSubType is Missing", "send Valid Json", res)
		// 	return false
		// }

		// if product.StandardHtmlDatasheet == nil {
		// 	sendError("Missing Attribute", "StandardHtmlDatasheet is Missing", "send Valid Json", res)
		// 	return false
		// }

	}
	return true

}

// func IdBySku(sku string) {
// 	Api := ApiSet.Api{}
// 	productId := Api.Get("https://firewallforce.se/wp-json/wc/v3/idbysku?sku=" + sku + "&consumer_key=" + ConsumerKey + "&consumer_secret=" + ConsumerSecret)

// 	intId, err := strconv.Atoi(productId)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//productIds = append(productIds, intId)
// 	waitGroup.Done()

// }

func sendError(Err string, Reason string, Solution string, res http.ResponseWriter) {
	fmt.Println(Err)
	respStruct := Error{Error: Err, Reason: Reason, Solution: Solution}
	json.NewEncoder(res).Encode(respStruct)
	return
}

func sendResp(resp string, res http.ResponseWriter) {
	fmt.Println(resp)
	respStruct := Response{Resp: resp}
	json.NewEncoder(res).Encode(respStruct)
	return
}

func sendFinalResponse(structure []StructSet.WoocommerceInfo, res http.ResponseWriter) {
	//json.NewEncoder(res).Encode((structure))
}
