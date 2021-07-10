package main

import (
	"math"
	"time"

	"github.com/robfig/cron"

	"nsaop/config"
	"nsaop/model"
	"nsaop/router"
	"nsaop/utils/email"
)

// @title NSAOP Documents for Developer
// @version 2.0.0
// @description Backend API of NSAOP
// @BasePath /v2

func cronTask() {
	c := cron.New()
	c.AddFunc("@every "+config.Router.GetString("jwt.clearFreq"), func() {
		t := time.Now().AddDate(0, 0, -int(math.Ceil(float64(config.Router.GetInt("jwt.refreshTime"))/3600.0)))
		model.DB.Where("create_at < ?", t).Delete(&model.Refresh{})
	})
	if *config.EnableEmail {
		c.AddFunc(config.Email.GetString("send_cron"), email.SendBillAll)
	}
	c.Start()
}

func main() {
	config.Parse()
	model.Run()
	cronTask()
	router.Run()
}
