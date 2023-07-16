package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	model "hqujwc/types"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var loginurl string

func GetGsSession(requestBody *model.LoginRequestBody) (gsSession string, err error) {
	client := &http.Client{}
	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("编码 JSON 失败:", err)
		return
	}
	req, err := http.NewRequest("POST", loginurl, bytes.NewReader(jsonValue))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	var resp *http.Response
	var body map[string]interface{}
	for i := 0; i < 4; i++ {
		resp, err = client.Do(req)
		fmt.Println(req)
		if err != nil {
			fmt.Println("发送请求失败:", err)
			return
		}
		//fmt.Println(resp.Header)
		respJson, _ := ioutil.ReadAll(resp.Body)
		body = JSONToMap(string(respJson))
		fmt.Println(body)
		code := body["code"].(float64)
		if code == 1 {
			break
		}
		msg := body["msg"].(string)
		if msg == "您提供的用户名或者密码有误" {
			return gsSession, errors.New("登录失败")
		}
		if msg == "模拟登录错误" {
			return gsSession, errors.New("由于数据问题，暂不提供研究生查成绩服务")
		}
		if i == 3 {
			return gsSession, errors.New("登录失败")
		}
	}
	defer resp.Body.Close()
	//jwapp, ok := body["msg"].(map[string]interface{})
	//if !ok {
	//	return gsSession, errors.New("gsSession解析失败")
	//}
	//_, ok = jwapp["jwapp"].(map[string]interface{})["cookie"].(string)
	//if !ok {
	//	fmt.Println(1)
	//	return gsSession, errors.New("gsSession解析失败")
	//}
	gsSession = fmt.Sprint(body["msg"].(map[string]interface{})["jwapp"].(map[string]interface{})["cookie"])
	var start int
	for k, v := range gsSession {
		if v == '=' {
			start = k
		}
		if v == ';' {
			gsSession = gsSession[start+1 : k]
			break
		}
	}
	//fmt.Println(gsSession)
	return
}

func GetEmaphome_WEU(gsSession string) (emaphome_WEU string, err error) {
	client := &http.Client{}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req, err := http.NewRequest("GET", "https://jwapp.hqu.edu.cn/jwapp/sys/emaphome/portal/index.do?forceCas=1", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//req.AddCookie(cookie)
	//cookie = &http.Cookie{
	//	Name:  "EMAP_LANG",
	//	Value: "zh",
	//}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Cookies())
	cookies := fmt.Sprint(resp.Cookies())
	var start int
	for k, v := range cookies {
		if v == '=' {
			start = k
		}
		if v == ';' {
			emaphome_WEU = cookies[start+1 : k]
			break
		}
	}
	//fmt.Println(emaphome_WEU)
	return
}

func GetCjcx_WEU(gsSession string, emaphome_WEU string) (cjcx_WEU string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://jwapp.hqu.edu.cn/jwapp/sys/cjcx/*default/index.do?EMAP_LANG=zh", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	cookie = &http.Cookie{
		Name:    "_WEU",
		Value:   emaphome_WEU,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Cookies())
	cookies := fmt.Sprint(resp.Cookies())
	var start int
	for k, v := range cookies {
		if v == '=' && k > 3 && cookies[k-4:k] == "_WEU" {
			start = k
		}
		if v == ';' && start != 0 {
			cjcx_WEU = cookies[start+1 : k]
			break
		}
	}
	//fmt.Println(cjcx_WEU)
	return
}

func GetInfo(gsSession string, cjcx_WEU string, semester string, year string, stuNum string) (grades []model.GradeInfo) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://jwapp.hqu.edu.cn/jwapp/sys/cjcx/modules/cjcx/xscjcx.do", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	cookie = &http.Cookie{
		Name:    "_WEU",
		Value:   cjcx_WEU,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()
	respJson, _ := ioutil.ReadAll(resp.Body)
	body := JSONToMap(string(respJson))
	gradeResponse := body["datas"].(map[string]interface{})["xscjcx"].(map[string]interface{})["rows"]
	stuYear, _ := strconv.Atoi(stuNum[:2])
	semesterYear, _ := strconv.Atoi(year)
	startAcademicYear := strconv.Itoa(stuYear + semesterYear - 1)
	endAcademicYear := strconv.Itoa(stuYear + semesterYear)
	//fmt.Println(startAcademicYear)
	//fmt.Println(endAcademicYear)
	var zcj float64
	var pscj string
	var qmcj string
	xq := "20" + startAcademicYear + "-20" + endAcademicYear + "-" + semester
	for _, v := range gradeResponse.([]interface{}) {
		if v.(map[string]interface{})["XNXQDM"] == xq {
			if _, ok := v.(map[string]interface{})["ZCJ"].(float64); !ok {
				zcj = 0
			} else {
				zcj = v.(map[string]interface{})["ZCJ"].(float64)
			}
			if _, ok := v.(map[string]interface{})["PSCJ"].(string); !ok {
				pscj = ""
			} else {
				pscj = fmt.Sprint(v.(map[string]interface{})["PSCJ"].(string))
			}
			if _, ok := v.(map[string]interface{})["QMCJ"].(string); !ok {
				qmcj = ""
			} else {
				qmcj = fmt.Sprint(fmt.Sprint(v.(map[string]interface{})["QMCJ"].(string)))
			}
			grades = append(grades, model.GradeInfo{
				XF:      v.(map[string]interface{})["XF"].(float64),
				XSKCM:   fmt.Sprint(v.(map[string]interface{})["XSKCM"].(string)),
				XSZCJMC: fmt.Sprint(v.(map[string]interface{})["DJCJMC"].(string)),
				ZCJ:     zcj,
				QMCJ:    qmcj,
				PSCJ:    pscj,
			})
		}
	}
	return
}

