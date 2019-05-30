package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"
)

var funcs = template.FuncMap{"TypeOf": reflect.TypeOf}

type JsonData map[string]interface{}
type Data struct {
	Package  string
	ModeName string
	Data     JsonData
}

const form = `
package {{.Package}}

type {{.ModeName}} struct { {{range $index, $value := .Data}}
  {{$index}} {{TypeOf $value}} {{end}}
}`

func main() {
	// log.WithFields(log.Fields{
	// 	"animal": "walrus",
	// }).Info("A walrus appears")
	// ac := accounting.Accounting{Symbol: "$", Precision: 2}
	// fmt.Println(ac.FormatMoney(123.123))
	// hello.Sayhello()
	var data Data
	data.Package = "first"
	data.ModeName = "ModeName"

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
	// fmt.Println(out.String())
}
