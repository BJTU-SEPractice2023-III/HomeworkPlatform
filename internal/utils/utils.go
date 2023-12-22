package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// GetTimeStamp ...
func GetTimeStamp() string { //获得时间戳
	return time.Now().Format("2006-01-02-15-04-05")
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStr(n int) string { //得到随机字符串
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// CopyFile ...
func CopyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// DeletePath
func DeletePath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

// CopyDir ...
func CopyDir(srcDir string, dstDir string) error {
	// // // fmt.Print([]byte(dstDir))
	err := os.Mkdir(dstDir, 0666)
	if err != nil {
		// log.Println(err)
	}
	fileInfoList, _ := ioutil.ReadDir(srcDir)
	for i := 0; i < len(fileInfoList); i++ {
		// // // fmt.Println("Copying: ", fileInfoList[i].Name(), fileInfoList[i].IsDir(), "...")
		if fileInfoList[i].IsDir() {
			CopyDir(path.Join(srcDir, fileInfoList[i].Name()), path.Join(dstDir, fileInfoList[i].Name()))
		} else {
			CopyFile(path.Join(srcDir, fileInfoList[i].Name()), path.Join(dstDir, fileInfoList[i].Name()))
		}
	}
	return nil
}

func ForwardStd(f io.ReadCloser, c chan string) {
	// for {
	// 	// log.Print("1")
	// }
	defer func() {
		recover()
	}()
	cache := ""
	buf := make([]byte, 1024)
	for {
		num, err := f.Read(buf)
		if err != nil && err != io.EOF { // 非EOF错误
			log.Panicln(err)
		}
		if num > 0 {
			str := cache + string(buf[:num])
			lines := strings.SplitAfter(str, "\n") // 按行分割
			for i := 0; i < len(lines)-1; i++ {
				c <- lines[i]
			}
			cache = lines[len(lines)-1] // 最后一行下次循环处理
		}
	}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func EncodePassword(password string, salt string) string { //一种单向加密算法
	//计算 Salt 和密码组合的SHA1摘要
	hash := sha1.New()
	hash.Write([]byte(password + salt))
	bs := hex.EncodeToString(hash.Sum(nil))
	//存储 Salt 值和摘要， ":"分割
	return salt + ":" + string(bs)
}

//下面是星火api

// 生成参数
func GenParams1(appid, question string) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名

	messages := []Message{
		{Role: "user", Content: question},
	}

	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": appid, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      "generalv3",  // 根据实际情况修改返回的数据结构和字段名
				"temperature": float64(0.8), // 根据实际情况修改返回的数据结构和字段名
				"top_k":       int64(6),     // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  int64(2048),  // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",    // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// 创建鉴权url  apikey 即 hmac username
func AssembleAuthUrl(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		// // fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// // // fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// // // fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func ReadResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ImageMessage struct {
	Role        string `json:"role"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

// 将multipart.FileHeader变成字节流
func FileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}
