package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	id := os.Args[1]

	resp, err := http.Get("http://localhost:8080?id=" + id)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)

	var jsonResp any
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, &jsonResp)
	if err != nil {
		panic(err)
	}

	fmt.Println(jsonResp)
}
