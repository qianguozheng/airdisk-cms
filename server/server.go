package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
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
}
func Run()  {
	http.HandleFunc("/", index)
	http.HandleFunc("/upgradeInfo", upgradeInfo)
	http.HandleFunc("/upgrade/create", upgradeCreateForm)
	http.HandleFunc("/upgrade/create/process", upgradeCreateProcess)
	http.HandleFunc("/upgrade/update", upgradeUpdateForm)
	http.HandleFunc("/upgrade/update/process", upgradeUpdateProcess)
	http.HandleFunc("/upgrade/delete/process", upgradeDeleteProcess)
	http.HandleFunc("/controlInfo", controlInfo)
	http.HandleFunc("/control/create", controlCreateForm)
	http.HandleFunc("/control/create/process", controlCreateProcesss)
	http.ListenAndServe(":8080", nil)

}

func index(w http.ResponseWriter, r *http.Request)  {
	//http.Redirect(w,r, "/upgradeInfo", http.StatusSeeOther)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
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
		fmt.Println(err.Error())
	}

	http.Redirect(w, r, "/upgradeInfo", http.StatusSeeOther)
}