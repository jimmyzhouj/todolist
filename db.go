// go wiki study program
//

package main

import (
    "fmt"
    "database/sql" 
    _ "github.com/mattn/go-sqlite3"
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
    stmt, err := db.Prepare("INSERT INTO items(title, content, created) values(?, ?, ?)")
    checkErr(err)

    //_, err = stmt.Exec("aa", "aa", "aa")
    _, err = stmt.Exec(item.Title, item.Body, item.Date)
    checkErr(err)
}

func dbGetAllItems() []*Item {

    rows, err := db.Query("SELECT * FROM items")
    checkErr(err)

    var list []*Item
    for rows.Next() {
        var uid uint64
        var title string
        var content []byte
        var created string
        err = rows.Scan(&uid, &title, &content, &created)
        checkErr(err)
        fmt.Println("list item")
        fmt.Println(uid)
        fmt.Println(title)
        fmt.Println(content)
        fmt.Println(created)
        item := &Item{Id:uid, Title:title, Body:content, Date:created}
        list = append(list, item)
    }
    return list
}

// query item
func dbGetItem() *Item {  
    //update
    // rows, err := db.Query("SELECT * FROM ")
    // checkErr(err)

    item := &Item{}
    return item
}

// update item
func dbUpdateItem(item *Item) {

    if item == nil {
        panic("item is nil")
        return 
    }

    // stmt, err := db.Prepare("update items set username=? where uid=?")
    // checkErr(err)

    // res, err := stmt.Exec("zjupdate", id)
    // checkErr(err)    
}


func dbDeleteItem(item *Item) {

    if item == nil {
        panic("item is nil")
        return 
    }

    // stmt, err = db.Prepare("delete from items where uid=?")
    // checkErr(err)

    // res, err = stmt.Exec(id)
    // checkErr(err)

    // affect, err = res.RowsAffected()
    // checkErr(err)
    // fmt.Println("delete affected rows ", affect)
}


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}



func init() {
    fmt.Println("init data base")
    db, err = sql.Open("sqlite3", "db/data.db")
    checkErr(err)  
}

