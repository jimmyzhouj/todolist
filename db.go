// go wiki study program
//

package main

import (
    "database/sql" 
    _ "github.com/mattn/go-sqlite3"
    "time"
    log "github.com/cihub/seelog"    
)

var db *sql.DB
var err error


func closeDb() {
    db.Close()
}

// insert new item
func dbInsertItem(item *Item) {

    if item == nil {
        panic("item is nil")
        return 
    }
    //insert data
    stmt, err := db.Prepare("INSERT INTO items(userid, title, content, done, created, duetime) values(?, ?, ?, ?, ?, ?)")
    defer stmt.Close()
    checkErr(err)

    //_, err = stmt.Exec("aa", "aa", "aa")
    _, err = stmt.Exec(item.UserId, item.Title, item.Body, item.Done, item.Created, item.DueTime)
    checkErr(err)
}

func dbGetAllItems() []*Item {

    rows, err := db.Query("SELECT * FROM items")
    defer rows.Close()
    checkErr(err)

    var list []*Item
    for rows.Next() {
        var uid uint64
        var userid uint64
        var title string
        var content []byte
        var created time.Time
        var duetime time.Time
        var done bool
        err = rows.Scan(&uid, &userid, &title, &content, &done, &created, &duetime)
        checkErr(err)
        item := &Item{Id:uid, UserId:userid, Title:title, Body:content, Created:created, DueTime:duetime, Done:done}
        list = append(list, item)
    }
    return list
}

func dbGetItemsByUserId(userid uint64) []*Item {

    stmt, err := db.Prepare("SELECT * FROM items where userid=?")
    defer stmt.Close()
    checkErr(err)
    rows, err := stmt.Query(userid)
    checkErr(err)      

    var list []*Item
    for rows.Next() {
        var uid uint64
        var userid uint64
        var title string
        var content []byte
        var created time.Time
        var duetime time.Time
        var done bool
        err = rows.Scan(&uid, &userid, &title, &content, &done, &created, &duetime)
        checkErr(err)
        item := &Item{Id:uid, UserId:userid, Title:title, Body:content, Created:created, DueTime:duetime, Done:done}
        list = append(list, item)
    }
    return list
}

// query item
func dbGetItem(uid uint64) *Item {  
    
    stmt, err := db.Prepare("SELECT * FROM items where uid=?")
    defer stmt.Close()
    checkErr(err)
    rows, err := stmt.Query(uid)
    checkErr(err)    

    var item *Item
    for rows.Next() {
        var uid uint64
        var userid uint64        
        var title string
        var content []byte
        var created time.Time
        var duetime time.Time
        var done bool        
        err = rows.Scan(&uid, &userid, &title, &content, &done, &created, &duetime)
        checkErr(err)
        item = &Item{Id:uid, UserId:userid, Title:title, Body:content, Created:created, DueTime:duetime, Done:done}
    }
    return item
}

// update item
func dbUpdateItem(item *Item) {

    if item == nil {
        panic("item is nil")
        return 
    }

    stmt, err := db.Prepare("update items set title=?, content=?, done=?, created=?, duetime=?  where uid=?")
    defer stmt.Close()
    checkErr(err)
    _, err = stmt.Exec(item.Title, item.Body, item.Done, item.Created, item.DueTime, item.Id)
    checkErr(err)    
}


func dbDeleteItem(id uint64) {

    stmt, err := db.Prepare("delete from items where uid=?")
    checkErr(err)

    _, err = stmt.Exec(id)
    checkErr(err)

    // affect, err = res.RowsAffected()
    // checkErr(err)
    // fmt.Println("delete affected rows ", affect)
}

// table users related

func dbInsertUser(user *User) bool {

    if user == nil {
        panic("user is nil")
        return false
    }
    //insert data
    stmt, err := db.Prepare("INSERT INTO users(name, password, created, active) values(?, ?, ?, ?)")
    defer stmt.Close()
    checkErr(err)

    _, err = stmt.Exec(user.Name, user.Password, user.Created, user.Active)
    checkErr(err)

    return true
}


func dbGetUser(name string) *User {  
    
    stmt, err := db.Prepare("SELECT * FROM users where name=?")
    defer stmt.Close()
    checkErr(err)
    rows, err := stmt.Query(name)
    checkErr(err)    

    var user *User
    for rows.Next() {
        var uid uint64
        var name string
        var password string
        var created time.Time
        var active bool        
        err = rows.Scan(&uid, &name, &password, &created, &active)
        checkErr(err)
        user = &User{Id:uid, Name:name, Password:password, Created:created, Active:active}
    }

    return user
}




func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}



func initDB() {
    log.Info("init data base")
    db, err = sql.Open("sqlite3", "db/data.db")
    checkErr(err)  
}

