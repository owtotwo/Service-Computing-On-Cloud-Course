package server

import (
	"database/sql"
	"go/build"
	"html/template"
)

// template files for sending to client
// html template files path
const (
	registerTemplate   = "register.html"
	addItemTemplate    = "add.html"
	deleteItemTemplate = "delete.html"
	ShowItemsTemplate  = "show.html"
)

var htmlPosition string = build.Default.GOPATH + "/src/github.com/owtotwo/Service-Computing-On-Cloud-Course/html/"

var htmlFilesNames = []string{
	htmlPosition + registerTemplate,
	htmlPosition + addItemTemplate,
	htmlPosition + deleteItemTemplate,
	htmlPosition + ShowItemsTemplate,
}

var templates = template.Must(template.ParseFiles(htmlFilesNames...))

// -------------------------------------------------

// database names and sqlstatements
var (
	dbName       = "cloudgo"     // database name
	dbTableName  = "todos"       // a table name
	dbPara       string          // database open parameter
)

// database execute statements
var dbStatements = map[string]string{
	"CREATETABLE": "CREATE TABLE IF NOT EXISTS " + dbTableName +
		" (username varchar(255) PRIMARY KEY, password varchar(255) NOT NULL, todos Text)",
	"REGISTER":  "INSERT INTO " + dbTableName + " (username, password, todos) values (?, ?, ?)",
	"EDITTODOS": "UPDATE " + dbTableName + " set todos=? WHERE username=? AND password=?",
	"SHOWTODOS": "SELECT todos FROM " + dbTableName + " WHERE username=? AND password=?",
	"QUERYUSER": "SELECT username FROM " + dbTableName + " WHERE username=? AND password=?",
}

// -------------------------------------------------------------

// message to client
var messages = map[string]string{
	"EmptyUsernameOrPassword": "username and password should be non-empty",
	"RegisterSuccess":         "register success",
	"RegisterFail":            "register fail: the username may have been used",
	"AddSuccess":              "add success",
	"AddFail":                 "add fail: please check username and password and the item should be non-empty",
	"DeleteSuccess":           "delete success",
	"DeleteFail":              "delete fail: please check username and password and the item index should be valid",
	"ShowSuccess":             "show success： you have %d todo items",
	"ShowFail":                "show fail: please check username and password",
}

// exec some simple sql statement
func dbExec(db *sql.DB, DBStatement string) {
	if _, err := db.Exec(DBStatement); err != nil {
		db.Close()
		panic(err)
	}
}

// create database and table when init package
func init() {
	db, err := openDB()
	if err != nil {
		panic(err)
	}

	dbExec(db, dbStatements["CREATETABLE"]) // create table in database

	db.Close()
}
