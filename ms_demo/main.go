// You can edit this code!
// Click here and start typing.
package main

import (
"fmt"
"reflect"
)

type Student struct {
	ID	int
	Name	int
}

func main() {
	st1 := Student{
		ID: 1,
		Name: 2,
	}
	
	var temp interface{} = st1
	temp2 := map[string]interface{}{}
	temp2["id"] = st1.ID
	temp2["name"] = st1.Name
		
	// fmt.Println(reflect.TypeOf(temp))
	
	structValue := reflect.New(reflect.TypeOf(temp)).Elem()
	// fmt.Println("==>", structValue)
	// structValue.Set(reflect.ValueOf(temp))
	// for i := 0; i < structValue.NumField(); i++ {
	// 	// field := structValue.Field(i)
	// 	// fmt.Println(field)
	// 	structValue.Field(i).SetInt(32)
	// }
	

	
	
	// fmt.Println(reflect.TypeOf(temp) == reflect.TypeOf(st1))

	// structTpye := reflect.TypeOf(temp2)
	// fmt.Println(structTpye)
	// for i := 0; i < structTpye.NumField(); i++ {
	// 	fmt.Println(structTpye.Field(i).Name)
	// }
	
	fmt.Println("==>", temp2)

	for key, value := range temp2 {
		fmt.Println(key, value)
		structValue.FieldByName(key).Set(reflect.ValueOf(value))
	}
	fmt.Println("Done: ", structValue.Interface().(Student))
	fmt.Println(structValue.Type())
}