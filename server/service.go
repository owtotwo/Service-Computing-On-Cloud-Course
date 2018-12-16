package server

import (
	"io"
	"fmt"
	"strconv"
	"crypto/md5"
)

type TodoListOperationResult struct {
	Operation string
	Username  string
	Todos     []string
	Result    string
	Message   string
}

// the details of user request results
var messages = map[string]string{
	// empty username and password
	"EmptyUsernameOrPassword": "username and password should be non-empty",
	// register
	"RegisterSuccess": "register success",
	"RegisterFail":    "register fail",
	// add item
	"AddSuccess": "add success",
	"EmptyItem":  "add fail: the item should be non-empty",
	"AddFail":    "add fail: please check username and password",
	// delete item
	"DeleteSuccess": "delete success",
	"DeleteFail":    "delete fail: please check username and password and the item index should be valid",
	// show items
	"ShowSuccess": "show success： you have %d todo items",
	"ShowFail":    "show fail: please check username and password",
}

// message to client
const (
	// the [result] of user's request
	SUCCESS = "success"
	FAIL    = "fail"
	// the [operation] of user's request
	REGISTER = "register"
	ADD      = "add"
	DELETE   = "delete"
	SHOW     = "show"
)

// judge if username or password is empty
func isEmptyUsernameOrPassword(username string, password string) bool {
	return username == "" || password == ""
}

// check if username or password is empty
// add user into datebase and check if it succeeded
// return result and details
func registerByUsernameAndPassword(username string, password string) TodoListOperationResult {
	if isEmptyUsernameOrPassword(username, password) { // if one of important parameters is empty
		return TodoListOperationResult{Operation: REGISTER, Username: username, Result: FAIL, Message: messages["EmptyUsernameOrPassword"]}
	} else if _, ok := operateOnDB(addUserIntoDB, username, tools.MD5Encryption(password), nil); ok {
		return TodoListOperationResult{Operation: REGISTER, Username: username, Result: SUCCESS, Message: messages["RegisterSuccess"]}
	}
	return TodoListOperationResult{Operation: REGISTER, Username: username, Result: FAIL, Message: messages["RegisterFail"]}
}

// check if item is empty
func isEmptyItem(item string) bool {
	return item == ""
}

// check if username or password is empty
// add item into datebase, check if item is empty and check if succeeded
// return result and details
func addItemByUsernamePasswordItem(username string, password string, item string) TodoListOperationResult {
	if isEmptyUsernameOrPassword(username, password) { // if one of important parameters is empty
		return TodoListOperationResult{Operation: ADD, Username: username, Result: FAIL, Message: messages["EmptyUsernameOrPassword"]}
	} else if isEmptyItem(item) {
		return TodoListOperationResult{Operation: ADD, Username: username, Result: FAIL, Message: messages["EmptyItem"]}
	} else if _, ok := operateOnDB(addItemIntoDB, username, tools.MD5Encryption(password), item); ok {
		return TodoListOperationResult{Operation: ADD, Username: username, Result: SUCCESS, Message: messages["AddSuccess"]}
	}
	return TodoListOperationResult{Operation: ADD, Username: username, Result: FAIL, Message: messages["AddFail"]}
}

// check if username or password is empty
// delete item from datebase and check if succeeded
// return result and details
func deleteItemByUsernameAndPasswordItemindex(
	username string, password string, itemIndexString string) TodoListOperationResult {
	if isEmptyUsernameOrPassword(username, password) { // if one of important parameters is empty
		return TodoListOperationResult{Operation: DELETE, Username: username, Result: FAIL, Message: messages["EmptyUsernameOrPassword"]}
	} else if itemIndex, err := strconv.Atoi(itemIndexString); err == nil {
		if _, ok := operateOnDB(deleteItemIntoDB, username, tools.MD5Encryption(password), itemIndex+1); ok {
			return TodoListOperationResult{Operation: DELETE, Username: username, Result: SUCCESS, Message: messages["DeleteSuccess"]}
		}
	}
	return TodoListOperationResult{Operation: DELETE, Username: username, Result: FAIL, Message: messages["DeleteFail"]}
}

// check if username or password is empty
// show items from datebase and check if succeeded
// return todoListOperationResult, result and details
func showItemsByUsernameAndPassword(username string, password string) TodoListOperationResult {
	if isEmptyUsernameOrPassword(username, password) { // if one of important parameters is empty
		return TodoListOperationResult{Operation: SHOW, Username: username, Result: FAIL,
			Message: messages["EmptyUsernameOrPassword"]}
	} else if todoListOperationResult, ok := operateOnDB(showItemsFromDB, username, tools.MD5Encryption(password), nil); ok {
		return TodoListOperationResult{
			Operation: SHOW, Username: username, Todos: todoListOperationResult[1:], Result: SUCCESS, // remove the first item(empty) of todoListOperationResult
			Message: fmt.Sprintf(messages["ShowSuccess"], len(todoListOperationResult)-1)}
	}
	return TodoListOperationResult{Operation: SHOW, Username: username, Result: FAIL, Message: messages["ShowFail"]}
}

// auxiliary functions - MD5 hash function
func MD5Encryption(text string) string {
	hash := md5.New()
	io.WriteString(hash, text)
	return fmt.Sprintf("%x", hash.Sum(nil))
}