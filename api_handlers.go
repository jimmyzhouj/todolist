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
)

// add 
func apiItemCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("post url: /api/v1/todos")  

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
    log.Debug("add item id is ", id)
    if id > 0 {
        msg := map[string]string {
            "uri": "/api/v1/todo/" + strconv.FormatUint(id, 10),
        }

        resp.Data = msg

        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusCreated)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            panic(err)
        }        
    }
}



// list all items
func apiItemList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("post url: /api/v1/todos")  

    list := dbGetAllItems()

    var resp RespBody
    resp.Data = list

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        panic(err)
    } 

}    