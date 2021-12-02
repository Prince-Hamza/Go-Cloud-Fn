package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	// "math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	//"github.com/gin-gonic/contrib/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"sync"

	ApiSet "Main.go/Core/Api"
	SuperJsonSet "Main.go/Core/SuperJson"
	WooApiSet "Main.go/Core/WoocommerceApi"
	// "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

var reqObjectCount int
var respCount int
var responseArray []string
var strReqObjectCount string

func mainx() {

	r := mux.NewRouter()
	r.HandleFunc("/", UpdateFirewallForce)
	fmt.Println("listening on localhost:5000")
	http.ListenAndServe(":5000", handlers.CORS()(r))

}

var wg sync.WaitGroup = sync.WaitGroup{}

func awaitThreading(respWriter http.ResponseWriter, req *http.Request) {
	wg.Add(1)
	go routine()
	wg.Add(1)
	go routine()
	wg.Add(1)
	go routine()

	wg.Wait()
	fmt.Println("End :) ")

}

func routine() {
	println("go routine")
	wg.Done()
}

func UpdateFirewallForce(respWriter http.ResponseWriter, req *http.Request) {

	// wooApi := WooApiSet.WooApi{Url: "firewallforce.se", ConsumerKey: "ck_42a75ce7a233bc1e341e33779723c304e6d820cc", ConsumerSecret: "cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526"}
	Cors(respWriter, req)
	fmt.Println("headers Set")

	type Resp struct {
		RespArray string
	}

	super := SuperJsonSet.SuperJson{}
	JsonString := super.Stringify(req.Body)
	JsonType := super.JsonStringToJson(JsonString)
	JsonMap := JsonType["products"]

	readJsonArray(JsonMap)

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				fmt.Println("req object count : ", reqObjectCount)
				fmt.Println("response array count : ", respCount)
				// fmt.Println("response array : ", responseArray)

				if reqObjectCount == respCount && reqObjectCount != 0 {
					fmt.Println("existing Products and responses are equal now")
					fmt.Println("Ticker stopped")
					fmt.Println("send Response")

					type Resp struct {
						RespArray []string
					}
					respStruct := Resp{RespArray: responseArray}
					json.NewEncoder(respWriter).Encode(respStruct)
					return

					// ticker.Stop()
					// done <- true

				}
			}
		}
	}()

	// send Response

}

func Cors(respWriter http.ResponseWriter, req *http.Request) {
	enableCors(&respWriter)

	if req.Method == http.MethodOptions {
		respWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		respWriter.Header().Set("Access-Control-Allow-Headers", "Authorization")
		respWriter.Header().Set("Access-Control-Allow-Methods", "POST")
		respWriter.Header().Set("Access-Control-Allow-Origin", "*")
		respWriter.Header().Set("Access-Control-Max-Age", "3600")
		respWriter.Header().Set("Content-Type", "application/json")
		respWriter.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	respWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	respWriter.Header().Set("Access-Control-Allow-Origin", "*")
	respWriter.Header().Set("Content-Type", "Application/Json")

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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

		go idBySku(obj)

	}

}

