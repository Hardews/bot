package main

import (
	"bot/check"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron"
	"github.com/xuri/excelize/v2"
)

var classmate = make(map[string]info, 33)
var names = make(map[string]info)
var mu sync.Mutex

type info struct {
	QNum string
	xh   string
}

func main() {
	fmt.Println("程序运行中....")

	// 初始化
	GetInfo()

	// 创建定时任务为 9 12 14 20点整查询
	c := cron.New()
	spec := "0 0 9,12,14,20 * * *"
	err := c.AddFunc(spec, func() {
		fmt.Println("server start")
		Server()
		time.Sleep(2 * time.Minute)
	})
	if err != nil {
		fmt.Println(err)
	}
	c.Start()
	select {}
}

func Server() {
	fmt.Println("check start")

	var num = 0

checks:
	// 用班上所有人的学号查询
	for s, i := range classmate {
		res := check.IsClock(i.xh)
		if !res {
			// 将未打卡的人加入名单
			names[s] = classmate[s]
		}
	}

	fmt.Println("check successful")
	if len(names) == 0 {
		DoNot()
	} else if len(names) == len(classmate) {
		// 如果全部人都要艾特，则删掉重来
		// 超过三次出bug发信息给我
		if num == 3 {
			IfBug("因为艾特全部人")
		}
		num++

		// 清空map
		for s := range names {
			delete(names, s)
		}

		// 重返
		goto checks
	}

	DoAt()
}

// GroupServer 给群组发送消息
func GroupServer(msg string) {
	// 机器人发送信息接口（http）
	url := "http://127.0.0.1:5700/send_group_msg?group_id=676416672&message=" + msg
	var req *http.Request
	var client = &http.Client{}
	var err error
	// 新建请求
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("new request failed,err:", err)
		return
	}

	// 发送请求，不太需要它的回应
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("do request failed,err:", err)
		return
	}

	// 成功
	fmt.Println("successful")

	// 清空map，让它不重复上次的艾特
	for s := range names {
		delete(names, s)
	}
}

func IfBug(m string) {
	msg := "您的程序" + m + "成功出问题了，赶紧去看看看看看看看看吧"
	url := "http://127.0.0.1:5700//send_private_msg?user_id=1225101127&message=" + msg

	var req *http.Request
	var client = &http.Client{}
	var err error
	// 新建请求
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("new request failed,err:", err)
		return
	}

	// 发送请求，不太需要它的回应
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("do request failed,err:", err)
		return
	}

	// 成功
	fmt.Println("successful")
}

func DoNot() {
	var msg string
	msg = "[CQ:face,id=0]居然全部打完卡了袜，那我岂不是失业了？"

	GroupServer(msg)
}

func DoAt() {
	var msg string

	msg = "[今日提醒]"

	// 先用CQ码将需要艾特的人补齐
	for name, _ := range names {
		msg += "[CQ:at,qq=" + names[name].QNum + "]"
	}
	// 提醒打卡
	msg += "[CQ:face,id=30]这些同学别忘记打卡嗷!"

	GroupServer(msg)
}

func GetInfo() {
	fileName := "./13002101.xlsx"
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取 Sheet1 上所有单元格
	for i := 1; i <= 33; i++ {
		str := strconv.Itoa(i)
		var Info info

		xh, err := f.GetCellValue("Sheet1", "A"+str)
		if err != nil {
			fmt.Println(err)
			return
		}
		name, err := f.GetCellValue("Sheet1", "B"+str)
		if err != nil {
			fmt.Println(err)
			return
		}
		QQNum, err := f.GetCellValue("Sheet1", "C"+str)
		if err != nil {
			fmt.Println(err)
			return
		}

		Info.xh = xh
		Info.QNum = QQNum
		classmate[name] = Info
	}
}
