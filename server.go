package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

	// 解析参数，填充到Form、PostForm
	r.ParseForm()
	id := r.Form["id"][0]
	if id == "" || len(id) <= 0 {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"id is empty\"}")
		return
	}

	if !VerifySignature(r.Header, string(bodyContent), config[id].Secret) {
		utils.Log2file("验证失败", config[id].Logfile)
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Signature error\"}")
	}
	fmt.Fprintln(w, "{\"code\":200, \"description\":\"OK\"}")
	utils.Log2file("验证通过,启动部署任务", config[id].Logfile)
	AddNewTask(id)
}

// 验证Signature
func VerifySignature(header http.Header, data string, secret string) bool {
	signature := header.Get("X-Hub-Signature")
	if signature != "" && len(signature) > 0 {
		signature = strings.Split(signature, "=")[0]
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
