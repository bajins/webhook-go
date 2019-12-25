package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"webhook-go/utils"
)

// 启动服务
func StartService(address string, port string) error {
	http.HandleFunc("/", index)
	http.HandleFunc("/webHooks", webHooks)

	utils.Log2file(fmt.Sprintf("service starting... %s:%s", address, port), "")
	return http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil)
}

// 首页
func index(w http.ResponseWriter, r *http.Request) {
	utils.Log2file(string(r.URL.Host), "")
	fmt.Fprintln(w, "{\"code\":200, \"description\":\"service running...\"}")
}

// 自动编译
func webHooks(w http.ResponseWriter, r *http.Request) {
	if err := LoadConfig(); err != nil {
		utils.Log2file(err.Error(), "")
		os.Exit(1)
	}
	if strings.ToUpper(r.Method) != "POST" {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Error Method or unknow request url\"}")
		return
	}
	if !VerifyEvent(r.Header, "push") {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unmatch x-github-event\"}")
		return
	}
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Error Response Body\"}")
		return
	}
	defer r.Body.Close()
	if len(bodyContent) <= 0 {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Response Body is empty\"}")
		return
	}

	var contentMap map[string]interface{}

	ct := r.Header.Get("Content-Type")
	// 针对`Gitea`请使用`v1.10.0-rc2`以下版本issue
	// https://github.com/go-gitea/gitea/issues/7700
	if ct == "" || len(ct) <= 0 {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		ct = r.Header.Get("Content-Type")
	}
	ct = strings.ToLower(ct)
	if ct == "application/x-www-form-urlencoded" {
		// 恢复Body内容
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyContent))
		// 解析参数，填充到Form、PostForm
		err = r.ParseForm()
		if r.Form == nil || len(r.Form) <= 0 {
			fmt.Fprintln(w, "{\"code\":200, \"error\":\"Response Body is empty\"}")
			return
		}
		payload := r.Form["payload"][0]
		json.Unmarshal([]byte(payload), &contentMap)
	} else if ct == "application/json" {
		err = json.Unmarshal(bodyContent, &contentMap)
		if err != nil {
			fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unmatch Response Body\"}")
			return
		}
	}
	if contentMap == nil || contentMap["repository"] == nil {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unmatch Response Body\"}")
		return
	}
	id := contentMap["repository"].(map[string]interface{})["full_name"].(string)
	log.Println("当前full_name：", id)
	//id = strings.ToLower(id)
	config := config[id]
	log.Println("当前配置：", config)
	if reflect.DeepEqual(config, Config{}) {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Config is not found\"}")
		return
	}

	if !VerifySignature(r.Header, string(bodyContent), config.Secret) {
		utils.Log2file("验证失败", config.Logfile)
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Signature error\"}")
		return
	}
	fmt.Fprintln(w, "{\"code\":200, \"description\":\"OK\"}")
	utils.Log2file("验证通过,启动部署任务", config.Logfile)
	AddNewTask(id)
}

// 验证Signature
func VerifySignature(header http.Header, data string, secret string) bool {
	signature := header.Get("X-Hub-Signature")
	if signature != "" && len(signature) > 0 {
		signature = strings.Split(signature, "=")[1]
		return signature == utils.ComputeHash1(data, secret)
	}
	signature = header.Get("X-Gitea-Signature")

	if signature == "" || len(signature) <= 0 {
		signature = header.Get("X-Gogs-Signature")
	}
	return signature == utils.ComputeHmacSha256(data, secret)
}

// 验证Event
func VerifyEvent(header http.Header, event string) bool {
	e := header.Get("X-Gitea-Event")

	if e == "" || len(e) <= 0 {
		e = header.Get("X-Gogs-Event")
	}
	if e == "" || len(e) <= 0 {
		e = header.Get("X-GitHub-Event")
	}
	return event == strings.Trim(e, "UTF-8")
}
