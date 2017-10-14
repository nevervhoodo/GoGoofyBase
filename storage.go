package main

import (
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

//Insert string to file in special place
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