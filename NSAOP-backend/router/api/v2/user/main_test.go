package user_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"nsaop/config"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

var deferBase = []string{constant.TableUser}

func TestMain(m *testing.M) {
	config.TestMode = true
	config.UpdateNeedReCAPTCHA()
	// setup a database for user related testing
	defer test.SetupDatabaseForTest(deferBase)()
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = ioutil.Discard
	code := m.Run()
	os.Exit(code)
}
