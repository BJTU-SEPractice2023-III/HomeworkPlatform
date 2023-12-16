package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"homework_platform/internal/utils"
	"io"
	"mime/multipart"
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
	_, err := client.CreateChatCompletion(
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

	// fmt.Println(resp.Choices[0].Message.Content)
	return nil, nil
}

// TODO:许一涵的TOKEN,就先放在这了
var (
	appid     = "fbc7bd6f"
	apiSecret = "ZWI4N2Q5M2NlZDQ3YzFmMzM4YmY0MGVk"
	apiKey    = "f81f4915c4d8fddaafe5a5755fcc0708"
)

var sseChanMap = make(map[uint]*chan string)

func addSSEChan(id uint, ch *chan string) {
	sseChanMap[id] = ch
}

func removeSSEChan(id uint) {
	sseChanMap[id] = nil
}

type ConnectSpark struct{}

func (s *ConnectSpark) Handle(c *gin.Context) (any, error) {
	// log.Println("ServerConsole")
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")

	id := c.GetUint("ID")
	messageChan := make(chan string)
	addSSEChan(id, &messageChan)
	// log.Printf("[service/ai]: %d connected", id)
	defer removeSSEChan(id)

	c.Stream(func(w io.Writer) bool {
		if message, ok := <-*sseChanMap[id]; ok {
			c.SSEvent("message", message)
			return true
		}
		return false
	})

	// log.Printf("[service/ai]: %d disconnected", id)
	return nil, nil
}

type SparkService struct {
	Context string `form:"context"`
}

func (s *SparkService) Handle(c *gin.Context) (any, error) {
	hostUrl := "wss://spark-api.xf-yun.com/v3.1/chat" //这里调模型是v几,然后调整utils里面的general字段
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	// 握手并建立websocket 连接
	conn, resp, err := d.Dial(utils.AssembleAuthUrl(hostUrl, apiKey, apiSecret), nil)
	if err != nil || resp.StatusCode != 101 {
		return nil, err
	}
	defer conn.Close()

	// 生成参数并发送请求
	data := utils.GenParams1(appid, s.Context)
	conn.WriteJSON(data)

	// 获取返回的数据
	id := c.GetUint("ID")
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("read message error:", err)
			break
		}

		var data map[string]interface{}
		err = json.Unmarshal(msg, &data)
		if err != nil {
			// fmt.Println("Error parsing JSON:", err)
			return nil, errors.New("连接失败")
		}
		// // fmt.Println(string(msg))

		// 解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			// fmt.Println(data["payload"])
			return nil, errors.New("连接失败")
		}
		status := choices["status"].(float64)
		// // fmt.Println(status)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)

		if sseChanMap[id] != nil {
			*sseChanMap[id] <- content
		} else {
			return nil, nil
		}
		if status == 2 {
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			_ = totalTokens
			// fmt.Println("total_tokens:", totalTokens)
			break
		}
	}

	return nil, nil
}

type SparkImageService struct {
	Content string               `form:"content"`
	Files   multipart.FileHeader `form:"file"`
}

func (s *SparkImageService) Handle(c *gin.Context) (any, error) {
	hostUrl := "wss://spark-api.cn-huabei-1.xf-yun.com/v2.1/image"
	if c.ContentType() != "multipart/form-data" {
		return nil, errors.New("not supported content-type")
	}

	var err error
	// 从 Form 获取其他数据
	err = c.ShouldBind(s)
	if err != nil {
		return nil, err
	}

	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(utils.AssembleAuthUrl(hostUrl, apiKey, apiSecret), nil)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 101 {
		return nil, errors.New("连接失败")
	}

	file, err := s.Files.Open()
	if err != nil {
		// fmt.Println("无法打开文件：", err)
		return nil, err
	}
	defer file.Close()
	image, err := io.ReadAll(file)
	// log.Println("content为" + s.Content)
	messages := []utils.ImageMessage{
		{Role: "user", Content: base64.StdEncoding.EncodeToString(image), ContentType: "image"}, // 首条必须是图片
		{Role: "user", Content: s.Content, ContentType: "text"},
	}

	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": appid,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain": "general",
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	conn.WriteJSON(data)

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("read message error:", err)
			break
		}

		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			// fmt.Println("Error parsing JSON:", err)
			return nil, err
		}
		// fmt.Println(string(msg))
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			// fmt.Println(data["payload"])
			return nil, errors.New(data["payload"].(string))
		}
		status := choices["status"].(float64)
		// fmt.Println(status)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content
		} else {
			// fmt.Println("收到最终结果")
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			_ = totalTokens
			// fmt.Println("total_tokens:", totalTokens)
			conn.Close()
			break
		}

	}
	//输出返回结果
	// fmt.Println(answer)
	return answer, nil
}
