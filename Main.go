//package p
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"io/ioutil"
	//"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	funcframework.RegisterHTTPFunctionContext(ctx, "/", HelloWorld)
	funcframework.Start("8080")
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {

	type Bird struct {
		Species     string
		Description string
	}

	birdJson := `{"species": "pig","description": "likes to perch on rocks"}`
	var bird Bird
	json.Unmarshal([]byte(birdJson), &bird)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bird)

}

func postApi() {

	httpposturl := "https://firewallforce.se/wp-json/wc/v3/idbysku?"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"sku": "morpheus",
	}`)

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

}