func GetJwpubapp_WEU(gsSession string, emaphome_WEU string) (jwpubapp_WEU string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://jwapp.hqu.edu.cn/jwapp/sys/jwpubapp/pub/setJwCommonAppRole.do", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	cookie = &http.Cookie{
		Name:    "_WEU",
		Value:   emaphome_WEU,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	cookies := fmt.Sprint(resp.Cookies())
	fmt.Println(cookies)
	var start int
	for k, v := range cookies {
		if v == '=' && k > 3 && cookies[k-4:k] == "_WEU" {
			start = k
		}
		if v == ';' && start != 0 {
			jwpubapp_WEU = cookies[start+1 : k]
			break
		}
	}
	fmt.Println(jwpubapp_WEU)
	return
}

func GetPYFADM(gsSession string, jwpubapp_WEU string) (pyfadm int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://jwapp.hqu.edu.cn/jwapp/sys/byshapp/api/grbg/queryXsjbxx.do", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	cookie = &http.Cookie{
		Name:    "_WEU",
		Value:   jwpubapp_WEU,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	respJson, _ := ioutil.ReadAll(resp.Body)
	body := JSONToMap(string(respJson))
	pyfadm, err = strconv.Atoi(fmt.Sprint(body["datas"].(map[string]interface{})["queryXsjbxx"].(map[string]interface{})["faArr"].([]interface{})[0].(map[string]interface{})["PYFADM"]))
	return
}

// 这个是对应学业监测统计的，等学业监测统计更新完成就可以使用
func GetGPA(gsSession string, jwpubapp_WEU string, pyfadm int) (xytjs []model.Xytj, err error) {
	baseURL := "https://jwapp.hqu.edu.cn/jwapp/sys/byshapp/modules/grbg/cxxsxqxf.do"
	params := url.Values{}
	params.Add("PYFADM", fmt.Sprintf("%d", pyfadm))
	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Println(url)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	client := &http.Client{}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	cookie = &http.Cookie{
		Name:    "_WEU",
		Value:   jwpubapp_WEU,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	respJson, _ := ioutil.ReadAll(resp.Body)
	body := JSONToMap(string(respJson))
	xqytjResponse := body["datas"].(map[string]interface{})["cxxsxqxf"].(map[string]interface{})["rows"]
	for _, v := range xqytjResponse.([]interface{}) {
		xytjs = append(xytjs, model.Xytj{
			XH:     fmt.Sprint(v.(map[string]interface{})["XH"].(string)),
			XQGPA:  v.(map[string]interface{})["XQGPA"].(float64),
			YXXF:   v.(map[string]interface{})["YXXF"].(float64),
			YHXF:   v.(map[string]interface{})["YHXF"].(float64),
			XNXQDM: fmt.Sprint(v.(map[string]interface{})["XNXQDM"].(string)),
			BJGXF:  v.(map[string]interface{})["BJGXF"].(float64),
			LJGPA:  v.(map[string]interface{})["LJGPA"].(float64),
			WLCJXF: v.(map[string]interface{})["WLCJXF"].(float64),
		})
	}
	return
}

func GetSessionId(gsSession string, cjcx_WEU string, XH int) (sessionId string, err error) {
	client := &http.Client{}
	body := model.GetSessionIBody{
		Reportlet: "cjcx/xscjpmtj.cpt",
		XH:        XH,
		BBWID:     "",
	}
	jsonValue, err := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://jwapp.hqu.edu.cn/jwapp/sys/frReport2/show.do", bytes.NewReader(jsonValue))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	cookie = &http.Cookie{
		Name:    "_WEU",
		Value:   cjcx_WEU,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析HTML文档失败:", err)
		return
	}
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptContent := s.Text()
		if strings.Contains(scriptContent, "currentSessionID") {
			currentSessionID := GetCurrentSessionID(scriptContent)
			sessionId = currentSessionID
		}
	})
	fmt.Println(sessionId)
	return
}

func GetAllGpa(gsSession string, cjcx_WEU string, sessionId string) (result []string, err error) {
	baseURL := "https://jwapp.hqu.edu.cn/jwapp/sys/frReport2/show.do"
	params := url.Values{}
	params.Add("sessionID", sessionId)
	params.Add("pn", "1")
	params.Add("op", "page_content")
	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Println(url)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	client := &http.Client{}
	cookie := &http.Cookie{
		Name:    "GS_SESSIONID",
		Value:   gsSession,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	cookie = &http.Cookie{
		Name:    "_WEU",
		Value:   cjcx_WEU,
		Expires: time.Now().Add(5 * time.Second),
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	reader := transform.NewReader(resp.Body, unicode.UTF8.NewDecoder())
	doc2, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("解析HTML文档失败:", err)
		return
	}
	//// 解析HTML代码为节点树
	//gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	//utf8Bytes, err := gbkDecoder.Bytes(doc2)
	doc, err := html.Parse(strings.NewReader(string(doc2)))
	if err != nil {
		log.Fatal(err)
	}
	// 开始遍历节点树
	Traverse(doc, &result)
	fmt.Println(result)
	return
}
