package server

import (
	"github.com/gorilla/sessions"
	"net/http"
	//"github.com/astaxie/beego/session"
	"strings"
)

const (
	SESSION_WEB = "SESSION_WEB"
	SESSION_FLASH_MESSAGE = "SESSION_FLASH_MESSAGE"
	SESSION_FLASH_OBJECT = "SESSION_FLASH_OBJECT"
	STORE_KEY = "AIRDISKCMS"
	TIMEFORMAT = "2017-02-17 117:57:48"
)

var (
	Store = sessions.NewCookieStore([]byte(STORE_KEY))
	SessionFlash *sessions.Session
	SessionWeb *sessions.Session
)

func SetFlashMessages(req *http.Request, w http.ResponseWriter, msg string)  {
	SessionFlash, _ = Store.Get(req, SESSION_FLASH_MESSAGE)
	SessionFlash.Options = &sessions.Options{
		Path:"/",
		MaxAge:300,
		HttpOnly:true,
	}

	SessionFlash.AddFlash(msg)
	SessionFlash.Save(req, w)
}

func GetFlashMessages(req *http.Request, w http.ResponseWriter) []interface{}{
	SessionFlash,_:=Store.Get(req, SESSION_FLASH_MESSAGE)
	var msg = SessionFlash.Flashes()
	SessionFlash.Save(req, w)
	return msg
}

func SetFlashObject(req *http.Request, w http.ResponseWriter, obj interface{})  {
	if obj != nil{
		SessionFlash,_ = Store.Get(req, SESSION_FLASH_OBJECT)
		SessionFlash.Options = &sessions.Options{
			Path:"/",
			MaxAge:300,
			HttpOnly:true,
		}
		SessionFlash.AddFlash(obj)
		SessionFlash.Save(req, w)
	}
}

func GetFlashObject(req *http.Request, w http.ResponseWriter) (bool, interface{}) {
	SessionFlash, _ = Store.Get(req, SESSION_FLASH_OBJECT)
	var obj = SessionFlash.Flashes()
	SessionFlash.Save(req, w)
	if size := len(obj); size > 0{
		return true, obj[size-1]
	}
	return false, nil
}

func SetSession(req *http.Request, w http.ResponseWriter, key string, value interface{}){
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	SessionWeb.Options = &sessions.Options{
		Path:"/",
		MaxAge: 60 * 10 * 24,
		HttpOnly: true,
	}
	SessionWeb.Values[key] = value
	SessionWeb.Save(req, w)
}

func GetSession(req *http.Request, key string) interface{} {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	return SessionWeb.Values[key]
}

//获取到key的value后，清除该session
func PopSession(req *http.Request, w http.ResponseWriter, key string) interface{} {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	v := SessionWeb.Values[key]
	SessionWeb.Values[key] = nil
	SessionWeb.Save(req, w)
	return v
}

func ClearSession(req *http.Request, w http.ResponseWriter, key string)  {
	SessionWeb, _ = Store.Get(req, SESSION_WEB)
	if nil != SessionWeb{
		flashes := SessionWeb.Flashes()
		size := len(flashes)
		if size > 0{
			for i := 0; i < size; i++{
				flashes[i] = nil
			}
			SessionWeb.Save(req, w)
		}
	}
}


//================== Controller ================= Redirect user to login
type CheckLogin struct {
	LoginUrl     string
	RememberNext bool
}

// NewLogger returns a new Logger instance
func NewCheckLogin() *CheckLogin {
	return &CheckLogin{"/account/login", true}
}

func (t *CheckLogin) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if GetSession(req, SESSION_WEB) == nil {
		var login_url = t.LoginUrl
		next := req.URL.Path
		next_lower := strings.ToLower(next)
		if t.RememberNext && (!strings.Contains(next_lower, "login") || !strings.Contains(next_lower, "logout")) {
			//login_url = login_url + "?next=" + next
			login_url = "/?next="+next
		}

		http.Redirect(rw, req, login_url, http.StatusFound)
		return
	} else {
		next(rw, req)
	}
}