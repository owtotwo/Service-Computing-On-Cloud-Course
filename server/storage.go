package server

import (
	"errors"
	"strings"
	"fmt"
	"time"
	"log"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)


// to create Bucket when program starts run
func init() {
	// Open the todolist.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("todolist.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Todos"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

/* Do some operations on database directly */

// open database and make sure to enable to be pinged
// to encrypt password by MD5
// do some operation on database and close connect after that
// return the result of operation (true/false)
func operateOnDB(function interface{}, username string, password string, parameter interface{}) ([]string, bool) {

	db, err := bolt.Open("todolist.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// judge function type
	// if is func addUserIntoDB(db *bolt.DB, username string, password string) error
	if addUserIntoDBFunc, ok := function.(func(*bolt.DB, string, string) error); ok {
		if addUserIntoDBFunc(db, username, password) != nil {
			return nil, false
		}
		return nil, true

	// if is addItemIntoDB(db *bolt.DB, username string, password string, item string) error
	} else if addItemIntoDBFunc, ok := function.(func(*bolt.DB, string, string, string) error); ok {
		var item string
		if item, ok = parameter.(string); ok { // get string parameter item
			if addItemIntoDBFunc(db, username, password, item) == nil {
				return nil, true
			}
		}
		return nil, false

	// if is deleteItemIntoDB(db *bolt.DB, username string, password string, itemIndex int) error
	} else if deleteItemIntoDBFunc, ok := function.(func(*bolt.DB, string, string, int) error); ok {
		var itemIndex int
		if itemIndex, ok = parameter.(int); ok {
			if deleteItemIntoDBFunc(db, username, password, itemIndex) == nil {
				return nil, true
			}
		}
		return nil, false

	// if is showItemsFromDB(db *bolt.DB, username string, password string) ([]string, error)
	} else if showItemsFromDBFunc, ok := function.(func(*bolt.DB, string, string) ([]string, error)); ok {
		todoList, err := showItemsFromDBFunc(db, username, password)
		if err != nil {
			return nil, false
		}
		return todoList, true

	// must not exist other kinds of function
	} else {
		panic(errors.New("wrong parameters"))
	}

}

// check if user exists in database by username and password, return bool result and user's id if exsits
func dbQueryUser(db *bolt.DB, username string, password string) (bool, string) {
	var id string

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		todosBucket := tx.Bucket([]byte("Todos"))

		return todosBucket.ForEach(func(k, _ []byte) error {
			v := todosBucket.Bucket(k)
			if string(v.Get([]byte("username"))) == username && 
			   string(v.Get([]byte("password"))) == password {
				id = string(k)
			}
			return nil
		})
	})

	if id != "" {
		return true, id
	}

	return false, ""
}

// check if user exists in database by username, return bool result
func dbQueryUserIfExist(db *bolt.DB, username string) bool {
	var result bool

	db.View(func(tx *bolt.Tx) error {
		todosBucket := tx.Bucket([]byte("Todos"))

		return todosBucket.ForEach(func(k, _ []byte) error {
			v := todosBucket.Bucket(k)
			if string(v.Get([]byte("username"))) == username {
				result = true
			}
			return nil
		})
	})

	return result
}

// get todolist string from database by username and password
func dbQueryTodos(db *bolt.DB, id string) string {
	var todos string // get todos

	db.View(func(tx *bolt.Tx) error {
		todos = string(tx.Bucket([]byte("Todos")).Bucket([]byte(id)).Get([]byte("todos")))
		return nil
	})

	return todos
}

// ---------------------------------------------------------------
// only do some main operations: add user, update or select todolist
// add user in database
func addUserIntoDB(db *bolt.DB, username string, password string) error {
	if dbQueryUserIfExist(db, username) == false {
		return db.Update(func(tx *bolt.Tx) error {

			userBucket, err := tx.Bucket([]byte("Todos")).CreateBucketIfNotExists([]byte(getUUID()))
			if err != nil {
				return err
			}
			
			if err := userBucket.Put([]byte("username"), []byte(username)); err != nil {
				return err
			}
			if err := userBucket.Put([]byte("password"), []byte(password)); err != nil {
				return err
			}
			if err := userBucket.Put([]byte("todos"), []byte("")); err != nil {
				return err
			}
			return nil
		})
	}

	return errors.New("This account is existed.")
}

func getItemsById(db *bolt.DB, id string) []string {
	return strings.Split(dbQueryTodos(db, id), ",")
}

// add item in database
func addItemIntoDB(db *bolt.DB, username string, password string, item string) error {
	// check if the user exists
	if result, id := dbQueryUser(db, username, password); result {
		// add item in todos
		var todosList []string = append(getItemsById(db, id), item)

		return db.Update(func(tx *bolt.Tx) error {
			userBucket := tx.Bucket([]byte("Todos")).Bucket([]byte(id))
			if err := userBucket.Put([]byte("todos"), []byte(strings.Join(todosList, ","))); err != nil {
				return err
			}
			return nil
		})
	}

	return errors.New("The account with this password was not found")
}

// parameter@todoListLength is len(strings.Split(dbQueryTodos(db, id), ","))-1
// to get actual todolistLength by "len-1" after spliting by "," because
// the first one string is empty after spliting by "," that should be ignored.
// check if item index is invalid; valid item index : from 0 to todolist actual length
func isInvalidItemIndex(itemIndex int, todoListLength int) bool {
	return itemIndex < 1 || itemIndex > todoListLength
}

// delete item by index in database
func deleteItemIntoDB(db *bolt.DB, username string, password string, itemIndex int) error {
	// check if the user exists
	if result, id := dbQueryUser(db, username, password); result {
		// get item in todos
		var todosList []string = getItemsById(db, id)
		// check if itemIndex is valid, 0 is invalid becase the 0th one is empty string after spliting by ","
		if isInvalidItemIndex(itemIndex, len(todosList)-1) {
			return errors.New("Could not find the item")
		}

		// delete an item by index
		todosList = append(todosList[:itemIndex], todosList[itemIndex+1:]...)

		return db.Update(func(tx *bolt.Tx) error {
			userBucket := tx.Bucket([]byte("Todos")).Bucket([]byte(id))
			if err := userBucket.Put([]byte("todos"), []byte(strings.Join(todosList, ","))); err != nil {
				return err
			}
			return nil
		})
	}

	return errors.New("The account with this password was not found")
}

// get all todo items from database
func showItemsFromDB(db *bolt.DB, username string, password string) ([]string, error) {
	// check if the user exists
	if result, id := dbQueryUser(db, username, password); result {
		// get item in todos
		return getItemsById(db, id), nil
	}

	return nil, errors.New("The account with this password was not found")
}

// auxiliary functions for getting uuid
func getUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return u.String() 
}