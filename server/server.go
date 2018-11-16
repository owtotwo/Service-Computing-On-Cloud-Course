package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/owtotwo/Service-Computing-On-Cloud-Course/tools"
)

/* Server handles some requests from client and log some neccessary information.
 * The correct URLs as follows:
 * user/registe?username=XXX&password=XXX
 * todo/add?username=XXX&&password=XXX&item=XXX
 * todo/delete?username=XXX&password=XXX&itemIndex=XXX
 * todo/show?username=XXX&&password=XXX
 * Some incorrect URLs may receive error message, some may get 404 page.
 */

// handle "user/register"
// correct URL: user/registe?username=XXX&password=XXX
// return infomation of registering fail/successfully
func RegisterHandler(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm() // parsing the parameters
	var (
		username = req.FormValue("username")
		password = req.FormValue("password")
		message  string
	)
	if username == "" || password == "" { // if one of important parameters is empty
		message = messages["EmptyUsernameOrPassword"]
	} else if dealAddUserIntoDBFn(username, password) { // succeed to execute
		message = messages["RegisterSuccess"]
	} else {
		message = messages["RegisterFail"]
	}
	// render html template
	renderTemplate(writer, registerTemplate, &TodoList{Username: username, Todos: []string{}, Message: message})
	// log infomation on server
	tools.LogOKInfo(req.Method, "register")
}

// handle "todo/add"
// correct URL: todo/add?username=XXX&&password=XXX&item=XXX
// return infomation of adding item fail/successfully
func AddItemHandler(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var (
		username = req.FormValue("username")
		password = req.FormValue("password")
		item     = req.FormValue("item")
		message  string
	)
	if username == "" || password == "" {
		message = messages["EmptyUsernameOrPassword"]
	} else if dealAddItemIntoDBFn(username, password, item) {
		message = messages["AddSuccess"]
	} else {
		message = messages["AddFail"]
	}
	renderTemplate(writer, addItemTemplate, &TodoList{Username: username, Todos: []string{}, Message: message})
	tools.LogOKInfo(req.Method, "add")
}

// handle "todo/delete"
// correct URL: todo/delete?username=XXX&password=XXX&itemIndex=XXX
// return infomation of deleting item fail/successfully
func DeleteItemHandler(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var (
		username        = req.FormValue("username")
		password        = req.FormValue("password")
		itemIndexString = req.FormValue("index")
		message         string
	)

	if username == "" || password == "" {
		message = messages["EmptyUsernameOrPassword"]
	} else {
		itemIdex, err := strconv.Atoi(itemIndexString)
		if err != nil || !dealDeleteItemIntoDBFn(username, password, itemIdex - 1) {
			message = messages["DeleteFail"]
		} else {
			message = messages["DeleteSuccess"]
		}
	}
	renderTemplate(writer, deleteItemTemplate, &TodoList{Username: username, Todos: []string{}, Message: message})
	tools.LogOKInfo(req.Method, "delete")
}

// handle "todo/show"
// correct URL: todo/show?username=XXX&&password=XXX
// return infomation of showing items fail/successfully
func ShowListHandler(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var (
		username = req.FormValue("username")
		password = req.FormValue("password")
		message  string
		todoList []string
		result   bool
	)
	if username == "" || password == "" {
		message = messages["EmptyUsernameOrPassword"]
	} else {
		todoList, result = dealShowItemsFromDBFn(username, password)
		if result {
			message = fmt.Sprintf(messages["ShowSuccess"], len(todoList))
		} else {
			message = messages["ShowFail"]
		}
	}
	todoList = append([]string{""}, todoList...)
	renderTemplate(writer, ShowItemsTemplate, &TodoList{Username: username, Todos: todoList, Message: message})
	tools.LogOKInfo(req.Method, "show")
}

// handle URL that could not be know the purpose of client
// return 404 page
func OtherHandler(writer http.ResponseWriter, req *http.Request) {
	http.NotFound(writer, req)
	tools.LogNoFound(req.Method)
}

// server listens requests from client
// add "" tools.LogPortListening " to log when server begins to work
func ListenAndServe(addr string, handler http.Handler) error {
	tools.LogPortListening(addr[1:])
	server := &http.Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
