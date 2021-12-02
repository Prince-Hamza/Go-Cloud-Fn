package Price

import (
	ApiSet "Main.go/Core/Api"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Prices struct{}

//  fun finalize ((int x , int y) *)
// hard coded
func (prices Prices) GetFinalPrice(price string) string {

	intPrice, err := strconv.Atoi(price)
	if err != nil {
		fmt.Println(err)
	}

	krona := getMemoApiPrice(intPrice)
	return krona

}

func getMemoApiPrice(inputEur int) string {
	// hard coded
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
