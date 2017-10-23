// 2017
//First database realisation stage

package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"time"
	"math/rand"
	//"sync"
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
	router.GET("/v2/records", returnAll)
	router.GET("/", homePage)
	router.GET("/v1/records/:id", returnSingleRecord)
	router.GET("/v2/records/:id", findRecord)
	router.POST("/v1/records", addRecord)
	router.POST("/v2/records", addRecord2Page)
	router.GET("/v1/update/:id/:val", updateRecord)
	router.GET("/v2/update/:id/:val", updateRecordInPage)
	router.GET("/v1/reset", resetTable)
	router.GET("/v2/reset", reset2Table)
	router.GET("/v1/delete/:id", deleteRecord)
	router.GET("/v2/delete/:id", delete2Record)
	http.ListenAndServe(":8484", router)

}
var mytable Table = InitTable()
var mytable2 Table = CollectTable()
//var requestChannel chan string
var global int = 0
var store Storage = addPages(config_dbdir)

var npage int = 0
func cococo () {
	fmt.Println("ticker ", npage)
	var result bool
	store, result = store.checkPage(npage)
	fmt.Println("ticker ", npage, result)
	if (npage == len(store)-1){
		npage = 0
	} else {
		npage += 1
	}
}

func checkcococo(){
	var n int
	rand.Seed(time.Now().Unix())
	n = rand.Intn(len(store))
	i := 0
	for i < n {
		cococo ()
		i += 1
	}
}

func main() {
	fmt.Println("Hello my dummy users")
	Opendb()
	//var muMain = &sync.RWMutex{}
	//muMain.Lock()
	//store, _ := addPages(config_dbdir)
	//muMain.Unlock()
	//mytable = InitTable()
	//requestChannel = make(chan string, 10)
	store.Print()
	//a := []int{0,1,2,3,4,5}
	//fmt.Println(a)
	//i := 2
	//a = append(a[0:i])
	//b := append(a[i:i*2])
	//c := append(a[i*2:i*3])
	//fmt.Println(a,b,c)
	//go func() {
	//	for {
	//		request, ok := <-requestChannel
	//
	//		if !ok{
	//			return
	//		}
	//
	//		fmt.Println("got request: ",request)
	//	}
	//}()
	handleRequests()
}
