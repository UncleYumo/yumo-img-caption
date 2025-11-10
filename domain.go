package main

type PostApiData struct {
	Model    string     `json:"model"`
	Messages []Messages `json:"messages"`
}

type Messages struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageUrl *ImageUrl `json:"image_url,omitempty"`
}

type ImageUrl struct {
	Url string `json:"url"`
}

type ApiResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint *string  `json:"system_fingerprint"` // 系统指纹，用于标识请求处理的系统版本
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	LogProbs     *any    `json:"logprobs"` // 日志概率信息，通常为null或包含详细的token概率数据
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens            int                      `json:"prompt_tokens"`
	CompletionTokens        int                      `json:"completion_tokens"`
	TotalTokens             int                      `json:"total_tokens"`
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`     // 提示token详细信息
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // 完成token详细信息
}

// PromptTokensDetails 提示token的详细信息
type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"` // 缓存的token数量
}

// CompletionTokensDetails 完成token的详细信息
type CompletionTokensDetails struct {
	TextTokens int `json:"text_tokens"` // 文本token数量
}
