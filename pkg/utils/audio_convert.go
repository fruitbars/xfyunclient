package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"os"
)

func ConvertRawToWav(rawFile, wavFile string, sampleRate, bitsPerSample, channels int) error {
	rawData, err := os.ReadFile(rawFile)
	if err != nil {
		return fmt.Errorf("failed to read raw file: %w", err)
	}

	outFile, err := os.Create(wavFile)
	if err != nil {
		return fmt.Errorf("failed to create wav file: %w", err)
	}
	defer outFile.Close()

	encoder := wav.NewEncoder(outFile, sampleRate, bitsPerSample, channels, 1) // 1 = PCM format
	buf := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  sampleRate,
			NumChannels: channels,
		},
		Data:           make([]int, len(rawData)/2),
		SourceBitDepth: bitsPerSample,
	}

	for i := 0; i < len(rawData); i += 2 {
		buf.Data[i/2] = int(int16(rawData[i]) | int16(rawData[i+1])<<8)
	}

	if err := encoder.Write(buf); err != nil {
		return fmt.Errorf("failed to write wav data: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close wav encoder: %w", err)
	}

	return nil
}

func PcmToWav(pcmFile, wavFile string, sampleRate, channels, bitsPerSample int) error {
	pcmData, err := os.ReadFile(pcmFile)
	if err != nil {
		return err
	}

	wavFileHandle, err := os.Create(wavFile)
	if err != nil {
		return err
	}
	defer wavFileHandle.Close()

	// WAV 头部
	var header bytes.Buffer
	byteRate := sampleRate * channels * bitsPerSample / 8
	blockAlign := channels * bitsPerSample / 8
	dataLen := uint32(len(pcmData))

	// 写 RIFF 头
	header.WriteString("RIFF")
	binary.Write(&header, binary.LittleEndian, uint32(36+dataLen)) // 文件总长
	header.WriteString("WAVE")

	// fmt 子块
	header.WriteString("fmt ")
	binary.Write(&header, binary.LittleEndian, uint32(16))            // 子块大小
	binary.Write(&header, binary.LittleEndian, uint16(1))             // 音频格式 PCM = 1
	binary.Write(&header, binary.LittleEndian, uint16(channels))      // 声道数
	binary.Write(&header, binary.LittleEndian, uint32(sampleRate))    // 采样率
	binary.Write(&header, binary.LittleEndian, uint32(byteRate))      // 字节率
	binary.Write(&header, binary.LittleEndian, uint16(blockAlign))    // 块对齐
	binary.Write(&header, binary.LittleEndian, uint16(bitsPerSample)) // 位深

	// data 子块
	header.WriteString("data")
	binary.Write(&header, binary.LittleEndian, dataLen)

	// 写入 WAV 文件
	wavFileHandle.Write(header.Bytes())
	wavFileHandle.Write(pcmData)

	return nil
}
