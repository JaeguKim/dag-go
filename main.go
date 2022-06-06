package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	dag "github.com/JaeguKim/dag-go/parser"
)


func check(e error) {
	if e != nil {
		fmt.Println("error while opening file")
		panic(e)
	}
}

func main() {
	xmlFilePath := flag.String("xmlFilePath", "", "xml file path for dag")
	flag.Parse()
	data,err := ioutil.ReadFile(*xmlFilePath)
	check(err)
	_, graph := dag.InitWithXML(data)
	graph.Start()
}
