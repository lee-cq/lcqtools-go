package main

import (
	"fmt"
	"reflect"
)

var myUint16Var uint16 = 100
var myUint32Var uint32 = 200

func main() {
	varName := "myUint16Var" // 或者 "myUint32Var"
	value := getVariableValue(varName)
	if value.IsValid() {
		v := value.Interface()
		fmt.Printf("Value of '%s': %v (type: %T)\n", varName, v, v)
	} else {
		fmt.Printf("Variable '%s' not found.\n", varName)
	}
}

func getVariableValue(name string) reflect.Value {
	// 在当前包的命名空间中查找变量
	val := reflect.ValueOf(main).Elem().FieldByName(name)
	return val
}
