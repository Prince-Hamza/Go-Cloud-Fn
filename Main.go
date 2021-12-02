package main
import (
	//"encoding/json"
	ITScopeProduct "Main.go/ITScopeProduct"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", ITScopePro)
	fmt.Println("listening on localhost:5000")
	http.ListenAndServe(":5000", handlers.CORS()(r))

}

func ITScopePro(res http.ResponseWriter, req *http.Request) {
	itScopePro := ITScopeProduct.ITScopePro{}
	itScopePro.ParseItScopeProduct(res, req)
}


