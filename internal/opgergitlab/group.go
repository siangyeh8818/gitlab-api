package opgergitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	//"github.com/coreos/go-etcd/etcd"

	"github.com/siangyeh8818/gitlab.api/pkg/setting"
)

type getGroupResponse struct {
	ID                             int           `json:"id"`
	WebURL                         string        `json:"web_url"`
	Name                           string        `json:"name"`
	Path                           string        `json:"path"`
	Description                    string        `json:"description"`
	Visibility                     string        `json:"visibility"`
	ShareWithGroupLock             bool          `json:"share_with_group_lock"`
	RequireTwoFactorAuthentication bool          `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod           int           `json:"two_factor_grace_period"`
	ProjectCreationLevel           string        `json:"project_creation_level"`
	AutoDevopsEnabled              interface{}   `json:"auto_devops_enabled"`
	SubgroupCreationLevel          string        `json:"subgroup_creation_level"`
	EmailsDisabled                 interface{}   `json:"emails_disabled"`
	MentionsDisabled               interface{}   `json:"mentions_disabled"`
	LfsEnabled                     bool          `json:"lfs_enabled"`
	DefaultBranchProtection        int           `json:"default_branch_protection"`
	AvatarURL                      interface{}   `json:"avatar_url"`
	RequestAccessEnabled           bool          `json:"request_access_enabled"`
	FullName                       string        `json:"full_name"`
	FullPath                       string        `json:"full_path"`
	CreatedAt                      time.Time     `json:"created_at"`
	ParentID                       interface{}   `json:"parent_id"`
	LdapCn                         interface{}   `json:"ldap_cn"`
	LdapAccess                     interface{}   `json:"ldap_access"`
	SharedWithGroups               []interface{} `json:"shared_with_groups"`
	Projects                       []interface{} `json:"projects"`
	SharedProjects                 []interface{} `json:"shared_projects"`
	SharedRunnersMinutesLimit      interface{}   `json:"shared_runners_minutes_limit"`
	ExtraSharedRunnersMinutesLimit interface{}   `json:"extra_shared_runners_minutes_limit"`
	PreventForkingOutsideGroup     interface{}   `json:"prevent_forking_outside_group"`
}

func checkGroup(groupname string) bool {
	result := false
	url := setting.GitlabAPISetting.Url + "api/v4/groups"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("PRIVATE-TOKEN", gitlabToken)

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

	response := []getGroupResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//fmt.Println(response)
	for va := range response {
		if response[va].Name == groupname {
			fmt.Println("group is already exist.")
			result = true
		}
	}
	return result
}

func addNewGroup(groupname string) (bool, []byte) {
	url := setting.GitlabAPISetting.Url + "api/v4/groups?name=" + groupname + "&visibility=private&path=" + strings.ToLower(groupname) + "&subgroup_creation_level=maintainer"
	var spaceResponse []byte
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

func checkSubGroup(groupid string, targetsub string) bool {

	result := false
	url := setting.GitlabAPISetting.Url + "api/v4/groups/" + groupid + "/subgroups"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return result
	}
	req.Header.Add("PRIVATE-TOKEN", gitlabToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return result
	}

	response := []getGroupResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return result
	}
	fmt.Println(response)
	for va := range response {
		if response[va].Name == targetsub {
			fmt.Printf("subgroup '" + targetsub + "'is already exist.\n")
			result = true
		}
	}

	return result
}

func addNewSubGroup(parentgroupid string, groupname string) (bool, []byte) {
	var spaceResponse []byte
	url := setting.GitlabAPISetting.Url + "api/v4/groups?parent_id=" + parentgroupid + "&name=" + groupname + "&path=" + groupname
	method := "POST"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

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
