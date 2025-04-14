package backup

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"xfyunclient/pkg/ase"
)

// handleAudioMessageToFile decodes base64 audio from a TTS response and writes it to the provided file.
func handleAudioMessageToFile(p []byte, file *os.File) {
	var resp V2TTSResponse
	if err := json.Unmarshal(p, &resp); err != nil {
		log.Printf("failed to unmarshal TTS response: %v", err)
		return
	}

	if resp.Data.Audio == "" {
		return
	}

	audioBytes, err := base64.StdEncoding.DecodeString(resp.Data.Audio)
	if err != nil {
		log.Printf("failed to decode base64 audio: %v", err)
		return
	}

	if _, err := file.Write(audioBytes); err != nil {
		log.Printf("failed to write audio to file: %v", err)
	}
}

// V2TextToSpeech is the basic version with callback handling
func V2TextToSpeech(appid, aseAppid, aseAPIKey, aseAPISecret, serverUrl, text, aue, vcn, tte string, cb ase.WSCB) (string, error) {
	v2ttsRequest := NewV2TTSRequestDefault(appid, text, aue, vcn, tte)
	client := ase.NewASEWebsocketClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, ase.DefaultASEHttpProto, ase.DefaultASEAlgorithm)
	return client.CallASEAPICallBack(v2ttsRequest, cb)
}

// V2TextToSpeechToFile provides full configuration and writes to file
func V2TextToSpeechToFile(
	aseAppid, aseAPIKey, aseAPISecret, serverUrl, text, aue string,
	sfl int, auf, vcn string, speed, volume, pitch, bgs int,
	tte, reg, rdn, fileName string,
) (string, error) {
	v2ttsRequest := NewV2TTSRequest(aseAppid, text, aue, sfl, auf, vcn, speed, volume, pitch, bgs, tte, reg, rdn)
	client := ase.NewASEWebsocketClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, ase.DefaultASEHttpProto, ase.DefaultASEAlgorithm)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	sid, err := client.CallASEAPICallBack(v2ttsRequest, func(p []byte) {
		handleAudioMessageToFile(p, file)
	})
	if err != nil {
		return sid, fmt.Errorf("TTS request failed: %w", err)
	}

	return sid, nil
}

// V2TextToSpeechToFileDefault uses default config for common use cases
func V2TextToSpeechToFileDefault(
	aseAppid, aseAPIKey, aseAPISecret, serverUrl, text, aue, vcn, tte, fileName string,
) (string, error) {
	v2ttsRequest := NewV2TTSRequestDefault(aseAppid, text, aue, vcn, tte)
	client := ase.NewASEWebsocketClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, ase.DefaultASEHttpProto, ase.DefaultASEAlgorithm)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	sid, err := client.CallASEAPICallBack(v2ttsRequest, func(p []byte) {
		handleAudioMessageToFile(p, file)
	})
	if err != nil {
		return sid, fmt.Errorf("TTS request failed: %w", err)
	}

	return sid, nil
}
