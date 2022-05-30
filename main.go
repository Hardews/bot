package main

import (
	"bot/check"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

var classmate = make(map[string]info, 33)
var names = make(map[string]info)

type info struct {
	QNum string
	xh   string
}

func main() {
	fmt.Println("程序运行中....")

	for true {
		if time.Now().Minute() == 0 {
			break
		}
		time.Sleep(20 * time.Second)
	}

	// 初始化
	GetInfo()

	// 进入服务页面
	Server()
}

func Server() {
	// 空转，到达五点或八点时开始下面的程序
	for true {
		if time.Now().Hour() == 16 || time.Now().Hour() == 20 || time.Now().Hour() == 10 {
			break
		} else {
			fmt.Println("睡眠中....")
			time.Sleep(58 * time.Minute)
		}
	}

	// 用班上所有人的学号查询
	for s, i := range classmate {
		res := check.IsClock(i.xh)
		if !res {
			// 将未打卡的人加入名单
			names[s] = classmate[s]
		}
	}

	if len(names) == 0 {
		DoNot()
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

	// 睡眠一小时，防止再次艾特
	time.Sleep(1 * time.Hour)
	Server()
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
