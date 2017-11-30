package gocodecc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type MsgPushConfig struct {
	Host  string `toml:"host"`
	SCKey string `toml:"sc-key"`
}

type msgPushRsp struct {
	ErrNo  int    `json:"errno"`
	ErrMsg string `json:"errmsg"`
}

var (
	globalMsgPushConfig *MsgPushConfig
)

func PushMessage(title, body string) error {
	if nil == globalMsgPushConfig {
		return errors.New("Msg push not config")
	}
	if "" == globalMsgPushConfig.Host ||
		"" == globalMsgPushConfig.SCKey {
		return errors.New("Invalid msg push config")
	}
	requestURL := fmt.Sprintf("%s/%v.send", globalMsgPushConfig.Host, globalMsgPushConfig.SCKey)
	statusCode, rspData, err := doFormPost(requestURL, url.Values{"text": {title}, "desp": {body}})
	if nil != err {
		return err
	}
	if http.StatusOK != statusCode {
		return errors.New("Msg push error")
	}

	var rsp msgPushRsp
	if err = json.Unmarshal(rspData, &rsp); nil != err {
		return err
	}

	if rsp.ErrNo != 0 {
		return errors.New(rsp.ErrMsg)
	}

	return nil
}
