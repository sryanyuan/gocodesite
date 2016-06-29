package gocodecc

type WebUser struct {
	Uid        uint32
	Permission uint32
	UserName   string
	Avatar     string
	Sex        int
	NickName   string
	EMail      string
}

func defaultWebUser() *WebUser {
	user := &WebUser{
		Uid:        0,
		Permission: kPermission_Guest,
		UserName:   "Guest",
		Avatar:     "",
	}
	return user
}
