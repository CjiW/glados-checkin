package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct{
	Cookie string `json:"cookie"`
}

type Msg struct {
	Code int `json:"code"`
	Message string `json:"message"`
	List []List `json:"list"`
}
type List struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Time int `json:"time"`
	Asset string `json:"asset"`
	Business string `json:"business"`
	Change string `json:"change"`
	Balance string `json:"balance"`
	Detail string `json:"detail"`
}

func (msg Msg)output(){
	fmt.Printf("%v\n", time.Now())
	fmt.Printf("Code: %d\n", msg.Code)
	fmt.Printf("Message: %s\n", msg.Message)
	if msg.Code < 0 {
		return
	}
	fmt.Printf("Balance: %s day\n", strings.Split(msg.List[0].Balance, ".")[0])
}

func main() {

	url := "https://glados.network/api/user/checkin"
	method := "POST"

	payload := strings.NewReader(`{"token": "glados.network"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	var cfg Config
	file, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return
	}
	ss := strings.Split(cfg.Cookie, ";")
	s1 := strings.Split(ss[2], "koa:sess=")[1]
	s2 := strings.Split(ss[3], "koa:sess.sig=")[1]
	req.AddCookie(&http.Cookie{
		Name: "koa:sess",
		Value: s1,
		Path: "/",
		Domain: "glados.network",
	})
	req.AddCookie(&http.Cookie{
		Name: "koa:sess.sig",
		Value: s2,
		Path: "/",
		Domain: "glados.network",
	})
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var msg Msg
	json.Unmarshal(body, &msg)
	msg.output()
}
