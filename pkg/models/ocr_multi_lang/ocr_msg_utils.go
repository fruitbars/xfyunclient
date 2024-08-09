package ocr_multi_lang

import (
	"aseclient/pkg/models/ocr_universal_2024/ocr_universal_2024_engine_text"
)

// GetBlockContents 返回每个Block合并后的文本内容
func GetBlockContents(data *OCRResponseText) map[int]string {
	blockContents := make(map[int]string)

	for _, page := range data.Pages {
		for _, block := range page.Blocks {
			var content string
			for _, lineID := range block.LineIDs {
				for _, line := range page.Lines {
					if line.ID == lineID {
						content += line.Content
						break
					}
				}
			}
			blockContents[block.ID] = content + "\n"

			//			log.Println(block.ID, ":", block.LineIDs)
		}
	}

	return blockContents
}
func GetUniversal2024BlockContents(data *ocr_universal_2024_engine_text.OCREngineText) map[int]string {
	blockContents := make(map[int]string)

	for _, page := range data.Pages {
		for _, block := range page.Blocks {
			var content string
			for _, lineID := range block.LineIDs {
				for _, line := range page.Lines {
					if line.ID == lineID {
						content += line.Content
						break
					}
				}
			}
			blockContents[block.ID] = content + "\n"
		}
	}

	return blockContents
}
