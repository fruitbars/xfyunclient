# 讯飞 API Go SDK 文档

## 概述

本SDK为科大讯飞语音识别、语音合成、OCR文字识别等服务提供了Go语言客户端实现。SDK遵循现代Go语言设计理念，采用接口驱动设计，提供优雅、易用且高度可扩展的API。

## 目录结构

```
xfyunclient/
├── pkg/
│   ├── ase/                  # 底层API通信模块
│   ├── utils/                # 工具函数
│   ├── v2tts/                # 语音合成模块
│   ├── ocr_layout/           # OCR版面分析模块
│   ├── ocr_universal/        # OCR通用识别模块
│   ├── translation/          # 机器翻译模块
│   ├── image_block/          # 图像块识别模块
│   └── models/               # 请求/响应模型定义
├── examples/                 # 示例代码
└── README.md                 # 项目说明
```

## 安装

```bash
go get github.com/yourname/xfyunclient
```

## 模块说明

### 1. 语音合成 (v2tts)

语音合成模块提供文本到语音的转换功能，支持多种发音人、语速、音高等参数设置。

#### 基本使用

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/v2tts"
)

func main() {
    // 创建客户端
    client := v2tts.NewV2TTSClient(
        "your_app_id",
        "your_ase_app_id",
        "your_api_key",
        "your_api_secret",
        "wss://tts-api.xfyun.cn/v2/tts",
    )
    
    // 设置默认参数
    client.SetDefaultVoice("xiaoyan")  // 默认使用小燕的声音
    
    // 简单文本转语音
    err := client.TextToSpeechToFile(
        "你好，这是一个测试文本。", 
        "output.wav",
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("语音合成完成!")
}
```

#### 高级使用

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/v2tts"
)

func main() {
    // 创建客户端
    client := v2tts.NewV2TTSClient(
        "your_app_id",
        "your_ase_app_id",
        "your_api_key",
        "your_api_secret",
        "wss://tts-api.xfyun.cn/v2/tts",
    )
    
    // 使用选项模式设置参数
    err := client.TextToSpeechToFile(
        "这是一个带有自定义选项的示例。",
        "custom_output.wav",
        v2tts.WithVoice("xiaofeng"),
        v2tts.WithSpeed(80),          // 语速: 0-100
        v2tts.WithVolume(90),         // 音量: 0-100
        v2tts.WithPitch(50),          // 音高: 0-100
        v2tts.WithAue("wav"),         // 音频格式
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("自定义语音合成完成!")
}
```

#### 处理长文本

```go
package main

import (
    "fmt"
    "os"
    "xfyunclient/pkg/v2tts"
)

func main() {
    // 创建客户端
    client := v2tts.NewV2TTSClient(
        "your_app_id",
        "your_ase_app_id",
        "your_api_key",
        "your_api_secret",
        "wss://tts-api.xfyun.cn/v2/tts",
    )
    
    // 长文本
    longText := `这是一个较长的文本示例，用于展示如何处理超出单次请求长度限制的文本。
    语音合成服务通常对单次请求的文本长度有限制，当需要转换的文本超过限制时，
    需要将文本分割成多个部分，分别进行转换，然后将多个音频文件合并成一个完整的文件。`
    
    // 分割长文本
    textParts := v2tts.SplitTextByLength(longText, 300)
    
    // 为每个部分生成音频
    var audioFiles []string
    for i, part := range textParts {
        fileName := fmt.Sprintf("part_%d.wav", i)
        err := client.TextToSpeechToFile(part, fileName)
        if err != nil {
            fmt.Printf("Error processing part %d: %v\n", i, err)
            // 清理已生成的文件
            for _, f := range audioFiles {
                os.Remove(f)
            }
            return
        }
        audioFiles = append(audioFiles, fileName)
    }
    
    // 合并音频文件
    err := v2tts.ConcatAudioFiles("complete.wav", audioFiles...)
    if err != nil {
        fmt.Printf("Error concatenating audio files: %v\n", err)
        return
    }
    
    // 清理临时文件
    for _, f := range audioFiles {
        os.Remove(f)
    }
    
    fmt.Println("长文本语音合成完成!")
}
```

### 2. OCR版面分析 (ocr_layout)

OCR版面分析模块提供图像中文字的版面识别和分析功能。

#### 基本使用

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/ocr_layout"
)

