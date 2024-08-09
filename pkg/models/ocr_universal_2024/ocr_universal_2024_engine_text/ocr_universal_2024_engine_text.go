package ocr_universal_2024_engine_text

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Word struct {
	Conf    float64      `json:"conf"`
	Content string       `json:"content"`
	Coord   []Coordinate `json:"coord"`
}

type Line struct {
	Angle     float64      `json:"angle"`
	Conf      float64      `json:"conf"`
	Content   string       `json:"content"`
	Coord     []Coordinate `json:"coord"`
	Exception int          `json:"exception"`
	ID        int          `json:"id"`
	Type      string       `json:"type"`
	Words     []Word       `json:"words"`
}

type Block struct {
	Coord   []Coordinate `json:"coord"`
	ID      int          `json:"id"`
	LineIDs []int        `json:"line_ids"`
}

type Graph struct {
	Coord []Coordinate `json:"coord"`
	ID    int          `json:"id"`
}

type PageNumberElement struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type PageNumber struct {
	Coord    []Coordinate        `json:"coord"`
	Elements []PageNumberElement `json:"elements"`
	ID       int                 `json:"id"`
}

type Page struct {
	Angle       float64      `json:"angle"`
	Blocks      []Block      `json:"blocks"`
	Exception   int          `json:"exception"`
	Graphs      []Graph      `json:"graphs"`
	Height      int          `json:"height"`
	Lines       []Line       `json:"lines"`
	PageNumbers []PageNumber `json:"page_numbers"`
	Width       int          `json:"width"`
}

type OCREngineText struct {
	Category string `json:"category"`
	Pages    []Page `json:"pages"`
	Protoc   string `json:"protoc"`
	Version  string `json:"version"`
}
