
package main

import (
    log "github.com/cihub/seelog"
    "fmt"
)

func loadAppConfig() {
    appConfig := `
<seelog type="sync">
    <outputs formatid="app">
        <console />
        <rollingfile type="date" filename="logs/roll.log" datepattern="02.01.2006" maxrolls="30" />        
        <filter levels="warn, error, critical">
            <file path="logs/error.log" formatid="critical"/>    
        </filter>
    </outputs>
    <formats>
        <format id="app" format="%Time %Date  [%LEV] %Msg%n" />
        <format id="critical" format="%Time %Date %RelFile %Func %Msg"/>
    </formats>
</seelog>
`
    logger, err := log.LoggerFromConfigAsBytes([]byte(appConfig))
    if err != nil {
        fmt.Println(err)
        return
    }
    log.ReplaceLogger(logger)
}


func initSeeLog() {
    defer log.Flush()
    loadAppConfig() 
    log.Info("App started")
    log.Info("seelog Config loaded")   
}

