# Yumo-File-Caption

Yumo-File-Caption 是一个基于通义千问多模态模型的命令行工具，可以自动为图片生成符合新闻/媒体平台要求的标题和描述。

## 功能特点

- 🤖 基于通义千问VL模型，准确识别图片内容
- 📝 自动生成符合媒体规范的标题和描述
- ⌨️ 命令行操作，简单易用
- 🖼️ 支持 JPG/PNG 格式图片
- 📦 支持 Windows 系统打包为 exe 可执行文件
- 🛠️ 可自定义生成字数限制

## 安装

### 方法一：下载预编译版本（推荐）

从 [Releases](https://github.com/your-username/yumo-file-caption/releases) 页面下载适用于 Windows 的预编译版本。

### 方法二：从源码构建

1. 安装 Go 1.21 或更高版本
2. 克隆本项目：
   ```
   git clone https://github.com/UncleYumo/yumo-img-caption.git
   ```
3. 进入项目目录并构建：
   ```
   cd yumo-file-caption
   go build -o ym-img-caption.exe .
   ```

## 使用方法

### 1. 设置环境变量

在使用之前，需要设置通义千问 API Key：

**Windows (PowerShell):**
```powershell
$env:YUMO_IMG_CAPTION_QWEN_API_KEY = "your_api_key_here"
```

**Windows (命令提示符):**
```cmd
set YUMO_IMG_CAPTION_QWEN_API_KEY=your_api_key_here
```

你可以在 [阿里云百炼平台](https://help.aliyun.com/zh/bailian) 获取 API Key。

### 2. 基本使用

```bash
ym-img-caption -file path/to/your/image.jpg
```

### 3. 高级选项

```bash
# 显示详细信息
ym-img-caption -file image.jpg -info

# 自定义标题和内容字数限制
ym-img-caption -file image.jpg -title_count 15 -content_count 50

# 保存图片的 base64 编码到文件
ym-img-caption -file image.jpg -base64
```

### 4. 命令行参数

| 参数 | 默认值 | 描述 |
|------|--------|------|
| `-file` | 无 | 需要处理的图片文件路径 |
| `-info` | false | 是否显示系统主要参数 |
| `-base64` | false | 是否保存图片的 base64 编码到文件 |
| `-content_count` | 30 | 生成内容的推荐字数 |
| `-title_count` | 20 | 生成标题的推荐字数 |
| `-model` | qwen-vl-plus | 视觉理解模型名称 |
| `-url` | https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions | API 请求地址 |

## 输出格式

工具会严格按照以下格式输出：

```
标题: [包含地点信息的标题]
描述: [包含时间、地点、主体、事件的描述]
```

例如：
```
标题: 北京迎来鼠年首场降雪
描述: 近日，北京迎来鼠年首场降雪。图片显示市民在雪中行走，街道被白雪覆盖。
```

## 开发指南

### 项目结构

```
yumo-file-caption/
├── main.go          # 主程序入口
├── utils.go         # 工具函数
├── domain.go        # 数据结构定义
├── build-win64.bat  # Windows编译脚本
└── README.md        # 说明文档
```

### 本地开发

1. 克隆项目
2. 设置 Go 环境
3. 运行程序：
   ```
   go run main.go -file path/to/image.jpg
   ```

## 打包

### Windows

运行 `build-win64.bat` 脚本：

```
build-win64.bat
```

这将在项目目录下生成 `ym-img-caption.exe` 文件。

## 注意事项

1. 图片大小限制：为保证API调用效率，程序会自动压缩大于4MB的图片
2. API 调用会产生费用，请根据实际需求使用
3. 生成的内容基于AI识别结果，可能需要人工审核

## 许可证

本项目采用 MIT 许可证，详情请见 [LICENSE](LICENSE) 文件。

## 免责声明

本工具由AI生成，生成的内容可能不完全准确，请在使用时进行人工审核。
