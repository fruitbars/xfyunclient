package ocr_multi_lang

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
	Words     []Word       `json:"words"`
}

type Block struct {
	Coord   []Coordinate `json:"coord"`
	ID      int          `json:"id"`
	LineIDs []int        `json:"line_ids"`
}

type Page struct {
	Angle     float64 `json:"angle"`
	Blocks    []Block `json:"blocks"`
	Exception int     `json:"exception"`
	Height    int     `json:"height"`
	Lines     []Line  `json:"lines"`
	Width     int     `json:"width"`
}

type OCRResponseText struct {
	Category string `json:"category"`
	Pages    []Page `json:"pages"`
	Protoc   string `json:"protoc"`
	Version  string `json:"version"`
}