func main() {
    // 创建客户端
    client := ocr_layout.NewOCRLayoutClient(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://api.xfyun.cn/v1/service/v1/ocr/layout",
    )
    
    // 启用调试输出
    client.SetDebug(true)
    
    // 识别图像文件
    text, err := client.ScanFile("path/to/document.jpg")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("OCR识别结果:")
    fmt.Println(text)
}
```

#### 从不同源读取图像

```go
package main

import (
    "fmt"
    "os"
    "xfyunclient/pkg/ocr_layout"
)

func main() {
    // 创建客户端
    client := ocr_layout.NewOCRLayoutClient(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://api.xfyun.cn/v1/service/v1/ocr/layout",
    )
    
    // 方法1: 从文件路径识别
    text1, err := client.ScanFile("path/to/image.jpg")
    if err != nil {
        fmt.Printf("从文件识别失败: %v\n", err)
    } else {
        fmt.Printf("文件识别结果: %s\n", text1)
    }
    
    // 方法2: 从文件对象识别
    file, err := os.Open("path/to/another_image.png")
    if err != nil {
        fmt.Printf("打开文件失败: %v\n", err)
        return
    }
    defer file.Close()
    
    text2, err := client.ScanReader(file)
    if err != nil {
        fmt.Printf("从文件对象识别失败: %v\n", err)
    } else {
        fmt.Printf("文件对象识别结果: %s\n", text2)
    }
    
    // 方法3: 从内存图像数据识别
    imageBytes, err := os.ReadFile("path/to/third_image.jpg")
    if err != nil {
        fmt.Printf("读取文件失败: %v\n", err)
        return
    }
    
    text3, err := client.ScanImage("jpg", imageBytes)
    if err != nil {
        fmt.Printf("从内存数据识别失败: %v\n", err)
    } else {
        fmt.Printf("内存数据识别结果: %s\n", text3)
    }
}
```

### 3. OCR通用识别 (ocr_universal)

OCR通用识别模块提供多语言文字识别功能，支持中英文、日韩文、欧洲语系等多种语言。

#### 基本使用

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/ocr_universal"
)

func main() {
    // 创建客户端
    client := ocr_universal.New(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://api.xfyun.cn/v1/service/v1/ocr/universal",
    )
    
    // 启用调试输出
    client.SetDebug(true)
    
    // 识别图像文件
    result, err := client.RecognizeFile(
        "path/to/image.jpg", 
        string(ocr_universal.LangChineseEnglish),
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("OCR识别结果:")
    fmt.Println(result.DecodedText)
}
```

#### 批量处理

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "xfyunclient/pkg/ocr_universal"
)

