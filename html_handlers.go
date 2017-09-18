// html handler

package main 

import (
    "fmt"
    "net/http"
    "html/template"
    "strconv"
    "time"
    "github.com/julienschmidt/httprouter"
    log "github.com/cihub/seelog"
    _ "github.com/jimmyzhouj/session/providers/memory"           
)


var templates = template.Must(template.ParseFiles("template/login.html", "template/logout.html", "template/edit.html", "template/list.html", 
                                                "template/write.html"))


func renderTemplate(w http.ResponseWriter, tmpl string, p *Item) {
    err := templates.ExecuteTemplate(w, tmpl, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func getUser(w http.ResponseWriter, r *http.Request) (*User, error) {
    sess := globalSessions.SessionStart(w, r)

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


func defaultHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("get / handler run start ")
    time.Sleep(7 * time.Second)
    log.Debug("get / handler run finished")    
    fmt.Fprint(w, "default, Welcome!\n")
} 

func listAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("get /list handler run")

    user, err := getUser(w, r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    list := dbGetItemsByUserId(user.Id)

    err = templates.ExecuteTemplate(w, "list.html", list)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }    
}

func writeItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("get /write handler run")      

    _, err := getUser(w, r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }
    renderTemplate(w, "write.html", &Item{})
}

func addItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("post /write handler run")  

    user, err := getUser(w, r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    r.ParseForm()
    body := r.FormValue("body")
    title := r.FormValue("title")
    //now := time.Now().Format("2006-01-02 15:04:05")
    now := time.Now()
    due := time.Date(2017, time.November, 10, 23, 0, 0, 0, time.UTC)

    item := &Item{UserId: user.Id, Title:title, Body:body, Created:now, DueTime:due, Done:false}
    if item.Validate() == false {
        renderTemplate(w, "write.html", item)
        return
    }

    dbInsertItem(item)
    http.Redirect(w, r, "/list", http.StatusFound) 
}
  
func viewItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("get /item/xxx handler run")      
    //fmt.Fprintf(w, "get one item, key is %s, value is %s\n", ps[0].Key, ps[0].Value)
    user, err := getUser(w, r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        fmt.Fprintf(w, "can not find item for id %s\n", ps[0].Value)
        return
    }

    item := dbGetItem(id)
    if item.UserId == user.Id {
        renderTemplate(w, "edit.html", item)        
    } else {
        log.Errorf("item user id %d not match cur user %s, user id %d \n", item.UserId, user.Name, user.Id)
        http.Redirect(w, r, "/login", http.StatusFound) 
    }

}

func editItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("post /item/xxx handler run")  

    user, err := getUser(w, r)
    if err != nil {
        log.Error("can not get user ")
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    //fmt.Fprintf(w, "get one item, key is %s, value is %s\n", ps[0].Key, ps[0].Value) 
    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        fmt.Fprintf(w, "can not find item for id %s\n", ps[0].Value)
        log.Errorf("can not find item for id %s\n", ps[0].Value)
        return
    }

    item := dbGetItem(id)
    if item.UserId != user.Id {
        log.Errorf("item user id %d not match cur user %s, user id %d \n", item.UserId, user.Name, user.Id)
        http.Redirect(w, r, "/login", http.StatusFound)
        return 
    }

    r.ParseForm()
    cmd := r.FormValue("remove")   // if cmd is remove
    body := r.FormValue("body")
    title := r.FormValue("title")
    now := time.Now()
    
    if cmd == "Remove" {
        log.Infof("remove item for id %s", id)
        dbDeleteItem(id)
        http.Redirect(w, r, "/list", http.StatusFound) 
    } else {

        item.Title = title
        item.Body = body
        item.Created = now

        if item.Validate() == false {
            log.Warn("validate item failed")
            renderTemplate(w, "edit.html", item)
            return
        }
        log.Debugf("update item for id %d, title is %s", item.Id, title)
        dbUpdateItem(item) 
        http.Redirect(w, r, "/list", http.StatusFound)      
    }
}


func showLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("get /login handler run")  

    user, err := getUser(w, r)
    // not login, need login
    if err != nil {
        log.Debug("no username bind to session, need login")
        err := templates.ExecuteTemplate(w, "login.html", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } 
    } else {
        
        log.Debugf("get user name %s from session\n", user.Name)
        log.Debug("redirect to list")       
        http.Redirect(w, r, "/list", http.StatusFound) 
    }

}  


func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("post /login handler run") 

    r.ParseForm()
    name := r.FormValue("name")
    password := r.FormValue("password")

    user, ok := userMgr.Process(name, password)
    if ok == false {
        err := templates.ExecuteTemplate(w, "login.html", user)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    // login in success, get session 
    //fmt.Fprint(w, "get user success!\n")
    sess := globalSessions.SessionStart(w, r)
    // connect session and user
    sess.Set("username", name)
    log.Debugf("bind username %s to curr session \n", name)

    http.Redirect(w, r, "/list", http.StatusFound) 
}    


func showLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("get /logout handler run")  
    err := templates.ExecuteTemplate(w, "logout.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }  
} 

func logoutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    log.Debug("post /logout handler run")  
    sess := globalSessions.SessionStart(w, r)

    tmp := sess.Get("username")
    if tmp != nil && tmp.(string) != "" {
        log.Debugf("remove user name  %s from session\n", tmp.(string))
        sess.Delete("username")
    }   

    globalSessions.SessionEnd(w, sess)

    http.Redirect(w, r, "/login", http.StatusFound) 
} 