package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"nsaop/config"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type FastHTTPResp struct {
	StatusCode int
	RespBody   gin.H
	Err        error
}

func MustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

/*
Unmarshal a byte array of json object into the target (some instance of struct)
Numbers of either float or int type will be processed correctly and independently
(instead of always float64)
*/
func UnmarshalWithNumber(jsonBytes []byte, target interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.UseNumber()
	return decoder.Decode(&target)
}

func NewBufferString(body string) io.Reader {
	return bytes.NewBufferString(body)
}

func Str2UUID(uuidStr string) (uuid.UUID, error) {
	return uuid.FromBytes([]byte(uuidStr))
}

func MustStr2UUID(uuidStr string) uuid.UUID {
	if id, err := Str2UUID(uuidStr); err != nil {
		return uuid.UUID{}
	} else {
		return id
	}
}

/*
Request url with method, headers, (query) args, and jsonBody
*/
func PlainRequest(
	method string, url string,
	headers map[string]string, args map[string]string, jsonBody gin.H) (
	fastHttpResp FastHTTPResp) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	// set headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	// append args to url
	if args != nil {
		var reqArgs fasthttp.Args
		for k, v := range args {
			reqArgs.Add(k, v)
		}
		url += "?"
		url += string(reqArgs.QueryString())
	}
	log.Printf("\tPlain Request to %s\n", url)
	req.Header.SetRequestURI(url)
	// set body
	if jsonBody != nil {
		req.SetBody(MustMarshal(jsonBody))
		req.Header.SetContentType("application/json")
	}
	resp := fasthttp.AcquireResponse()
	client := fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if fastHttpResp.Err = client.Do(req, resp); fastHttpResp.Err != nil {
		return
	}
	if fastHttpResp.Err = json.Unmarshal(resp.Body(), &fastHttpResp.RespBody); fastHttpResp.Err != nil {
		return
	}
	fastHttpResp.StatusCode = resp.StatusCode()
	return
}

/*
Perform a http request by fasthttp with given params and timeout
*/
func Request(
	method string, url string,
	headers map[string]string, args map[string]string, jsonBody gin.H,
	timeout time.Duration) (
	resp FastHTTPResp) {
	RespCh := make(chan FastHTTPResp, 1)
	go func() {
		RespCh <- PlainRequest(method, url, headers, args, jsonBody)
	}()
	for {
		select {
		case resp = <-RespCh:
			return
		case <-time.After(timeout * time.Millisecond):
			resp.StatusCode = http.StatusRequestTimeout
			resp.Err = errors.Errorf("request timeout (%s to %s), exceeds limit %dms", method, url, timeout.Milliseconds())
			return
		}
	}
}

// Ref: https://developers.google.com/recaptcha/docs/v3
func ReCAPTCHA(response string, expectedAction string, ip string) (bool, error) {
	if !config.NeedReCAPTCHA {
		return true, nil
	}
	if len(response) == 0 {
		return false, errors.Errorf("reCAPTCHA failed: failed to get response from client")
	}
	args := map[string]string{
		"secret":   config.ReCAPTCHA.GetString("reCAPTCHA_secret"),
		"response": response,
		"remoteip": ip,
	}
	resp := Request(
		http.MethodPost,
		"https://recaptcha.net/recaptcha/api/siteverify",
		nil,
		args,
		nil,
		config.NCETimeout,
	) // assert(StatusCode == 200)
	if resp.Err != nil {
		if resp.Err.Error()[:15] == "request timeout" {
			return true, resp.Err
		} else {
			return false, resp.Err
		}
	}
	// invalid token for this site
	// or unexpected action or low score
	if resp.RespBody["success"] == false {
		return false, errors.Errorf("Unsuccessful reCAPTCHA challenge")
	} else if resp.RespBody["action"] != expectedAction {
		return false, errors.Errorf("Unexpected action: %s", resp.RespBody["action"])
	} else if resp.RespBody["score"].(float64) < config.MinReCAPTCHAScore {
		return false, errors.Errorf("Low score: %s", resp.RespBody["score"])
	} else {
		return true, nil
	}
}
