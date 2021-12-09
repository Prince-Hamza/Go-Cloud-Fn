package Api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Api struct{}

func (api Api) Get(url string) string {

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("GET:", url)
	//fmt.Println("Response:", string(body))

	return string(body)

}

func (api Api) Post(url string, jsonString string) string {

	//fmt.Println("POST : ", url)

	json_data := []byte(jsonString)
	fmt.Println(bytes.NewBuffer(json_data))

	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(string(body))
	return string(body)

}

func (api Api) sendJsonResp(w http.ResponseWriter, r *http.Request) {
	//Api.postApi()

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

func (api Api) ParseReqJson(reqBody io.ReadCloser) interface{} {
	var target interface{}
	body, _ := ioutil.ReadAll(reqBody)
	json.Unmarshal(body, &target)
	return target
}

func (api Api) JsonToMap() map[string]interface{} {
	coronaVirusJSON := `{
        "name" : "covid-11",
        "country" : "China",
        "city" : "Wuhan",
        "reason" : "Non vedge Food"
    }`

	// Declared an empty map interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(coronaVirusJSON), &result)
	//  fmt.Println(result)

	return result
}
