package main

import (
	"fmt"
	"reflect"
)

var print = fmt.Println


func  typeof(i interface{}) string{
	return fmt.Sprintf("%T",i)
}


func main (){

	var obj interface {}

	obj = [] string{"1","2"}

	print(reflect.TypeOf(obj))
	print(typeof(obj))

	obj = map[string]string{"a":"b","c":"1"}
	print(reflect.TypeOf(obj))

	print(typeof(obj))
	return 
}
