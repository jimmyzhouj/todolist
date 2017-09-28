// user manage code
//

package main

import (
    "time"
    "fmt"
    "strings"
    "errors"    
    "golang.org/x/crypto/bcrypt"
    log "github.com/cihub/seelog"        
)

type UserMgr struct {

}


func NewUserMgr() (*UserMgr, error) {
    return &UserMgr{}, nil
}



func (manager *UserMgr) AuthUser(name, password string) (*User, error) {

    var msg string

    err := manager.Validate(name, password)
    if err != nil {
        return nil, err
    }

    var user *User
    for _, s := range userCache {
        if s.Name == name {
            user = s
        }
    }

    if user == nil {
        msg = fmt.Sprintf("user name %s not exit exist, auth fail", name)
        log.Warn(msg)
        return nil, errors.New(msg)

    } else {
        // compare password
        log.Debug("user exist, verify password")        
        pass := []byte(password)
        hash := []byte(user.Password)
        if bcrypt.CompareHashAndPassword(hash, pass) != nil {
            log.Infof("%s not match %s!", hash, pass)
            msg = fmt.Sprintf("password not match user name")
            log.Warn(msg)
            return nil, errors.New(msg)            
        }        

        log.Debug("pass verify, get user data")        
    }
    return user, nil
}    

func (manager *UserMgr) CreateUser(name, password string) (*User, error) {
    var msg string

    err := manager.Validate(name, password)
    if err != nil {
        return nil, err
    }

    for _, s := range userCache {
        if s.Name == name {
            msg = fmt.Sprintf("user name %s already exist, register fail", name)
            log.Warn(msg)
            return nil, errors.New(msg)
        }
    }

    pass := []byte(password)
    hp, err := bcrypt.GenerateFromPassword(pass, 10)
    if err != nil {
        msg = fmt.Sprintf("GenerateFromPassword error: %s", err)
        log.Warn(msg)
        return nil, errors.New(msg)
    }

    if bcrypt.CompareHashAndPassword(hp, pass) != nil {
        msg = fmt.Sprintf("%s should hash %s correctly", hp, pass)        
        log.Warn(msg)
        return nil, errors.New(msg)
    }
    passwd := string(hp)
    log.Debugf("salted password is %s", passwd)
    
    user := &User{Name:name, Password:passwd, Created:time.Now(), Active:true}
    ok := dbInsertUser(user)
    if ok {
        // append new user
        userCache = append(userCache, user)
    }

    log.Debug("create user")
    return user, nil   
}    



// validate if user name and password fit our rule, like not empyt, only a-zA-Z0-9
func (manager *UserMgr) Validate(name string, password string) error {

    name = strings.TrimSpace(name)
    password = strings.TrimSpace(password)

    if name == "" {
        return errors.New("user name can not be empty")
    }
    if password == "" {
        return errors.New("password can not be empty")
    }

    return nil
}