func main() {
    // 创建客户端
    client := ocr_universal.New(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://api.xfyun.cn/v1/service/v1/ocr/universal",
    )
    
    // 图像目录
    imageDir := "path/to/images"
    
    // 输出目录
    outputDir := "ocr_results"
    os.MkdirAll(outputDir, 0755)
    
    // 遍历目录
    err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // 跳过目录
        if info.IsDir() {
            return nil
        }
        
        // 检查是否为图像文件
        ext := strings.ToLower(filepath.Ext(path))
        if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".bmp" {
            return nil
        }
        
        fmt.Printf("处理文件: %s\n", path)
        
        // 识别文本
        result, err := client.RecognizeFile(path, string(ocr_universal.LangChineseEnglish))
        if err != nil {
            fmt.Printf("  Error: %v\n", err)
            return nil // 继续处理下一个文件
        }
        
        // 创建输出文件名
        baseName := filepath.Base(path)
        outputPath := filepath.Join(
            outputDir, 
            baseName[:len(baseName)-len(ext)] + ".txt",
        )
        
        // 保存结果
        err = os.WriteFile(outputPath, []byte(result.DecodedText), 0644)
        if err != nil {
            fmt.Printf("  保存结果失败: %v\n", err)
            return nil
        }
        
        fmt.Printf("  结果已保存到: %s\n", outputPath)
        return nil
    })
    
    if err != nil {
        fmt.Printf("遍历目录失败: %v\n", err)
        return
    }
    
    fmt.Println("批量处理完成!")
}
```

### 4. 机器翻译 (translation)

机器翻译模块提供文本翻译功能，支持多种语言间的互译。

#### 基本使用

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/translation"
)

func main() {
    // 创建客户端
    client := translation.New(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://itrans.xfyun.cn/v1/its",
    )
    
    // 简单英译中
    text := "Hello, how are you today?"
    translated, err := client.Translate(text, string(translation.LangEN), string(translation.LangZH))
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("原文: %s\n", text)
    fmt.Printf("译文: %s\n", translated)
}
```

#### 自动检测语言

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/translation"
)

func main() {
    // 创建客户端
    client := translation.New(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://itrans.xfyun.cn/v1/its",
    )
    
    // 多语言文本
    texts := []string{
        "Hello, how are you?",                 // 英语
        "你好，最近怎么样？",                    // 中文
        "こんにちは、お元気ですか？",              // 日语
        "안녕하세요, 어떻게 지내세요?",           // 韩语
        "Bonjour, comment ça va?",             // 法语
    }
    
    // 自动检测并翻译为中文
    for _, text := range texts {
        translated, err := client.TranslateAuto(text, string(translation.LangZH))
        if err != nil {
            fmt.Printf("Error translating '%s': %v\n", text, err)
            continue
        }
        
        fmt.Printf("原文: %s\n", text)
        fmt.Printf("译文: %s\n\n", translated)
    }
}
```

#### 批量翻译

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/translation"
)

func main() {
    // 创建客户端
    client := translation.New(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://itrans.xfyun.cn/v1/its",
    )
    
    // 要翻译的文本列表
    texts := []string{
        "Hello",
        "Good morning",
        "Thank you",
        "Goodbye",
        "How are you?",
    }
    
    // 批量翻译为中文
    results, err := client.BatchTranslate(texts, string(translation.LangEN), string(translation.LangZH))
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // 打印结果
    fmt.Println("批量翻译结果:")
    for i, result := range results {
        fmt.Printf("%d. %s -> %s\n", i+1, texts[i], result)
    }
}
```

### 5. 图像块识别 (image_block)

图像块识别模块提供从图像中识别不同文本块的功能，支持多语言OCR和通用OCR 2024两种模式。

#### 基本使用

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/image_block"
)

