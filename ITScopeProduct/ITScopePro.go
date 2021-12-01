package ITScopeProduct

import (
	"encoding/json"
	"fmt"
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
	//PriceSet "Main.go/ITScopeProduct/Prices"
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

}

var productsWithIds []interface{}

func parallelUpdate(itScopeJson StructSet.ITScopeInfo, res http.ResponseWriter) {
	Api := ApiSet.Api{}

	for _, product := range itScopeJson.Products {
		waitGroup.Add(1)
		productId := Api.Get("https://firewallforce.se/wp-json/wc/v3/idbysku?sku=" + product.ManufacturerSKU + "&consumer_key=" + ConsumerKey + "&consumer_secret=" + ConsumerSecret)
		product.Id = productId
		productsWithIds = append(productsWithIds, product)
	}

	fmt.Println("products With Ids", productsWithIds)

}

func bundleInformation() {

	// prices := PriceSet.Prices{}
	// prices.GetFinalPrices()

	// for each id :
	//        price.getPrices()
	//        stock.getStock()
	//        attribs.getAttr()
	//        categories.getCategories()
	//        images.getImages()
	// Post to Woocommerce
	// sendResp("success", res)

	waitGroup.Done()

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
}

func sendResp(resp string, res http.ResponseWriter) {
	fmt.Println(resp)
	respStruct := Response{Resp: resp}
	json.NewEncoder(res).Encode(respStruct)
	return
}
