package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/VictoriaMetrics-Community/victorialogs-tool/cfgs"
)

func getCfgByToml() *cfgs.Config {
	// generated configuration file path
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := user.HomeDir
	p := filepath.Join(homeDir, "vtool.json")
	// Read the toml file
	byt, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	var base struct {
		Base string `json:"base"`
	}
	if err := json.Unmarshal(byt, &base); err != nil {
		panic(err)
	}

	cfg := new(cfgs.Config)
	_, err = toml.DecodeFile(base.Base, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func splitRequest(cfg *cfgs.Config) []*cfgs.Config {
	// split the request by start and end
	var res []*cfgs.Config
	start, _ := time.Parse(time.RFC3339, cfg.Start)
	end, _ := time.Parse(time.RFC3339, cfg.End)

	var num int
	// split by 1 hour
	for start.Before(end) {
		// copy the cfg
		c := *cfg
		c.Num = num
		num++
		// case: 不足一小时的情况
		if end.Sub(start) < time.Hour {
			res = append(res, &c)
			break
		}

		c.Start = start.Format(time.RFC3339)
		start = start.Add(time.Hour)
		c.End = start.Format(time.RFC3339)
		res = append(res, &c)
	}
	return res
}

func buildParams(cfg *cfgs.Config) string {
	// build query params
	var query string
	query += "_time:" + cfg.LastDuration
	query += " " + cfg.Query
	query += " topic:" + cfg.Topic
	if cfg.Caller != "" {
		query += " caller:" + cfg.Caller
	}
	query += " _stream:" + "{service=" + `"` + cfg.Stream.Service + `"}`
	query += " level:" + cfg.Level

	if len(cfg.Fileds) > 0 {
		query += " | fields " + strings.Join(cfg.Fileds, ",")
	} else {
		query += " | fields *"
	}

	if cfg.Sort != "" {
		query += " | sort by (_time) " + string(cfg.Sort)
	}

	// build custom pipes
	if len(cfg.CustomPipes) > 0 {
		query += " | " + strings.Join(cfg.CustomPipes, " | ")
	}

	fmt.Println("Query:", query)
	// url encode
	query = strings.ReplaceAll(url.QueryEscape(query), "+", "%20")

	// time format to ts

	// manually construct the form
	res := fmt.Sprintf("query=%s&limit=%d", query, cfg.Limit)

	// default start time is 5 minutes ago
	start := time.Now().Add(-time.Minute * 5)
	if cfg.Start != "" {
		// time format to ts
		start, _ = time.Parse(time.RFC3339, cfg.Start)
	}
	res += fmt.Sprintf("&start=%d", start.Unix())

	if cfg.End != "" {
		// time format to ts
		end, _ := time.Parse(time.RFC3339, cfg.End)
		res += fmt.Sprintf("&end=%d", end.Unix())
	}

	return res
}
