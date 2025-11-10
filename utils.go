package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// GetImageUrl 根据输入的本地图片地址（可能是相对地址），获取图片的绝对地址
func GetImageUrl(imagePath string) string {
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		// 如果转换失败，返回原始路径
		return imagePath
	}
	return absPath
}

// PrintBeautiful 打印带分隔符的输出
func PrintBeautiful(fn func()) {
	fmt.Printf("==========================================\n")
	fn()
	fmt.Printf("==========================================\n")
}

// PerformObfuscationToString 对输入的字符串进行模糊处理，返回处理后的字符串
func PerformObfuscationToString(input string, scope int) string {
	runes := []rune(input)
	length := len(runes)
	if length <= scope {
		for i := range runes {
			runes[i] = '*'
		}
		return string(runes)
	}
	for i := length - scope; i < length; i++ {
		runes[i] = '*'
	}
	return string(runes)
}

// ImagecompressionAndencoding 处理成类似格式：data:image/jpeg;base64,CEAAUDBAQEAw...
func ImagecompressionAndencoding(imagePath string, maxSize int64) (string, error) {
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %v", err)
	}

	ext := filepath.Ext(imagePath)
	var prefix string
	switch ext {
	case ".jpg", ".jpeg":
		prefix = "data:image/jpeg;base64,"
	case ".png":
		prefix = "data:image/png;base64,"
	default:
		return "", fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if int64(len(fileData)) <= maxSize {
		return prefix + base64.StdEncoding.EncodeToString(fileData), nil
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return "", fmt.Errorf("解码图片失败: %v", err)
	}

	quality := 95
	for quality > 5 {
		var buf bytes.Buffer

		// 根据图片格式进行压缩
		switch ext {
		case ".jpg", ".jpeg":
			options := &jpeg.Options{Quality: quality}
			err = jpeg.Encode(&buf, img, options)
		case ".png":
			// PNG是无损压缩，这里简化处理
			err = png.Encode(&buf, img)
		default:
			return "", fmt.Errorf("不支持的图片格式: %s", ext)
		}

		if err != nil {
			return "", fmt.Errorf("图片编码失败: %v", err)
		}

		compressedData := buf.Bytes()
		if int64(len(compressedData)) <= maxSize {
			// 满足大小要求，返回base64编码
			return prefix + base64.StdEncoding.EncodeToString(compressedData), nil
		}

		// 降低质量继续尝试
		quality -= 10
	}

	return "", fmt.Errorf("无法将图片压缩到 %d 字节以下", maxSize)
}

func GenPromptTemplateByCount(titleCount int, contentCount int) string {
	prompt := fmt.Sprintf(`
请根据图片内容生成符合要求的文本描述，必须严格按照以下格式输出：

标题: [不超过%d个字，需包含具体地点信息]
描述: [不超过%d个字，需包含时间、地点、主要人物/主体、行为/事件等核心要素]

输出要求：
1. 内容必须完全基于图片信息，不得虚构或猜测
2. 时间描述应具体（如：2023年春天、上周三下午等），不确定时可写"近日"或"某天、秋季等笼统的时间词汇"
3. 地点需具体到城市或明确的场所（如：北京天安门广场、某小区花园等），若无法分辨地点则用环境特征代替（如：拥挤的集市、某片荒芜的平原等）
4. 人物描述仅限图片中清晰可辨的部分，不得虚构身份信息
5. 不得添加任何主观评价、修辞或额外解释
6. 严格遵守上述字数限制
7. 如果判断为艺术价值较高的摄影作品或绘画等艺术作品，需要更深入的进行鉴赏和剖析

示例:
标题: 北京迎来鼠年首场降雪
描述: 近日，北京迎来鼠年首场降雪。图片显示市民在雪中行走，街道被白雪覆盖。

请仅输出标题和描述两行内容，不要包含其他任何文字:`, titleCount, contentCount)
	return prompt
}
