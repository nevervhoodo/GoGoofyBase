// 2017
//First database realisation stage

package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

//additional function to catch error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//User reuests handlers
func handleRequests() {

	router := httprouter.New()
	router.GET("/v1/records", returnTable)
	router.GET("/", homePage)
	router.GET("/v1/records/:id", returnSingleRecord)
	router.POST("/v1/records", addRecord)
	//router.PUT("/v1/records/:id", updateRecord)
	//router.DELETE("/api/v1/records/:id", deleteRecord)
	http.ListenAndServe(":8080", router)

}

func main() {
	fmt.Println("Hello my dummy users")
	handleRequests()
}
