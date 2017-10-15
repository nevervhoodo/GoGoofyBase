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

/*
func signalCatch(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			closedb()
		}
	}()
}
*/

//User reuests handlers
func handleRequests() {

	router := httprouter.New()
	router.GET("/v1/records", returnTable)
	router.GET("/", homePage)
	router.GET("/v1/records/:id", returnSingleRecord)
	router.POST("/v1/records", addRecord)
	router.GET("/v1/update/:id/:val", updateRecord)
	router.GET("/v1/reset", resetTable)
	router.GET("/v1/delete/:id", deleteRecord)
	http.ListenAndServe(":8080", router)

}

func main() {
	fmt.Println("Hello my dummy users")
	Opendb()
	handleRequests()
}
