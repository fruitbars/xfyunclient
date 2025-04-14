package ocr_layer

type ASEOCRLayerEngineResult struct {
	Category string `json:"category"`
	Pages    []struct {
		Angle  float64 `json:"angle"`
		Blocks []struct {
			Coord []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"coord"`
			ID int `json:"id"`
		} `json:"blocks"`
		Columns []struct {
			Coord []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"coord"`
			Elements []struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			} `json:"elements"`
			ID int `json:"id"`
		} `json:"columns"`
		Exception int `json:"exception"`
		Footers   []struct {
			Coord []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"coord"`
			ID int `json:"id"`
		} `json:"footers"`
		Graphs []struct {
			Coord []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"coord"`
			ID       int `json:"id"`
			Elements []struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			} `json:"elements,omitempty"`
		} `json:"graphs"`
		Headers []struct {
			Coord []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"coord"`
			Elements []struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			} `json:"elements"`
			ID int `json:"id"`
		} `json:"headers"`
		Height      int `json:"height"`
		PageNumbers []struct {
			Coord []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"coord"`
			Elements []struct {
				ID   int    `json:"id"`
				Type string `json:"type"`
			} `json:"elements"`
			ID int `json:"id"`
		} `json:"page_numbers"`
		Width int `json:"width"`
	} `json:"pages"`
	Protoc  string `json:"protoc"`
	Version string `json:"version"`
}
