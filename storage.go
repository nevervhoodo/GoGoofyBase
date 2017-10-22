package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"io"
	"bufio"
)

//Special subfunc to represent file as an array of strings
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

//Special function to read file and return array of strings
func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

//array with the whole file splitted into lines
var dblines[] string

func Opendb(){
	var er error
	dblines, er = File2lines(config_dbpath)
	if er != nil {
		fmt.Print("db err")
	}
}

//Insert string to file in special place
func InsertStringToFile(path, str string, index int) error {
	fmt.Printf("%#v \n", dblines)
	lines := dblines
	/*
	if err != nil {
		return err
	}
	*/
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

	ioutil.WriteFile(path, []byte(fileContent), 0644)
	var er error
	dblines, er = File2lines(config_dbpath)
	return er
}

func changeLine(line int, got_val string) bool {
	var id int
	var key, val string
	_, err := fmt.Sscanf(dblines[line], config_format, &id, &key, &val)
	if err != nil{
		return false
	}
	if val == got_val{
		return false
	}
	str := fmt.Sprintf(config_format,id, key, got_val)
	dblines[line] = str

	fileContent := ""
	for _, l := range dblines {
		fileContent += l
		fileContent += "\n"
	}

	if ioutil.WriteFile(config_dbpath, []byte(fileContent), 0644) != nil{
		return false
	}else{
		return true
	}
}

func emptydb(){
	ioutil.WriteFile(config_dbpath, []byte(""), 0644)
	dblines, _ = File2lines(config_dbpath)
}

func removeline(ind int){
	fileContent := ""
	for i, l := range dblines {
		if (i != ind){
			fileContent += l
			fileContent += "\n"
		}
	}
	ioutil.WriteFile(config_dbpath, []byte(fileContent), 0644)
	dblines, _ = File2lines(config_dbpath)
}
