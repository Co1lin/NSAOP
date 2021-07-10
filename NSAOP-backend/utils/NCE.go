package utils

import (
	"fmt"
	"net/http"
	"nsaop/config"
	"nsaop/utils/constant"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	lastTimeGotToken = time.Now()
	NCEToken         = ""
)

const NCESiteURL = "/controller/campus/v3/sites"

func GetNCEToken() (string, error) {
	if len(NCEToken) == 0 || time.Since(lastTimeGotToken) > 10*time.Minute {
		url := config.NCEHost + "/controller/v2/tokens"
		body := gin.H{
			"userName": config.NCEUsername,
			"password": config.NCEPassword,
		}
		resp := Request(http.MethodPost, url, nil, nil, body, config.NCETimeout)
		if resp.Err != nil || resp.StatusCode != http.StatusOK {
			return "", errors.New(fmt.Sprintf("get NCE Token Error with HTTP Code: %d", resp.StatusCode))
		} else {
			lastTimeGotToken = time.Now()
			NCEToken = resp.RespBody["data"].(map[string]interface{})["token_id"].(string)
			return NCEToken, nil
		}
	} else {
		return NCEToken, nil
	}
}

func RequestNCEWithToken(
	method string, path string,
	headers map[string]string, params map[string]string, body gin.H) (
	resp FastHTTPResp) {
	token, err := GetNCEToken()
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Err = err
		return
	} else {
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["X-AUTH-TOKEN"] = token
		resp = Request(method, config.NCEHost+path, headers, params, body, config.NCETimeout)
		if NCERespondWithFail(resp.RespBody) {
			resp.StatusCode = http.StatusInternalServerError
		}
		return
	}
}

/*
create a NCE site with given siteName
return the ID of created NCE site
*/
func CreateNCESite(siteName string) (NCESiteID string, err error) {
	resp := RequestNCEWithToken(
		http.MethodPost,
		NCESiteURL,
		nil,
		nil,
		gin.H{
			"sites": []gin.H{
				{
					"name": siteName,
				},
			},
		},
	)
	if resp.StatusCode == http.StatusOK {
		NCESiteID = resp.RespBody["success"].([]interface{})[0].(map[string]interface{})["id"].(string)
		return
	} else {
		err = errors.Errorf("Got HTTP Code: %d when creating NCE Site", resp.StatusCode)
		return
	}
}

func GetNCESites() (sites []string, err error) {
	resp := RequestNCEWithToken(
		http.MethodGet,
		NCESiteURL,
		nil,
		nil,
		nil,
	)
	if resp.StatusCode == http.StatusOK {
		sitesList := resp.RespBody["data"].([]interface{})
		for _, site := range sitesList {
			sites = append(sites, site.(map[string]interface{})["id"].(string))
		}
		return
	} else {
		err = errors.Errorf("Got HTTP Code: %d when getting NCE Sites", resp.StatusCode)
		return
	}
}

func DeleteNCESites(sites []string) (err error) {
	resp := RequestNCEWithToken(
		http.MethodDelete,
		NCESiteURL,
		nil,
		nil,
		gin.H{
			"ids": sites,
		},
	)
	if resp.StatusCode == http.StatusOK {
		return
	} else {
		err = errors.Errorf("Got HTTP Code: %d when deleting NCE Sites", resp.StatusCode)
		return
	}
}

func getSSIDsByNCESiteID(siteID string) (ssids []string, err error) {
	if ssidsResp := RequestNCEWithToken(
		http.MethodGet,
		fmt.Sprintf(constant.SsidFormatUrl, siteID),
		nil, nil, nil); ssidsResp.Err != nil || ssidsResp.StatusCode != http.StatusOK ||
		len(ssidsResp.RespBody["data"].([]interface{})) == 0 {
		err = errors.Errorf("Failed to get SSIDs with HTTP Code: %d", ssidsResp.StatusCode)
		return
	} else {
		// get IDs of SSIDs
		ssidsList := ssidsResp.RespBody["data"].([]interface{})
		for _, ssid := range ssidsList {
			ssids = append(ssids, ssid.(map[string]interface{})["id"].(string))
		}
		return
	}
}

/*
change status of all of the SSIDs of the given NCE Site to the target
*/
func ChangeSSIDStatus(siteID string, targetStatus bool) error {
	if ssids, err := getSSIDsByNCESiteID(siteID); err != nil {
		return err
	} else {
		url := fmt.Sprintf("/controller/campus/v3/networkconfig/site/%s/apssid/", siteID)
		var err error
		for _, ssid := range ssids {
			resp := RequestNCEWithToken(
				http.MethodPut,
				url+ssid,
				nil, nil,
				gin.H{"enable": targetStatus})
			if resp.Err != nil || resp.StatusCode != http.StatusOK {
				err = errors.Errorf("%s\nError occurred when changing SSID %s; HTTP Code: %d",
					resp.Err.Error(), ssid, resp.StatusCode)
			}
		}
		return err
	}
}

func NCERespondWithFail(resp gin.H) bool {
	fail := resp["fail"]
	if fail != nil && len(fail.([]interface{})) > 0 {
		return true
	} else {
		return false
	}
}
