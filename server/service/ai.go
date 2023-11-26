package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework_platform/internal/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sashabaranov/go-openai"
)

type GPTService struct {
	Context string `form:"context"`
}

func (s *GPTService) Handle(c *gin.Context) (any, error) {
	// 需要设置代理,不然访问不到
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")
	os.Setenv("ALL_PROXY", "socks5://127.0.0.1:7890")
	client := openai.NewClient("?") //TODO:有钱人买个再说
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: s.Context,
				},
			},
		},
	)
	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("HTTPS_PROXY")
	os.Unsetenv("ALL_PROXY")
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.Choices[0].Message.Content)
	return nil, nil
}

// TODO:许一涵的TOKEN,就先放在这了
var (
	hostUrl   = "wss://aichat.xf-yun.com/v1.1/chat"
	appid     = "fbc7bd6f"
	apiSecret = "ZWI4N2Q5M2NlZDQ3YzFmMzM4YmY0MGVk"
	apiKey    = "f81f4915c4d8fddaafe5a5755fcc0708"
)

type SparkService struct {
	Context string `form:"context"`
}

func (s *SparkService) Handle(c *gin.Context) (any, error) {
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(utils.AssembleAuthUrl1(hostUrl, apiKey, apiSecret), nil)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 101 {
		panic(utils.ReadResp(resp) + err.Error())
	}

	go func() {
		data := utils.GenParams1(appid, s.Context)
		conn.WriteJSON(data)

	}()

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}

		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			fmt.Println("Error parsing JSON:", err)
			return nil, errors.New("连接失败")
		}
		fmt.Println(string(msg))
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			fmt.Println(data["payload"])
			return nil, errors.New("连接失败")
		}
		status := choices["status"].(float64)
		fmt.Println(status)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content
		} else {
			fmt.Println("收到最终结果")
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			fmt.Println("total_tokens:", totalTokens)
			conn.Close()
			break
		}

	}
	//输出返回结果
	fmt.Println(answer)
	return answer, nil
}
