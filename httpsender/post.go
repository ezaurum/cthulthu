package httpsender

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 그냥 메세지 전송용 유틸
func PostJSON(sendURL string,
	requestObject interface{},
	responseObject interface{}, header map[string]string) error {

	if _, b := header["Content-Type"]; !b {
		header["Content-Type"] = "application/json;charset=UTF-8"
	}

	if _, b := header["Accept"]; !b {
		header["Accept"] = "application/json"
	}

	marshal, e := json.Marshal(requestObject)
	if e != nil {
		return e
	}
	buffer := bytes.NewBuffer(marshal)
	req, err := http.NewRequest("POST", sendURL, buffer)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	resBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(resBody, responseObject)
	if nil != err {
		return err
	}

	return nil
}

// 그냥 메세지 전송용 유틸
func PostForm(sendURL string,
	form url.Values,
	responseObject interface{}) error {

	encode := form.Encode()
	reader := strings.NewReader(encode)
	req, err := http.NewRequest("POST", sendURL, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if nil != resp {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	code := resp.StatusCode
	resBody, _ := ioutil.ReadAll(resp.Body)
	s := string(resBody)
	if 400 <= code {
		return errors.WithMessage(errors.Errorf("status error %d", code), s)
	}
	err = xml.Unmarshal(resBody, responseObject)
	if nil != err {
		return err
	}
	return nil
}
