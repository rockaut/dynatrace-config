package main

import (
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"
    "net/http"
    "bytes"
    "github.com/yudai/gojsondiff"
    "github.com/yudai/gojsondiff/formatter"
)

func main() {    
    var config map[string]interface{}
    var uri string
    var method string
    var payload map[string]interface{}
    var success int
    var apiUrl string
    var apiToken string

    // API_URL and API_TOKEN must be set as env variables, e.g.:
    // export DYNATRACE_API_URL="https://fdlfd1132.live.dynatrace.com/api/config"
    // export DYNATRACE_API_TOKEN="t-nJAasfsafsf"
    // API_TOKEN must have access to "Read Configuration" and "Write Configuration"
    // you can create an API_TOKEN in the dynatrace web console "Settings -> Integration -> Dynatrace API"
    apiUrl = os.Getenv("DYNATRACE_API_URL")
    apiToken = os.Getenv("DYNATRACE_API_TOKEN")

    // commands can be get, validate or apply
    command := os.Args[1]
 
    // configuration json is set via command line argument
    jsonFilePath := os.Args[2]

    // open the dynatrace configuration json
    jsonFile, err := os.Open(jsonFilePath)
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }

    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &config)

    uri = config["uri"].(string)
    method = config["method"].(string)
    payload = config["payload"].(map[string]interface {})
    success = int(config["success"].(float64))

    payloadJson, err := json.Marshal(payload)	
	if err != nil {
		fmt.Println(err.Error())
		return
	}
    
    switch command {
    case "get":
        success, response := execHttpRequest("GET", apiUrl + "/" + uri, nil, apiToken, 200) 
        if success {
            bodyBytes, err := ioutil.ReadAll(response.Body)
            if err != nil {
                fmt.Println(err.Error())
                return
            }
            bodyString := string(bodyBytes)
            fmt.Println(bodyString)
        }
    case "diff":
        success, response := execHttpRequest("GET", apiUrl + "/" + uri, nil, apiToken, 200) 
        if success {
            currentState, err := ioutil.ReadAll(response.Body)
            if err != nil {
                fmt.Println(err.Error())
                return
            }

            // thats experimantal quick approach using github.com/yudai/gojsondiff
            // configurationVersions and clusterVersion is obviously always different.
            var aJson map[string]interface{}
            json.Unmarshal(currentState, &aJson)
            config := formatter.AsciiFormatterConfig{
                ShowArrayIndex: true,
                Coloring:       true,
            }
            differ := gojsondiff.New()
            d, err := differ.Compare(currentState, payloadJson)
            formatter := formatter.NewAsciiFormatter(aJson, config)
            diffString, err := formatter.Format(d)
            fmt.Print(diffString)
        }
    case "validate":
        success, _ := execHttpRequest("POST", apiUrl + "/" + uri + "/validator", payloadJson, apiToken, 204) 
        if success {
            fmt.Println("Successfully validated configuration!")
        }
    case "apply":
        success, _ := execHttpRequest(method, apiUrl + "/" + uri, payloadJson, apiToken, success)
        if success {
            fmt.Println("Successfully applied configuration!")
        }
    default:
        fmt.Println("Please specify a valid command as first parameter [get|diff|validate|apply]")
    }
    defer jsonFile.Close()
}

func execHttpRequest(method string, url string, data []byte, apiToken string, expectedReturncode int) (bool, *http.Response) {
    // execute http request as defined in the configuration file
    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
    if err != nil {
        fmt.Println("Error occured: ", err)
        return false, nil
    }
    req.Header.Add("accept","application/json")
    req.Header.Add("Authorization","Api-Token " + apiToken)
    req.Header.Add("Content-Type", "application/json; charset=utf-8")
    resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false, resp
	}
    // check if http statuscode is the same as in the success field in the configuration file
    if resp.StatusCode == expectedReturncode {
        return true, resp
    } else {
        fmt.Println("Returncode not as expected: " + string(resp.StatusCode) + resp.Status)
    }
    return false, resp
}