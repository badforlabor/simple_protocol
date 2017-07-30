package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	exportGoFile(nil)
	return

	// 结果解析过程中的异常
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	filename := "test.proto"

	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		datas, err := ioutil.ReadAll(file)
		if err == nil {
			p := NewParser(string(datas))
			p.DoParse()
		}
	}
}
