// response struct
//

package main

import (

)


type RespBody struct {
    Meta map[string]string `json:"meta,omitempty"`
    Data interface{} `json:"data,omitempty"`
    URI string `json:"uri,omitempty"`
}


