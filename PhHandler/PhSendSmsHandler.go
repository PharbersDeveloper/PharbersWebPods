package PhHandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type PhSendSmsHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

type sms struct {
	Phone string `json:"phone"`
	Code string `json:"code"`
}

type responseSms struct {
	Status string	`json:"status"`
	Msg		string 	`json:"msg"`
}

func (h PhSendSmsHandler) NewSmsHandler(args ...interface{}) PhSendSmsHandler {
	var m *BmMongodb.BmMongodb
	var r *BmRedis.BmRedis
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				} else if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
			}
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}

	return PhSendSmsHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

func (h PhSendSmsHandler) SendSms(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	response := map[string]interface{}{}

	if err != nil {
		log.Printf("解析Body出错：%v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	sms := sms{}
	err = json.Unmarshal(body, &sms)

	if err != nil {
		log.Printf("解析Json出错：%v", err)
		http.Error(w, "can't convert Sms struct", http.StatusBadRequest)
		return 1
	}

	// 拼接转发的URL
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	resource := fmt.Sprint(h.Args[1], "/", h.Args[0], "/", h.Args[2])
	mergeURL := strings.Join([]string{scheme, resource}, "")

	fmt.Println(mergeURL)

	requestSms := []byte(`{
		"phone": "`+ sms.Phone +`"
	}`)

	// 转发
	client := &http.Client{}
	req, _ := http.NewRequest("POST", mergeURL, bytes.NewBuffer(requestSms))
	req.Header.Set("Content-Type", "application/json")

	for k, v := range r.Header {
		req.Header.Add(k, v[0])
	}

	smsResponse, err := client.Do(req)

	smsBody, err := ioutil.ReadAll(smsResponse.Body)

	smsStatus := responseSms{}

	json.Unmarshal(smsBody, &smsStatus)

	if smsStatus.Status == "error" {
		enc := json.NewEncoder(w)
		enc.Encode(smsStatus)
		return 1
	}

	response["status"] = "success"
	response["msg"] = "短信发送成功！"

	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 0
}

func (h PhSendSmsHandler) VerifyCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	response := map[string]interface{}{}

	if err != nil {
		log.Printf("解析Body出错：%v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	sms := sms{}
	err = json.Unmarshal(body, &sms)

	if err != nil {
		log.Printf("解析Json出错：%v", err)
		http.Error(w, "can't convert Sms struct", http.StatusBadRequest)
		return 1
	}

	// 拼接转发的URL
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	resource := fmt.Sprint(h.Args[1], "/", h.Args[0], "/", h.Args[2])
	mergeURL := strings.Join([]string{scheme, resource}, "")

	fmt.Println(mergeURL)

	requestSms := []byte(`{
		"phone": "`+ sms.Phone +`",
		"code": "`+ sms.Code +`"
	}`)

	// 转发
	client := &http.Client{}
	req, _ := http.NewRequest("POST", mergeURL, bytes.NewBuffer(requestSms))
	req.Header.Set("Content-Type", "application/json")

	for k, v := range r.Header {
		req.Header.Add(k, v[0])
	}

	smsResponse, err := client.Do(req)

	smsBody, err := ioutil.ReadAll(smsResponse.Body)

	smsStatus := responseSms{}

	json.Unmarshal(smsBody, &smsStatus)

	if smsStatus.Status == "error" {
		enc := json.NewEncoder(w)
		enc.Encode(smsStatus)
		return 1
	}

	response["status"] = "success"
	response["msg"] = "验证成功！"

	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 0
}

func (h PhSendSmsHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhSendSmsHandler) GetHandlerMethod() string {
	return h.Method
}
