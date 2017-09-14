package main

import (
	"log"
)
type Bike struct {
	Name string
	Age int
}
var bikes []Bike
func main() {
	bikes=append(bikes,Bike{Name:"aaa",Age:15})
	bikes=append(bikes,Bike{Name:"bbb",Age:15})
	bikes=append(bikes,Bike{Name:"ccc",Age:15})
	for k, v := range bikes {
		if v.Name == "bbb" {
			kk := k + 1
			bikes = append(bikes[:k], bikes[kk:]...)
		}
	}
	log.Println(bikes)
}