package main

import (
	"fmt"
	"io/ioutil"
)

//Table initialization realisation and declaration
func initTable()(Table){
	strs, err := File2lines(dbpath)
	if err != nil{
		fmt.Print("db error ")
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
var mytable Table = initTable()

