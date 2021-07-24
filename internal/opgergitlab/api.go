package opgergitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/siangyeh8818/gitlab.api/pkg/setting"
)

var gitlabToken = ""
var organizationID = ""
var cicdGroupID = ""
var provisionerProjectID = ""
var provisionerTriggerToken = ""
var provisionerUrlEncode = ""

func ImportOpger() {

	fmt.Println(setting.EtcdSetting.Endpoint)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{setting.EtcdSetting.Endpoint},
		Username:    setting.EtcdSetting.User,
		Password:    setting.EtcdSetting.Password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}

	defer cli.Close()

	_, err = cli.Put(context.TODO(), "/gitlab/auth", "") // 記得塞token
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	gitlabToken = getGitlabToken(cli)

	log.Printf("getGitlabToken is %s \n", gitlabToken)

	// 公司組織在conf中定義
	if checkGroup(setting.GitlabAPISetting.Organization) { // 檢查公司組織的Group 是否有被創建
		log.Printf("We don't need to create group , %s group is exist.\n", setting.GitlabAPISetting.Organization)
	} else { // 若不存在 , 則創建該公司組織的Group
		log.Printf("Creating group , name is %s \n", setting.GitlabAPISetting.Organization)
		newgResult, newgResponse := addNewGroup(setting.GitlabAPISetting.Organization)
		if newgResult {
			// 若創建成功後 , 其info 會被存在etcd 中
			_, err = cli.Put(context.TODO(), "/gitlab/"+setting.GitlabAPISetting.Organization+"/info", string(newgResponse)) //etcd存放路徑為 /gitlab/組織名/info
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Println("Get Organization group info from etcd")
	resp, err := cli.Get(ctx, "/gitlab/"+setting.GitlabAPISetting.Organization+"/info")
	cancel() //超出timeout便取消
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)
	response := getGroupResponse{}
	for _, ev := range resp.Kvs {
		//fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &response)
		if err != nil {
			fmt.Println(err)
		}
	}
	organizationID = strconv.Itoa(response.ID)
	//檢查組織group中 ，是否有cicd的subGroup 存在
	if checkSubGroup(organizationID, "cicd") {
		log.Printf("We don't need to create subgroup 'cicd' , %s own 'cicd' subgroup.\n", setting.GitlabAPISetting.Organization)
	} else {
		log.Printf("Creating 'cicd' subgroup into group %s \n", setting.GitlabAPISetting.Organization)
		newgResult, newgResponse := addNewSubGroup(organizationID, "cicd")
		if newgResult {
			// 若subGroup創建成功後 , 其info 會被存在etcd 中
			_, err = cli.Put(context.TODO(), "/gitlab/"+setting.GitlabAPISetting.Organization+"/cicd/info", string(newgResponse)) //etcd存放路徑為 /gitlab/組織名/cicd/info
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	fmt.Println("Get cicd subgroup info from etcd")
	subresp, err := cli.Get(ctx2, "/gitlab/"+setting.GitlabAPISetting.Organization+"/cicd/info")
	cancel() //超出timeout便取消
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)
	subresponse := getGroupResponse{}
	for _, ev := range subresp.Kvs {
		//fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &subresponse)
		if err != nil {
			fmt.Println(err)
		}
	}
	cicdGroupID = strconv.Itoa(subresponse.ID)

	//檢查cicd下的opger project存不存在
	identifyProjectInit("deploy", cli)
	identifyProjectInit("environment", cli)
	identifyProjectInit("base", cli)
	identifyProjectInit("configuration", cli)
	identifyProjectInit("network", cli)
	identifyProjectInit("release", cli)
	identifyProjectInit("scripts", cli)
	identifyProjectInit("provisioner", cli)

	ctx3, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	fmt.Println("Get cicd subgroup info from etcd")
	triggerresp, err := cli.Get(ctx3, "/gitlab/"+setting.GitlabAPISetting.Organization+"/cicd/provisioner/info")
	cancel() //超出timeout便取消
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)
	provisionerInfo := getGProjectResponse{}
	for _, ev := range triggerresp.Kvs {
		//fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &provisionerInfo)
		if err != nil {
			fmt.Println(err)
		}
	}
	provisionerProjectID = strconv.Itoa(provisionerInfo.ID)
	provisionerUrlEncode = UrlEncoding(provisionerInfo.WebURL)

	fmt.Println(provisionerUrlEncode)

	// 檢查trigger是否存在
	if checkProjectTrigger(UrlEncoding(provisionerInfo.WebURL), provisionerProjectID) {
		log.Printf("We don't need to create trigger into '%s/cicd/provisioner'.\n", setting.GitlabAPISetting.Organization)
	} else {
		log.Printf("Creating trigger into '%s/cicd/provisioner'.\n", setting.GitlabAPISetting.Organization)
		newgResult, newgResponse := addNewProjectTrigger(UrlEncoding(provisionerInfo.WebURL), provisionerInfo.Name, provisionerProjectID)
		if newgResult {
			// 若subGroup創建成功後 , 其info 會被存在etcd 中
			_, err = cli.Put(context.TODO(), "/gitlab/"+setting.GitlabAPISetting.Organization+"/cicd/"+provisionerInfo.Name+"/trigger/info", string(newgResponse))
			//etcd存放路徑為 /gitlab/組織名/cicd/provisioner/trigger/info
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// 檢查variable是否存在

	log.Printf("--- organizationID : %s\n", organizationID)
	log.Printf("--- cicdGroupID : %s\n", cicdGroupID)
	log.Printf("--- provisionerProjectID : %s\n", provisionerProjectID)
}

func getGitlabToken(client *clientv3.Client) string {

	var result string
	requestTimeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := client.Get(ctx, "/gitlab/auth")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		//fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		result = string(ev.Value)
	}
	return result
}

func identifyProjectInit(projectname string, client *clientv3.Client) {
	if checkProject(projectname, cicdGroupID) {
		log.Printf("We don't need to create project '%s'.\n", projectname)
	} else {
		log.Printf("Creating '%s' project into subgroup '%s'\n", projectname, setting.GitlabAPISetting.Organization+"/cicd")
		newgResult, newgResponse := addNewProject(projectname, cicdGroupID)

		if newgResult {
			log.Printf("Creating '%s' project into subgroup '%s' is success.\n", projectname, setting.GitlabAPISetting.Organization+"/cicd")
			// 若subGroup創建成功後 , 其info 會被存在etcd 中
			_, err := client.Put(context.TODO(), "/gitlab/"+setting.GitlabAPISetting.Organization+"/cicd/"+projectname+"/info", string(newgResponse)) //etcd存放路徑為 /gitlab/組織名/cicd/專案名稱/info
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("Creating '%s' project into subgroup '%s' is failed.", projectname, setting.GitlabAPISetting.Organization+"/cicd")
		}
	}
}

func identifyVariableInit(variable string, client *clientv3.Client) {

	if checkGroupVariable(organizationID, variable) {
		log.Printf("We don't need to create group-variable into group '%s'.\n", setting.GitlabAPISetting.Organization)
	} else {
		log.Printf("Creating group-variable into '%s' group.\n", setting.GitlabAPISetting.Organization)
		newgResult, newgResponse := addNewGroupVariable(provisionerUrlEncode, organizationID, variable, UrlEncoding(""))
		if newgResult {
			// 若subGroup創建成功後 , 其info 會被存在etcd 中
			_, err := client.Put(context.TODO(), "/gitlab/"+setting.GitlabAPISetting.Organization+"/variable/"+variable, string(newgResponse))
			//etcd存放路徑為 /gitlab/組織名/variable/變數名
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
