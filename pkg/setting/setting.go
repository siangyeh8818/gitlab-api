package setting

import (
	"log"
	"os"

	"github.com/go-ini/ini"
)

type App struct {
	//JwtSecret string
	//PageSize  int
	//PrefixUrl string

	RuntimeRootPath string

	//ImageSavePath  string
	//ImageMaxSize   int
	//ImageAllowExts []string

	//ExportSavePath string
	//QrCodeSavePath string
	//FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Etcd struct {
	Endpoint string
	User     string
	Password string
}

var EtcdSetting = &Etcd{}

type GitlabAPI struct {
	Url          string
	Organization string
}

var GitlabAPISetting = &GitlabAPI{}

type NatsStreaming struct {
	Address    string
	ClusterId  string
	QueueGroup string
}

var NatsStreamingSetting = &NatsStreaming{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	if os.Getenv("CONFIG-PATH") == "" {
		os.Setenv("CONFIG-PATH", "conf/app.ini")
	}
	cfg, err = ini.Load(os.Getenv("CONFIG-PATH"))
	//cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	/*

		mapTo("server", ServerSetting)
		mapTo("database", DatabaseSetting)
		mapTo("redis", RedisSetting)
	*/
	mapTo("app", AppSetting)
	mapTo("etcd", EtcdSetting)
	mapTo("gitlab", GitlabAPISetting)
	mapTo("nats-streaming", NatsStreamingSetting)
	/*
		AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
		ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
		ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
		RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	*/
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
