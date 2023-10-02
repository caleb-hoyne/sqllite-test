package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/caleb-hoyne/sqllite-test/handler"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	id := os.Args[1]
	idStr, _ := strconv.Atoi(id)
	operation := os.Args[2]

	switch operation {
	case "get":
		resp, err := http.Get("http://localhost:8080/name?id=" + id)
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
	case "post":
		req := handler.PostNameReq{
			ID:   idStr,
			Name: "some name " + id,
		}
		data, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		httpReq, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/name", bytes.NewReader(data))
		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.StatusCode)
	case "put":
		newName := os.Args[3]
		req := handler.PostNameReq{
			ID:   idStr,
			Name: newName,
		}
		data, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		httpReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080/name/%s", newName), bytes.NewReader(data))
		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.StatusCode)
	}
}
