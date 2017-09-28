// a cloud to do list app 

package main 

import (
    "net/http"
    log "github.com/cihub/seelog"
    "github.com/jimmyzhouj/session"
    _ "github.com/jimmyzhouj/session/providers/memory"

)



//session manager
var globalSessions *session.Manager

var userMgr *UserMgr
var userCache []*User

// init all resoureses
func init() {
    globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
    //go globalSession.GC()
    userMgr, _ = NewUserMgr()
}


func main() {

    initSeeLog()
    initDB()

    userCache = dbGetAllUsers()

    router := NewRouter()
    //http.ListenAndServe(":8082", router)
    http.ListenAndServeTLS(":8082", "cert/server.crt", "cert/server.key", router)
    log.Critical ("listen and serve panic")
} 