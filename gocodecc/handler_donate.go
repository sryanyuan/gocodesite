package gocodecc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

type DonateRsp struct {
	Result int
	Msg    string
}

type orderCreateInfo struct {
	OrderID     string
	PpayOrderID string
	PpayURL     string
	ApiID       string
	ApiKey      string
	Uid         int
	Num         int
	NumFloat    float64
	QRUrl       string
	// Config field
	CallHost   string
	CallSecret string
}

const (
	payMethodAlipayQR = iota
	payMethodWxQR
	payMethodUnion
	payMethodWxPpay
	payMethodAliPpay
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
	ctx.r.ParseForm()
	account := ctx.r.FormValue("account")
	valueStr := ctx.r.FormValue("value")
	value := 0
	var err error
	if "" != valueStr {
		value, err = strconv.Atoi(valueStr)
		if nil != err {
			value = 0
		}
	}
	dataCtx := map[string]interface{}{
		"active":  "donate",
		"account": account,
		"value":   value,
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

func requestForPaymentURL(info *orderCreateInfo, config *AppConfig) (string, error) {
	urlBase := "https://yun.maweiwangluo.com/pay/union/submit.php"
	if config.DonateUnionURL != "" {
		urlBase = config.DonateUnionURL
	}
	args := map[string]string{
		"addnum":  info.OrderID,
		"total":   fmt.Sprintf("%.2f", info.NumFloat),
		"apiid":   info.ApiID,
		"apikey":  info.ApiKey,
		"showurl": fmt.Sprintf("%s/ajax/zfbqr_pay_confirm", config.Domain),
		"uid":     strconv.Itoa(info.Uid),
	}
	rspData, err := doGet(urlBase, args)
	if nil != err {
		return "", err
	}
	res := string(rspData)
	res = strings.Trim(res, "\xEF\xBB\xBF")
	return res, nil
}

type ppayResp struct {
	Code    int                    `json:"code"`
	Message string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
}

func createDonateOrderPpay(cfg *AppConfig, pm int, order *orderCreateInfo) (string, error) {
	if "" == cfg.Ppay.PayURL ||
		"" == cfg.Ppay.PayKey {
		return "", errors.New("Bad config of ppay")
	}

	sign := order.OrderID + strconv.Itoa(order.Uid)
	if pm == payMethodAlipayQR {
		sign += "2"
	} else {
		sign += "1"
	}
	sign += fmt.Sprintf("%.2f", order.NumFloat)
	sign += cfg.Ppay.PayKey

	payMethod := "1"
	if pm == payMethodAlipayQR {
		payMethod = "2"
	}

	tm := time.Now().UnixNano() / 1e6
	h := md5.New()
	h.Write([]byte(sign))
	sign = hex.EncodeToString(h.Sum(nil))

	requestURL := fmt.Sprintf("%s/createOrder?t=%v&sign=%v&type=%v&payId=%v&price=%.2f&param=%v",
		cfg.Ppay.PayURL, tm, sign, payMethod, order.OrderID, order.NumFloat, order.Uid)
	rspData, err := doPost(requestURL, nil, nil)
	if nil != err {
		return "", err
	}

	var rsp ppayResp
	if err = json.Unmarshal(rspData, &rsp); nil != err {
		return "", err
	}
	if nil == rsp.Data {
		if "" != rsp.Message {
			return "", errors.New(rsp.Message)
		}
		return "", errors.New("Null data")
	}
	oi, ok := rsp.Data["orderId"]
	if !ok {
		return "", errors.New("OrderId not found")
	}
	orderId, ok := oi.(string)
	if !ok {
		return "", errors.New("OrderId not string")
	}
	return orderId, nil
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

var (
	orderConfirmedMap  = make(map[string]struct{})
	orderConfirmedLock sync.Mutex
)

func confirmDonateOrder(uid string, orderID string, apikey string, total float64) error {
	if "" == donateCall {
		return errors.New("Donate not enabled")
	}

	requestURL := fmt.Sprintf("%s/ctrl?cmd=insertdonatecb&secret=%v&addnum=%v&total=%.2f&apikey=%v&uid=%v", donateCall, callSecret, orderID, total, apikey, uid)
	rspData, err := doGet(requestURL, nil)
	if nil != err {
		return err
	}

	if string(rspData) == "success" {
		pushMsg := fmt.Sprintf("%v_%v_pay_%v", orderID, uid, total)
		var pushed bool
		orderConfirmedLock.Lock()
		_, pushed = orderConfirmedMap[pushMsg]
		if !pushed {
			orderConfirmedMap[pushMsg] = struct{}{}
		}
		orderConfirmedLock.Unlock()
		if !pushed {
			PushMessage("OrderConfirmed", pushMsg)
		}
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
