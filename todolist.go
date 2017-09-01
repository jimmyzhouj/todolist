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
)


var templates = template.Must(template.ParseFiles("template/view.html", "template/edit.html", "template/list.html", "template/write.html"))

// item definition
type Item struct {
    Id uint64
    Title string
    Body  []byte
    Date string
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Item) {
    err := templates.ExecuteTemplate(w, tmpl, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func defaultHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprint(w, "default, Welcome!\n")
} 

func listAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    list := dbGetAllItems()

    err := templates.ExecuteTemplate(w, "list.html", list)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }    
}

func writeItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    //fmt.Fprint(w, "add an item\n")
    renderTemplate(w, "write.html", &Item{})
}

func addItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    body := r.FormValue("body")
    title := r.FormValue("title")
    now := time.Now().Format("2006-01-02 15:04:05")
    fmt.Println("add item time is ", now)
    item := &Item{Title:title, Body:[]byte(body), Date:now}
    dbInsertItem(item)

    http.Redirect(w, r, "/list", http.StatusFound) 
    //fmt.Fprintf(w, "add one  item, title is %s, body is %s\n", title, body)
}
  
func viewItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println(ps[0])
    //fmt.Fprintf(w, "get one item, key is %s, value is %s\n", ps[0].Key, ps[0].Value)
    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        fmt.Fprintf(w, "can not find item for id %s\n", ps[0].Value)
        return
    }

    item := dbGetItem(id)
    renderTemplate(w, "edit.html", item)
}

func editItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    //fmt.Fprintf(w, "get one item, key is %s, value is %s\n", ps[0].Key, ps[0].Value)
    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        fmt.Fprintf(w, "can not find item for id %s\n", ps[0].Value)
        return
    }

    body := r.FormValue("body")
    title := r.FormValue("title")
    now := time.Now().Format("2006-01-02 15:04:05")
    fmt.Println("add item time is ", now)
    item := &Item{Id:id, Title:title, Body:[]byte(body), Date:now}
    dbUpdateItem(item)

    //fmt.Fprintf(w, "edit item success, key is %s, value is %s\n", ps[0].Key, ps[0].Value)
    //http.Redirect(w, r, "/item/" + ps[0].Value, http.StatusFound) 
    http.Redirect(w, r, "/list", http.StatusFound) 

}

func getUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println(ps)
    fmt.Fprint(w, "get user info\n")
}



/*
func init() {

}  
*/  


func main() {

    router := httprouter.New()

    router.GET("/", defaultHandler)
    router.GET("/list", listAll)     // list all items           list.html
    router.GET("/write", writeItem)  //edit.html
    router.POST("/write", addItem)     // add one item to list    write.html
    router.GET("/item/:id", viewItem)     // get one item detail view.html
    router.POST("/item/:id", editItem)   // edit one item  edit.html
    router.GET("/people/:name", getUserInfo) // get user info

    log.Fatal(http.ListenAndServe(":8082", router))
} 