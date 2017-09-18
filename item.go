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
    Id uint64               `json:"id"`
    UserId uint64           `json:"userid"`
    Title string            `json:"title"`
    Body  string            `json:"body"`
    Created time.Time       `json:"created"`
    DueTime time.Time       `json:"duetime"`
    Done bool               `json:"done"`
    Errors map[string]string `json:"-"` // show field errors
}

type Items []Item

// validate the received item
func (item *Item) Validate() bool {
    item.Errors = make(map[string]string)

    if strings.TrimSpace(item.Title) == "" {
        item.Errors["Title"] = "title can not be empty"
    }

    if strings.TrimSpace(item.Body) == "" {
        item.Errors["Body"] = "item body can not be empty"
    }

    return len(item.Errors) == 0
}

