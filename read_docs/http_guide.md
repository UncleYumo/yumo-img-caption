n# Go语言HTTP请求完全指南

本指南将详细介绍如何在Go语言中发起HTTP请求，包括POST请求、设置Header、构建请求体、处理JSON数据等核心操作。无论你是初学者还是有一定经验的开发者，都能通过本文掌握Go语言中的HTTP操作。

## 目录
1. [基础HTTP请求](#基础http请求)
2. [构建POST请求](#构建post请求)
3. [设置请求Header](#设置请求header)
4. [构建JSON请求体](#构建json请求体)
5. [处理响应数据](#处理响应数据)
6. [解析JSON响应](#解析json响应)
7. [完整示例](#完整示例)
8. [错误处理](#错误处理)
9. [进阶技巧](#进阶技巧)

## 基础HTTP请求

Go语言标准库提供了强大的[net/http](https://pkg.go.dev/net/http)包来处理HTTP请求。最基本的GET请求如下：

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    // 发起GET请求
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // 读取响应体
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // 输出响应内容
    fmt.Printf("Status: %s\n", resp.Status)
    fmt.Printf("Body: %s\n", string(body))
}
```

关键点说明：
- 使用[http.Get](https://pkg.go.dev/net/http#Get)函数发起GET请求
- 必须关闭响应体([resp.Body.Close()](https://pkg.go.dev/net/http#Response.Body))，通常使用defer语句
- 使用[io.ReadAll](https://pkg.go.dev/io#ReadAll)读取响应体内容

## 构建POST请求

要发送POST请求，我们需要使用更灵活的[http.NewRequest](https://pkg.go.dev/net/http#NewRequest)方法或[http.Post](https://pkg.go.dev/net/http#Post)函数。

### 方法一：使用[http.Post](https://pkg.go.dev/net/http#Post)函数

```go
package main

import (
    "fmt"
    "strings"
    "net/http"
    "io"
)

func main() {
    // 准备POST数据
    postData := strings.NewReader("name=张三&age=25")

    // 发送POST请求
    resp, err := http.Post(
        "https://api.example.com/users", 
        "application/x-www-form-urlencoded", 
        postData)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // 读取响应
    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response: %s\n", string(body))
}
```

### 方法二：使用[http.NewRequest](https://pkg.go.dev/net/http#NewRequest)（推荐）

这种方法更加灵活，可以更好地控制请求：

```go
package main

import (
    "fmt"
    "strings"
    "net/http"
    "io"
)

func main() {
    // 创建请求体
    postData := strings.NewReader("name=张三&age=25")

    // 创建请求对象
    req, err := http.NewRequest("POST", "https://api.example.com/users", postData)
    if err != nil {
        panic(err)
    }

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // 读取响应
    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response: %s\n", string(body))
}
```

## 设置请求Header

在使用[http.NewRequest](https://pkg.go.dev/net/http#NewRequest)创建请求后，可以通过[req.Header.Set()](https://pkg.go.dev/net/http#Header.Set)方法设置Header：

```go
package main

import (
    "fmt"
    "strings"
    "net/http"
    "io"
)

func main() {
    postData := strings.NewReader("name=张三&age=25")
    
    req, err := http.NewRequest("POST", "https://api.example.com/users", postData)
    if err != nil {
        panic(err)
    }

    // 设置请求头
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", "MyApp/1.0")
    req.Header.Set("Authorization", "Bearer your-token-here")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response: %s\n", string(body))
}
```

常用Header设置：
- [req.Header.Set("Content-Type", "application/json")](https://pkg.go.dev/net/http#Header.Set) - 设置内容类型为JSON
- [req.Header.Set("Authorization", "Bearer token")](https://pkg.go.dev/net/http#Header.Set) - 设置认证Token
- [req.Header.Set("User-Agent", "MyApp/1.0")](https://pkg.go.dev/net/http#Header.Set) - 设置用户代理

## 构建JSON请求体

在实际开发中，我们经常需要发送JSON格式的数据。以下是几种构建JSON请求体的方法：

### 方法一：使用结构体和json.Marshal

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

// 定义数据结构
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
    Email string `json:"email"`
}

func main() {
    // 创建数据对象
    user := User{
        Name: "张三",
        Age: 25,
        Email: "zhangsan@example.com",
    }

    // 将结构体转换为JSON
    jsonData, err := json.Marshal(user)
    if err != nil {
        panic(err)
    }

    // 创建请求
    req, err := http.NewRequest("POST", "https://api.example.com/users", bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }

    // 设置正确的Content-Type
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Printf("Status: %s\n", resp.Status)
}
```

### 方法二：直接使用JSON字符串

```go
package main

import (
    "bytes"
    "fmt"
    "net/http"
)

func main() {
    // 直接定义JSON字符串
    jsonStr := `{"name":"张三","age":25,"email":"zhangsan@example.com"}`

    // 创建请求
    req, err := http.NewRequest("POST", "https://api.example.com/users", bytes.NewBufferString(jsonStr))
    if err != nil {
        panic(err)
    }

    // 设置Content-Type
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Printf("Status: %s\n", resp.Status)
}
```

## 处理响应数据

HTTP响应包含状态码、Header和响应体等多个部分：

```go
package main

import (
    "fmt"
    "net/http"
    "io"
)

func main() {
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // 检查状态码
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
        return
    }

    // 读取响应体
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // 输出响应Header
    fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
    fmt.Printf("响应内容: %s\n", string(body))
}
```

重要字段说明：
- [resp.StatusCode](https://pkg.go.dev/net/http#Response.StatusCode) - HTTP状态码（如200、404、500等）
- [resp.Header](https://pkg.go.dev/net/http#Response.Header) - 响应头信息
- [resp.Body](https://pkg.go.dev/net/http#Response.Body) - 响应体内容

## 解析JSON响应

接收JSON响应并解析是常见的需求，我们可以使用[json.Unmarshal](https://pkg.go.dev/encoding/json#Unmarshal)函数：

### 方法一：解析为结构体

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
)

// 定义响应结构体
type APIResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    } `json:"data"`
}

func main() {
    resp, err := http.Get("https://api.example.com/user/1")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // 解析JSON到结构体
    var result APIResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        panic(err)
    }

    // 使用解析后的数据
    fmt.Printf("Success: %t\n", result.Success)
    fmt.Printf("Message: %s\n", result.Message)
    fmt.Printf("User ID: %d\n", result.Data.ID)
    fmt.Printf("User Name: %s\n", result.Data.Name)
}
```

### 方法二：解析为map

当JSON结构不确定时，可以解析为map：

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
)

func main() {
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // 解析为map
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        panic(err)
    }

    // 访问数据
    fmt.Printf("解析结果: %+v\n", result)
    // 访问特定字段
    if name, ok := result["name"].(string); ok {
        fmt.Printf("Name: %s\n", name)
    }
}
```

## 完整示例

结合前面的知识，我们来看一个完整的示例，向API发送用户数据并处理响应：

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "io"
)

// 请求数据结构
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// 响应数据结构
type APIResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    struct {
        ID    int    `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    } `json:"data"`
}

func main() {
    // 1. 准备请求数据
    requestData := CreateUserRequest{
        Name:  "张三",
        Email: "zhangsan@example.com",
        Age:   25,
    }

    // 2. 将数据序列化为JSON
    jsonData, err := json.Marshal(requestData)
    if err != nil {
        fmt.Printf("JSON序列化失败: %v\n", err)
        return
    }

    // 3. 创建请求
    req, err := http.NewRequest("POST", "https://api.example.com/users", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("创建请求失败: %v\n", err)
        return
    }

    // 4. 设置请求头
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer your-api-token")
    req.Header.Set("User-Agent", "MyApp/1.0")

    // 5. 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("发送请求失败: %v\n", err)
        return
    }
    defer resp.Body.Close()

    // 6. 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("读取响应失败: %v\n", err)
        return
    }

    // 7. 检查状态码
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        fmt.Printf("请求失败，状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
        return
    }

    // 8. 解析JSON响应
    var apiResp APIResponse
    err = json.Unmarshal(body, &apiResp)
    if err != nil {
        fmt.Printf("解析响应JSON失败: %v\n", err)
        return
    }

    // 9. 使用响应数据
    if apiResp.Success {
        fmt.Printf("用户创建成功!\n")
        fmt.Printf("用户ID: %d\n", apiResp.Data.ID)
        fmt.Printf("用户名: %s\n", apiResp.Data.Name)
        fmt.Printf("用户邮箱: %s\n", apiResp.Data.Email)
    } else {
        fmt.Printf("用户创建失败: %s\n", apiResp.Message)
    }
}
```

## 错误处理

在实际应用中，良好的错误处理非常重要：

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

func makeRequest() error {
    // 创建带超时的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 创建请求
    req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com/data", nil)
    if err != nil {
        return fmt.Errorf("创建请求失败: %w", err)
    }

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("发送请求失败: %w", err)
    }
    defer resp.Body.Close()

    // 检查HTTP状态码
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP错误，状态码: %d", resp.StatusCode)
    }

    // 解析JSON
    var data map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return fmt.Errorf("解析JSON失败: %w", err)
    }

    fmt.Printf("数据: %+v\n", data)
    return nil
}

func main() {
    if err := makeRequest(); err != nil {
        fmt.Printf("错误: %v\n", err)
    }
}
```

关键改进点：
1. 使用[context.WithTimeout](https://pkg.go.dev/context#WithTimeout)设置请求超时
2. 使用[json.NewDecoder](https://pkg.go.dev/encoding/json#NewDecoder)直接解码，避免额外内存分配
3. 返回错误时包装原始错误，保留错误链

## 进阶技巧

### 1. 自定义HTTP客户端

```go
client := &http.Client{
    Timeout: 30 * time.Second, // 设置超时时间
    Transport: &http.Transport{
        MaxIdleConns:        100,               // 最大空闲连接数
        MaxIdleConnsPerHost: 10,                // 每个主机最大空闲连接数
        IdleConnTimeout:     90 * time.Second,  // 空闲连接超时时间
    },
}
```

### 2. 处理重定向

```go
client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        // 限制重定向次数
        if len(via) >= 10 {
            return fmt.Errorf("重定向次数过多")
        }
        return nil
    },
}
```

### 3. 添加中间件/拦截器

```go
// 日志中间件
func loggingMiddleware(next http.RoundTripper) http.RoundTripper {
    return &loggingRoundTripper{next}
}

type loggingRoundTripper struct {
    next http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    start := time.Now()
    resp, err := l.next.RoundTrip(req)
    duration := time.Since(start)
    
    fmt.Printf("[%s] %s %s (耗时: %v)\n", 
        req.Method, req.URL, resp.Status, duration)
    
    return resp, err
}

// 使用中间件
client := &http.Client{
    Transport: loggingMiddleware(http.DefaultTransport),
}
```

## 总结

通过本文的学习，你应该掌握了以下关键知识点：

1. 如何使用[http.Get](https://pkg.go.dev/net/http#Get)和[http.Post](https://pkg.go.dev/net/http#Post)发送基本请求
2. 如何使用[http.NewRequest](https://pkg.go.dev/net/http#NewRequest)创建更灵活的请求
3. 如何设置请求Header
4. 如何构建JSON请求体
5. 如何处理和解析JSON响应
6. 如何进行错误处理和超时控制

这些技能足以应对大多数HTTP通信场景。在实际项目中，建议封装HTTP客户端以便复用，并根据具体业务需求进行适当的扩展。