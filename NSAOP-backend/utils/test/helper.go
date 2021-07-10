package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"

	"nsaop/config"
	"nsaop/router"
	"nsaop/router/resp"
	"nsaop/utils"
)

type TestCase struct {
	Url         string            // request Url
	Method      string            // request Method
	Param       interface{}       // request param (in body)
	ContentType string            // request content-type
	Headers     map[string]string // request headers
	Cookies     []http.Cookie     // cookies

	Code     int         // expected HTTP status Code
	Msg      string      // expected message
	Data     interface{} // expected Data
	Desc     string      // description of the case
	ShowBody bool        // whether to show the body of the response in log
}

func performRequest(url, method, contentType string, param interface{}, headers map[string]string, cookies []http.Cookie) (writer *httptest.ResponseRecorder) {
	var body []byte
	if method == "GET" {
		if v, err := query.Values(param); err != nil {
			log.Fatal("param parse to query fail")
		} else if q := v.Encode(); q != "" {
			url += "?" + strings.ToLower(q)
		}
	} else {
		body = utils.MustMarshal(param)
	}
	baseUrl := "/" + config.Router.GetString("version")
	url = baseUrl + url
	fmt.Println(url)
	// setup router
	rt := router.InitRouter()
	router.SetupRouter(rt)
	// perform request
	writer = httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}
	context.Request = req
	context.Request.Header.Set("Content-Type", contentType)
	rt.ServeHTTP(writer, req)
	return
}

func GetResponse(testCase TestCase, k int, logging bool) (w *httptest.ResponseRecorder, rsp resp.Response) {
	w = performRequest(testCase.Url, testCase.Method, testCase.ContentType, testCase.Param, testCase.Headers, testCase.Cookies)
	if logging {
		fmt.Printf("Test case No. %d：%s\n", k, testCase.Desc)
		if testCase.ShowBody {
			fmt.Printf("API returns: %s\n", w.Body.String())
		}
	}
	rsp = resp.Response{}
	err := json.Unmarshal(w.Body.Bytes(), &rsp)
	if err != nil {
		log.Fatalf("GetResponse Body Unmarshal error != nil: %v\nwhich body is: %v", err.Error(), w.Body.String())
	}
	return
}

func RequestOnly(testCase TestCase) {
	performRequest(testCase.Url, testCase.Method, testCase.ContentType, testCase.Param, testCase.Headers, testCase.Cookies)
}

func PerformTest(t *testing.T, testCases []TestCase, log bool) {
	assert := assert.New(t)
	// expected: testCase; actual: w / resp
	for k, testCase := range testCases {
		k := k
		copyCase := testCase
		t.Run(copyCase.Desc, func(t *testing.T) {
			t.Parallel()
			w, rsp := GetResponse(copyCase, k, log)
			if copyCase.Desc == "wtf" {
				t.Log("*################################ what the fuck here ###############################")
				t.Log(copyCase.Code, w.Code)
				t.Log(copyCase.Msg, rsp.Msg)
				t.Log(copyCase.Data, rsp.Data)
			}
			// compare code
			assert.Equal(copyCase.Code, w.Code, "Got unexpected HTTP status code.")
			// compare msg
			assert.Equal(copyCase.Msg, rsp.Msg, "Got unexpected msg in response.")
			// compare data
			if copyCase.Data != nil {
				tmp := reflect.New(reflect.TypeOf(copyCase.Data))
				actualData := tmp.Interface() // get a naive struct with the same type of copyCase.Data
				tmp = reflect.New(reflect.TypeOf(copyCase.Data))
				tmp.Elem().Set(reflect.ValueOf(copyCase.Data))
				expectData := tmp.Interface()

				decoderConfig := &mapstructure.DecoderConfig{TagName: "json", Result: &actualData}
				decoder, err := mapstructure.NewDecoder(decoderConfig)
				err = decoder.Decode(rsp.Data)
				assert.Equal(err, nil, "Got unexpected data in response")
				assert.Equal(expectData, actualData, "Got unexpected data in response.")
			}
		}) // end Run
	} // end for
}
