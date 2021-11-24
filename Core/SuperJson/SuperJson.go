package SuperJson

import (
	"encoding/json"
	"fmt"
	"io"

	//"io/ioutil"
	//"fmt"
	"net/http"
)

type SuperJson struct{}

func (super SuperJson) SendJsonResp(w http.ResponseWriter, r *http.Request) {
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

// func (super SuperJson) ParseReqJson(reqBody io.ReadCloser) interface{} {
// 	var target interface{}
// 	body, _ := ioutil.ReadAll(reqBody)
// 	json.Unmarshal(body, &target)
// 	return target
// }

func (super SuperJson) ReqBodyToString(body io.ReadCloser) string {
	b, err := io.ReadAll(body)
	if err != nil {
		print(err)
	}
	return string(b)
}

func (super SuperJson) JsonToMap(jsonString string) map[string]interface{} {

	// Declared an empty map interface
	//bird := `{"birdList" : ["sparrow", "parrot"]}`
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(jsonString), &result)
	//fmt.Println("result", result)

	return result

}

func (super SuperJson) ReqBodyToMap(body io.ReadCloser) map[string]interface{} {
	jsnString := super.ReqBodyToString(body)
	mapify := super.JsonToMap(jsnString)
	return mapify
}

func (super SuperJson) JsonStringToIntMap() {
	type employee struct {
		Name string
	}

	j := `{"Name":"John"}`
	var b map[int]employee
	json.Unmarshal([]byte(j), &b)
	fmt.Println(b)
}



func (super SuperJson) JsonStringToJson(jsonString string) map[string]interface{} {

	// in := []byte(`{ "votes": { "option": "3" } , "c2":"7" , "list":["A" , "B" , "C"] }`)
	in := []byte(jsonString)
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(in, &jsonMap); err != nil {
		panic(err)
	}

	//opt := raw["votes"].(map[string]interface{})["option"]
	//opt := product.Nested(raw["votes"], "option")

	return jsonMap
	//fmt.Println(opt)

}

func (super SuperJson) Nested(input interface{}, attribute string) interface{} {
	x := input.(map[string]interface{})[attribute]
	return x
}
