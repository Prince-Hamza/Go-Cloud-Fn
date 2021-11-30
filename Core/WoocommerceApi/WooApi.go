package WooApi

import (
	//	"fmt"
	"fmt"
	"log"
	"strconv"
	"sync"

	ApiSet "Main.go/Core/Api"
	SuperJson "Main.go/Core/SuperJson"
)

type WooApi struct {
	Url            string
	ConsumerKey    string
	ConsumerSecret string
}

var woo WooApi

func (wooApi WooApi) ApiFlow(productJson string) {
	superJson := SuperJson.SuperJson{}
	productMap := superJson.JsonToMap(productJson)
	fmt.Println(productMap)

	// wooApi.IdBySku(productMap["manufacturerSKU"])
}

func (wooApi WooApi) IdBySku(sku string , wg sync.WaitGroup) int {
	Api := ApiSet.Api{}
	productId := Api.Get("https://firewallforce.se/wp-json/wc/v3/idbysku?sku=" + sku + "&consumer_key=" + wooApi.ConsumerKey + "&consumer_secret=" + wooApi.ConsumerSecret)
	intId, err := strconv.Atoi(productId)
	if err != nil {
		log.Fatal(err)
	}
	return intId
}

func (wooApi WooApi) CreateOrUpdate(productsJson string) {

	//products := productsJson["products"]

	//fmt.Print(products)

	// wooApi.IdBySku(productMap["manufacturerSKU"])
}

// func (wooApi WooApi) CreateProduct (productJson string) interface {

// }

// func (wooApi WooApi) UpdateProduct (productJson string) interface {

// }
