package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/rockaut/dynatrace-config/cmd/calls"
	"github.com/rockaut/dynatrace-config/cmd/config"
)

var (
	GetCommand = &cobra.Command{
		Use:   "get",
		Short: "Fetches the actual config from the given API call.",
		Run:   runGet,
	}
)

func init() {
	GetCommand.Flags().StringP("outfile", "o", "stdout", "The file to write to.")
	GetCommand.Flags().StringP("file", "f", "", "API configuration to load/read")
	GetCommand.MarkFlagRequired("file")
}

func runGet(cmd *cobra.Command, args []string) {
	var currentState []byte
	var err error
	var loadedConfig map[string]interface{}
	var uri string

	filePath, _ := cmd.Flags().GetString("file")
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &loadedConfig)

	uri = loadedConfig["uri"].(string)

	success, response := calls.ExecHttpRequest("GET", config.ApiUrl+"/"+uri, nil, config.ApiToken, 200)
	if success {
		currentState, err = ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	fmt.Println(string(currentState))
}
