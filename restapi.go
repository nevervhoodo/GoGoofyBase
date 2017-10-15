package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"log"
)

//Return full table as json response
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

//Special function to get id GET parametr
func getID(w http.ResponseWriter, ps httprouter.Params) (int, bool) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		return 0, false
	}
	return id, true
}

//Special function to get key GET parametr
func getKey(w http.ResponseWriter, ps httprouter.Params) (string, bool){
	return ps.ByName("id"), true
}

//Return one record by id or key as json response
func returnSingleRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params){

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)
	fmt.Print("val",id,ok)
	if !ok {
		rec, ires := mytable.searchByKey(ps.ByName("id"))
		fmt.Println(rec,ires)
		if ires == -1{
			json.NewEncoder(w).Encode("No record ith that key")
		} else {
			json.NewEncoder(w).Encode(rec)
		}
	} else {
		rec, ires := mytable.searchById(id)
		fmt.Println(rec,ires)
		if ires == -1 {
			json.NewEncoder(w).Encode("No value with that id")

		} else {
			json.NewEncoder(w).Encode(rec)
		}
	}
}

//Change record
func updateRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)
	fmt.Print("val",id,ok)

	if !ok {
		rec, ires := mytable.searchByKey(ps.ByName("id"))
		fmt.Println(rec,ires)
		if ires == -1 {
			json.NewEncoder(w).Encode("No record with that key")
		} else {
			json.NewEncoder(w).Encode("Record found")
			if !changeLine(ires, ps.ByName("val")){
				json.NewEncoder(w).Encode("Error")
			}
		}
	} else {
		rec, ires := mytable.searchById(id)
		fmt.Println(rec,ires)
		if ires == -1 {
			json.NewEncoder(w).Encode("No value with that id")

		} else {
			json.NewEncoder(w).Encode("Record found")
			if !changeLine(ires, ps.ByName("val")){
				json.NewEncoder(w).Encode("Error")
			}
		}
	}


}

//Clear table and database
func resetTable(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	mytable.reset()
	emptydb()
}

//Delete one record from database
func deleteRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)
	fmt.Print("val",id,ok)

	var rec Field
	var ires int
	if !ok {
		rec, ires = mytable.searchByKey(ps.ByName("id"))
	} else {
		rec, ires = mytable.searchById(id)
	}

	fmt.Println(rec,ires)
	if ires == -1 {
		if !ok {
			json.NewEncoder(w).Encode("No record with that key")
		} else{
			json.NewEncoder(w).Encode("No record with that id")
		}
	} else {
		json.NewEncoder(w).Encode("Record found")

		mytable.delete(ires)

		if !changeLine(ires, ps.ByName("val")){
			json.NewEncoder(w).Encode("Error")
		}
	}

}

//Add ne record to table
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

	var ok int
	_, ok = mytable.searchByKey(t.Key)
	if ok != -1 {
		json.NewEncoder(w).Encode("Record with that key exists")
	} else {
		num := len(mytable)
		fstr := fmt.Sprintf(config_format+"\n", num, t.Key, t.Value)
		fmt.Println(num, fstr)
		InsertStringToFile(config_dbpath, fstr, num)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode("added")
	}
}


//return HomePage. Print README.MD
func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprintf(w, "<h1>Welcome to the HomePage!<h1>")
	lines, err := File2lines(config_readmepath)
	if err != nil {
		fmt.Println("Problem with Readme")
	} else {
		for _, line := range lines{
			fmt.Fprintf(w, "<div style=\"font-size: 14px\"><p>"+line+"<p></div>")
		}
		fmt.Println("Endpoint Hit: homePage")
	}
}
