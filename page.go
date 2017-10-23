package main

import (
	"sync"
	"io/ioutil"
	"regexp"
	"log"
	"fmt"
	"os"
	"strconv"
)

type Page struct {
	sync.RWMutex
	pageNum int `json:"pageNum"`
	fileName string `json:"fileName"`
	offset int `json:"offset"`
	content []string `json:"content"`
	need_to_write bool
}

type Storage []Page

func (store Storage) Print() {
	fmt.Printf("\n###########################################\n")
	fmt.Printf("Current storage state\n")
	fmt.Printf("%#v \n", store)
	fmt.Printf("###########################################\n\n")
}

func _get_all_dbs (dir string) []string {
	var result []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		matched, _ := regexp.MatchString("^(dat-)([0-9])+$", f.Name())
		if matched {
			result = append(result, f.Name())
		}
	}
	return result
}

func _add_db_file (dir string) string {
	files := _get_all_dbs(dir)
	var i int
	var max int = 0
	for _, file := range files{
		_, err := fmt.Sscanf(file, "dat-%d", &i)
		if err != nil{
			fmt.Println("problems with new file at ",file)
		}
		if i > max {
			max = i
		}
	}
	newname := dir+"dat-"+strconv.Itoa(max+1)
	_, er := os.OpenFile(newname, os.O_CREATE, 0666)
	if er != nil{
		fmt.Println("cannot create ",newname)
	}
	return newname
}

//TODO: remove page space!!
func addPages (dir string) (Storage) {
	//init empty Store
	store := Storage {}
	//all files with db found
	dbs := _get_all_dbs(config_dbdir)

	page_num := 0
	for _, db_name := range (dbs){
		//got lines from each file
		lines, er := File2lines(config_dbdir+db_name)
		if er != nil{
			log.Fatal(er)
		}
		//check if file_size is correct
		if len(lines) > config_page_size*config_file_size {
			log.Fatal("too much strings in file!!\n")
		}
		//get number of pages in current file
		file_size := int(len(lines)/config_page_size)
		if len(lines) % config_page_size != 0 {
			file_size += 1
		}
		//page counter inside file
		cur_page := 0
		//for all pages in file
		for cur_page < file_size {
			//append page to Storage
			if cur_page != file_size -1{
				store = append(store, Page {
					pageNum: page_num,
					fileName: db_name,
					offset: cur_page*config_page_size,
					content: append(lines[cur_page*config_page_size:(cur_page+1)*config_page_size]),
					need_to_write: false})
			} else{
				store = append(store, Page {
					pageNum: page_num,
					fileName: db_name,
					offset: cur_page*config_page_size,
					content: append(lines[cur_page*config_page_size:len(lines)]),
					need_to_write: false})
			}
			cur_page += 1
			page_num += 1
		}

	}
	return store
}

//create table from page
func (page Page) Page2Table () Table {
	//create empty table
	var table Table
	//create field
	var id int
	var key,val string
	var counter int = 0
	for i, str := range (page.content) {
		n, err := fmt.Sscanf(str, config_format, &id, &key, &val)
		fmt.Println(str)
		if n < 3 || err!= nil{
			fmt.Printf("cannot parse line %d\n", i)
		} else {
			counter += 1
			table = append(table, Field{ID: id, Key: key, Value: val})
		}
	}
	//if page is empty
	if counter == 0{
		fmt.Print("empty page\n")
	}
	return table
}

//add line to page
func (store Storage) InsertStringToPage(fstr string, id int) Storage{
	for i, page := range store{
		fmt.Println("InsertPage ", page)
		page.Lock()
		fmt.Println("sdadasd")
		if len(page.content) < config_page_size{
			page.content = append(page.content, fstr)
			page.need_to_write = true
			store[i] = page
			page.Unlock()
			return store
		}
		page.Unlock()
	}
	var lines []string
	store = append(store, Page {
		pageNum: len(store),
		fileName: _add_db_file (config_dbdir),
		offset: 0,
		content: append(lines, fstr),
		need_to_write: true})
	return store
}

//update line in page
func (store Storage) updatePage(id int, val string) Storage {
	for i, page := range store {
		page.Lock()
		defer page.Unlock()
		for _, line := range page.content {
			var cur_id int
			var key, cur_val string
			fmt.Sscanf(line, config_format, &cur_id, &key, &cur_val)
			if cur_id == id {
				str := fmt.Sprintf(config_format,id, key, val)
				line = str
				page.need_to_write = true
				store[i] = page
				fmt.Print("result", page)
				return store
			}
		}
	}
	return store
}

//remove one line from page
func (store Storage) rmFromPage (id int) Storage {
	var resi int = -1
	var resp Page
	var pi int = -1
	for n, page := range store {
		page.Lock()
		defer page.Unlock()
		for i, line := range page.content {
			var cur_id int
			var key, cur_val string
			fmt.Sscanf(line, config_format, &cur_id, &key, &cur_val)
			if cur_id == id {
				resi = i
				resp = page
				pi = n
				break
			}
		}
	}
	if resi != -1 {
		resp.need_to_write = true
		resp.content = append(resp.content[0:resi], resp.content[resi+1:len(resp.content)]...)
		store[pi] = resp
	}
	return store
}

//Delete all db
func (store Storage) emptyPages() Storage{
	files := _get_all_dbs (config_dbdir)
	for _, file := range files{
		os.Remove(config_dbdir+file)
	}
	return Storage{}
}

func (page Page) writePage () {
	fmt.Println("write ",page)
	lines, err := File2lines(page.fileName)
	if (err != nil){
		log.Fatal(err)
	}
	fileContent := ""
	for i,line := range lines {
		fmt.Print(i,line)
		if i < page.offset {
			fileContent += line
			fileContent += "\n"
		} else if i < page.offset + len(page.content) {
			fileContent += page.content[i-page.offset]
			fileContent += "\n"
		} else if i< page.offset + config_page_size {
			fileContent += "\n"
		} else	{
			fileContent += line
			fileContent += "\n"
		}
	}
	ioutil.WriteFile(config_dbdir+page.fileName, []byte(fileContent), 0644)
}

//pull page to file if needed
func (store Storage) checkPage (npage int) (Storage, bool) {
	page := store[npage]
	st_n := npage
	fmt.Println("m ",npage)
	if page.pageNum != npage {
		for i,p := range store {
			fmt.Println("m ", p.pageNum, npage)
			if p.pageNum == npage {
				page = p
				st_n = i
			}
		}
	}

	var result bool
	fmt.Println("before check", page)
	page.RLock()
	if page.need_to_write {
		result = true
		fmt.Println("swap")
		page.writePage()
		page.need_to_write = false
		store[st_n] = page
	} else {
		result = false
	}
	page.RUnlock()
	return store, result
}

