//package p
package main



import (
	ccc "Main.go/Basic"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {


	// bsc.Mnemo()
	ccc.Mnemo()

	println("Main")
	ctx := context.Background()
	funcframework.RegisterHTTPFunctionContext(ctx, "/", HelloWorld)
	funcframework.Start("8080")
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {

	println("Hello World")
	Api := Book{"firewallforce", ".se"}
	Api.postApi()

	type Bird struct {
		Species     string
		Description string
	}

	birdJson := `{"species": "pigeoness","description": "likes to perch on rocks"}`
	var bird Bird
	json.Unmarshal([]byte(birdJson), &bird)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bird)

}

type Book struct {
	url string
	ext string
}

func (book Book) postApi() {

	values := map[string]string{"Sku": "J9780A"}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("https://firewallforce.se/wp-json/wc/v3/idbysku", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// var res map[string]interface{}
	// json.NewDecoder(resp.Body).Decode(&res)
	// fmt.Println(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	fmt.Println(string(body))

}

func (book Book) getApi() {

	println(`GET Api`)

	resp, err := http.Get("https://firewallforce.se/wp-json/wc/v3/products/44780?consumer_key=ck_42a75ce7a233bc1e341e33779723c304e6d820cc&consumer_secret=cs_6e5a683ab5f08b62aa1894d8d2ddc4ad69ff0526")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

}
