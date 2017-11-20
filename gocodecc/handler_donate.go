package gocodecc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cihub/seelog"
)

type DonateRsp struct {
	Result int
	Msg    string
}

type orderCreateInfo struct {
	OrderID string
	ApiID   string
	ApiKey  string
	Uid     int
	Num     int
}

var (
	donateCall string
	callSecret string
)

func initDonateCall(addr string, secret string) {
	donateCall = addr
	callSecret = secret
}

func donateHander(ctx *RequestContext) {
	dataCtx := map[string]interface{}{
		"active": "donate",
	}
	dataHTML := renderTemplate(ctx, []string{"template/donate.html"}, dataCtx)
	ctx.w.Write(dataHTML)
}

func createDonateOrder(user string, num int) (*orderCreateInfo, error) {
	if "" == donateCall {
		return nil, errors.New("Donate not enabled")
	}

	requestURL := fmt.Sprintf("%s/ctrl?cmd=preinsertdonate&secret=%v&userid=%v&num=%v", donateCall, callSecret, user, num)
	rspData, err := doGet(requestURL, nil)
	if nil != err {
		return nil, err
	}

	var rsp DonateRsp
	if err = json.Unmarshal(rspData, &rsp); nil != err {
		return nil, err
	}

	if 0 != rsp.Result {
		return nil, errors.New(rsp.Msg)
	}

	seelog.Info(rsp.Msg)
	var orderInfo orderCreateInfo
	if err = json.Unmarshal([]byte(rsp.Msg), &orderInfo); nil != err {
		return nil, err
	}
	return &orderInfo, nil
}

func confirmDonateOrder(orderID string, apikey string, total int) error {
	if "" == donateCall {
		return errors.New("Donate not enabled")
	}

	requestURL := fmt.Sprintf("%s/ctrl?cmd=insertdonatecb&secret=%v&addnum=%v&total=%v&apikey=%v", donateCall, callSecret, orderID, total, apikey)
	rspData, err := doGet(requestURL, nil)
	if nil != err {
		return err
	}

	var rsp DonateRsp
	if err = json.Unmarshal(rspData, &rsp); nil != err {
		return err
	}

	if 0 != rsp.Result {
		return errors.New(rsp.Msg)
	}

	return nil
}

func doGet(reqUrl string, args map[string]string) ([]byte, error) {
	u, _ := url.Parse(strings.Trim(reqUrl, "/"))
	q := u.Query()
	if nil != args {
		for arg, val := range args {
			q.Add(arg, val)
		}
	}

	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf("Http statusCode:%d", res.StatusCode))
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
