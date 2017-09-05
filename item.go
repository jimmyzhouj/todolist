// item
//

package main

import (
    //"fmt"
    "time"
    "strings"
)

// item definition
type Item struct {
    Id uint64
    UserId uint64
    Title string
    Body  []byte
    Created time.Time
    DueTime time.Time
    Done bool
    Errors map[string]string  // show field errors
}

// validate the received item
func (item *Item) Validate() bool {
    item.Errors = make(map[string]string)

    if strings.TrimSpace(item.Title) == "" {
        item.Errors["Title"] = "title can not be empty"
    }

    if len(item.Body) == 0 {
        item.Errors["Body"] = "item body can not be empty"
    }

    return len(item.Errors) == 0
}

