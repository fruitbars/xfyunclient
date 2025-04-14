package ost

// models.go - 数据结构定义
// ====================================

// 请求和响应结构定义
type CreateTaskRequest struct {
	Common struct {
		AppID string `json:"app_id"`
	} `json:"common"`
	Business struct {
		Language    string `json:"language"`
		Domain      string `json:"domain"`
		Accent      string `json:"accent"`
		CallbackURL string `json:"callback_url,omitempty"`
		// 其他可选参数
		VsppOn         string `json:"vspp_on,omitempty"`
		SpeakerNum     string `json:"speaker_num,omitempty"`
		OutputType     string `json:"output_type,omitempty"`
		PostprocOn     string `json:"postproc_on,omitempty"`
		Pd             string `json:"pd,omitempty"`
		Duration       string `json:"duration,omitempty"`
		EnableSubtitle string `json:"enable_subtitle,omitempty"`
		Smoothproc     string `json:"smoothproc,omitempty"`
		Colloqproc     string `json:"colloqproc,omitempty"`
		LanguageType   string `json:"language_type,omitempty"`
		Vto            string `json:"vto,omitempty"`
		Dhw            string `json:"dhw,omitempty"`
	} `json:"business"`
	Data struct {
		AudioURL  string `json:"audio_url"`
		AudioSrc  string `json:"audio_src"`
		AudioSize string `json:"audio_size,omitempty"`
		Format    string `json:"format,omitempty"`
		Encoding  string `json:"encoding"`
	} `json:"data"`
}

type CreateTaskResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Data    struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

type QueryTaskRequest struct {
	Common struct {
		AppID string `json:"app_id"`
	} `json:"common"`
	Business struct {
		TaskID string `json:"task_id"`
	} `json:"business"`
}

type Word struct {
	Text         string  `json:"w"`
	Confidence   float64 `json:"sc,omitempty"`
	WordProperty string  `json:"wp,omitempty"`
	WordSeg      string  `json:"wordseg,omitempty"`
	Role         string  `json:"rl,omitempty"`
}

type Sentence struct {
	Begin     string `json:"begin"`
	End       string `json:"end"`
	Words     []Word `json:"words"`
	SpeakerID string `json:"spk,omitempty"`
}

type UserFriendlyRecognitionResult struct {
	Code int `json:"code"`
	Data struct {
		TaskID     string     `json:"task_id"`
		TaskStatus string     `json:"task_status"`
		ErrorCode  string     `json:"error_code,omitempty"`
		ErrorDesc  string     `json:"error_desc,omitempty"`
		Result     ResultData `json:"result,omitempty"`
	} `json:"data"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
}

type ResultData struct {
	FileLength int        `json:"file_length"`
	Sentences  []Sentence `json:"sentences,omitempty"`
	FullText   string     `json:"full_text,omitempty"` // 解析后的完整文本
	Sentences2 []Sentence `json:"sentences2,omitempty"`
	FullText2  string     `json:"full_text2,omitempty"` // 解析后的完整文本
}

// 上传相关结构
type InitUploadRequest struct {
	RequestID string `json:"request_id"`
	AppID     string `json:"app_id"`
	CloudID   string `json:"cloud_id,omitempty"`
}

type InitUploadResponse struct {
	Code int    `json:"code"`
	Sid  string `json:"sid"`
	Data struct {
		UploadID string `json:"upload_id"`
	} `json:"data"`
	Message string `json:"message"`
}

type UploadCompleteRequest struct {
	RequestID string `json:"request_id"`
	AppID     string `json:"app_id"`
	UploadID  string `json:"upload_id"`
}

type UploadResponse struct {
	Code int    `json:"code"`
	Sid  string `json:"sid"`
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
	Message string `json:"message"`
}

// 音频源接口
type AudioSource interface {
	// Read 读取音频数据
	Read() ([]byte, error)
	// Size 返回音频数据的大小
	Size() (int64, error)
	// Name 返回音频数据的名称
	Name() string
}

// 文件音频源
type FileAudioSource struct {
	FilePath string
}

// 内存音频源
type MemoryAudioSource struct {
	Data     []byte
	FileName string
}

// 回调函数类型
type StatusCallback func(status string, result *RawQueryResponse)
type ResultCallback func(text string, result *RawQueryResponse)

// 识别选项
type RecognitionOptions struct {
	Language       string         // 识别语言，默认 "zh_cn"
	Accent         string         // 方言，默认 "mandarin"
	Domain         string         // 领域，默认 "pro_ost_ed"
	Encoding       string         // 编码方式，默认 "raw"
	Format         string         // 音频格式，可选
	CallbackURL    string         // 回调URL，可选
	StatusCallback StatusCallback // 状态回调函数
	ResultCallback ResultCallback // 结果回调函数
	Timeout        int            // 超时时间（秒）
	MaxRetries     int            // 最大重试次数
	RetryInterval  int            // 重试间隔（秒）
}

type RawLattice struct {
	Begin     string `json:"begin"`
	End       string `json:"end"`
	JSON1Best struct {
		St struct {
			Bg string `json:"bg"`
			Ed string `json:"ed"`
			Pa string `json:"pa"`
			Pt string `json:"pt"`
			Rl string `json:"rl"`
			Rt []struct {
				Nb string `json:"nb"`
				Nc string `json:"nc"`
				Ws []struct {
					Cw []struct {
						W  string  `json:"w"`
						Wc string  `json:"wc,omitempty"`
						Wp string  `json:"wp,omitempty"`
						Sc float64 `json:"sc,omitempty"`
					} `json:"cw"`
					Wb int `json:"wb"`
					We int `json:"we"`
				} `json:"ws"`
			} `json:"rt"`
			Sc string `json:"sc"`
			Si string `json:"si"`
		} `json:"st"`
	} `json:"json_1best"`
	Lid string `json:"lid"`
	Spk string `json:"spk"`
}

// 原始返回结构
type RawQueryResponse struct {
	Code int `json:"code"`
	Data struct {
		ForceRefresh string `json:"force_refresh"`
		Result       struct {
			FileLength int          `json:"file_length"`
			Lattice    []RawLattice `json:"lattice"`
			Lattice2   []RawLattice `json:"lattice2"`
		} `json:"result"`
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
		TaskType   string `json:"task_type"`
	} `json:"data"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
}
