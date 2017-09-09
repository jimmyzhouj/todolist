// Copyright 2017 Jing zhou. All rights reserved.
// Based on the path package, Copyright 2009 The Go Authors.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package main

import (
    "testing"
)

// type User struct {
//     Id uint64
//     Name string
//     Password string
//     Created time.Time
//     Active bool

//     Errors map[string]string  // show field errors
// }

var testUsers = []struct {
    user User
    valid bool
}{
    {User{Id:0, Name:"hello", Password:"111"}, true},

}


func TestUser(t *testing.T) {

    for _, test := range testUsers {
        if s := test.user.Validate(); s != test.valid {
            t.Errorf("function result %v not match, should be %v\n", s, test.valid)
            t.Errorf("user data is %v\n", test.user)
        }       
    }
}