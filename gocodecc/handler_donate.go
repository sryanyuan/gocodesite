package gocodecc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

type DonateRsp struct {
	Result int
	Msg    string
}

type orderCreateInfo struct {
	OrderID  string
	ApiID    string
	ApiKey   string
	Uid      int
	Num      int
	NumFloat float64
	// Config field
	CallHost   string
	CallSecret string
}

const (
	payMethodAlipayQR = iota
	payMethodWxQR
)

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

func donateCheckHandler(ctx *RequestContext) {
	var rsp DonateRsp
	rsp.Result = 1

	defer func() {
		jsonBytes, _ := json.Marshal(&rsp)
		ctx.w.Write(jsonBytes)
	}()

	vars := mux.Vars(ctx.r)
	orderID := vars["orderid"]

	if "" == orderID {
		rsp.Msg = "Invalid order id"
		return
	}
	orderStatus, err := checkDonateOrder(orderID)
	if nil != err {
		rsp.Msg = err.Error()
		return
	}

	rsp.Result = 0
	rsp.Msg = orderStatus
}

func createDonateOrder(user string, num int, pm int, debug bool) (*orderCreateInfo, error) {
	if "" == donateCall {
		return nil, errors.New("Donate not enabled")
	}

	requestURL := fmt.Sprintf("%s/ctrl?cmd=preinsertdonate&secret=%v&userid=%v&num=%v&paymethod=%v", donateCall, callSecret, user, num, pm)
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
	orderInfo.CallHost = donateCall
	orderInfo.CallSecret = callSecret
	orderInfo.NumFloat = float64(num)
	if debug {
		orderInfo.NumFloat = 0.01
	}
	return &orderInfo, nil
}

func confirmDonateOrder(uid string, orderID string, apikey string, total int) error {
	if "" == donateCall {
		return errors.New("Donate not enabled")
	}

	requestURL := fmt.Sprintf("%s/ctrl?cmd=insertdonatecb&secret=%v&addnum=%v&total=%v&apikey=%v&uid=%v", donateCall, callSecret, orderID, total, apikey, uid)
	rspData, err := doGet(requestURL, nil)
	if nil != err {
		return err
	}

	if string(rspData) == "success" {
		PushMessage("OrderConfirmed", fmt.Sprintf("%v_%v_pay_%v", orderID, uid, total))
		return nil
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

func checkDonateOrder(orderID string) (string, error) {
	if "" == donateCall {
		return "", errors.New("Donate not enabled")
	}

	requestURL := fmt.Sprintf("%s/ctrl?cmd=donatecheck&orderid=%v&secret=%v", donateCall, orderID, callSecret)
	rspData, err := doGet(requestURL, nil)
	if nil != err {
		return "", err
	}

	var rsp DonateRsp
	if err = json.Unmarshal(rspData, &rsp); nil != err {
		return "", err
	}

	if 0 != rsp.Result {
		return "", errors.New(rsp.Msg)
	}

	return rsp.Msg, nil
}
