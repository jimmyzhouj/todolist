// api handlers

package main 

import (
    "net/http"
    "time"
    log "github.com/cihub/seelog"
    "github.com/julienschmidt/httprouter"
    "io"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "fmt"                    
)

type ErrorResponse struct {
    Error string `json:"error"`
}

type SessionTokenResponse struct {
    Token string `json:"sessionToken"`
}

func writeJSON(w http.ResponseWriter, status int, response interface{}) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(status)
    err := json.NewEncoder(w).Encode(response)
    if err != nil {
        panic(err)
    }     
}


func apiGetUser(r *http.Request) (*User, error) {

    sess := globalSessions.ApiSessionStart(r)
    if sess == nil {
        log.Info("no valid session, need login first")
        return nil, fmt.Errorf("no valid session")        
    }

    tmp := sess.Get("username")
    if tmp == nil || tmp.(string) == "" {
        log.Info("no user bind to session, need login first")
        return nil, fmt.Errorf("no user bind to session")
    }

    name := tmp.(string)
    user := dbGetUser(name)
    if user != nil {
        return user, nil
    }
    return nil, fmt.Errorf("can not get user id from db") 
}


func apiRegister(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ POST url: /api/v1/user")  
    //log.Debug(r.Header)
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    var resp RespBody
    resp.Meta = make(map[string]string)
    var status int
    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        resp.Meta["message"] = err.Error()
        status = http.StatusUnprocessableEntity
        writeJSON(w, status, resp)        
        return
    }    

    role, err := userMgr.CreateUser(user.Name, user.Password)
    if err != nil {
        resp.Meta["message"] = err.Error()
        status = http.StatusBadRequest                    
    } else {
        // register success, create session
        sess := globalSessions.ApiSessionStart(r)
        if sess != nil {
            sess.Set("username", role.Name)
            log.Debugf("bind username %s to curr session \n", role.Name)
            resp.Data = SessionTokenResponse{Token: sess.SessionID()}
            status = http.StatusOK          
        } else {
            resp.Meta["message"] = "no valid session token"
            status = http.StatusInternalServerError 
        }
              
    }
    writeJSON(w, status, resp)
  
}


func apiLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ apiLogin Handler, POST url: /api/v1/user/login")  
    //log.Debug(r.Header)
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    var resp RespBody
    resp.Meta = make(map[string]string)
    var status int
    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        resp.Meta["message"] = err.Error()
        status = http.StatusUnprocessableEntity
        writeJSON(w, status, resp)        
        return
    }    


    role, err := userMgr.AuthUser(user.Name, user.Password)
    if err != nil {
        resp.Meta["message"] = err.Error()
        status = http.StatusBadRequest                    
    } else {
        // register success, create session
        sess := globalSessions.ApiSessionStart(r)
        if sess != nil {
            sess.Set("username", role.Name)
            log.Debugf("bind username %s to curr session \n", role.Name)
            resp.Data = SessionTokenResponse{Token: sess.SessionID()}
            status = http.StatusOK          
        } else {
            resp.Meta["message"] = "no valid session token"
            status = http.StatusInternalServerError 
        }
              
    }
    writeJSON(w, status, resp)
  
}


func apiLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ apiLogout Handler, POST url: /api/v1/user/logout")  

    var resp RespBody
    resp.Meta = make(map[string]string)

    user, err := apiGetUser(r)
    if err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnauthorized, resp)        
        return
    }
    log.Debugf("logout user : %s", user.Name)
    sess := globalSessions.ApiSessionStart(r)
    tmp := sess.Get("username")
    if tmp != nil && tmp.(string) != "" {
        log.Debugf("remove user name  %s from session\n", tmp.(string))
        sess.Delete("username")
    }   
    globalSessions.ApiSessionEnd(sess)

    //resp.URI = fmt.Sprintf("%s", r.URL)

    w.WriteHeader(http.StatusNoContent)

    //writeJSON(w, http.StatusNoContent, resp)  
}


// list all items
func apiItemList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ GET url: /api/v1/todos")  

    var resp RespBody
    resp.Meta = make(map[string]string)    
    var status int

    user, err := apiGetUser(r)
    if err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnauthorized, resp)        
        return
    }

    list := dbGetItemsByUserId(user.Id)
    //list := dbGetAllItems()
    resp.Data = list
    status = http.StatusOK
    writeJSON(w, status, resp)

} 

func apiItemShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ GET url: ", r.URL) 
    log.Debug(ps[0])

    var resp RespBody
    resp.Meta = make(map[string]string)
    var status int

    user, err := apiGetUser(r)
    if err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnauthorized, resp)        
        return
    }
    log.Debug(user.Id)

    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        resp.Meta["message"] = "the id part of url is not valid digit"
        status = http.StatusNotFound
    } else {
        item := dbGetItem(id)
        if item != nil {
            if item.UserId != user.Id {
                status = http.StatusUnauthorized
                log.Errorf("item user id %d not match cur user. cur user name %s, user id %d \n", item.UserId, user.Name, user.Id)                
                resp.Meta["message"] = "this item does not belongs to you"              
            } else {
                resp.Data = item
                status = http.StatusOK                
            }

        } else {
            status = http.StatusNotFound
            resp.Meta["message"] = "can not get todo for id " + ps[0].Value
        }     
    }

    writeJSON(w, status, resp)
}


// add 
func apiItemCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ POST url: /api/v1/todos")  

    var resp RespBody
    resp.Meta = make(map[string]string)    

    user, err := apiGetUser(r)
    if err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnauthorized, resp)        
        return
    }

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    var item Item
    if err := json.Unmarshal(body, &item); err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnprocessableEntity, resp)
        return
    }

    if item.Validate() == false {
        for name, value := range item.Errors {
            resp.Meta["message"] += name + " : " + value + " || " 
        }
        writeJSON(w, http.StatusNotAcceptable, resp)    
        return
    }

    now := time.Now()
    due := time.Date(2017, time.November, 10, 23, 0, 0, 0, time.UTC)
    item.UserId = user.Id
    item.Created = now
    item.DueTime = due
    item.Done = false
    id := dbInsertItem(&item)
    item.Id = id
    log.Debug("add item id is ", id)
    if id > 0 {
        log.Trace("add item success ")
        log.Trace(item)
        resp.Data = item
        resp.URI = "/api/v1/todo/" + strconv.FormatUint(id, 10)

        writeJSON(w, http.StatusCreated, resp)       
    }
}


func apiItemUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ PATCH url: ", r.URL) 

    var resp RespBody
    resp.Meta = make(map[string]string)

    user, err := apiGetUser(r)
    if err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnauthorized, resp)        
        return
    }

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    var item Item
    if err := json.Unmarshal(body, &item); err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnprocessableEntity, resp)        
        return        
    }
    log.Trace("received item ")
    log.Trace(item)

    if item.Validate() == false {
        for name, value := range item.Errors {
            resp.Meta["message"] += name + " : " + value + " || " 
        }
        writeJSON(w, http.StatusNotAcceptable, resp)            
        return
    }

    var status int
    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        msg := "the id part of url is not valid digit"
        log.Info(msg)
        resp.Meta["message"] = msg         
    } else {
        curr := dbGetItem(id)

        if curr != nil {
            if curr.UserId != user.Id {
                log.Errorf("item user id %d not match cur user %s, user id %d \n", item.UserId, user.Name, user.Id)
                status = http.StatusUnauthorized
                resp.Meta["message"] = "Not your todo, can not modify it"

            } else {
                curr.Title = item.Title
                curr.Body = item.Body
                curr.Done = item.Done
                log.Debugf("update item for id %d, title is %s", curr.Id, curr.Title)
                dbUpdateItem(curr)             
                log.Trace("item after update")
                log.Trace(curr)

                resp.Data = curr
                resp.URI = fmt.Sprintf("%s", r.URL)
                status = http.StatusOK
            }

        } else {
            status = http.StatusNotFound
            msg := "todo not exist for id " + ps[0].Value
            log.Info(msg)
            resp.Meta["message"] = msg               
        }     
    }

    writeJSON(w, status, resp)
}


func apiItemDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ DELETE url: ", r.URL) 


    var resp RespBody
    resp.Meta = make(map[string]string)
    var status int
    
    user, err := apiGetUser(r)
    if err != nil {
        resp.Meta["message"] = err.Error()
        writeJSON(w, http.StatusUnauthorized, resp)        
        return
    }
    log.Debug(user.Id)

    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        status = http.StatusNotFound
        msg := "the id part of url is not valid digit"
        log.Info(msg)
        resp.Meta["message"] = msg        
    } else {

        curr := dbGetItem(id)
        if curr != nil && (curr.UserId != user.Id) {
            log.Errorf("item user id %d not match cur user. cur user name %s, user id %d \n", curr.UserId, user.Name, user.Id)
            status = http.StatusUnauthorized
            resp.Meta["message"] = "Not your todo, can not delete"            
        } else {
            log.Infof("remove item for id %v", id)
            num := dbDeleteItem(id)
            if num > 0 {
                resp.URI = fmt.Sprintf("%s", r.URL)            
                status = http.StatusOK
                //resp.Meta["message"] = "delete todo success, id " + ps[0].Value
            } else {
                status = http.StatusNotFound
                msg := "todo not exist,  id " + ps[0].Value
                log.Info(msg)
                resp.Meta["message"] = msg
            }        
        }
    }

    writeJSON(w, status, resp)
}