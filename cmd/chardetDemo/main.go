package main

import (
	"fmt"
	"xfyunclient/pkg/utils"
)

func main() {
	// 假设你从某种来源拿到非 UTF-8 的字符串
	// 这里我们直接用 GB2312 编码的字节作为示例
	//inputBytes := []byte{0xc4, 0xe3, 0xba, 0xc3} // "你好" in GB2312
	inputString := "你好"

	charset, confidence, err := utils.DetectEncoding(inputString)
	if err != nil {
		fmt.Println("检测失败:", err)
		return
	}

	fmt.Printf("检测结果: %s（置信度: %d%%）\n", charset, confidence)
}
