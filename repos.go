package main

import (
	"encoding/json"
	"io/ioutil"
)

// repository interface for User info
type userRepo interface {
	init()
	close()
	validateUser(username, passowrd string) bool
}

// struct to hold user info
type user struct {
	Name     string
	Password string
}

// struct user repository from file
type userRepoFile struct {
	// TO DO => Replace slice with map
	users map[string]user
}

func (ur *userRepoFile) init() {
	bytes, err := ioutil.ReadFile(config().UsersDb)
	checkerr(err)
	var us []user
	err = json.Unmarshal(bytes, &us)
	ur.users = make(map[string]user)
	checkerr(err)
	for _, u := range us {
		ur.users[u.Name] = u
	}
	//log.Printf("init @ users = %v\n", ur.users)
}

func (ur *userRepoFile) close() {
	// nothing to do fo this implementation
}

func (ur *userRepoFile) validateUser(username, passowrd string) bool {
	isValid := false
	for _, u := range ur.users {
		//fmt.Printf("user: %s, pwd: %s \n", u.Name, u.Password)
		if username == u.Name {
			isValid = passowrd == u.Password
			break
		}
	}
	return isValid
}
