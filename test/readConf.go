package main

import (
	"io/ioutil"
	"fmt"
)
func LoadConf()string{
	data, err := ioutil.ReadFile("c:\\server.conf")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(data))
	return string(data)
}