// user manage code
//

package main

import (
    "time"
    "strings"
    "golang.org/x/crypto/bcrypt"
    log "github.com/cihub/seelog"        
)

type UserMgr struct {

}


func NewUserMgr() (*UserMgr, error) {

    return &UserMgr{}, nil
}



func (manager *UserMgr) Process(name, password string) (*User, bool) {

    user := &User{Name:name, Password:password}

    flag := manager.Validate(user)
    //do not send back password plaintext
    user.Password = ""

    if flag == false {
        return user, false
    }

    tmp := dbGetUser(name)
    if tmp == nil {
        log.Debug("user not exist, create one")
        tmp, ok := manager.CreateUser(name, password)
        if ok {
            log.Infof("Created user for name %s success", name)            
            user = tmp
        }
    } else {
        // compare password
        log.Debug("user exist, verify password")        
        pass := []byte(password)
        hash := []byte(tmp.Password)
        if bcrypt.CompareHashAndPassword(hash, pass) != nil {
            log.Infof("%s not match %s!", hash, pass)
            user.Errors["Password"] = "password not match!"
            return user, false
        }        
        user = tmp
        log.Debug("pass verify, get user data")        
    }
    return user, true
}    

func (manager *UserMgr) CreateUser(name, password string) (*User, bool) {

    pass := []byte(password)
    hp, err := bcrypt.GenerateFromPassword(pass, 10)
    if err != nil {
        log.Warnf("GenerateFromPassword error: %s", err)
        return nil, false
    }

    if bcrypt.CompareHashAndPassword(hp, pass) != nil {
        log.Warnf("%s should hash %s correctly", hp, pass)
        return nil, false
    }
    passwd := string(hp)
    log.Debugf("salted password is %s", passwd)
    
    user := &User{Name:name, Password:passwd, Created:time.Now(), Active:true}
    ok := dbInsertUser(user)
    
    log.Debug("create user")
    return user, ok   
}    



// validate if user name and password fit our rule, like not empyt, only a-zA-Z0-9
func (manager *UserMgr) Validate(user *User) bool {

    name := strings.TrimSpace(user.Name)
    password := strings.TrimSpace(user.Password)

    user.Errors = make(map[string]string)

    if name == "" {
        user.Errors["Name"] = "usr name can not be empty"
    }
    if password == "" {
        user.Errors["Password"] = "password can not be empty"
    }

    return len(user.Errors) == 0
}



