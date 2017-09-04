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


var templates = template.Must(template.ParseFiles("template/edit.html", "template/list.html", "template/write.html"))


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
    //fmt.Println("show write html \n")
    renderTemplate(w, "write.html", &Item{})
}

func viewItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    //fmt.Println(ps[0])
    //fmt.Fprintf(w, "get one item, key is %s, value is %s\n", ps[0].Key, ps[0].Value)
    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        fmt.Fprintf(w, "can not find item for id %s\n", ps[0].Value)
        return
    }

    item := dbGetItem(id)
    renderTemplate(w, "edit.html", item)
}

func addItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    //fmt.Println("add one item \n")
    r.ParseForm()
    body := r.FormValue("body")
    title := r.FormValue("title")
    now := time.Now().Format("2006-01-02 15:04:05")
    due := time.Date(2017, time.November, 10, 23, 0, 0, 0, time.UTC)

    item := &Item{Title:title, Body:[]byte(body), Created:now, DueTime:due, Done:false}
    if item.Validate() == false {
        renderTemplate(w, "write.html", item)
        return
    }

    dbInsertItem(item)
    http.Redirect(w, r, "/list", http.StatusFound) 
}
  


func editItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    //fmt.Fprintf(w, "get one item, key is %s, value is %s\n", ps[0].Key, ps[0].Value)
    id, err := strconv.ParseUint(ps[0].Value, 10, 64)
    if err != nil {
        fmt.Fprintf(w, "can not find item for id %s\n", ps[0].Value)
        return
    }

    r.ParseForm()
    cmd := r.FormValue("remove")   // if cmd is remove
    body := r.FormValue("body")
    title := r.FormValue("title")
    now := time.Now().Format("2006-01-02 15:04:05")
    
    if cmd == "Remove" {
        item := &Item{Id:id}
        dbDeleteItem(item)
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

func removeItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


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
    router.GET("/write", writeItem)  //write.html
    router.POST("/write", addItem)     // add one item to list
    router.GET("/item/:id", viewItem)     // get one item detail ,edit.html
    router.POST("/item/:id", editItem)   // edit one item
    router.GET("/people/:name", getUserInfo) // get user info

    log.Fatal(http.ListenAndServe(":8082", router))
} 