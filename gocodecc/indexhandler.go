package gocodecc

func indexHandler(ctx *RequestContext) {
	ctx.w.Write([]byte("hello index"))
}
