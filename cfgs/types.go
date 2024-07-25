package cfgs

type Config struct {
	URL    string `toml:"url" json:"url"`
	Topic  string `toml:"topic" json:"topic"`
	Caller string `toml:"caller" json:"caller"`

	LastDuration string   `toml:"last_duration" json:"last_duration"`
	Start        string   `toml:"start" json:"start"`
	End          string   `toml:"end" json:"end"`
	Limit        int64    `toml:"limit" json:"limit"`
	Sort         SortType `toml:"sort" json:"sort"`

	StartTs int64 `json:"-"`
	EndTs   int64 `json:"-"`
	Num     int   `json:"-"`

	Query string `toml:"query" json:"query"`
	Level string `toml:"level" json:"level"`

	Stream Stream   `toml:"stream" json:"_stream"`
	Fileds []string `toml:"fileds" json:"fileds"`
}

type SortType string

const (
	SortTypeAsc  SortType = "asc"
	SortTypeDesc SortType = "desc"
)

type Stream struct {
	Service string `toml:"service"`
}

type RequestParams struct {
	// 查询语句
	Query string `json:"query"`
	Limit int64  `json:"limit"`
	Start string `json:"start"`
	End   string `json:"end"`
}