func main() {
    // 创建客户端
    client := image_block.NewClient(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://api.xfyun.cn/v1/service/v1/ocr/block",
    )
    
    // 启用调试模式
    client.SetDebug(true)
    
    // 设置结果保存目录
    client.SetOutputDirectory("./ocr_results")
    
    // 识别图像文件中的文本块
    blocks, err := client.RecognizeBlocksFromFile(
        "path/to/document.jpg", 
        string(image_block.LangChineseEnglish),
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("识别到 %d 个文本块\n", len(blocks))
    
    // 打印识别结果
    for blockID, text := range blocks {
        fmt.Printf("Block %s: %s\n", blockID, text)
    }
}
```

#### 使用通用OCR 2024

```go
package main

import (
    "fmt"
    "xfyunclient/pkg/image_block"
)

func main() {
    // 创建客户端
    client := image_block.NewClient(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://api.xfyun.cn/v1/service/v1/ocr/universal_2024",
    )
    
    // 启用通用OCR 2024
    client.UseUniversalOCR(true)
    
    // 设置输出目录
    client.SetOutputDirectory("./ocr_results/universal")
    
    // 识别图像文件中的文本块
    blocks, err := client.RecognizeBlocksFromFile(
        "path/to/document.jpg", 
        string(image_block.LangUniversalChineseEnglish),
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("识别到 %d 个文本块\n", len(blocks))
    
    // 打印识别结果
    for blockID, text := range blocks {
        fmt.Printf("Block %s: %s\n", blockID, text)
    }
}
```

#### 批量处理目录

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "xfyunclient/pkg/image_block"
)

func main() {
    // 创建客户端
    client := image_block.NewClient(
        "your_app_id",
        "your_api_key",
        "your_api_secret",
        "https://api.xfyun.cn/v1/service/v1/ocr/block",
    )
    
    // 设置输出目录
    client.SetOutputDirectory("./ocr_results/batch")
    
    // 处理目录中的所有图像文件
    dirPath := "path/to/images"
    fmt.Printf("处理目录: %s\n", dirPath)
    
    results, err := client.ProcessDirectory(dirPath, string(image_block.LangChineseEnglish))
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("已处理 %d 个文件\n", len(results))
    
    // 将结果保存到文本文件
    outputDir := "./ocr_results/text"
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        fmt.Printf("Error creating output directory: %v\n", err)
        return
    }
    
    for fileName, blocks := range results {
        // 创建输出文件名
        baseName := filepath.Base(fileName)
        outputName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".txt"
        outputPath := filepath.Join(outputDir, outputName)
        
        // 打开输出文件
        file, err := os.Create(outputPath)
        if err != nil {
            fmt.Printf("Error creating output file %s: %v\n", outputPath, err)
            continue
        }
        
        // 写入文件
        for blockID, text := range blocks {
            fmt.Fprintf(file, "Block %s:\n%s\n\n", blockID, text)
        }
        
        file.Close()
        fmt.Printf("结果已保存到: %s\n", outputPath)
    }
    
    fmt.Println("批量处理完成")
}
```

## 常见问题

### 1. 如何处理超过限制的文本？

对于语音合成，超过长度限制的文本可以使用`SplitTextByLength`函数分割后分别处理：

```go
// 分割长文本，最大长度300字符
textParts := v2tts.SplitTextByLength(longText, 300)

// 处理每个部分
for _, part := range textParts {
    // 处理每个文本片段
}
```

### 2. 如何设置自定义参数？

所有模块都支持选项模式设置参数，例如：

```go
// 语音合成参数
client.TextToSpeechToFile(
    text,
    outputFile,
    v2tts.WithVoice("xiaofeng"),
    v2tts.WithSpeed(80),
    v2tts.WithVolume(90),
)

// OCR参数
client.RecognizeImage(
    format,
    imageData,
    language,
    ocr_universal.WithDetectOrientation(true),
    ocr_universal.WithOutputFormat("json"),
)
```

### 3. 如何处理错误？

所有API都返回明确的错误信息，建议使用标准的错误处理方式：

```go
result, err := client.SomeFunction()
if err != nil {
    // 检查是否为特定类型的错误
    if ocr_universal.IsNetworkError(err) {
        // 处理网络错误
    } else if ocr_universal.IsServerError(err) {
        // 处理服务器错误
    } else {
        // 处理其他错误
    }
    return
}

// 使用结果
```

## 总结

讯飞API Go SDK提供了多种语音和文字处理服务的客户端实现，包括语音合成、OCR识别、机器翻译等功能。SDK设计遵循现代Go语言最佳实践，提供了易用、灵活且功能强大的API。

通过使用这些模块，开发者可以轻松地将讯飞的人工智能能力集成到自己的应用中，实现文本到语音的转换、图像文字识别、多语言翻译等功能。