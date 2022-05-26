// MIT License
//
// Copyright (c) 2022 Grzegorz Piechnik @gpiechnik2
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package httpagg

import (
	"encoding/json"
	"log"
	"os"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/httpagg", new(Httpagg))
}

// Httpsaver is the k6 extension
type Httpagg struct{}

// constants
type response struct {
	// remote_ip        int
	// remote_port      int
	// url              string
	status int
	// proto            string
	// headers          interface{}
	// cookies          interface{}
	// body             string
	// timings          interface{}
	// tls_version      string
	// tls_cipher_suite string
	// ocsp             interface{}
	// error            string
	// error_code       int
	// request          interface{}
}

type options struct {
	errorFileName     string
	successesFileName string
	allReqFileName    string
	// htmlReportFileName string
}

// parses a config JSON file and using json. Unmarshal stores its data in a struct
func ParseJSON(jsonData interface{}) string {
	out, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func AppendJSONToFile(file string, jsonData interface{}) {
	// get string from JSON
	stringData := ParseJSON(jsonData)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(stringData + "\n")); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func CreateErrorsOutput(res response, errorFileName string) {
	// if it is a client-side or server-side error
	if res.status > 399 {
		AppendJSONToFile(errorFileName, res)
	}
}

func CreateSuccessesOutput(res response, successesFileName string) {
	// if it is not a client-side or server-side error
	if res.status < 400 {
		AppendJSONToFile(successesFileName, res)
	}
}

func CreateAllOutput(res response, allReqFileName string) {
	AppendJSONToFile(allReqFileName, res)
}

func (*Httpagg) CheckRequest(res response, level string, options options) bool {
	switch level {
	case "errors":
		CreateErrorsOutput(res, options.errorFileName)
		return true
	case "successes":
		CreateSuccessesOutput(res, options.successesFileName)
		return true
	case "all":
		CreateAllOutput(res, options.allReqFileName)
		return true
	default:
		CreateAllOutput(res, options.allReqFileName)
		return true
	}
	return false
}

// generate html raport
// func (*Httpsaver) GenerateRaport(htmlReportFileName string) {

// 	// pass
// }
