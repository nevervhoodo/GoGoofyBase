package main

import (
	"fmt"
    //"regexp"
	"net/http"
	"encoding/json"
	//"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"io/ioutil"
	"os"
	"io"
	"bufio"
	"log"
)
// this is a comment

const dbpath string = "/tmp/dat"
const format string = "%d: %s = %s"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Field struct {
	ID int `json:"Id"`
	Key string `json:"Key"`
	Value string `json:"value"`
}

type gotField struct {
	Key string `json:"Key"`
	Value string `json:"value"`
}

type Table []Field

func initTable()(Table){

	strs, err := File2lines(dbpath)
	if err != nil{
		fmt.Print("db error")
		table := Table{
			Field{ID: 0, Key: "test", Value: "testme"},
		}
		var str string
		str = fmt.Sprintf(format+"\n",0,"test","testme")
		ioutil.WriteFile(dbpath, []byte(str), 0644)
		return table
	}
	if len(strs) == 0{
		table := Table{
			Field{ID: 0, Key: "test", Value: "testme"},
		}
		var str string
		str = fmt.Sprintf(format+"\n",0,"test","testme")
		ioutil.WriteFile(dbpath, []byte(str), 0644)
		return table
	}
	var i,id int
	var key,val string
	var table Table
	for i=0; i<len(strs); i++{
		//var validID = regexp.MustCompile(`"^([a-z]+): (['a'-'z'],['A'-'Z']+)=(['a'-'z'],['A'-'Z']+)$"`)
		n, err := fmt.Sscanf(strs[i], format, &id, &key, &val)
		fmt.Println(strs[i])
		//fmt.Println(validID.MatchString(strs[i]))
		if n < 3 || err!= nil{
			fmt.Printf("problems %d: %s = %s\n", id, val)
		}
		table = append(table,Field{ID: id, Key: key, Value: val})
	}
	return table
}

func searchById(table Table, id int)(Field, bool){

	if len(table) == 0 {
		return Field{0,"",""}, false
	}
	var i int
	for i = 0; i<len(table); i++{
		if table[i].ID == id{
			return table[i], true
		}
	}
	return Field{0,"",""}, false
}

func searchByKey(table Table, key string)(Field, bool){
	if len(table) == 0 {
		return Field{0,"",""}, false
	}
	var i int
	for i = 0; i<len(table); i++{
		if table[i].Key == key{
			return table[i], true
		}
	}
	return Field{0,"",""}, false
}

var mytable Table = initTable()

func returnTable(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	fmt.Println("Endpoint Hit: returnTable")

	// json.NewEncoder(w).Encode(table)

	var str string
	if len(r.URL.RawQuery) > 0 {
		str = r.URL.Query().Get("Key")
		if str == "" {
			w.WriteHeader(400)
			return
		}
		fmt.Println(str)
		fmt.Println(json.NewEncoder(w).Encode(mytable))
	}
	/*recs, err := read(str)
	if err != nil {
		w.WriteHeader(500)
		return
	}*/
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(mytable); err != nil {
		w.WriteHeader(501)
	}

}

func getID(w http.ResponseWriter, ps httprouter.Params) (int, bool) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		return 0, false
	}
	return id, true
}

func getKey(w http.ResponseWriter, ps httprouter.Params) (string, bool){
	return ps.ByName("id"), true
}

func returnSingleRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params){

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)
	fmt.Print("val",id,ok)
	if !ok {
		rec, boolres := searchByKey(mytable, ps.ByName("id"))
		fmt.Println(rec,boolres)
		if !boolres {
			json.NewEncoder(w).Encode("No record ith that key")
		} else {
			json.NewEncoder(w).Encode(rec)
		}
	} else {
		rec, boolres := searchById(mytable, id)
		fmt.Println(rec,boolres)
		if !boolres {
			json.NewEncoder(w).Encode("No value with that id")

		} else {
			json.NewEncoder(w).Encode(rec)
		}
	}



}

func addRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var t gotField
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()
	log.Println(t.Key,t.Value,"ok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var ok bool
	_, ok = searchByKey(mytable,t.Key)
	if ok {
		json.NewEncoder(w).Encode("Record ith that key exists")
	} else {
		num := len(mytable)
		fstr := fmt.Sprintf(format+"\n", num, t.Key, t.Value)
		fmt.Println(num, fstr)
		InsertStringToFile(dbpath, fstr, num)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode("added")
	}
}


func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	//http.HandleFunc("/", )
	fmt.Println("Endpoint Hit: homePage")
}


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

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}
	var inserted bool = false
	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
			inserted = true
		}
		fileContent += line
		fileContent += "\n"
	}
	if !inserted{
		fileContent += str
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}


func main() {
	fmt.Println("Hello my dummy users")
	//err := InsertStringToFile("/tmp/dat", "hello world!\n", 2)
	//if err != nil {
	//	fmt.Println("Database died")
	//}



	//table := initTable()
	handleRequests()

}
