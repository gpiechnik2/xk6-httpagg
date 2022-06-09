package httpagg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/js/modules/k6/http"
)

func init() {
	modules.Register("k6/x/httpagg", new(Httpagg))
}

// Httpagg is the k6 extension
type Httpagg struct{}

func AppendJSONToFile(fileName string, jsonData http.Response) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0666)
	check(err)
	defer f.Close()

	file, _ := json.MarshalIndent(jsonData, "", " ")
	falseContent, err := f.Write(file)
	check(err)

	if false {
		fmt.Println(falseContent)
	}
}

func getJSONAggrResults(fileName string) []http.Response {
	jsonFile, err := os.Open(fileName)
	check(err)

	var responses []http.Response
	byteValue, _ := ioutil.ReadAll(jsonFile)
	responsesCoded := json.NewDecoder(strings.NewReader(string(byteValue[:])))

	for {
		var response http.Response

		err := responsesCoded.Decode(&response)
		if err == io.EOF {
			// all done
			break
		}

		check(err)
		responses = append(responses, response)
	}
	return responses
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Custom function must have only 1 return value, or 1 return value and an error
func formatDate(timeStamp time.Time) string {
	// Define layout for formatting timestamp to string
	// return timeStamp.Format("01-02-2006")
	return timeStamp.Format("Mon, 02 Jan 2006")

}

// Map name formatDate to formatDate function above
var funcMap = template.FuncMap{
	"formatDate": formatDate,
}

func (*Httpagg) CheckRequest(response http.Response, status bool, fileName string, aggregateLevel string) {
	if fileName == "" {
		fileName = "httpagg.json"
	}

	switch aggregateLevel {
	case "onError":
		if !status {
			AppendJSONToFile(fileName, response)
		}
	case "onSuccess":
		if status {
			AppendJSONToFile(fileName, response)
		}
	case "all":
		AppendJSONToFile(fileName, response)
	default:
		// by default, aggregate only invalid http responses
		if !status {
			AppendJSONToFile(fileName, response)
		}
	}
}

func (*Httpagg) GenerateRaport(httpaggResultsFileName string, httpaggReportFileName string) {
	temp := template.Must(template.New("index.txt").Funcs(funcMap).ParseFiles("index.txt"))

	if httpaggResultsFileName == "" {
		httpaggResultsFileName = "httpagg.json"
	}

	if httpaggReportFileName == "" {
		httpaggReportFileName = "httpaggReport.html"
	}

	responses := getJSONAggrResults(httpaggResultsFileName)

	file, err := os.Create(httpaggReportFileName)
	check(err)

	err = temp.Execute(file, responses)
	check(err)
}
