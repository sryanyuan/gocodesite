package gocodecc

import (
	"net/http"
)

const (
	apiMetaKeyPerm = "perm"
)

func apiWrapper(handler HttpHandler) HttpHandler {
	return func(ctx *RequestContext) {
		user := ctx.user
		// Permission check, default is guest
		pem, ok := ctx.ri.Meta.GetInt(apiMetaKeyPerm)
		if !ok {
			pem = kPermission_Guest
		}
		if pem == kPermission_None ||
			user.Permission < uint32(pem) {
			// Need login
			if user.Uid == 0 {
				ctx.WriteAPIRspBadNeedLogin("")
				return
			} else {
				// Permission denied
				ctx.WriteAPIRsp(http.StatusForbidden, nil)
				return
			}
		}
		// Pass
		handler(ctx)
	}
}

func registerApi(path string, pem uint32, handler HttpHandler, methods []string) {
	ri := RouterItem{
		Url:        path,
		Permission: kPermission_Guest,
		Handler:    apiWrapper(handler),
		Methods:    methods,
	}
	if kPermission_Guest != pem {
		meta := map[string]interface{}{
			apiMetaKeyPerm: int(pem),
		}
		ri.Meta.Init(meta)
	}
	routerItems = append(routerItems, ri)
}
