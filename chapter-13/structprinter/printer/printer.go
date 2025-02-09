package printer

import (
	"fmt"
	"reflect"
)

func PrintStructFields(data interface{}) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		fmt.Println("Provided input is not a struct")
		return
	}
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		fmt.Printf("Field: %s, Type: %v, Value: %v\n", field.Name, field.Type, value)
		if value.Kind() == reflect.Struct {
			printNestedFields(value, "\t")
		}
	}
}

func printNestedFields(val reflect.Value, prefix string) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		fmt.Printf("%sNested Field: %s, Nested Type: %v, Nested Value: %v\n", prefix, field.Name, field.Type, value)
		if value.Kind() == reflect.Struct {
			printNestedFields(value, prefix+"\t")
		}
	}
}
