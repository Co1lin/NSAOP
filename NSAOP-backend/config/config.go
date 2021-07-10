package config

import (
	"flag"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

func GetConfig(name string) *viper.Viper {
	var (
		_, b, _, _ = runtime.Caller(0)
		basePath   = filepath.Dir(b)
	)
	v := viper.New()
	v.AddConfigPath(basePath)
	v.SetConfigName(name)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Read config failed: %v", err)
	}
	return v
}

var Config, Site, Model, Router, Email, NCE, ReCAPTCHA *viper.Viper
var Secret = make(map[string]string)
var (
	IfCORS      *bool
	EnableEmail *bool
	Dev         *bool
	TestMode    bool
)
var (
	NeedReCAPTCHA    bool
	ReCAPTCHATimeout time.Duration
	MinReCAPTCHAScore		 float64
)
var (
	NCEHost    	string
	NCETimeout 	time.Duration
	NCEUsername	string
	NCEPassword	string
)
var AdminRegisterCode 	string
var FrontendDomain		string

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// argument flag
	IfCORS = flag.Bool("cors", false, "whether to use CORS")
	Dev = flag.Bool("dev", false, "whether dev deploy")
	EnableEmail = flag.Bool("email", false, "whether to send bill email")
	// get sub items in the config file
	Config = GetConfig("config")
	Site = Config.Sub("site")
	Model = Config.Sub("database")
	Router = Config.Sub("router")
	Email = Config.Sub("email")
	NCE = Config.Sub("NCE")
	ReCAPTCHA = Config.Sub("recaptcha")

	AdminRegisterCode = Site.GetString("admin_register_code")
	Secret["access"] = Site.GetString("access_secret")
	Secret["refresh"] = Site.GetString("refresh_secret")

	NCEHost = NCE.GetString("host")
	NCETimeout = time.Duration(NCE.GetInt64("timeout")) * time.Millisecond
	NCEUsername = NCE.GetString("username")
	NCEPassword = NCE.GetString("password")

	FrontendDomain = Config.GetString("domain")
}

func Parse() {
	flag.Parse()
	if *Dev {
		Model.Set("db", Model.GetString("db")+"-dev")
	}
	// setup google reCAPTCHA
	UpdateNeedReCAPTCHA()
	ReCAPTCHATimeout = time.Duration(ReCAPTCHA.GetInt64("timeout")) * time.Millisecond
	MinReCAPTCHAScore = ReCAPTCHA.GetFloat64("min_reCAPTCHA_score")
}

func UpdateNeedReCAPTCHA() {
	NeedReCAPTCHA = ReCAPTCHA.GetBool("enable") && !(*Dev) && !TestMode
}
