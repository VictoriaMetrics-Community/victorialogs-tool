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

	// IgnoreOriQuery Whether to ignore the original statement, default is false
	IgnoreOriQuery bool `toml:"ignore_ori_query" json:"ignore_ori_query"`

	Stream      map[string]any `toml:"stream" json:"_stream"`
	Fileds      []string       `toml:"fileds" json:"fileds"`
	CustomPipes []string       `toml:"custom_pipes" json:"custom_pipes"`
}

type SortType string

const (
	SortTypeAsc  SortType = "asc"
	SortTypeDesc SortType = "desc"
)

type RequestParams struct {
	// 查询语句
	Query string `json:"query"`
	Limit int64  `json:"limit"`
	Start string `json:"start"`
	End   string `json:"end"`
}
