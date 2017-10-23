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
	//requestChannel <- "List"
	fmt.Println("Endpoint Hit: returnTable")
	fmt.Println(global)
	mytable.Print()
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
	if len(mytable) == 0 {
		fmt.Fprintf(w, "<h1>Empty table</h1>")
	} else{
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(mytable); err != nil {
			w.WriteHeader(501)
		}
	}
	//requestChannel <- "Table"

}

func returnAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Println("all")
	checkcococo()
	if len(mytable2) == 0 {
		fmt.Fprintf(w, "<h1>Empty table</h1>")
	} else{
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(mytable2); err != nil {
			w.WriteHeader(501)
		}
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

	//msg := <- requestChannel
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

func findRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	checkcococo()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)
	fmt.Print("val",id,ok)
	if !ok {
		rec, ires := mytable2.searchByKey(ps.ByName("id"))
		fmt.Println(rec,ires)
		if ires == -1{
			json.NewEncoder(w).Encode("No record ith that key")
		} else {
			json.NewEncoder(w).Encode(rec)
		}
	} else {
		rec, ires := mytable2.searchById(id)
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
	fmt.Print("Update val ",id,ok)

	got_val := ps.ByName("val")
	//search for item. Ires = line with item
	if !ok {
		rec, ires := mytable.searchByKey(ps.ByName("id"))
		fmt.Println(rec,ires)
		if ires == -1 {
			json.NewEncoder(w).Encode("No record with that key")
		} else {
			json.NewEncoder(w).Encode("Record found")
			mytable = mytable.updateRecord(ires, got_val)
			if !changeLine(ires, got_val){
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
			mytable = mytable.updateRecord(ires, got_val)
			mytable.Print()
			if !changeLine(ires, got_val){
				json.NewEncoder(w).Encode("Error")
			}
		}
	}
}

func updateRecordInPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	checkcococo()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)

	got_val := ps.ByName("val")
	//search for item. Ires = line with item
	if !ok {
		rec, ires := mytable2.searchByKey(ps.ByName("id"))
		fmt.Println(rec,ires)
		if ires == -1 {
			json.NewEncoder(w).Encode("No record with that key")
		} else {
			json.NewEncoder(w).Encode("Record found")
			mytable2 = mytable2.updateRecord(ires, got_val)
			store.Print()
			store = store.updatePage(mytable2.getKeyID(ps.ByName("id")), got_val)
			store.Print()
		}
	} else {
		rec, ires := mytable2.searchById(id)
		fmt.Println(rec,ires)
		if ires == -1 {
			json.NewEncoder(w).Encode("No value with that id")
		} else {
			json.NewEncoder(w).Encode("Record found")
			mytable2 = mytable2.updateRecord(ires, got_val)
			mytable2.Print()
			store.Print()
			store = store.updatePage(id, got_val)
			store.Print()
		}
	}
}

//Clear table and database
func resetTable(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	mytable = mytable.reset()
	emptydb()
}

func reset2Table(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	checkcococo()
	mytable2 = mytable2.reset()
	store = store.emptyPages()
}

//Delete one record from database
func deleteRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	fmt.Println("dalete")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)
	fmt.Println("val",id,ok)

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
		fmt.Println("DELETE")
		mytable.Print()
		ok, mytable = mytable.delete(ires)
		if ok {
			removeline(ires)
		}

		//if !changeLine(ires, ps.ByName("val")){
		//	json.NewEncoder(w).Encode("Error")
		//}
	}

}

func delete2Record(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	checkcococo()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id, ok := getID(w, ps)

	//var rec Field
	var ires int
	var ID int
	if !ok {
		_, ires = mytable2.searchByKey(ps.ByName("id"))
		ID = mytable2.getKeyID(ps.ByName("id"))
	} else {
		_, ires = mytable2.searchById(id)
		ID = id
	}

	if ires == -1 {
		if !ok {
			json.NewEncoder(w).Encode("No record with that key")
		} else{
			json.NewEncoder(w).Encode("No record with that id")
		}
	} else {
		json.NewEncoder(w).Encode("Record found")
		fmt.Println("DELETE")
		mytable2.Print()
		ok, mytable2 = mytable2.delete(ires)
		if ok {
			store = store.rmFromPage(ID)
		}

		//if !changeLine(ires, ps.ByName("val")){
		//	json.NewEncoder(w).Encode("Error")
		//}
	}

}

//Add new record to table
func addRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//parse got parametrs into gotField structure
	decoder := json.NewDecoder(r.Body)
	var t gotField
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()
	log.Println(t.Key, t.Value,"ok")

	mytable.Print()
	var ok int
	_, ok = mytable.searchByKey(t.Key)
	if ok != -1 {
		json.NewEncoder(w).Encode("Record with that key exists")
	} else {

		_, mytable = mytable.updateTable(t, false)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode("added")
	}
}

func addRecord2Page(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	checkcococo()
	decoder := json.NewDecoder(r.Body)
	var t gotField
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()
	log.Println(t.Key, t.Value,"ok")

	mytable2.Print()
	var ok int
	_, ok = mytable2.searchByKey(t.Key)
	if ok != -1 {
		json.NewEncoder(w).Encode("Record with that key exists")
	} else {

		_, mytable2 = mytable2.updateTable(t, true)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	//requestChannel <- "Home"
}
