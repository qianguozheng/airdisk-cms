package server

import (
	"database/sql"
	"html/template"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"encoding/gob"
	//"github.com/google/martian/log"
	//log "github.com/Sirupsen/logrus"
	"os"
	"fmt"
)

var db *sql.DB
var tpl *template.Template

type Upgrade struct {
	Mac string //Should Remove it later
	Url string
	Version string
	Md5 string
}

type AirDisk struct {
	Mac string
	Upgrade int
	Control int
}

type Account struct {
	Id        int64
	UserName  string
	Password  string
}
func init() {
	var err error
	db, err = sql.Open("sqlite3", "./airdisk.db")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

	//Session 使用前需要注册数据结构
	gob.Register(&Account{})
}

func InitLog(file string) *os.File{

	f, err := os.OpenFile(file, os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Open Logfile failed")
		return nil
	}
	//log.SetOutput(f)
	//defer f.Close()
	return f
}
func Run()  {
	//fileHandle := InitLog("airdisk-cms.log")
	//defer fileHandle.Close()
	fmt.Println("WebServer Start")

	router := mux.NewRouter()
	adminRoutes := mux.NewRouter()

	router.HandleFunc("/", index)
	router.HandleFunc("/register", register)
	router.HandleFunc("/account/register", registerCheck)

	adminRoutes.HandleFunc("/admin/", adminIndex)
	adminRoutes.HandleFunc("/admin/upgradeInfo", upgradeInfo)
	adminRoutes.HandleFunc("/admin/upgrade/create", upgradeCreateForm)
	adminRoutes.HandleFunc("/admin/upgrade/create/process", upgradeCreateProcess)
	adminRoutes.HandleFunc("/admin/upgrade/update", upgradeUpdateForm)
	adminRoutes.HandleFunc("/admin/upgrade/update/process", upgradeUpdateProcess)
	adminRoutes.HandleFunc("/admin/upgrade/delete/process", upgradeDeleteProcess)
	adminRoutes.HandleFunc("/admin/controlInfo", controlInfo)
	adminRoutes.HandleFunc("/admin/control/create", controlCreateForm)
	adminRoutes.HandleFunc("/admin/control/create/process", controlCreateProcesss)

	router.PathPrefix("/admin").Handler(negroni.New(
		NewCheckLogin(),
		negroni.Wrap(adminRoutes),
	))
	// /account/login
	router.HandleFunc("/account/login", login)
	//router.Handle("/static/", http.FileServer(http.Dir("/Users/weeds/Documents/go/weeds/airdisk-cms/static/")))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":8080", router)

}

