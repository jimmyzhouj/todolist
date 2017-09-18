// create my rooter

package main 

import (
    "github.com/julienschmidt/httprouter"
)



func NewRouter() *httprouter.Router {

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

    //router.POST("/api/login", apiLoginHandler)  // api user login
    router.POST("/api/v1/todos", apiItemCreate)     
    router.GET("/api/v1/todos", apiItemList)     

    return router
} 