package opgergitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/siangyeh8818/gitlab.api/pkg/setting"
)

//"github.com/coreos/go-etcd/etcd"

type getGroupVariableResponse struct {
	VariableType     string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	EnvironmentScope string `json:"environment_scope"`
}

func checkGroupVariable(groupid string, variable string) bool {

	var result = false
	url := setting.GitlabAPISetting.Url + "api/v4/groups/" + groupid + "/variables"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("PRIVATE-TOKEN", gitlabToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//fmt.Println(string(body))
	response := []getGroupVariableResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//fmt.Println(response)
	for va := range response {
		if response[va].Key == variable {
			fmt.Println("variable already exist.")
			result = true
		}
	}

	return result
}

func addNewGroupVariable(urlId string, groupid string, key string, valueencoding string) (bool, []byte) {
	var spaceResponse []byte

	url := setting.GitlabAPISetting.Url + "api/v4/groups/" + groupid + "/variables?id=" + urlId + "&key=" + key + "&value=" + valueencoding + "&protected=true&masked=false"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return false, spaceResponse
	}
	req.Header.Add("PRIVATE-TOKEN", gitlabToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, spaceResponse
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, spaceResponse
	}
	fmt.Println(string(body))

	return true, body
}
