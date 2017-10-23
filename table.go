package main

import (
	"fmt"
	"sync"
)

//One base record struct
type Field struct {
	ID int `json:"Id"` //id for enumeration of records
	Key string `json:"Key"` //key field for database
	Value string `json:"value"` //value
}

//Array of records
type Table []Field
var muForTable = &sync.RWMutex{}


//Record struct used for adding new key-value pair
type gotField struct {
	Key string `json:"Key"`
	Value string `json:"value"`
}

//Table initialization realisation and declaration
func InitTable()(Table){
	strs, err := File2lines(config_dbpath)

	//create empty table
	table := Table {}
	//catch error while opening
	if err != nil{
		fmt.Print("db error",err)
		return table
	}

	var id int
	var key,val string
	var counter int = 0
	for i, str := range strs {
		n, err := fmt.Sscanf(str, config_format, &id, &key, &val)
		fmt.Println(str)
		if n < 3 || err!= nil{
			fmt.Printf("cannot parse line %d\n", i)
		} else {
			counter += 1
			table = append(table, Field{ID: id, Key: key, Value: val})
		}
	}

	//if table is empty
	if counter == 0{
		fmt.Print("empty db\n")
	}

	return table
}

func CollectTable() Table {
	fulltable := Table {}
	fmt.Println(store)
	for _, page := range (store){
		fmt.Println(page)
		addFields := page.Page2Table()
		fmt.Println(addFields)
		fulltable = append(fulltable, addFields...)
	}
	return fulltable
}

func (table Table) Print() {
	fmt.Printf("\n###########################################\n")
	fmt.Printf("Current table state\n")
	fmt.Printf("%#v \n", table)
	fmt.Printf("###########################################\n\n")
}

func (table Table) reset() Table{
	muForTable.Lock()
	table = append(table[:0], table[len(table):]...)
	muForTable.Unlock()
	return table
}

func (table Table) delete (ires int) (bool, Table){
	muForTable.Lock()
	defer muForTable.Unlock()
	for i, _ := range table{
		if i == ires{
			table = append(table[:i], table[i+1:len(table)]...)
			//table[len(table)]
			return true, table
		}
	}
	return false, table
}

//Find record in table by its id
func (table Table) searchById(id int)(Field, int){
	muForTable.RLock()
	defer muForTable.RUnlock()
	fmt.Println("search by id ", id)
	if len(table) == 0 {
		return Field{0,"",""}, -1
	}
	var i int
	for i = 0; i<len(table); i++{
		if table[i].ID == id{
			fmt.Println("try ", i, " : ", table[i].ID)
			return table[i], i
		}
	}
	return Field{0,"",""}, -1
}

//Find record in table by its Key name
func (table Table) searchByKey(key string)(Field, int){
	muForTable.RLock()
	defer muForTable.RUnlock()
	if len(table) == 0 {
		return Field{0,"",""}, -1
	}
	var i int
	for i = 0; i<len(table); i++{
		if table[i].Key == key{
			return table[i], i
		}
	}
	return Field{0,"",""}, -1
}

func (table Table) getKeyID (key string) int {
	muForTable.RLock()
	defer muForTable.RUnlock()
	if len(table) == 0 {
		return -1
	}
	var i int
	for i = 0; i<len(table); i++{
		if table[i].Key == key{
			return table[i].ID
		}
	}
	return -1
}

//add new record to table. Return ID
func (table Table) updateTable(rec gotField, updatePage bool) (int, Table) {
	muForTable.Lock()
	defer muForTable.Unlock()
	//new id = find free one
	id := -1
	for i,_ := range (table){
		if table[i].ID != i{
			id = i
			break
		}
	}
	if id == -1{
		id = len(table)
	}
	fmt.Println("add record",rec.Key,len(table))
	fi := Field{ID: id, Key: rec.Key, Value: rec.Value}
	fmt.Printf("%#v \n", fi)
	table = append(table, fi)
	fmt.Println("added",len(table))
	fmt.Printf("%#v \n", table)
	//update db file
	fstr := fmt.Sprintf(config_format+"\n", id, rec.Key, rec.Value)
	if !updatePage{
		fmt.Println(id, fstr)
		InsertStringToFile(config_dbpath, fstr, id)
	} else {
		store = store.InsertStringToPage(fstr, id)
	}
	return id, table
}

func (table Table) updateRecord(line int, val string) Table{
	muForTable.Lock()
	defer muForTable.Unlock()
	table[line].Value = val
	return table
}


