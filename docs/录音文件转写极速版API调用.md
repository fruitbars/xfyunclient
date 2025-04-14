# 讯飞语音识别(xfspeech) Go SDK 使用文档

## 目录

- [简介](#简介)
- [安装](#安装)
- [快速开始](#快速开始)
- [核心功能](#核心功能)
    - [客户端初始化](#客户端初始化)
    - [语音识别](#语音识别)
    - [回调和异步处理](#回调和异步处理)
    - [高级用法](#高级用法)
- [API参考](#api参考)
    - [Client](#client)
    - [AudioSource](#audiosource)
    - [RecognitionOptions](#recognitionoptions)
    - [RecognitionResult](#recognitionresult)
- [常见问题](#常见问题)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)
- [示例代码](#示例代码)

## 简介

一个轻量级、高性能的科大讯飞语音识别Go SDK，提供了简洁易用的API接口，方便开发者快速集成讯飞语音识别服务到Go应用中。SDK支持多种音频输入方式，同步和异步处理，以及灵活的回调机制。

主要特点：

- 支持文件和内存中的音频数据
- 支持同步和异步识别
- 提供状态和结果回调
- 灵活的配置选项
- 自动处理文件分片上传
- 强类型支持

## 安装

使用Go模块安装:

```bash
go get github.com/yourusername/xfspeech
```

## 快速开始

以下是一个简单的示例，展示如何使用SDK进行语音识别：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/fruitbars/xfyunclient/ost"
)

func main() {
    // 初始化客户端
    client := ost.NewClient(
        "your_app_id",     // 讯飞开放平台应用ID
        "your_api_key",    // 讯飞开放平台API Key
        "your_api_secret", // 讯飞开放平台API Secret
    )
    
    // 识别音频文件
    result, err := client.RecognizeFile("audio_sample.wav", nil)
    if err != nil {
        log.Fatalf("识别失败: %v", err)
    }
    
    // 输出识别结果
    fmt.Printf("识别结果: %s\n", result.Data.Result.FullText)
}
```

## 核心功能

### 客户端初始化

创建一个新的客户端实例：

```go
// 基本初始化
client := ost.NewClient("your_app_id", "your_api_key", "your_api_secret")

// 设置超时时间
client = ost.NewClient("your_app_id", "your_api_key", "your_api_secret").WithTimeout(60)

// 设置自定义HTTP客户端
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    // 其他自定义配置...
}
client = ost.NewClient("your_app_id", "your_api_key", "your_api_secret").WithHTTPClient(httpClient)
```

### 语音识别

SDK提供多种方式进行语音识别：

#### 识别文件

```go
// 使用默认选项
result, err := client.RecognizeFile("audio_sample.wav", nil)

// 使用自定义选项
options := ost.CreateOptions().
    WithLanguage("zh_cn").
    WithAccent("mandarin").
    WithDomain("pro_ost_ed")
    
result, err := client.RecognizeFile("audio_sample.wav", options)
```

#### 识别内存中的音频数据

```go
// 读取音频数据
audioData, err := ioutil.ReadFile("audio_sample.wav")
if err != nil {
    log.Fatalf("读取音频文件失败: %v", err)
}

// 识别内存中的数据
result, err := client.RecognizeBytes(audioData, "audio_sample.wav", nil)
```

#### 使用自定义音频源

```go
// 创建自定义音频源
source := ost.NewFileAudioSource("audio_sample.wav")
// 或者
source := ost.NewMemoryAudioSource(audioData, "audio_sample.wav")

// 识别自定义音频源
result, err := client.RecognizeAudio(source, nil)
```

### 回调和异步处理

#### 使用回调监控识别过程

```go
options := ost.CreateOptions().
    WithStatusCallback(func(status string, result *ost.RecognitionResult) {
        fmt.Printf("状态更新: %s\n", ost.StatusToString(status))
    }).
    WithResultCallback(func(text string, result *ost.RecognitionResult) {
        fmt.Printf("中间结果: %s\n", text)
    })

result, err := client.RecognizeFile("audio_sample.wav", options)
```

#### 异步识别

```go
// 创建通道接收结果
done := make(chan struct{})

// 开始异步识别
client.RecognizeFileAsync("audio_sample.wav", nil, func(result *ost.RecognitionResult, err error) {
    if err != nil {
        fmt.Printf("识别失败: %v\n", err)
    } else {
        fmt.Printf("识别结果: %s\n", result.Data.Result.FullText)
    }
    close(done)
})

// 继续执行其他任务...

// 等待识别完成
<-done
```

### 高级用法

#### 分步骤处理

```go
// 1. 上传文件
uploader := ost.NewUploader(
    "your_app_id",
    "your_api_key",
    "your_api_secret", 
    ost.NewFileAudioSource("audio_sample.wav"),
)
fileURL, err := uploader.Upload()
if err != nil {
    log.Fatalf("上传失败: %v", err)
}

// 2. 创建识别任务
taskID, err := client.CreateTask(fileURL, nil)
if err != nil {
    log.Fatalf("创建任务失败: %v", err)
}

// 3. 等待结果
result, err := client.GetResultWithTimeout(taskID, 5*time.Minute)
if err != nil {
    log.Fatalf("获取结果失败: %v", err)
}
```

#### 设置回调URL

```go
options := ost.CreateOptions().
    WithCallbackURL("https://your-server.com/callback")
    
taskID, err := client.CreateTask(fileURL, options)
```

## API参考

### Client

主要客户端接口，提供所有识别功能。

#### 方法

- `NewClient(appID, apiKey, apiSecret string) *Client` - 创建新客户端
- `WithTimeout(timeout int) *Client` - 设置超时时间
- `WithHTTPClient(client interface{}) *Client` - 设置自定义HTTP客户端
- `RecognizeFile(filePath string, options *RecognitionOptions) (*RecognitionResult, error)` - 识别文件
- `RecognizeBytes(data []byte, fileName string, options *RecognitionOptions) (*RecognitionResult, error)` - 识别内存中数据
- `RecognizeAudio(source AudioSource, options *RecognitionOptions) (*RecognitionResult, error)` - 识别音频源
- `RecognizeFileAsync(filePath string, options *RecognitionOptions, resultCallback func(*RecognitionResult, error))` - 异步识别文件
- `RecognizeBytesAsync(data []byte, fileName string, options *RecognitionOptions, resultCallback func(*RecognitionResult, error))` - 异步识别内存中数据
- `GetResult(taskID string) (*RecognitionResult, error)` - 获取结果
- `GetResultWithCallback(taskID string, statusCallback StatusCallback, resultCallback ResultCallback, options *RecognitionOptions) error` - 获取结果并使用回调
- `GetResultWithTimeout(taskID string, timeout time.Duration) (*RecognitionResult, error)` - 带超时获取结果
- `CreateTask(fileURL string, options *RecognitionOptions) (string, error)` - 创建任务

### AudioSource

音频源接口，定义了如何获取音频数据。

```go
type AudioSource interface {
    Read() ([]byte, error)  // 读取音频数据
    Size() (int64, error)   // 返回音频数据大小
    Name() string           // 返回音频名称
}
```

#### 实现

- `FileAudioSource` - 从文件读取音频数据
- `MemoryAudioSource` - 从内存读取音频数据

#### 工厂函数

- `NewFileAudioSource(filePath string) AudioSource` - 创建文件音频源
- `NewMemoryAudioSource(data []byte, fileName string) AudioSource` - 创建内存音频源

### RecognitionOptions

配置识别选项。

#### 字段

- `Language string` - 识别语言，如 "zh_cn"
- `Accent string` - 方言，如 "mandarin"
- `Domain string` - 领域，如 "pro_ost_ed"
- `Encoding string` - 编码方式，如 "raw"
- `Format string` - 音频格式，可选
- `CallbackURL string` - 回调URL，可选
- `StatusCallback StatusCallback` - 状态回调函数
- `ResultCallback ResultCallback` - 结果回调函数
- `Timeout int` - 超时时间（秒）
- `MaxRetries int` - 最大重试次数
- `RetryInterval int` - 重试间隔（秒）

#### 链式设置方法

- `WithLanguage(language string) *RecognitionOptions`
- `WithAccent(accent string) *RecognitionOptions`
- `WithDomain(domain string) *RecognitionOptions`
- `WithEncoding(encoding string) *RecognitionOptions`
- `WithFormat(format string) *RecognitionOptions`
- `WithCallbackURL(callbackURL string) *RecognitionOptions`
- `WithStatusCallback(callback StatusCallback) *RecognitionOptions`
- `WithResultCallback(callback ResultCallback) *RecognitionOptions`
- `WithTimeout(timeout int) *RecognitionOptions`
- `WithMaxRetries(maxRetries int) *RecognitionOptions`
- `WithRetryInterval(retryInterval int) *RecognitionOptions`

### RecognitionResult

识别结果结构。

```go
type RecognitionResult struct {
    Code int
    Data struct {
        TaskID     string
        TaskStatus string
        ErrorCode  string
        ErrorDesc  string
        Result     ResultData
    }
    Message string
    Sid     string
}

type ResultData struct {
    FileLength int
    Sentences  []Sentence
    FullText   string  // 解析后的完整文本
}

type Sentence struct {
    Begin     string
    End       string
    Words     []Word
    SpeakerID string
}

type Word struct {
    Text       string
    Confidence float64
    WordPos    int
    WordSeg    string
    Role       string
}
```

## 常见问题

### 1. 支持哪些音频格式？

SDK支持多种音频格式，包括WAV、MP3、PCM等。建议使用采样率为16kHz的WAV格式获得最佳识别效果。

### 2. 如何处理长音频文件？

对于大于30MB的音频文件，SDK会自动使用分片上传功能，无需开发者手动处理。

### 3. 如何处理识别超时？

可以使用`WithTimeout`方法设置超时时间，或者使用`GetResultWithTimeout`方法指定获取结果的超时时间。

### 4. 如何获取每个词的置信度？

可以通过`result.Data.Result.Sentences`数组中的`Words`字段获取每个词的置信度信息。

## 错误处理

SDK使用标准的Go错误处理机制，所有可能失败的操作都会返回错误。错误信息包含详细的失败原因，方便排查问题。

```go
result, err := client.RecognizeFile("audio_sample.wav", nil)
if err != nil {
    if strings.Contains(err.Error(), "上传失败") {
        // 处理上传错误
    } else if strings.Contains(err.Error(), "创建任务失败") {
        // 处理任务创建错误
    } else {
        // 处理其他错误
    }
}
```

常见错误码及其含义：

| 错误码 | 含义 |
|-------|------|
| 10101 | 参数不合法 |
| 10102 | 系统繁忙 |
| 10103 | 未开通服务权限 |
| 10104 | 未授权或IP不在白名单内 |
| 10105 | 授权过期 |
| 10106 | 签名错误 |
| 10107 | 应用ID不存在 |
| 10108 | 语音文件为空或过大 |
| 10109 | 语音文件解码失败 |
| 10110 | 语音文件无有效数据 |

## 最佳实践

### 性能优化

1. **复用客户端**：创建Client实例是一个相对昂贵的操作，应尽量复用同一个客户端实例。

```go
// 创建一个全局客户端
var client = ost.NewClient("your_app_id", "your_api_key", "your_api_secret")

func handler1() {
    // 使用全局客户端
    client.RecognizeFile(...)
}

func handler2() {
    // 使用同一个客户端
    client.RecognizeFile(...)
}
```

2. **异步处理**：对于Web应用，使用异步处理避免阻塞请求处理线程。

```go
func handleUpload(w http.ResponseWriter, r *http.Request) {
    // 解析上传的文件...
    
    // 异步处理识别
    go func() {
        result, err := client.RecognizeFile(filePath, nil)
        if err != nil {
            log.Printf("识别失败: %v", err)
            return
        }
        
        // 处理结果...
    }()
    
    // 立即返回响应
    w.Write([]byte("文件已接收，正在处理..."))
}
```

3. **设置合理的超时时间**：根据音频长度设置合理的超时时间。

```go
// 对于短音频（<1分钟）
options := ost.CreateOptions().WithTimeout(30)

// 对于中等长度音频（1-5分钟）
options := ost.CreateOptions().WithTimeout(60)

// 对于长音频（>5分钟）
options := ost.CreateOptions().WithTimeout(300)
```

### 错误处理最佳实践

实现指数退避重试机制：

```go
func recognizeWithRetry(client *ost.Client, filePath string, maxRetries int) (*ost.RecognitionResult, error) {
    var result *ost.RecognitionResult
    var err error
    
    for i := 0; i < maxRetries; i++ {
        result, err = client.RecognizeFile(filePath, nil)
        if err == nil {
            return result, nil
        }
        
        // 检查是否是临时错误
        if strings.Contains(err.Error(), "系统繁忙") || strings.Contains(err.Error(), "网络错误") {
            // 指数退避
            waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
            time.Sleep(waitTime)
            continue
        }
        
        // 非临时错误，立即返回
        return nil, err
    }
    
    return nil, fmt.Errorf("达到最大重试次数(%d): %v", maxRetries, err)
}
```

## 示例代码

### 完整示例：处理多个音频文件

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sync"
    
    "github.com/yourusername/xfspeech"
)

func main() {
    client := ost.NewClient("your_app_id", "your_api_key", "your_api_secret")
    
    // 获取音频文件列表
    audioFiles, err := filepath.Glob("audio/*.wav")
    if err != nil {
        log.Fatalf("查找音频文件失败: %v", err)
    }
    
    // 创建结果通道
    type Result struct {
        FileName string
        Text     string
        Error    error
    }
    resultChan := make(chan Result, len(audioFiles))
    
    // 使用WaitGroup等待所有识别完成
    var wg sync.WaitGroup
    
    // 同时处理的最大文件数
    maxConcurrent := 5
    semaphore := make(chan struct{}, maxConcurrent)
    
    for _, file := range audioFiles {
        wg.Add(1)
        
        // 获取信号量
        semaphore <- struct{}{}
        
        go func(file string) {
            defer wg.Done()
            defer func() { <-semaphore }()
            
            baseName := filepath.Base(file)
            
            // 设置选项
            options := ost.CreateOptions().
                WithStatusCallback(func(status string, _ *ost.RecognitionResult) {
                    log.Printf("[%s] 状态: %s", baseName, ost.StatusToString(status))
                })
            
            // 执行识别
            result, err := client.RecognizeFile(file, options)
            
            if err != nil {
                resultChan <- Result{FileName: baseName, Error: err}
                return
            }
            
            resultChan <- Result{
                FileName: baseName,
                Text:     result.Data.Result.FullText,
            }
        }(file)
    }
    
    // 等待所有识别完成
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    // 处理结果
    outputFile, err := os.Create("recognition_results.txt")
    if err != nil {
        log.Fatalf("创建输出文件失败: %v", err)
    }
    defer outputFile.Close()
    
    for result := range resultChan {
        if result.Error != nil {
            fmt.Fprintf(outputFile, "文件 %s 识别失败: %v\n\n", result.FileName, result.Error)
            continue
        }
        
        fmt.Fprintf(outputFile, "文件: %s\n文本: %s\n\n", result.FileName, result.Text)
    }
    
    fmt.Println("所有文件处理完成，结果已保存到 recognition_results.txt")
}
```

### 示例：会议记录转写

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/yourusername/xfspeech"
)

func main() {
    client := ost.NewClient("your_app_id", "your_api_key", "your_api_secret")
    
    // 设置选项，使用讯飞的会议场景优化
    options := ost.CreateOptions().
        WithDomain("pro_ost_ed").         // 增强版
        WithAccent("mandarin").           // 普通话
        WithFormat(ost.FormatWAV).   // WAV格式
        WithMaxRetries(150).              // 增加重试次数，适合长音频
        WithStatusCallback(func(status string, _ *ost.RecognitionResult) {
            fmt.Printf("当前状态: %s\n", ost.StatusToString(status))
        })
    
    // 开始时间
    startTime := time.Now()
    
    // 执行识别
    result, err := client.RecognizeFile("meeting_record.wav", options)
    if err != nil {
        log.Fatalf("识别失败: %v", err)
    }
    
    // 计算耗时
    duration := time.Since(startTime)
    
    // 保存结果
    outputFile, err := os.Create("meeting_transcript.txt")
    if err != nil {
        log.Fatalf("创建输出文件失败: %v", err)
    }
    defer outputFile.Close()
    
    // 写入标题
    fmt.Fprintf(outputFile, "会议记录转写\n")
    fmt.Fprintf(outputFile, "转写时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Fprintf(outputFile, "处理耗时: %s\n\n", duration)
    
    // 写入完整文本
    fmt.Fprintf(outputFile, "== 完整文本 ==\n\n")
    fmt.Fprintf(outputFile, "%s\n\n", result.Data.Result.FullText)
    
    // 按句子输出详细信息
    fmt.Fprintf(outputFile, "== 详细信息 ==\n\n")
    for i, sentence := range result.Data.Result.Sentences {
        // 计算开始和结束时间（毫秒转换为时:分:秒格式）
        beginMs, _ := time.ParseDuration(sentence.Begin + "ms")
        endMs, _ := time.ParseDuration(sentence.End + "ms")
        
        beginStr := fmt.Sprintf("%02d:%02d:%02d", 
            int(beginMs.Hours()),
            int(beginMs.Minutes()) % 60,
            int(beginMs.Seconds()) % 60,
        )
        
        endStr := fmt.Sprintf("%02d:%02d:%02d", 
            int(endMs.Hours()),
            int(endMs.Minutes()) % 60,
            int(endMs.Seconds()) % 60,
        )
        
        // 输出句子信息
        fmt.Fprintf(outputFile, "[%d] %s - %s", i+1, beginStr, endStr)
        
        if sentence.SpeakerID != "" {
            fmt.Fprintf(outputFile, " (说话人: %s)", sentence.SpeakerID)
        }
        
        fmt.Fprintf(outputFile, "\n")
        
        // 输出文本
        var sentenceText string
        for _, word := range sentence.Words {
            sentenceText += word.Text
        }
        
        fmt.Fprintf(outputFile, "%s\n\n", sentenceText)
    }
    
    fmt.Println("会议记录转写完成，结果已保存到 meeting_transcript.txt")
}
```