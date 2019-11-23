package main

import (
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"
    "net/http"
    "bytes"
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

    // configuration json is set via command line argument
    jsonFilePath := os.Args[1]
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
    
    // validate new configuration
    if execHttpRequest("POST", apiUrl + "/" + uri + "/validator", payloadJson, apiToken, 204) {
        fmt.Println("Successfully validated configuration!")
        // apply configuration
        if execHttpRequest(method, apiUrl + "/" + uri, payloadJson, apiToken, success) {
            fmt.Println("Successfully applied configuration!")
        }
    } else {
        fmt.Println("Validation was not successful. I am sorry.")
    }

    defer jsonFile.Close()
}

func execHttpRequest(method string, url string, data []byte, apiToken string, expectedReturncode int) bool {
    // execute http request as defined in the configuration file
    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
    if err != nil {
        fmt.Println("Error occured: ", err)
        return false
    }
    req.Header.Add("accept","application/json")
    req.Header.Add("Authorization","Api-Token " + apiToken)
    req.Header.Add("Content-Type", "application/json; charset=utf-8")
    resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
    // check if http statuscode is the same as in the success field in the configuration file
    if resp.StatusCode == expectedReturncode {
        return true
    } else {
        fmt.Println("Returncode not as expected: " + string(resp.StatusCode) + resp.Status)
    }
    return false
}