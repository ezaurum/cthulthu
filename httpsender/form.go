package httpsender

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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
