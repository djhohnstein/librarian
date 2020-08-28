package main

import "C"
import (
	"fmt"
	"io/ioutil"
	"github.com/djhohnstein/librarian/pkg/memfile"
	"github.com/djhohnstein/librarian/pkg/memlibrary"
)

func main() {
	ramfsName := "djhtest"
	memFile, err := memfile.New(ramfsName)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile("/gosharedlib.so")
	if err != nil {
		panic(err)
	}
	bytesWritten := memFile.Write(b)
	fmt.Printf("Wrote %d bytes to memory file.\n", bytesWritten)
	lib, err := memlibrary.New(memFile,
		"helloworld",
		"helloworld",
		"helloworldCallback")
	goStr := "This is data from App Core!"
	inputString := []byte(goStr)
	results := lib.Invoke(inputString)
	fmt.Printf("App Core got results: %s\n", results)
	lib.SendDataToModule([]byte("Data from Core"))
	secondData := []byte("App Core second invocation!")
	res, err := memlibrary.InvokeLoadedModule("helloworld", secondData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("App Core got results: %s\n", res)
}
