package opgergitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/siangyeh8818/gitlab.api/pkg/setting"
	//"github.com/coreos/go-etcd/etcd"
)

type getGProjectResponse struct {
	ID                                        int                             `json:"id"`
	Description                               interface{}                     `json:"description"`
	Name                                      string                          `json:"name"`
	NameWithNamespace                         string                          `json:"name_with_namespace"`
	Path                                      string                          `json:"path"`
	PathWithNamespace                         string                          `json:"path_with_namespace"`
	CreatedAt                                 time.Time                       `json:"created_at"`
	DefaultBranch                             interface{}                     `json:"default_branch"`
	TagList                                   []interface{}                   `json:"tag_list"`
	SSHURLToRepo                              string                          `json:"ssh_url_to_repo"`
	HTTPURLToRepo                             string                          `json:"http_url_to_repo"`
	WebURL                                    string                          `json:"web_url"`
	ReadmeURL                                 interface{}                     `json:"readme_url"`
	AvatarURL                                 interface{}                     `json:"avatar_url"`
	ForksCount                                int                             `json:"forks_count"`
	StarCount                                 int                             `json:"star_count"`
	LastActivityAt                            time.Time                       `json:"last_activity_at"`
	Namespace                                 GitlabNamespace                 `json:"namespace"`
	ContainerRegistryImagePrefix              string                          `json:"container_registry_image_prefix"`
	Links                                     GitlabLinks                     `json:"_links"`
	PackagesEnabled                           bool                            `json:"packages_enabled"`
	EmptyRepo                                 bool                            `json:"empty_repo"`
	Archived                                  bool                            `json:"archived"`
	Visibility                                string                          `json:"visibility"`
	ResolveOutdatedDiffDiscussions            bool                            `json:"resolve_outdated_diff_discussions"`
	ContainerRegistryEnabled                  bool                            `json:"container_registry_enabled"`
	ContainerExpirationPolicy                 gitlabContainerExpirationPolicy `json:"container_expiration_policy"`
	IssuesEnabled                             bool                            `json:"issues_enabled"`
	MergeRequestsEnabled                      bool                            `json:"merge_requests_enabled"`
	WikiEnabled                               bool                            `json:"wiki_enabled"`
	JobsEnabled                               bool                            `json:"jobs_enabled"`
	SnippetsEnabled                           bool                            `json:"snippets_enabled"`
	ServiceDeskEnabled                        bool                            `json:"service_desk_enabled"`
	ServiceDeskAddress                        interface{}                     `json:"service_desk_address"`
	CanCreateMergeRequestIn                   bool                            `json:"can_create_merge_request_in"`
	IssuesAccessLevel                         string                          `json:"issues_access_level"`
	RepositoryAccessLevel                     string                          `json:"repository_access_level"`
	MergeRequestsAccessLevel                  string                          `json:"merge_requests_access_level"`
	ForkingAccessLevel                        string                          `json:"forking_access_level"`
	WikiAccessLevel                           string                          `json:"wiki_access_level"`
	BuildsAccessLevel                         string                          `json:"builds_access_level"`
	SnippetsAccessLevel                       string                          `json:"snippets_access_level"`
	PagesAccessLevel                          string                          `json:"pages_access_level"`
	OperationsAccessLevel                     string                          `json:"operations_access_level"`
	AnalyticsAccessLevel                      string                          `json:"analytics_access_level"`
	EmailsDisabled                            interface{}                     `json:"emails_disabled"`
	SharedRunnersEnabled                      bool                            `json:"shared_runners_enabled"`
	LfsEnabled                                bool                            `json:"lfs_enabled"`
	CreatorID                                 int                             `json:"creator_id"`
	ImportStatus                              string                          `json:"import_status"`
	OpenIssuesCount                           int                             `json:"open_issues_count"`
	CiDefaultGitDepth                         int                             `json:"ci_default_git_depth"`
	CiForwardDeploymentEnabled                bool                            `json:"ci_forward_deployment_enabled"`
	PublicJobs                                bool                            `json:"public_jobs"`
	BuildTimeout                              int                             `json:"build_timeout"`
	AutoCancelPendingPipelines                string                          `json:"auto_cancel_pending_pipelines"`
	BuildCoverageRegex                        interface{}                     `json:"build_coverage_regex"`
	CiConfigPath                              interface{}                     `json:"ci_config_path"`
	SharedWithGroups                          []interface{}                   `json:"shared_with_groups"`
	OnlyAllowMergeIfPipelineSucceeds          bool                            `json:"only_allow_merge_if_pipeline_succeeds"`
	AllowMergeOnSkippedPipeline               interface{}                     `json:"allow_merge_on_skipped_pipeline"`
	RestrictUserDefinedVariables              bool                            `json:"restrict_user_defined_variables"`
	RequestAccessEnabled                      bool                            `json:"request_access_enabled"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool                            `json:"only_allow_merge_if_all_discussions_are_resolved"`
	RemoveSourceBranchAfterMerge              bool                            `json:"remove_source_branch_after_merge"`
	PrintingMergeRequestLinkEnabled           bool                            `json:"printing_merge_request_link_enabled"`
	MergeMethod                               string                          `json:"merge_method"`
	SuggestionCommitMessage                   interface{}                     `json:"suggestion_commit_message"`
	AutoDevopsEnabled                         bool                            `json:"auto_devops_enabled"`
	AutoDevopsDeployStrategy                  string                          `json:"auto_devops_deploy_strategy"`
	AutocloseReferencedIssues                 bool                            `json:"autoclose_referenced_issues"`
	RepositoryStorage                         string                          `json:"repository_storage"`
	RequirementsEnabled                       interface{}                     `json:"requirements_enabled"`
	SecurityAndComplianceEnabled              bool                            `json:"security_and_compliance_enabled"`
	ComplianceFrameworks                      []interface{}                   `json:"compliance_frameworks"`
}

type GitlabNamespace struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Kind      string      `json:"kind"`
	FullPath  string      `json:"full_path"`
	ParentID  int         `json:"parent_id"`
	AvatarURL interface{} `json:"avatar_url"`
	WebURL    string      `json:"web_url"`
}

type GitlabLinks struct {
	Self          string `json:"self"`
	Issues        string `json:"issues"`
	MergeRequests string `json:"merge_requests"`
	RepoBranches  string `json:"repo_branches"`
	Labels        string `json:"labels"`
	Events        string `json:"events"`
	Members       string `json:"members"`
}

type gitlabContainerExpirationPolicy struct {
	Cadence       string      `json:"cadence"`
	Enabled       bool        `json:"enabled"`
	KeepN         int         `json:"keep_n"`
	OlderThan     string      `json:"older_than"`
	NameRegex     string      `json:"name_regex"`
	NameRegexKeep interface{} `json:"name_regex_keep"`
	NextRunAt     time.Time   `json:"next_run_at"`
}

func addNewProject(projectname string, groupid string) (bool, []byte) {

	var spaceResponse []byte
	url := setting.GitlabAPISetting.Url + "api/v4/projects?path=" + projectname + "&name=" + projectname + "&namespace_id=" + groupid
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

func checkProject(projectname string, groupid string) bool {

	result := false
	url := setting.GitlabAPISetting.Url + "api/v4/groups/" + groupid + "/projects"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return result
	}
	req.Header.Add("PRIVATE-TOKEN", gitlabToken)
	req.Header.Add("Content-Type", "application/json")

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
	//fmt.Println(string(body))

	response := []getGProjectResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//fmt.Println(response)
	for va := range response {
		if response[va].Name == projectname && strconv.Itoa(response[va].Namespace.ID) == groupid {
			fmt.Printf("project '%s' is already exist.\n", projectname)
			result = true
		}
	}

	return result
}
