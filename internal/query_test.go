package internal

import (
	"testing"

	"github.com/VictoriaMetrics-Community/victorialogs-tool/cfgs"
)

func TestBuildParams(t *testing.T) {
	cfg := &cfgs.Config{
		URL:          "www.google.com",
		Topic:        "prod*",
		Caller:       "*",
		LastDuration: "5m",
		Start:        "2024-06-30T11:25:13+08:00",
		End:          "2024-06-31T11:25:13+08:00",
		Limit:        10,
		Query:        "_msg:'leslie'",
		Level:        "*",
		Stream: cfgs.Stream{
			Service: "search",
		},
		Fileds: []string{"_msg", "_time", "*"},
	}
	_ = buildParams(cfg)
}
