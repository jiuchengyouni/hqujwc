package httpUtil

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func DoPost(Url string, cookies map[string]string, jsonValue []byte) (resp *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", Url, bytes.NewReader(jsonValue))
	for k, v := range cookies {
		cookie := &http.Cookie{
			Name:    k,
			Value:   v,
			Expires: time.Now().Add(5 * time.Second),
		}
		req.AddCookie(cookie)
	}
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return nil, errors.New("创建请求失败")
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

func DoGet(loginUrl string, cookies map[string]string) (resp *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://jwapp.hqu.edu.cn/jwapp/sys/cjcx/*default/index.do?EMAP_LANG=zh", nil)
	for k, v := range cookies {
		cookie := &http.Cookie{
			Name:    k,
			Value:   v,
			Expires: time.Now().Add(5 * time.Second),
		}
		req.AddCookie(cookie)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}
