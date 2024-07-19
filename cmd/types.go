package cmd

type Config struct {
	URL    string `toml:"url"`
	Topic  string `toml:"topic"`
	Caller string `toml:"caller"`

	Start string `toml:"start"`
	End   string `toml:"end"`
	Limit int64  `toml:"limit"`

	Query string `toml:"query"`

	Stream stream   `toml:"stream"`
	Fileds []string `toml:"fileds"`
}

type stream struct {
	Service string `toml:"service"`
}
