package installerevent

type InstallerEvents struct {
	Companyid         string        `json:"companyID"`
	Publishername     string        `json:"publisherName"`
	Publishertimestap string        `json:"publisherTimestap"`
	Executemodule     string        `json:"executeModule"`
	Executeaction     string        `json:"executeAction"`
	Excuteargs        []interface{} `json:"excuteArgs"`
}


