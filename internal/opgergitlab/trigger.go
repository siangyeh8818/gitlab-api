package opgergitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/siangyeh8818/gitlab.api/pkg/setting"
	//"github.com/coreos/go-etcd/etcd"
)

type getTriggerResponse struct {
	ID          int          `json:"id"`
	Token       string       `json:"token"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	LastUsed    interface{}  `json:"last_used"`
	Owner       TriggerOwner `json:"owner"`
}

type TriggerOwner struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func addNewProjectTrigger(urlId string, projectname string, projectid string) (bool, []byte) {

	var description = "opger-trigger-" + projectname
	var spaceResponse []byte
	url := setting.GitlabAPISetting.Url + "api/v4/projects/" + projectid + "/triggers?id=" + urlId + "&description=" + description
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
	//fmt.Println(string(body))
	return true, body
}

func checkProjectTrigger(urlencoding string, projectid string) bool {
	var result = false
	url := setting.GitlabAPISetting.Url + "api/v4/projects/" + projectid + "/triggers?id=" + urlencoding
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

	response := []getTriggerResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//fmt.Println(response)
	for va := range response {
		if response[va].Owner.State == "active" {
			fmt.Println("triggeris already exist.")
			result = true
		}
	}
	return result
}

func UrlEncoding(path string) string {
	return url.PathEscape(path)
}