func adminIndex(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("access adminIndex")
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func index(w http.ResponseWriter, r *http.Request)  {
	//http.Redirect(w,r, "/upgradeInfo", http.StatusSeeOther)
	fmt.Println("access index")
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func register(w http.ResponseWriter, r *http.Request)  {
	//http.Redirect(w,r, "/upgradeInfo", http.StatusSeeOther)
	fmt.Println("access register")
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

func registerCheck(w http.ResponseWriter, r *http.Request)  {
	//http.Redirect(w,r, "/upgradeInfo", http.StatusSeeOther)
	fmt.Println("request:", r.FormValue("username"))
	fmt.Println("request:", r.FormValue("password"))
	fmt.Println("request:", r.FormValue("email"))

	//tpl.ExecuteTemplate(w, "login.gohtml", nil)
	//w.Write([]byte("OK"))
	type V struct {
		Result bool
		//Test bool
	}
	var v V
	v.Result = false
	//v.Test = false
	tpl.ExecuteTemplate(w, "register_status.gohtml", &v)
	//http.RedirectHandler("http://localhost:8080/acount/login", http.StatusPermanentRedirect)
}

func login(w http.ResponseWriter, req *http.Request)  {
	//http.Redirect(w,r, "/upgradeInfo", http.StatusSeeOther)
	//tpl.ExecuteTemplate(w, "login.gohtml", nil)
	fmt.Println("access login function")

	ctx := make(map[string]interface{})

	var next = req.FormValue("next")
	ctx["next"] = next
	//
	//ctx[csrf.TemplateTag] = csrf.TemplateField(req)
	//var account = Account{}
	var account = Account{Id:1, UserName:"hello",Password:"admin"}
	if req.Method == "POST" {
		var username = req.FormValue("username")
		var password = req.FormValue("password")

		if username != "" && password != "" {
			//err := dbMap.SelectOne(&account, "select * from account where username=$1 or email=$2 limit 1", username, username)
			//
			//if err != nil && account.Id > 0 {
			//	SetFlashMessages(req, w, "用户不存在！")
			//	http.Redirect(w, req, "/account/login", http.StatusFound)
			//	return
			//}
			//
			//if account.Password == lib.MD5(password) {
			if password == "admin" {
				SetSession(req, w, SESSION_WEB, account)
			//	dbMap.Exec("update account set lastlogin=$1 where $2", time.Now(), account.Id)
			//
				if next != "" {
					http.Redirect(w, req, next, http.StatusFound)
					return
				} else {
					//http.Redirect(w, req, "/admin/upgradeInfo", http.StatusFound)
					//return
				}
			} else {
				SetFlashMessages(req, w, "账号和密码不匹配")
				http.Redirect(w, req, "/account/login", http.StatusFound)
				return
			}
			//fmt.Println("account:", username, " password:", password)
			fmt.Println("account:", username, " password:", password)
			tpl.ExecuteTemplate(w, "index.gohtml", nil)
		}
	} else {
		//dbMap.SelectOne(&account, "select * from account limit 1")
		//if account.Id <= 0 {
		//	account.UserName = "admin"
		//	account.Password = lib.MD5("admin")
		//	account.Status = true
		//	dbMap.Insert(&account)
		//	SetFlashMessages(req, w, "初始化账号和密码：admin / admin，请及时更改密码！")
		//}
		//
		ctx["flashes"] = GetFlashMessages(req, w)

		//if currentUser := GetSession(req, SESSION_WEB); currentUser != nil {
		//	if next == "" {
		//		next = "/"
		//	}
		//	http.Redirect(w, req, next, http.StatusFound)
		//	return
		//}
		//
		//r.HTML(w, http.StatusOK, "account/login", ctx, render.HTMLOptions{Layout: ""})
		//return
	}
}

func controlInfo(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET"{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("select * from airdisk")
	if err != nil{
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	airs := make([]AirDisk, 0)

	for rows.Next(){
		air := AirDisk{}
		err := rows.Scan(&air.Mac, &air.Upgrade, &air.Control)
		if err != nil{
			http.Error(w, http.StatusText(500), 500)
			return
		}
		airs =append(airs, air)
	}

	if err = rows.Err(); err != nil{
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpl.ExecuteTemplate(w, "controlInfo.gohtml", airs)
}
func controlCreateForm(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "create_control.gohtml", nil)
}
func controlCreateProcesss(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	bk := AirDisk{}
	bk.Mac = r.FormValue("mac")
	bk.Control = 1
	//bk.Control = r.FormValue("control")

	// validate form values
	if bk.Mac == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}


	// Update airdisk table values.
	_, err:= db.Exec("update airdisk set control=1 where mac = $1", bk.Mac)
	if err != nil{
		fmt.Println(err.Error())
	}

	// confirm insertion
	tpl.ExecuteTemplate(w, "created_control.gohtml", bk)
}
func upgradeInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET"{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("select * from upgrade")
	if err != nil{
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	upgs := make([]Upgrade, 0)

	for rows.Next(){
		upg := Upgrade{}
		err := rows.Scan(&upg.Url, &upg.Version, &upg.Md5, &upg.Mac)
		if err != nil{
			http.Error(w, http.StatusText(500), 500)
			return
		}
		upgs =append(upgs, upg)
	}

	if err = rows.Err(); err != nil{
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tpl.ExecuteTemplate(w, "upgradeInfo.gohtml", upgs)
}

func upgradeCreateForm(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "create.gohtml", nil)
}

func upgradeCreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	bk := Upgrade{}
	bk.Mac = r.FormValue("mac")
	bk.Url = r.FormValue("url")
	bk.Version = r.FormValue("version")
	bk.Md5 = r.FormValue("md5")

	// validate form values
	if bk.Mac == "" || bk.Url == "" || bk.Version == "" || bk.Md5 == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}


	// insert values
	_, err := db.Exec("INSERT INTO upgrade (Url, Version, Md5, Mac) VALUES ($1, $2, $3, $4)", bk.Url, bk.Version, bk.Md5, bk.Mac)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Update airdisk table values.
	_, err= db.Exec("update airdisk set upgrade=1 where mac = $1", bk.Mac)
	if err != nil{
		fmt.Println(err.Error())
	}

	// confirm insertion
	tpl.ExecuteTemplate(w, "created.gohtml", bk)
}

func upgradeUpdateForm(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	mac := r.FormValue("mac")
	if mac == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM upgrade WHERE mac = $1", mac)

	bk := Upgrade{}
	err := row.Scan(&bk.Url, &bk.Version, &bk.Md5, &bk.Md5)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(w, "update.gohtml", bk)
}
func upgradeUpdateProcess(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	bk := Upgrade{}
	bk.Mac = r.FormValue("mac")
	bk.Version = r.FormValue("version")
	bk.Md5 = r.FormValue("md5")
	bk.Url = r.FormValue("url")

	// validate form values
	if bk.Mac == "" || bk.Version == "" || bk.Md5 == "" || bk.Url == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// insert values
	_, err := db.Exec("UPDATE upgrade SET mac = $1, url=$2, version=$3, md5=$4 WHERE mac=$1;", bk.Mac, bk.Url, bk.Version, bk.Md5)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Update airdisk table values.
	_, err= db.Exec("update airdisk set upgrade=1 where mac = $1", bk.Mac)
	if err != nil{
		fmt.Println(err.Error())
	}

	// confirm insertion
	tpl.ExecuteTemplate(w, "updated.gohtml", bk)
}
func upgradeDeleteProcess(w http.ResponseWriter, r*http.Request)  {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	mac := r.FormValue("mac")
	if mac == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// delete book
	_, err := db.Exec("DELETE FROM upgrade WHERE mac=$1;", mac)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Update airdisk table values.
	_, err= db.Exec("update airdisk set upgrade=0 where mac = $1", mac)
	if err != nil{
		//fmt.Println(err.Error())
		fmt.Println(err.Error())
	}

	http.Redirect(w, r, "/upgradeInfo", http.StatusSeeOther)
}