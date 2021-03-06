// user struct
//

package main

import (
    //"fmt"
    "time"
    "strings"
)


type Users []*User


// user definition
type User struct {
    Id uint64
    Name string     `json:"name"`
    Password string `json:"password"`
    Created time.Time
    Active bool
    SessionId string `json:sessionid`
    Errors map[string]string  // show field errors
}

// validate the received user
func (user *User) Validate() bool {
    user.Errors = make(map[string]string)

    if strings.TrimSpace(user.Name) == "" {
        user.Errors["Name"] = "usr name can not be empty"
    }

    if strings.TrimSpace(user.Password) == "" {
        user.Errors["Password"] = "password can not be empty"
    }

    return len(user.Errors) == 0
}

