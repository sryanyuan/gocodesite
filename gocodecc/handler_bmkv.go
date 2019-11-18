package gocodecc

type bmkvResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func bmkvGetHandler(ctx *RequestContext) {
	var rsp bmkvResp
	rsp.Code = 1

	ctx.r.ParseForm()
	key := ctx.r.FormValue("key")
	if "" == key {
		rsp.Message = "Key required"
		ctx.WriteJSONResponse(&rsp)
		return
	}

	value, err := modelBmkvGet(key)
	if nil != err {
		rsp.Message = err.Error()
		ctx.WriteJSONResponse(&rsp)
		return
	}

	rsp.Code = 0
	rsp.Message = value
	ctx.WriteJSONResponse(&rsp)
}
