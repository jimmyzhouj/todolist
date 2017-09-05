// a cloud to do list app 

package main 

import (
    "fmt"
    "net/http"
    "html/template"
    "strconv"
    "time"
    "log"
    "github.com/julienschmidt/httprouter"
    "github.com/jimmyzhouj/session"
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
        fmt.Println("no user bind to session, need login first")
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
    fmt.Println("get / handler run")
    fmt.Fprint(w, "default, Welcome!\n")
} 

func listAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("get /list handler run")

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
    fmt.Println("get /write handler run")      

    _, err := getUser(w, r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }
    renderTemplate(w, "write.html", &Item{})
}

func addItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("post /write handler run")  

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

    item := &Item{UserId: user.Id, Title:title, Body:[]byte(body), Created:now, DueTime:due, Done:false}
    if item.Validate() == false {
        renderTemplate(w, "write.html", item)
        return
    }

    dbInsertItem(item)
    http.Redirect(w, r, "/list", http.StatusFound) 
}
  
func viewItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("get /item/xxx handler run")      
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
        fmt.Errorf("item user id %d not match cur user %s, user id %d \n", item.UserId, user.Name, user.Id)
        http.Redirect(w, r, "/login", http.StatusFound) 
    }

}

func editItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("post /item/xxx handler run")  

    user, err := getUser(w, r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    //fmt.Fprintf(w, "get one item, key is %s, value is %s\n", ps[0].Key, ps[0].Value) 
    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        fmt.Fprintf(w, "can not find item for id %s\n", ps[0].Value)
        return
    }

    item := dbGetItem(id)
    if item.UserId != user.Id {
        fmt.Errorf("item user id %d not match cur user %s, user id %d \n", item.UserId, user.Name, user.Id)
        http.Redirect(w, r, "/login", http.StatusFound)
        return 
    }

    r.ParseForm()
    cmd := r.FormValue("remove")   // if cmd is remove
    body := r.FormValue("body")
    title := r.FormValue("title")
    now := time.Now()
    
    if cmd == "Remove" {
        dbDeleteItem(id)
        http.Redirect(w, r, "/list", http.StatusFound) 
    } else {

        item := &Item{Title:title, Body:[]byte(body), Created:now}
        if item.Validate() == false {
            renderTemplate(w, "edit.html", item)
            return
        }

        dbUpdateItem(item) 
        http.Redirect(w, r, "/list", http.StatusFound)      
    }
}


func showLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("get /login handler run")  

    user, err := getUser(w, r)
    // not login, need login
    if err != nil {
        fmt.Println("no username bind to session, need login")
        err := templates.ExecuteTemplate(w, "login.html", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } 
    } else {
        
        fmt.Printf("get user name %s from session\n", user.Name)
        fmt.Println("redirect to list")       
        http.Redirect(w, r, "/list", http.StatusFound) 
    }


    // sess := globalSessions.SessionStart(w, r)    
    // tmp := sess.Get("username")
    // if tmp != nil && tmp.(string) != "" {
    //     fmt.Println("get user name from session " , tmp)
    //     fmt.Println("redirect to list")       
    //     http.Redirect(w, r, "/list", http.StatusFound)         
    // } else {
    //     fmt.Println("no username bind to session, need login")
    //     err := templates.ExecuteTemplate(w, "login.html", nil)
    //     if err != nil {
    //         http.Error(w, err.Error(), http.StatusInternalServerError)
    //     } 
    // } 
}  


func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("post /login handler run")      
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
    fmt.Printf("bind username %s to curr session \n", name)

    http.Redirect(w, r, "/list", http.StatusFound) 
}    

func showLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("get /logout handler run")  
    err := templates.ExecuteTemplate(w, "logout.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }  
} 


func logoutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println("post /logout handler run")  
    sess := globalSessions.SessionStart(w, r)

    tmp := sess.Get("username")
    if tmp != nil && tmp.(string) != "" {
        fmt.Printf("remove user name  %s from session\n", tmp.(string))
        sess.Delete("username")
    }   

    globalSessions.SessionEnd(w, sess)

    http.Redirect(w, r, "/login", http.StatusFound) 
} 




//session manager
var globalSessions *session.Manager

var userMgr *UserMgr

// init all resoureses
func init() {
    globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
    //go globalSession.GC()
    userMgr, _ = NewUserMgr()
}


func main() {

    router := httprouter.New()

    router.GET("/", defaultHandler)
    router.GET("/list", listAll)     // list all items           list.html
    router.GET("/write", writeItem)  //write.html
    router.POST("/write", addItem)     // add one item to list
    router.GET("/item/:id", viewItem)     // get one item detail ,edit.html
    router.POST("/item/:id", editItem)   // edit one item
    router.GET("/login", showLogin)      // show login ui
    router.POST("/login", loginHandler)  // user login
    router.GET("/logout", showLogout)    
    router.POST("/logout", logoutHandler)

    log.Fatal(http.ListenAndServe(":8082", router))
} 