package check

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var now = time.Now().Format("2006-01-02 03:04:05")

func IsClock(xh string) bool {
	var res bool
	res = false
	// we重邮是否打卡的接口地址
	URL := "https://we.cqupt.edu.cn/api/mrdk/get_mrdk_list_test.php"

	// 当前时间戳
	stamp := strconv.FormatInt(time.Now().Unix(), 10)

	// 向we重邮请求时需要的参数key的原型json
	var u = "{\"xh\":\"" + xh + "\",\"openid\":\"3234\",\"timestamp\":" + stamp + "}"

	// 对其进行base64加密后成为key
	key := base64.StdEncoding.EncodeToString([]byte(u))

	//// 形成表单
	postData := url.Values{}
	postData.Add("key", key)

	var resp *http.Response
	var client = &http.Client{}
	var err error

	// 这里是直接发送表单内容的写法，而不需要写进body
	resp, err = client.PostForm(URL, postData)
	if err != nil {
		fmt.Println("do request failed,err:", err)
		return res
	}

	// 请求延时
	time.Sleep(2 * time.Second)

	defer resp.Body.Close()

	// we重邮返回的数据处理，返回的是一个json 格式是 一个大map里 一个字段data 对应以前的三次打卡信息map
	// map[data:map[] map[] map[]]
	var Info map[string]interface{}
	// 这里是对获得的响应进行反序列化
	if err = json.NewDecoder(resp.Body).Decode(&Info); err != nil {
		// 这里是因为如果没打卡，we重邮不会返回任何内容，所以就会导致读不到任何东西
		if err == io.EOF {
			err = nil
			return false
		}
		fmt.Println("write info in failed,err:", err)
		return res
	}

	// 先拿出第一层 三个map
	var s = Info["data"].([]interface{})

	// 只需要最新的打卡记录
	userMap := s[0].(map[string]interface{})
	latestTime := userMap["created_at"]

	// 再次确认是否打卡后返回
	// 如果we重邮有返回值但未打卡，可以通过判断打卡时间来确认是否打卡
	if now[:10] == latestTime.(string)[:10] {
		res = true
		return res
	} else {
		return res
	}
}
