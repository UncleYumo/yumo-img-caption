package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	ApiUrl        = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	ApiModel      = "qwen-vl-plus"
	ApiKeyEnvName = "YUMO_IMG_CAPTION_QWEN_API_KEY"
)

func main() {
	// 获取环境变量的值
	qwenApiKey := os.Getenv(ApiKeyEnvName)

	if qwenApiKey == "" {
		fmt.Println("系统未检索到有效的环境变量：YUMO_IMG_CAPTION_QWEN_API_KEY，请在阿里云控制台创建并获取 API_KEY。")
		return
	}

	// 解析命令行参数
	modelUrl := flag.String("url", ApiUrl, "API 请求地址")
	inputFile := flag.String("file", "", "需要处理的图片文件路径")
	apiModel := flag.String("model", ApiModel, "视觉理解模型名称")
	showInfoFlag := flag.Bool("info", false, "是否显示系统主要参数")
	saveBase64Flag := flag.Bool("base64", false, "是否保存图片的 base64 编码到文件")
	recommendedContentCount := flag.Int("content_count", 30, "生成内容的推荐字数")
	recommendedTitleCount := flag.Int("title_count", 20, "生成标题的推荐字数")
	modelPrompt := flag.String("prompt", "", "模型提示语")
	flag.Parse()

	// 在解析命令行参数后再生成提示词模板
	if *modelPrompt == "" {
		*modelPrompt = GenPromptTemplateByCount(*recommendedTitleCount, *recommendedContentCount)
	}

	imgAbsFilePath := GetImageUrl(*inputFile)

	if *showInfoFlag {
		PrintBeautiful(func() {
			fmt.Printf("API 请求地址：%s\n", *modelUrl)
			fmt.Printf("API 密钥：%s\n", PerformObfuscationToString(qwenApiKey, 12))
			fmt.Printf("视觉理解模型：%s\n", *apiModel)
			fmt.Printf("图片地址相对地址文件：%s\n", *inputFile)
			fmt.Printf("图片地址绝对地址文件：%s\n", imgAbsFilePath)
			fmt.Printf("视觉理解提示词：\n%s\n", *modelPrompt)
			fmt.Printf("生成内容的推荐字数：%d\n", *recommendedContentCount)
			fmt.Printf("生成标题的推荐字数：%d\n", *recommendedTitleCount)
		})
	}

	if *inputFile == "" {
		PrintBeautiful(func() {
			fmt.Println("系统未检索到有效的图片文件路径，请使用 -file 指定图片文件路径。")
		})
		return
	}

	PrintBeautiful(func() {
		fmt.Printf("正在解析图片文件：%s\n", *inputFile)

		if *apiModel != "" {
			fmt.Printf("正在使用指定模型：%s\n", *apiModel)
		} else {
			fmt.Printf("正在使用默认模型：%s\n", ApiModel)
		}
	})

	base64Img, err := ImagecompressionAndencoding(imgAbsFilePath, 1024*1024*4)
	if err != nil {
		PrintBeautiful(func() {
			fmt.Println("图片转换失败：", err)
		})
		return
	}

	if *saveBase64Flag {
		PrintBeautiful(func() {
			err := os.WriteFile("base64.txt", []byte(base64Img), 0644) // 0644表示权限为读写
			if err != nil {
				fmt.Println("写入文件失败：", err)
				return
			}
			fmt.Println("图片 base64 编码已保存到 base64.txt 文件中。")
		})
	}

	postApiData := PostApiData{
		Model: *apiModel,
		Messages: []Messages{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "image_url",
						ImageUrl: &ImageUrl{
							Url: base64Img,
						},
					},
					{
						Type: "text",
						Text: *modelPrompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(postApiData)
	if err != nil {
		PrintBeautiful(func() {
			fmt.Println("JSON请求体转换失败：", err)
		})
		return
	}
	request, err := http.NewRequest("POST", ApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		PrintBeautiful(func() {
			fmt.Println("创建 HTTP 请求失败：", err)
		})
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+qwenApiKey)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		PrintBeautiful(func() {
			fmt.Println("发送 HTTP 请求失败：", err)
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭 HTTP 响应失败：", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		PrintBeautiful(func() {
			fmt.Println("读取 HTTP 响应失败：", err)
		})
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("请求失败，状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
		return
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		PrintBeautiful(func() {
			fmt.Println("解析 JSON 响应失败：", err)
		})
		return
	}

	PrintBeautiful(func() {
		fmt.Println("模型生成描述成功")
		fmt.Println("\n" + apiResponse.Choices[0].Message.Content + "\n")
		fmt.Printf("提示词用量：\t%d\ttokens\n", apiResponse.Usage.PromptTokens)
		fmt.Printf("生成用量：\t%d\ttokens\n", apiResponse.Usage.CompletionTokens)
		fmt.Printf("总用量：\t%d\ttokens\n", apiResponse.Usage.TotalTokens)
	})
}
