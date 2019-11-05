package main

import (
	"fmt"
	"os/exec"
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
		_, err := exec.Command("/bin/sh", v).Output()
		if err == nil {
			utils.Log2file(fmt.Sprintf("部署成功：%s", v), GetLogName(task.Id))
		} else {
			utils.Log2file(fmt.Sprintf("部署失败：%s %s", v, err), GetLogName(task.Id))
		}
	}
	queue = queue[:0]
	running = false
	checkoutTaskStatus()
}
