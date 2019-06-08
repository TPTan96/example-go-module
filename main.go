package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"TypeOf":  reflect.TypeOf,
	"SliceOf": reflect.SliceOf,
	"GetType": GetType,
}

type JsonData map[string]interface{}

type Data struct {
	Package  string
	ModeName string
	Data     JsonData
	Structs  map[string]JsonData
}

const form = `
package {{.Package}}

type {{.ModeName}} struct { 
{{range $key, $value := .Data}}	{{$key}} {{GetType $value $key}}
{{end}}}

{{range $structName, $data := .Structs}}
type {{$structName}} struct { 
{{range $key, $value := $data}}	{{$key}} {{GetType $value $key}}
{{end}}}
{{end}}
`

var data Data

func main() {
	data.Package = "first"
	data.ModeName = "ModeName"
	data.Structs = make(map[string]JsonData)

	rawData, err := ioutil.ReadFile("./tmp/data.json")
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(rawData, &data.Data)

	temp := template.New("temp").Funcs(funcs)
	temp, err = temp.Parse(form)
	if err != nil {
		panic(err)
	}
	// var out bytes.Buffer

	f, _ := os.Create("./tmp/example.go")
	defer f.Close()
	err = temp.Execute(f, data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
	// fmt.Println(GetType(data.Data["array"]))
}

func GetType(value interface{}, key string) string {
	val := reflect.ValueOf(value)
	kind := val.Kind()
	if kind == reflect.Slice {
		fmt.Println(kind)
		if val.Len() == 0 {
			return "[]interface{}"
		}
		return GetTypeOfSlice(value.([]interface{}), key)
	} else if kind == reflect.Map {
		out := strings.ToUpper(key)
		temStruct := value.(map[string]interface{})
		data.Structs[out] = temStruct
		return out
	}
	return kind.String()
}

func GetTypeOfSlice(value []interface{}, key string) string {
	kind := reflect.ValueOf(value[0]).Kind()
	if kind == reflect.Map {
		for _, val := range value {
			if !DoMapsHaveSameField(value[0].(map[string]interface{}), val.(map[string]interface{})) {
				return "[]interface{}"
			}
		}
		nameOfMap := strings.ToUpper(key) + "S"
		data.Structs[nameOfMap] = value[0].(map[string]interface{})
		return "[]" + nameOfMap
	}
	for _, v := range value {
		if reflect.ValueOf(v).Kind() != kind {
			return "[]interface{}"
		}
	}
	return "[]" + kind.String()
}

func DoMapsHaveSameField(map1, map2 map[string]interface{}) bool {
	if len(map1) == len(map2) {
		for key := range map1 {
			if _, ok := map2[key]; !ok {
				return false
			}
		}
	}
	return true
}
