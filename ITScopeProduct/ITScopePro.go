package ITScopeProduct

import (
	"encoding/json"
	"fmt"
	"strings"

	//"log"
	//"strconv"
	"sync"

	// "github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	"net/http"

	ApiSet "Main.go/Core/Api"
	CorsSet "Main.go/Core/Cors"
	StructSet "Main.go/Core/Structures"
	SuperJsonSet "Main.go/Core/SuperJson"
	PriceSet "Main.go/ITScopeProduct/Prices"
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

func (ITS ITScopePro) ParseItScopeProduct(res http.ResponseWriter, req *http.Request) {

	fmt.Println("parsing It Scope Product")

	cors := CorsSet.CorsAccess{}
	super := SuperJsonSet.SuperJson{}

	cors.Cors(res, req)
	JsonString := super.Stringify(req.Body)

	var itScopeJson StructSet.ITScopeInfo
	err := json.Unmarshal([]byte(JsonString), &itScopeJson) // Json to struct
	if err != nil {
		sendError("Failed to Parse Req Body", "Invalid Json", "Valid Json", res)
		return
	}

	waitGroup.Add(1)
	go parallelUpdate(itScopeJson, res)
	waitGroup.Wait()

	// send response

	json.NewEncoder(res).Encode(finalInfo)
	return

}

func parallelUpdate(itScopeJson StructSet.ITScopeInfo, res http.ResponseWriter) {

	valid := Validate(itScopeJson, res)
	if valid != true {
		fmt.Println("Error in condition")
		waitGroup.Done()
		return
	}

	for i := 1; i <= len(itScopeJson.Products)-1; i++ {
		awaitUpdate.Add(1)
		go updateRoutine(itScopeJson, i, res)
	}

	awaitUpdate.Wait()
	waitGroup.Done()

}

func updateRoutine(itScopeJson StructSet.ITScopeInfo, index int, res http.ResponseWriter) {

	Api := ApiSet.Api{}
	pricePackage := PriceSet.Prices{}

	product := itScopeJson.Products[index]

	productId := Api.Get("https://firewallforce.se/wp-json/wc/v3/idbysku?sku=" + product.ManufacturerSKU + "&consumer_key=" + ConsumerKey + "&consumer_secret=" + ConsumerSecret)
	setResponseFields()

	if productId != "0" {
		finalPrice := pricePackage.GetFinalPrice(product.SupplierPriceInfo.Price)
		WooProduct := StructSet.WoocommerceInfo{
			FieldsInResponse: fieldsInResponse,
			Price:            finalPrice,
			RegularPrice:     finalPrice,
			Type:             "simple",
			ManageStock:      true,
			StockQuantity:    *product.AggregatedStock,
			StockStatus:      *product.AggregatedStockStatusText,
			//Attributes: []struct{ID string "jsonid"  Name string "json:\"name\""; Options []string "json:\"options\""}{},
		}

		wooCommerceJson, _ := json.Marshal(WooProduct)
		wooResponse := Api.Post("https://firewallforce.se/wp-json/wc/v3/products/"+productId+"?"+"consumer_key="+ConsumerKey+"&consumer_secret="+ConsumerSecret, string(wooCommerceJson))

		finalInfo = append(finalInfo, wooResponse)
		awaitUpdate.Done()
	}

}

func setResponseFields() {
	fieldsInResponse[0] = "id"
	fieldsInResponse[1] = "sku"
	fieldsInResponse[2] = "price"
}
func Validate(itScopeJson StructSet.ITScopeInfo, res http.ResponseWriter) bool {

	for _, product := range itScopeJson.Products {

		if product.AggregatedStock == nil {

			sendError("Missing Attribute", "Aggregated Stock is Missing", "send Valid Json", res)
			return false
		}

		if product.AggregatedStockStatusText == nil {
			sendError("Missing Attribute", "Aggregated Stock Status text is Missing", "send Valid Json", res)
			return false
		} else {
			if strings.Contains(*product.AggregatedStockStatusText, "Not") {
				*product.AggregatedStockStatusText = "outofstock"
			} else {
				*product.AggregatedStockStatusText = "instock"
			}
		}

		// if product.SupplierStockInfo == nil {
		// 	sendError("Missing Attribute", "supplierStockInfo is Missing", "send Valid Json", res)
		// 	return false
		// }

		if product.SupplierPriceInfo == nil {
			sendError("Missing Attribute", "supplierPriceInfo is Missing", "send Valid Json", res)
			return false
		}

		if product.ProductSubType == nil {
			sendError("Missing Attribute", "productSubType is Missing", "send Valid Json", res)
			return false
		}

		if product.StandardHtmlDatasheet == nil {
			sendError("Missing Attribute", "StandardHtmlDatasheet is Missing", "send Valid Json", res)
			return false
		}
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
