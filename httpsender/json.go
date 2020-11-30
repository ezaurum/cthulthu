package httpsender

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func SendJSON(sendURL string, requestObject interface{}, responseObject interface{}, header map[string]string, method string) error {
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
	req, err := http.NewRequest(method, sendURL, buffer)
	for k, v := range header {
		req.Header.Set(k, v)
	}
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

// 그냥 메세지 전송용 유틸
func GetJSON(sendURL string,
	responseObject interface{}, header map[string]string) error {
	return SendJSON(sendURL, nil, responseObject, header, "GET")
}

// 그냥 메세지 전송용 유틸
func PostJSON(sendURL string,
	requestObject interface{},
	responseObject interface{}, header map[string]string) error {
	return SendJSON(sendURL, requestObject, responseObject, header, "POST")
}

// 그냥 메세지 전송용 유틸
func PatchJSON(sendURL string,
	requestObject interface{},
	responseObject interface{}, header map[string]string) error {
	return SendJSON(sendURL, requestObject, responseObject, header, "PATCH")
}

// 그냥 메세지 전송용 유틸
func PutJSON(sendURL string,
	requestObject interface{},
	responseObject interface{}, header map[string]string) error {
	return SendJSON(sendURL, requestObject, responseObject, header, "PUT")
}

// 그냥 메세지 전송용 유틸
func DeleteJSON(sendURL string,
	requestObject interface{},
	responseObject interface{}, header map[string]string) error {
	return SendJSON(sendURL, requestObject, responseObject, header, "DELETE")
}
