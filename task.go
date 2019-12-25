package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"webhook-go/utils"
)

var running = false
var queue []*TaskQueue

type TaskQueue struct {
	Id string
	//Payload string
}

// AddNewTask add new task
func AddNewTask(id string) {
	queue = append(queue, &TaskQueue{id})
	checkoutTaskStatus()
}

func checkoutTaskStatus() {
	if running {
		return
	}
	if len(queue) > 0 {
		go startTask(queue[0])
	}
}

func startTask(task *TaskQueue) {
	commands := config[task.Id].Commands
	running = true
	for _, v := range commands {
		filePath, err := filepath.Abs(v)
		if err != nil {
			utils.Log2file(fmt.Sprintf("部署失败：%s", err), GetLogName(task.Id))
			return
		}
		out, err := exec.Command("/bin/sh", filePath).Output()
		log.Println(filePath, "执行结果：", string(out))
		utils.Log2file(fmt.Sprintf("%s执行结果：%s", filePath, string(out)), GetLogName(task.Id))
		if err == nil {
			utils.Log2file(fmt.Sprintf("部署成功：%s", filePath), GetLogName(task.Id))
		} else {
			log.Fatal(filePath, "执行错误：", err)
			utils.Log2file(fmt.Sprintf("部署失败：%s %s", filePath, err), GetLogName(task.Id))
			// 如果执行失败则添加到队列末尾，等待稍后继续执行
			//queue = append(queue, task)
		}
	}
	queue = queue[:0]
	running = false
	checkoutTaskStatus()
}
