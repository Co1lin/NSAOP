package test

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"nsaop/config"
	"nsaop/model"
	"nsaop/router/api/v2/user"
	"nsaop/utils/constant"
)

var CommonUsersArray []*model.User
var CommonUsers map[string]*model.User
var CommonServices []*model.Service
var CommonLocation []*model.Location
var (
	CommonNCESiteName string
	CommonNCESiteID   string
	CommonNCESSIDPath string
)

func init() {
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("TestPasswd@nsaop2021"),
		bcrypt.DefaultCost)
	passwd := string(hash)

	CommonUsers = make(map[string]*model.User)
	for _, key := range []string{constant.RoleCustomer, constant.RoleOperator, constant.RoleEngineer, constant.RoleAdmin, "other"} {
		tmp := &model.User{
			Role:     key,
			Username: "NSAOPTest-" + key,
			Password: passwd,
			Company:  "NSAOP",
			Phone:    "12345678901",
			Email:    "nsaop@nsaop.com",
		}
		CommonUsersArray = append(CommonUsersArray, tmp)
		CommonUsers[key] = tmp
	}

	for i := 1; i <= 40; i++ { // mysql autoincrement starts from 1
		CommonLocation = append(CommonLocation, &model.Location{
			Comment: fmt.Sprintf("Location#%v", i),
			Address: fmt.Sprintf("Beijing Tsinghua Zijing #2 6%vB", i),
			Contact: fmt.Sprintf("ZWL%v", i),
			Phone:   "13416219802",
			UserID:  1, // when inserted into database, CommonUsers[constant.RoleCustomer].ID,
		})
	}

	for i := 0; i < 8; i++ {
		var paytype string = constant.PaytypeMonth
		if i&1 != 0 {
			paytype = constant.PaytypeYear
		}
		for j, status := range []string{constant.StatusWaiting, constant.StatusPass, constant.StatusOn, constant.StatusPause, constant.StatusRetrieve, constant.StatusSuspend, constant.StatusCanceled} {
			k := uint(1 + i*5 + j)
			CommonServices = append(CommonServices, &model.Service{
				Comment:    fmt.Sprintf("comment%v", k),
				Detail:     fmt.Sprintf("detail%v", k),
				PayType:    paytype,
				Status:     status,
				Require:    i,
				LocationID: k,
				Users: []*model.User{
					CommonUsers[constant.RoleCustomer],
					CommonUsers[constant.RoleOperator],
					CommonUsers[constant.RoleEngineer],
				},
			})
		}
	}

	// get NCE related test info
	NCETestConfig := config.NCE.Sub("test_site")
	CommonNCESiteID = NCETestConfig.GetString("id")
	CommonNCESiteName = NCETestConfig.GetString("name")
	CommonNCESSIDPath = fmt.Sprintf(constant.SsidFormatUrl, CommonNCESiteID)
}

func LoginAndGetToken(info *model.User) (string, string) {
	loginCase := TestCase{
		Desc:   "",
		Url:    "/user/login",
		Method: http.MethodPost,
		Param: user.UserLogin{
			Username: info.Username,
			Password: "TestPasswd@nsaop2021",
		},
		Code: http.StatusBadRequest,
		Msg:  constant.MsgOK,
	}
	writer, rsp := GetResponse(loginCase, -1, false)
	token := rsp.Data.(map[string]interface{})["token"].(string)
	refreshTokenCookie := writer.Result().Cookies()[0]
	refreshToken := refreshTokenCookie.Value
	return token, refreshToken
}

// SetupDatabaseForTest setups a in-memory db for testing purpose.
// Shouldn't be called out of test.
func SetupDatabaseForTest(base []string) func() {
	oldDB := model.DB
	tmpDB, err := gorm.Open(
		sqlite.Open(":memory:"),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	if err != nil {
		log.Fatalf("Fail to connect to database: %v", err)
	}
	model.DB = tmpDB
	model.Build()
	sqlDB, e := model.DB.DB()
	if e != nil {
		log.Fatalf("Fail to initialize MySQL: %v", e)
	}
	sqlDB.SetMaxOpenConns(1)

	for _, which := range base {
		switch which {
		case constant.TableUser:
			for _, user := range CommonUsersArray {
				model.DB.Create(user)
			}
		case constant.TableLocation:
			for _, location := range CommonLocation {
				model.DB.Create(location)
			}
		case constant.TableService:
			for _, service := range CommonServices {
				model.DB.Create(service)
			}
		}
	}

	return func() {
		model.DB = oldDB
	}
}
