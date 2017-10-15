package main

import (
	"fmt"
	"io/ioutil"
)

//One base record struct
type Field struct {
	ID int `json:"Id"` //id for enumeration of records
	Key string `json:"Key"` //key field for database
	Value string `json:"value"` //value
}

//Array of records
type Table []Field

//Record struct used for adding new key-value pair
type gotField struct {
	Key string `json:"Key"`
	Value string `json:"value"`
}

//Table initialization realisation and declaration
func initTable()(Table){
	strs, err := File2lines(config_dbpath)
	if err != nil{
		fmt.Print("db error ")
		table := Table{
			Field{ID: 0, Key: "test", Value: "testme"},
		}
		var str string
		str = fmt.Sprintf(config_format+"\n",0,"test","testme")
		ioutil.WriteFile(config_dbpath, []byte(str), 0644)
		return table
	}
	if len(strs) == 0{
		table := Table{
			Field{ID: 0, Key: "test", Value: "testme"},
		}
		var str string
		str = fmt.Sprintf(config_format+"\n",0,"test","testme")
		ioutil.WriteFile(config_dbpath, []byte(str), 0644)
		return table
	}
	var i,id int
	var key,val string
	var table Table
	for i=0; i<len(strs); i++{
		//var validID = regexp.MustCompile(`"^([a-z]+): (['a'-'z'],['A'-'Z']+)=(['a'-'z'],['A'-'Z']+)$"`)
		n, err := fmt.Sscanf(strs[i], config_format, &id, &key, &val)
		fmt.Println(strs[i])
		//fmt.Println(validID.MatchString(strs[i]))
		if n < 3 || err!= nil{
			fmt.Printf("problems %d: %s = %s\n", id, val)
		}
		table = append(table,Field{ID: id, Key: key, Value: val})
	}
	return table
}

var mytable Table = initTable()

func (table Table) reset(){
	table = append(table[:0], table[len(table):]...)
}

func (table Table) delete (ires int) bool{
	for i, _ := range table{
		if i == ires{
			table = append(table[:i], table[i+1:len(table)-1]...)
			return true
		}
	}
	return true
}

//Find record in table by its id
func (table Table) searchById(id int)(Field, int){

	if len(table) == 0 {
		return Field{0,"",""}, -1
	}
	var i int
	for i = 0; i<len(table); i++{
		if table[i].ID == id{
			return table[i], i
		}
	}
	return Field{0,"",""}, -1
}

//Find record in table by its Key name
func (table Table) searchByKey(key string)(Field, int){
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


