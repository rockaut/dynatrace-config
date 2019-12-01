package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// some configuration objects are simple, some are lists and
	ListConfigurations = map[string]string{
		"v1/alertingProfiles":              "displayName",
		"v1/anomalyDetection/diskEvents":   "name",
		"v1/anomalyDetection/metricEvents": "name",
		"v1/applicationDetectionRules":     "applicationIdentifier",
		"v1/autoTags":                      "name",
		"v1/aws/credentials":               "label",
		"v1/cloudFoundry/credentials":      "name",
		//"v1/dashboards": "dashboardMetadata": {"name"} // currently nested path to nameIdentifier is not allowed
		"v1/kubernetes/credentials": "label",
		"v1/customMetric/log":       "displayName",
		"v1/maintenanceWindows":     "name",
		"v1/managementZones":        "name",
		"v1/notifications":          "name",
		//"v1/service/customServices": "name" // customServices is split into different technologies, not supported right now
		//"v1/service/ibmMQTracing/imsEntryQueue": "??" // document says name but in configuration payload there is no name
		"v1/service/ibmMQTracing/queueManager": "name",
		"v1/service/requestAttributes":         "name",
		"v1/service/requestNaming":             "namingPattern", // not sure ...
	}
)

func GetListItem(name string, apiUrl string, uri string, configurationJson []byte, apiToken string) []byte {
	var config map[string]interface{}
	json.Unmarshal(configurationJson, &config)
	nameIdentifier := config[name].(string)
	id := GetIdByName(nameIdentifier, apiUrl+"/"+uri, apiToken)
	if id == "" {
		return nil
	}

	success, response := ExecHttpRequest("GET", apiUrl+"/"+uri+"/"+id, nil, apiToken, 200)
	if success {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		return bodyBytes
	}

	return nil
}

func CreateOrConfigureListItem(name string, apiUrl string, uri string, configurationJson []byte, apiToken string) {
	var config map[string]interface{}
	json.Unmarshal(configurationJson, &config)
	nameIdentifier := config[name].(string)
	id := GetIdByName(nameIdentifier, apiUrl+"/"+uri, apiToken)
	if id == "" {
		success, _ := ExecHttpRequest("POST", apiUrl+"/"+uri, configurationJson, apiToken, 201)
		if success {
			fmt.Println("Successfully created new configuration!")
		}
	} else {
		success, _ := ExecHttpRequest("PUT", apiUrl+"/"+uri+"/"+id, configurationJson, apiToken, 204)
		if success {
			fmt.Println("Successfully configured existing configuration!")
		}
	}
}

func GetIdByName(name string, url string, apiToken string) string {
	success, response := ExecHttpRequest("GET", url, nil, apiToken, 200)
	if success {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}

		var responseMap map[string]interface{}
		json.Unmarshal(bodyBytes, &responseMap)
		for _, element := range responseMap["values"].([]interface{}) {
			elementMap := element.(map[string]interface{})
			if elementMap["name"] == name {
				return elementMap["id"].(string)
			}
		}
	}

	return ""
}

func ExecHttpRequest(method string, url string, data []byte, apiToken string, expectedStatusCode int) (bool, *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error occured: ", err)
		return false, nil
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Api-Token "+apiToken)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false, resp
	}

	if resp.StatusCode == expectedStatusCode {
		return true, resp
	} else {
		fmt.Println("StatusCode was not expected: " + string(resp.StatusCode) + resp.Status)
	}

	return false, resp
}
