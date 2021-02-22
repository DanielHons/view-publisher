package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"
)

var dataUrl *string
var reportFile *string
var outFile *string
var creation time.Time

func init() {
	creation = time.Now()
	creationNanos := strconv.Itoa(creation.Nanosecond())
	dataUrl = flag.String("dataUrl", "", "URL returning a json array")
	reportFile = flag.String("template", "sample-timesheet.gohtml", "the template file")
	outFile = flag.String("out", "report-out-"+creationNanos+".html", "the output file")
	flag.Parse()
}

func main() {
	tf := template.FuncMap{
		"add": func(i float64, j float64) float64 {
			return i + j
		},
		"now": func() string { return time.Now().Format("2006-01-02 15:04:05") },
		"isInt": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				return true
			default:
				return false
			}
		},
		"isString": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.String:
				return true
			default:
				return false
			}
		},
		"isSlice": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Slice:
				return true
			default:
				return false
			}
		},
		"isArray": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Array:
				return true
			default:
				return false
			}
		},
		"isMap": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Map:
				return true
			default:
				return false
			}
		},
	}
	if reportFile == nil {
		panic("no report file given")
	}

	t := template.New(path.Base(*reportFile)).Funcs(tf)
	tt, err := t.ParseFiles(*reportFile)
	if err != nil {
		log.Fatal(err)
	}
	if dataUrl == nil || *dataUrl == "" {
		log.Fatal("no data source url given")
	}
	data := loadListDataFromUrl(*dataUrl)

	f, err := os.Create(*outFile)
	check(err)
	defer f.Close()

	if err = tt.Execute(f, &data); err != nil {
		panic(err)
	}
	println(*outFile)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadListDataFromUrl(url string) []map[string]interface{} {
	resp, err := http.Get(url)
	check(err)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	check(readErr)

	var m []map[string]interface{}
	err = json.Unmarshal(body, &m)
	check(err)
	return m

}
