package e2e

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Azure/agentbaker/pkg/agent"
	"github.com/Azure/agentbaker/pkg/agent/datamodel"
	//"github.com/Azure/go-autorest/autorest/to"
)

// TODO 1: How to get the most accurate url links/image links for the currently hardcoded ones for eg CustomKubeBinaryURL, Pause Image etc
// TODO 2: Update --rotate-certificate (true for TLS enabled, false otherwise, small nit)
// TODO 3: Seperate out the certificate encode/decode
// TODO 4: Investigate CloudSpecConfig and its need. Without it, the bootstrapping struct breaks.

func createFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}

		defer f.Close()
	}
}

func TestE2EBasic(t *testing.T) {
	entry := "Generating CustomData and cseCmd"
	fmt.Println(entry)

	fields, err := os.Open("fields.json")
	if err != nil {
		fmt.Println(err)
	}

	defer fields.Close()
	fieldsByteValue, _ := ioutil.ReadAll(fields)

	values := customDataFields{}
	json.Unmarshal([]byte(fieldsByteValue), &values)

	createFile("../e2e/cloud-init.txt")
	createFile("../e2e/cseCmd")

	// caCertDecoded, _ := base64.URLEncoding.DecodeString(values.Cacert)
	// apiServerCertDecoded, _ := base64.URLEncoding.DecodeString(values.Apiservercert)
	// clientKeyDecoded, _ := base64.URLEncoding.DecodeString(values.Clientkey)
	// clientCertDecoded, _ := base64.URLEncoding.DecodeString(values.Clientcert)

	nbc, _ := ioutil.ReadFile("nodebootstrapping_config.json")
	config := &datamodel.NodeBootstrappingConfiguration{}
	json.Unmarshal([]byte(nbc), config)

	// customData
	baker := agent.InitializeTemplateGenerator()
	base64EncodedCustomData := baker.GetNodeBootstrappingPayload(config)
	customDataBytes, _ := base64.StdEncoding.DecodeString(base64EncodedCustomData)
	customData := string(customDataBytes)
	err = ioutil.WriteFile("cloud-init.txt", []byte(customData), 0644)
	if err != nil {
		fmt.Println("couldnt write to file", err)
	}

	// cseCmd
	cseCommand := baker.GetNodeBootstrappingCmd(config)
	err = ioutil.WriteFile("csecmd", []byte(cseCommand), 0644)
	if err != nil {
		fmt.Println("couldnt write to file", err)
	}
}
