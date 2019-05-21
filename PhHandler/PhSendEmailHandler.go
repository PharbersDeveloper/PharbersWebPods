package PhHandler

import (
	"Web/PhModel"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type PhSendMailHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

type eMail struct {
	Email string `json:"email"`
}

type responseMail struct {
	Status string	`json:"status"`
	Msg		string 	`json:"msg"`
}

func (h PhSendMailHandler) NewSmsHandler(args ...interface{}) PhSendMailHandler {
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

	return PhSendMailHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

// 等重构吧，太烂了
func (h PhSendMailHandler) SendMail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	response := map[string]interface{}{}

	if err != nil {
		log.Printf("解析Body出错：%v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	mail := eMail{}
	err = json.Unmarshal(body, &mail)

	if err != nil {
		log.Printf("解析Json出错：%v", err)
		http.Error(w, "can't convert Sms struct", http.StatusBadRequest)
		return 1
	}

	requestAccountMail := []byte(`{
		"email": "`+ mail.Email +`",
		"subject": "体验账号",
		"content": "测试内容，等调试结束，会替换",
		"content-type": "text/plain; charset=UTF-8"}`)

	mailResponse, _ := h.sendMail(r, requestAccountMail)

	mailBody, err := ioutil.ReadAll(mailResponse.Body)

	mailStatus := responseMail{}

	json.Unmarshal(mailBody, &mailStatus)

	if mailStatus.Status == "error" {
		enc := json.NewEncoder(w)
		enc.Encode(mailStatus)
		return 1
	}

	applyuser := PhModel.Applyuser{}
	var out PhModel.Applyuser

	condition := bson.M{"email": mail.Email}
	err = h.db.FindOneByCondition(&applyuser, &out, condition)

	requestUserMail := []byte(`{
		"email": "targetuser@pharbers.com",
		"subject": "申请人员",
		"content": "称呼：`+ out.Name +`<br>所在公司与团队：`+ out.Company +`<br>电子邮件：`+ out.Email +`<br>联系电话：`+ out.Phone +`",
		"content-type": "text/html; charset=UTF-8"}`)

	mailResponse, _ = h.sendMail(r, requestUserMail)

	mailBody, err = ioutil.ReadAll(mailResponse.Body)

	mailStatus = responseMail{}

	json.Unmarshal(mailBody, &mailStatus)

	if mailStatus.Status == "error" {
		enc := json.NewEncoder(w)
		enc.Encode(mailStatus)
		return 1
	}

	response["status"] = "success"
	response["msg"] = "邮件发送成功！"

	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 0
}

func (h PhSendMailHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhSendMailHandler) GetHandlerMethod() string {
	return h.Method
}

func (h PhSendMailHandler) sendMail(r *http.Request, content []byte) (*http.Response, error){
	// 拼接转发的URL
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	resource := fmt.Sprint(h.Args[1], "/", h.Args[0], "/", h.Args[2])
	mergeURL := strings.Join([]string{scheme, resource}, "")

	fmt.Println(mergeURL)

	// 转发
	client := &http.Client{}
	req, _ := http.NewRequest("POST", mergeURL, bytes.NewBuffer(content))
	req.Header.Set("Content-Type", "application/json")

	for k, v := range r.Header {
		req.Header.Add(k, v[0])
	}

	mailResponse, err := client.Do(req)

	return mailResponse, err
}

