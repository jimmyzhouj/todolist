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

// list all items
func apiItemList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ GET url: /api/v1/todos")  

    list := dbGetAllItems()

    var resp RespBody
    resp.Data = list

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        panic(err)
    } 
} 

func apiItemShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ GET url: ", r.URL) 
    log.Debug(ps[0])

    var resp RespBody
    resp.Meta = make(map[string]string)

    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        resp.Meta["message"] = "the id part of url is not valid digit"
    } else {
        item := dbGetItem(id)
        if item != nil {
            resp.Data = item
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusNotFound)
            resp.Meta["message"] = "can not get todo for id " + ps[0].Value
        }     
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        panic(err)
    } 
}


// add 
func apiItemCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ POST url: /api/v1/todos")  

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    var item Item
    if err := json.Unmarshal(body, &item); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusUnprocessableEntity)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
        return
    }

    var resp RespBody
    if item.Validate() == false {
        resp.Meta = make(map[string]string)

        for name, value := range item.Errors {
            resp.Meta["message"] += name + " : " + value + " || " 
        }

        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusNotAcceptable)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            panic(err)
        }        
        return
    }

    now := time.Now()
    due := time.Date(2017, time.November, 10, 23, 0, 0, 0, time.UTC)
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
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusCreated)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            panic(err)
        }        
    }
}


func apiItemUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ PATCH url: ", r.URL) 
    log.Debug(ps[0])

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    var item Item
    if err := json.Unmarshal(body, &item); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusUnprocessableEntity)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
        return
    }
    log.Trace("received item ")
    log.Trace(item)

    var resp RespBody
    resp.Meta = make(map[string]string)

    if item.Validate() == false {

        for name, value := range item.Errors {
            resp.Meta["message"] += name + " : " + value + " || " 
        }

        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusNotAcceptable)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            panic(err)
        }        
        return
    }


    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        msg := "the id part of url is not valid digit"
        log.Info(msg)
        resp.Meta["message"] = msg         
    } else {
        curr := dbGetItem(id)
        if curr != nil {
            curr.Title = item.Title
            curr.Body = item.Body
            curr.Done = item.Done
            log.Debugf("update item for id %d, title is %s", curr.Id, curr.Title)
            dbUpdateItem(curr)             
            log.Trace("item after update")
            log.Trace(curr)

            resp.Data = curr
            resp.URI = fmt.Sprintf("%s", r.URL)
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusNotFound)
            msg := "todo not exist for id " + ps[0].Value
            log.Info(msg)
            resp.Meta["message"] = msg               
        }     
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        panic(err)
    } 
}


func apiItemDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("############ DELETE url: ", r.URL) 


    var resp RespBody
    resp.Meta = make(map[string]string)

    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        msg := "the id part of url is not valid digit"
        log.Info(msg)
        resp.Meta["message"] = msg        
    } else {
        log.Infof("remove item for id %v", id)
        num := dbDeleteItem(id)
        if num > 0 {
            resp.URI = fmt.Sprintf("%s", r.URL)            
            w.WriteHeader(http.StatusOK)
            //resp.Meta["message"] = "delete todo success, id " + ps[0].Value
        } else {
            w.WriteHeader(http.StatusNotFound)
            msg := "todo not exist,  id " + ps[0].Value
            log.Info(msg)
            resp.Meta["message"] = msg
        }     
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        panic(err)
    } 
}