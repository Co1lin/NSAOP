package location_test

import (
	"io/ioutil"
	"nsaop/utils/constant"
	"nsaop/utils/test"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var deferBase = []string{constant.TableUser, constant.TableLocation}

func TestMain(m *testing.M) {
	// setup a database for location related testing
	defer test.SetupDatabaseForTest(deferBase)()
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = ioutil.Discard
	code := m.Run()
	os.Exit(code)
}