func idBySku(Product map[string]interface{}) {

	wooApi := WooApiSet.WooApi{Url: "", ConsumerKey: "ck_42a75ce7a233bc1e341e33779723c304e6d820cc", ConsumerSecret: "cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526"}
	super := SuperJsonSet.SuperJson{}

	Sku := Product["manufacturerSKU"].(string)

	// get id

	id := wooApi.IdBySku(Sku , wg)

	if id != 0 {

		// bundle information

		reqObjectCount += 1
		strReqObjectCount = strconv.Itoa(reqObjectCount)

		var finalPrice string
		var price interface{}

		priceInfo, ok := Product["supplierPriceInfo"]
		if ok {
			price = super.Nested(priceInfo, "price")
			finalPrice = getFinalPrice(price.(string))
		} else {
			priceInfo := Product["productPriceInfo"]
			price = super.Nested(priceInfo, "price")
			finalPrice = getFinalPrice(price.(string))
		}

		var stockQuantity string
		var stockStatus string
		var nestedObjectStockInfo interface{}

		nestedObjectStockInfo, ok = Product["supplierStockInfo"]

		if ok {
			stock_interface := super.Nested(nestedObjectStockInfo, "stock")
			stockQuantity = stock_interface.(string)
			status_interface := super.Nested(nestedObjectStockInfo, "stockStatusText")
			stockStatus = status_interface.(string)

		} else {
			nestedObjectStockInfo, ok := Product["productStockInfo"]
			if ok {
				stock_interface := super.Nested(nestedObjectStockInfo, "stock")
				stockQuantity = stock_interface.(string)
				status_interface := super.Nested(nestedObjectStockInfo, "stockStatusText")
				stockStatus = status_interface.(string)
			} else {
				stock, ok := Product["aggregatedStock"]
				if ok {
					stockQuantity = stock.(string)
				}
				status, ok := Product["aggregatedStockStatusText"]
				if ok {
					stockStatus = status.(string)
				}
			}
		}

		// stockQuantity := super.Nested(stockInfo, "stock")
		// roundStockQuantity = math.Round(stockQuantity.(float64))
		// strStockQuantity := strconv.FormatFloat(stockQuantity.(float64), 'f', 6, 64)

		stockStatusText := stockStatus

		if strings.Contains(stockStatusText, "Not") {
			stockStatusText = "outofstock"
		} else {
			stockStatusText = "instock"
		}

		Brand := Product["productSubType"]

		fmt.Println("id", id)
		fmt.Println("price", finalPrice)
		fmt.Println("stock Quantity", stockQuantity)
		fmt.Println("stock Status", stockStatusText)
		fmt.Println("Brand", Brand)

		// prepare json string

		productJsonString := `{ "price":` + `"` + finalPrice + `"` + "," +
			`"regular_price":` + `"` + finalPrice + `"` + "," +
			`"manage_stock":"true" ` + "," +
			`"stock_quantity":` + `"` + stockQuantity + `"` + "," +
			`"stock_status":` + `"` + stockStatusText + `"` + "," +
			//`"attributes": [{"id":4 , "name":"Brands" , "options":[ ` + Brand.(string) + `]` + `}] }` + "," +
			`"fields_in_response": ["id", "sku", "stock_quantity", "stock_status", "price"]}`

			//	fmt.Println("jsn str", productJsonString)
		stringId := strconv.Itoa(id)
		updateProduct(productJsonString, stringId)

	}

	//Select := wooApi.CreateOrUpdate()
	// sfun(select)  create() || update()
}

func updateProduct(jsonString string, productId string) {
	Api := ApiSet.Api{}
	resp := Api.Post(`https://firewallforce.se/wp-json/wc/v3/Products/`+productId+"?consumer_key=ck_42a75ce7a233bc1e341e33779723c304e6d820cc&consumer_secret=cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526", jsonString)
	responseArray = append(responseArray, resp)
	respCount += 1
}

func getFinalPrice(price string) string {

	intPrice, err := strconv.Atoi(price)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Eur Price : ", price)
	krona := getMemoApiPrice(intPrice)
	fmt.Println("SEK Price : ", krona)
	return krona

	// content := getDateFromFile()
	// currentDate := getCurrentDate()
	// storedDay := getDayFromDate(content)
	// currentDay := getDayFromDate(currentDate)
	// intCurrentDay, err := strconv.Atoi(currentDay)
	// if err != nil {
	// 	println(err)
	// }
	// intStoredDay, err := strconv.Atoi(storedDay)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if intStoredDay == intCurrentDay {
	// 	getStoredPrice()
	// } else {
	// 	getApiPrice()
	// }

	// return "currentTime.String()"

}

func getMemoApiPrice(inputEur int) string {
	apiEurRate := 0.89173
	apiSekRate := 9.072
	cnvbase := 1 / apiEurRate // Eur * by cnvbase = usd value
	USdollar := float64(inputEur) * cnvbase
	SEK := apiSekRate
	Krona := USdollar * SEK
	str := strconv.FormatFloat(Krona, 'f', 2, 64)
	return str
}

func getStoredPrice() string {
	fileData, err := ioutil.ReadFile("./priceStore.txt")
	if err != nil {
		fmt.Println(err)
	}
	content := string(fileData)
	list := strings.SplitAfter(content, ";")

	return list[0]
}

func getApiPrice() {
	api := ApiSet.Api{}
	resp := api.Get("'https://currencyapi.net/api/v1/rates?key=McRbxJQKvXlfe5D6EHIv2Q8qtSxTD37zEq9m&output=JSON")

	fmt.Println(resp)
	// overwrite

}

func getCurrentDate() string {
	now := time.Now()
	splity := strings.SplitAfter(now.String(), " ")
	return splity[0]
}

func getDayFromDate(date string) string {
	dateCompact := strings.Join(strings.Fields(strings.TrimSpace(date)), " ")
	dateparts := strings.Split(dateCompact, "-")
	return dateparts[2]
}

func getDateFromFile() string {
	fileData, err := ioutil.ReadFile("./priceStore.txt")
	if err != nil {
		fmt.Println(err)
	}
	content := string(fileData)
	list := strings.SplitAfter(content, ";")

	return list[1]
}

