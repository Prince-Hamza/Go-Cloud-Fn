//package p

package main

import (
	ApiSet "Main.go/Core/Api"
	SuperJsonSet "Main.go/Core/SuperJson"
	WooApiSet "Main.go/Core/WoocommerceApi"
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

//Api.Get("https://firewallforce.se/wp-json/wc/v3/products/44780?consumer_key=ck_42a75ce7a233bc1e341e33779723c304e6d820cc&consumer_secret=cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526&_fields=id,sku")

func main() {
	println("Main")
	ctx := context.Background()
	funcframework.RegisterHTTPFunctionContext(ctx, "/", updateFirewallForce)
	funcframework.Start("8080")
}

func updateFirewallForce(res http.ResponseWriter, req *http.Request) {

	// wooApi := WooApiSet.WooApi{Url: "firewallforce.se", ConsumerKey: "ck_42a75ce7a233bc1e341e33779723c304e6d820cc", ConsumerSecret: "cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526"}

	super := SuperJsonSet.SuperJson{}
	JsonString := super.ReqBodyToString(req.Body)

	// make json string readable

	JsonType := super.JsonStringToJson(JsonString)
	JsonMap := JsonType["products"]
	readJsonArray(JsonMap)

	// get id by sku

	// Bundle Information

	// update Product

}

func readJsonArray(JsonMap interface{}) {

	JsonArray, ok := JsonMap.([]interface{})
	if !ok {
		log.Fatal("expected an array of objects")
	}
	for i, obj := range JsonArray {
		obj, ok := obj.(map[string]interface{})
		if !ok {
			log.Fatalf("expected type map[string]interface{}, got %s", reflect.TypeOf(JsonArray[i]))
		}
		// priceInfo := obj["productPriceInfo"]
		// price := super.Nested(priceInfo, "price")
		// fmt.Println(price)
		go idBySku(obj)

	}

}

func idBySku(Product map[string]interface{}) {
	wooApi := WooApiSet.WooApi{Url: "", ConsumerKey: "ck_42a75ce7a233bc1e341e33779723c304e6d820cc", ConsumerSecret: "cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526"}
	Sku := Product["manufacturerSKU"].(string)
	id := wooApi.IdBySku(Sku)
	fmt.Println(id)
}

func HelloWorld(res http.ResponseWriter, req *http.Request) {

	Api := ApiSet.Api{}

	productId := Api.Get("https://firewallforce.se/wp-json/wc/v3/idbysku?sku=R0M67A")
	fmt.Println(`productId: `, productId)

	currentTime := time.Now()
	timeF12 := currentTime.Format("2017-09-07 2:3:5 PM")
	fmt.Println("Short Hour Minute Second: ", timeF12)

	Api.Post(`https://firewallforce.se/wp-json/wc/v3/Products/`+productId+"?consumer_key=ck_42a75ce7a233bc1e341e33779723c304e6d820cc&consumer_secret=cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526", `{"name":"HP Enterprise ROM67A"}`)

	return
}
