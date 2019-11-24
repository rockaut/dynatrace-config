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
    var configuration map[string]interface{}
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

    // some configuration objects are simple, some are lists and
    listConfigurations := map[string]string {
        "v1/alertingProfiles": "displayName",
        "v1/anomalyDetection/diskEvents": "name",
        "v1/anomalyDetection/metricEvents": "name",
        "v1/applicationDetectionRules": "applicationIdentifier",
        "v1/autoTags": "name",
        "v1/aws/credentials": "label",
        "v1/cloudFoundry/credentials": "name",
        //"v1/dashboards": "dashboardMetadata": {"name"} // currently nested path to nameIdentifier is not allowed
        "v1/kubernetes/credentials": "label",
        "v1/customMetric/log": "displayName",
        "v1/maintenanceWindows": "name",
        "v1/managementZones": "name",
        "v1/notifications": "name",
        //"v1/service/customServices": "name" // customServices is split into different technologies, not supported right now
        //"v1/service/ibmMQTracing/imsEntryQueue": "??" // document says name but in configuration payload there is no name
        "v1/service/ibmMQTracing/queueManager": "name",
        "v1/service/requestAttributes": "name",
        "v1/service/requestNaming": "namingPattern", // not sure ...
    }

    // open the dynatrace configuration json
    jsonFile, err := os.Open(jsonFilePath)
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }

    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &config)

    uri = config["uri"].(string)
    configuration = config["configuration"].(map[string]interface {})

    configurationJson, err := json.Marshal(configuration)	
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
            var currentStateJson map[string]interface{}
            json.Unmarshal(currentState, &currentStateJson)

            // ignore metadata when diffing currentState with desiredState
            delete(currentStateJson,"metadata")
            currentStateWithoutMetadata, _ := json.Marshal(currentStateJson)
            
            config := formatter.AsciiFormatterConfig{
                ShowArrayIndex: true,
                Coloring:       true,
            }
            differ := gojsondiff.New()
            d, err := differ.Compare(currentStateWithoutMetadata, configurationJson)
            if d.Modified() {
                formatter := formatter.NewAsciiFormatter(currentStateJson, config)
                diffString, _ := formatter.Format(d)
                fmt.Print(diffString)
            } else {
                fmt.Print("no difference found.")
            }
        }
    case "validate":
        success, _ := execHttpRequest("POST", apiUrl + "/" + uri + "/validator", configurationJson, apiToken, 204) 
        if success {
            fmt.Println("Successfully validated configuration!")
        }
    case "apply":
        // depending on the uri we need to apply the configuration directly or find the corresponding list item
        if name, found := listConfigurations[uri]; found {
            createOrConfigureListItem(name, apiUrl, uri, configurationJson, apiToken)
        } else {
            success, _ := execHttpRequest("PUT", apiUrl + "/" + uri, configurationJson, apiToken, 204)
            if success {
                fmt.Println("Successfully applied configuration!")
            }
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

func getIdByName(name string, url string, apiToken string) (string) {
    success, response := execHttpRequest("GET", url, nil, apiToken, 200)
    if success {
        bodyBytes, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Println(err.Error())
            return ""
        }

        var responseMap map[string]interface{}
        json.Unmarshal(bodyBytes, &responseMap)
        for _, element := range responseMap["values"].([]interface {}) {
            elementMap := element.(map[string]interface {})
            if elementMap["name"] == name {
                return elementMap["id"].(string)
            }
        }
    }
    return ""
}

// used for list items
func createOrConfigureListItem(name string, apiUrl string, uri string, configurationJson []byte, apiToken string) {
    var config map[string]interface{}
    json.Unmarshal(configurationJson, &config)
    nameIdentifier := config[name].(string)
    id := getIdByName(nameIdentifier, apiUrl + "/" + uri, apiToken)
    // if id is empty, we need to create a new element with POST method
    if id == "" {
        success, _ := execHttpRequest("POST", apiUrl + "/" + uri, configurationJson, apiToken, 201)
        if success {
            fmt.Println("Successfully created new configuration!")
        }       
    // if id already exists, then execute a PUT on the specific ID element        
    } else {
        success, _ := execHttpRequest("PUT", apiUrl + "/" + uri + "/" + id, configurationJson, apiToken, 204)
        if success {
            fmt.Println("Successfully configured existing configuration!")
        }                 
    }
}