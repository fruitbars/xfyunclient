package ocr_universal

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type WordUnit struct {
	CenterPoint Coordinate   `json:"center_point"`
	Coord       []Coordinate `json:"coord"`
	Conf        float64      `json:"conf"`
	Content     string       `json:"content"`
}

type Word struct {
	Coord   []Coordinate `json:"coord"`
	Conf    float64      `json:"conf"`
	Content string       `json:"content"`
}

type Line struct {
	Exception int          `json:"exception"`
	Coord     []Coordinate `json:"coord"`
	Words     []Word       `json:"words"`
	Angle     int          `json:"angle"`
	Conf      float64      `json:"conf"`
	WordUnits []WordUnit   `json:"word_units"`
}

type Page struct {
	Exception int    `json:"exception"`
	Width     int    `json:"width"`
	Angle     int    `json:"angle"`
	Lines     []Line `json:"lines"`
	Height    int    `json:"height"`
}

type OCRResponseText struct {
	Pages    []Page `json:"pages"`
	Category string `json:"category"`
	Version  string `json:"version"`
}
