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

type options struct {
	FileName       string `js:"fileName"`
	AggregateLevel string `js:"aggregateLevel"`
}

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
	if err != nil {
		fmt.Println("[httpagg] The result file named " + fileName + " does not exist")
		var responses []http.Response
		return responses
	}

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

func (*Httpagg) CheckRequest(response http.Response, status bool, options options) {
	if options.FileName == "" {
		options.FileName = "httpagg.json"
	}

	if options.AggregateLevel == "" {
		options.AggregateLevel = "onError"
	}

	switch options.AggregateLevel {
	case "onError":
		if !status {
			AppendJSONToFile(options.FileName, response)
		}
	case "onSuccess":
		if status {
			AppendJSONToFile(options.FileName, response)
		}
	case "all":
		AppendJSONToFile(options.FileName, response)
	default:
		// by default, aggregate only invalid http responses
		if !status {
			AppendJSONToFile(options.FileName, response)
		}
	}
}

func (*Httpagg) GenerateRaport(httpaggResultsFileName string, httpaggReportFileName string) {
	const tpl = `
	<html lang="en">

<head>
    <meta charset="utf-8" />
    <title>httpagg html report</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="stylesheet" href="/css/demo.css" />
    <link rel="preconnect" href="https://fonts.gstatic.com" />
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter&family=Source+Code+Pro&display=swap" />
    <script src="https://code.jquery.com/jquery-3.5.1.js"></script>
    <script src="https://cdn.datatables.net/1.12.1/js/jquery.dataTables.min.js"></script>
    <style>
        .container {
            display: flex;
            /* Misc */
            width: 96%;
            height: 100%;
            margin-left: 2%;
            margin-right: 2%;
        }

        .container__left {
            /* Initially, the left takes 3/4 width */
            width: 65%;
            min-width: 30%;
            max-height: 100%;
            border: 1px solid #ece8f1;
            padding: 2%;
            overflow-y: scroll;
            font-family: Helvetica, sans-serif;
        }

        .container__right {
            /* Scroll */
            max-height: 100%;
            overflow-y: scroll;
            border: 1px solid #ece8f1;
            flex: 1;
            padding: 2%;
        }

        table {
            color: #3c3c64;
            font-size: 14px;
            line-height: 25px;
            border-collapse: collapse;
            width: 100%;
            max-height: 50px;
            border: 1px solid #ece8f1;
        }

        th {
            background: #f9f8fc;
            color: #5a5c87;
            font-size: 10px;
            letter-spacing: .5px;
            line-height: 18px;
            padding: 10px 20px;
            text-align: left;
            text-transform: uppercase;
            border-bottom: 1px solid #ece8f1;
            box-sizing: border-box;
            border-collapse: collapse;
            cursor: pointer;
        }

        td {
            padding: 20px;
            vertical-align: baseline;
            border-bottom: 1px solid #ece8f1;
            box-sizing: border-box;
        }

        tr {
            cursor: pointer;
        }

        a {
            border-bottom: 1px solid rgba(125, 100, 255, 0);
            color: #7d64ff;
            cursor: pointer;
            font-weight: 500;
            padding-bottom: 1px;
            position: relative;
            text-decoration: none;
            transition: all .3s;
            outline-color: #00cdff;
            background-color: transparent;
            box-sizing: border-box;
            font-size: 1em;
            line-height: 25px;
            border-collapse: collapse;
        }

        .error {
            color: #fa3287;
        }

        .success {
            color: #c3e88d;
        }

        .requestContainer {
            background: #3c3c64;
            margin: 0;
            padding: 15px;
            overflow-y: auto;
            text-align: left;
            transition: max-height .2s ease-in-out;
            font-family: monospace, monospace;
            font-size: 1em;
            box-sizing: border-box;
            color: #3c3c64;
            line-height: 25px;
        }

        .purple {
            color: #00cdff;
        }

        .white {
            color: white;
        }

        h2 {
            font-size: 25px;
            font-weight: 400;
            line-height: 35px;
            margin-top: 50px;
            margin-bottom: 15px;
            position: relative;
            box-sizing: border-box;
            color: #3c3c64;
            font-family: Helvetica, sans-serif;
        }

        input {
            color: #5a5c87;
            font-weight: 400;
            appearance: none;
            border: 1px solid #5a5c87;
            border-radius: 0;
            box-shadow: 0 1px 5px rgba(60, 60, 100, .05);
            color: #3c3c64;
            flex: 1 1;
            font-size: 15px;
            font-weight: 500;
            line-height: 20px;
            outline: none;
            overflow-x: auto;
            padding: 0 40px 0 15px;
            text-align: left;
            overflow: visible;
            font-family: inherit;
            margin: 0;
            box-sizing: border-box;
            width: 100%;
            padding: 12px;
            padding-left: 20px;
            margin-bottom: 30px;
            margin-top: 20px;
        }

        select {
            float: right;
            border-style: none;
            background-color: transparent;
            border: none;
            color: black;
            cursor: pointer;
            font-size: 12px;
            font-weight: 700;
            position: relative;
            transition: color .3s ease;
            text-transform: none;
            overflow: visible;
            line-height: 1.15;
            margin-right: fill;
            align-items: center;
            display: flex;
            flex-direction: column;
            position: relative;
            padding-right: 5px;
            font-size: 14px;
            border-color: blue;
            position: relative;
            -moz-appearance: none;
            -webkit-appearance: none;
            appearance: none;
            border: none;
            background: white url("data:image/svg+xml;utf8,<svg width='10' height='10' viewBox='0 0 10 10' fill='none' xmlns='http://www.w3.org/2000/svg' ><path d='M9 3 5 7 1 3' stroke='black' stroke-width='1.6'></path></svg>") no-repeat;
            background-position: right 0px top 50%;
            font-family: Helvetica, sans-serif;
        }

        .dataTables_info {
            margin-top: 30px;
            color: #3c3c64;
            font-size: 14px;
            line-height: 25px;
            box-sizing: border-box;
            margin-bottom: 20px;
        }

        #example_paginate {
            display: flex;
        }

        #example_previous {
            align-items: flex-start;
            margin-right: auto;
            padding-left: 0px;
            padding-right: 0;
            color: #7d64ff;
            font-size: 12px;
            font-weight: 700;
            line-height: 18px;
            text-transform: uppercase;
            cursor: pointer;
            display: flex;
            flex-direction: column;
            position: relative;
            text-decoration: none;
            
        }

        #example_next {
            align-items: flex-start;
            margin-left: auto;
            padding-left: 0px;
            padding-right: 0px;
            color: #7d64ff;
            font-size: 12px;
            font-weight: 700;
            line-height: 18px;
            text-transform: uppercase;
            cursor: pointer;
            display: flex;
            flex-direction: column;
            position: relative;
            text-decoration: none;
        }

        .paginate_button {
            margin-left: 2px;
            margin-right: 2px;
            padding-left: 0px;
            padding-right: 0px;
            color: #3c3c64;
            font-size: 13px;
            font-weight: 700;
            line-height: 18px;
            text-transform: uppercase;
            cursor: pointer;
            flex-direction: column;
            position: relative;
            text-decoration: none;
        }

        .paginate_button.current {
            color: #7d64ff;
        }

        .invisible_req {
        }

        textarea {
            color: white;
            background-color: #3c3c64;
            resize: none;
            width: 100%;
            border: none;
            outline: none;
          }

        code {
            line-height: 16px;
        }

        #example_filter {
            padding-top: 40px;
        }
        
        #example_wrapper {
            padding-top: 30px;
        }
    </style>
</head>

<body>
    <div class="container">
        <div class="container__left">
            <table id="example">
                <thead>
                    <tr>
                        <th>Response timestamp</th>
                        <th>Status</th>
                        <th>Method</th>
                        <th>URL</th>
                        <th>Duration (ms)</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .}}
                        <tr>
                            <td><a>{{.Headers.Date}}</a></td>
                            <td>{{.Status}}</td>
                            <td>{{.Request.Method}}</td>
                            <td>{{.Request.URL}}</td>
                            <td>{{.Timings.Duration}}</td>
                        </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        <div></div>
        <div class="container__right">
            {{range .}}
                <div class="invisible_req">
                    <h2>Response</h2>
                    <div class="requestContainer">
                        <span class="purple"><code>{{.Proto}}</code></span>
                        <span class="white"><code>{{.Status}}</code></span>
                        <span class="purple"><code>{{.StatusText}}</code></span></br></br>

                        {{ if eq (len .Headers) 0 }}
                            <span class="white"><code>
                                No headers
                            </code></span></br></br>
                        {{ else }}
                            <span class="white"><code>
                                {{ range $key, $value := .Headers }}
                                    {{ $key }}: {{ index $value }}</br>
                                {{ end }}
                            </code></span></br>
                        {{ end }}
                       
                        {{ if eq (len .Cookies) 0 }}
                            <span class="white"><code>
                                No cookies
                            </code></span></br></br>
                        {{ else }}
                            <span class="white"><code>
                                {{ range $key, $value := .Cookies }}
                                    {{ $key }}: {{ index $value }}</br>
                                {{ end }}
                            </code></span></br>
                        {{ end }}

                        {{ if .Body }}
                            <span class="white"><code><textarea readonly>{{.Body}}</textarea></code></span>
                        {{ else if .Body }}
                            <span class="white"><code>
                                No body
                            </code></span>
                        {{ end }}

                    </div>


                    <h2>Request</h2>
                    <div class="requestContainer">
                        <span class="purple"><code>{{.Request.Method}}</code></span>
                        <span class="white"><code>{{.Request.URL}}</code></span>
                        <span class="purple"><code>{{.Proto}}</code></span></br></br>
                        
                        {{ if eq (len .Request.Headers) 0 }}
                            <span class="white"><code>
                                No headers
                            </code></span></br></br>
                        {{ else }}
                            <span class="white"><code>
                                {{ range $key, $value := .Request.Headers }}
                                    {{ $key }}: {{ index $value 0 }}</br>
                                {{ end }}
                            </code></span></br>
                        {{ end }}
                        
                        {{ if eq (len .Request.Cookies) 0 }}
                            <span class="white"><code>
                                No cookies
                            </code></span></br></br>
                        {{ else }}
                            <span class="white"><code>
                                {{ range $key, $value := .Request.Cookies }}
                                    {{ $key }}: {{ (index $value 0).Value }}</br>
                                {{ end }}
                            </code></span></br>
                        {{ end }}

                        {{ if eq (len .Request.Body) 0 }}
                            <span class="white"><code>
                                No body
                            </code></span>
                        {{ else }}
                            <span class="white"><code><textarea readonly>{{.Request.Body}}</textarea></code></span>
                        {{ end }}
                    </div>

                    <h2>Additional data</h2>
                    <div class="requestContainer">
                        <span class="white"><code>
                            <p class="purple">Timings:</p>
                            Duration: {{.Timings.Duration}} ms</br>
                            Blocked: {{.Timings.Blocked}} ms</br>
                            Connecting: {{.Timings.Connecting}} ms</br>
                            LookingUp: {{.Timings.LookingUp}} ms</br>
                            Receiving: {{.Timings.Receiving}} ms</br>
                            Sending: {{.Timings.Sending}} ms</br>
                            TLSHandshaking: {{.Timings.TLSHandshaking}} ms</br>
                            Waiting: {{.Timings.Waiting}} ms</br></br>
                            
                            <p class="purple">TLS:</p>
                            tls_version: {{.TLSVersion}}</br>
                            tls_cipher_suite: {{.TLSCipherSuite}}</br></br>
                           
                            <p class="purple">OCSP:</p>
                            NextUpdate: {{.OCSP.NextUpdate}}</br>
                            ProducedAt: {{.OCSP.ProducedAt}}</br>
                            RevocationReason: {{.OCSP.RevocationReason}}</br>
                            RevokedAt: {{.OCSP.RevokedAt}}</br>
                            Status: {{.OCSP.Status}}</br>
                            ThisUpdate: {{.OCSP.ThisUpdate}}</br></br>

                            <p class="purple">ERROR DATA:</p>
                            Error: {{.Error}}</br>
                            Error_code: {{.ErrorCode}}
                        </code></span>
                   </div>
                </div>
            {{end}}
        </div>
    </div>

    <script type="module">
       

        $(document).ready(function () {
            $('#example').DataTable({
                "language": {
                    "lengthMenu": '_MENU_',
                    "search": '<i class="search"></i>',
                    "searchPlaceholder": "Search",

                },
                order: []
            });

            // change testform
            document.querySelectorAll("textarea").forEach(element => {
                function autoResize(el) {
                    el.style.height = el.scrollHeight + 'px';
                }
                autoResize(element);
                element.addEventListener('input', () => autoResize(element));    
            });

            $(document).on("click", 'table tr', function() {
                $('table tr').css('background','#ffffff');
                $(this).css('background','#f9f8fc');

                var data = $('table').DataTable().cells( selectedRow, '' ).render( 'display' );
                var selectedRow = data.row(this).index();

                $('.invisible_req').css('display','none');
                $('.invisible_req').eq(selectedRow).css('display','block');
            });

            $(document).on("click", 'thead tr', function() {
                $('table tr').eq(1).trigger('click');
            });

            $(document).on("click", '.paginate_button', function() {
                $('table tr').eq(1).trigger('click');
            });

            $('table tr').eq(1).trigger('click');
        });
    </script>
</body>

</html>
	`

	// temp := template.Must(template.New("index.txt").Funcs(funcMap).ParseFiles("index.txt"))
	temp, err := template.New("index.txt").Parse(tpl)
	check(err)

	if httpaggResultsFileName == "" {
		httpaggResultsFileName = "httpagg.json"
	}

	if httpaggReportFileName == "" {
		httpaggReportFileName = "httpaggReport.html"
	}

	responses := getJSONAggrResults(httpaggResultsFileName)
	if len(responses) != 0 {
		file, err := os.Create(httpaggReportFileName)
		check(err)

		err = temp.Execute(file, responses)
		check(err)
	}
}

var index string
