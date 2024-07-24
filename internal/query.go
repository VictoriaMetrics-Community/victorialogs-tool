package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/here-Leslie-Lau/victorialogs-tool/cfgs"
)

const baseJSON = "cfgs/base.json"

// QueryLogs is a function that queries logs from the victoriametrics database
func QueryLogs() ([]byte, error) {
	// Read the toml file
	byt, err := os.ReadFile(baseJSON)
	if err != nil {
		return nil, err
	}
	var base struct {
		Base string `json:"base"`
	}
	if err := json.Unmarshal(byt, &base); err != nil {
		return nil, err
	}

	cfg := new(cfgs.Config)
	_, err = toml.DecodeFile(base.Base, cfg)
	if err != nil {
		return nil, err
	}

	// send request to victoria
	b, err := reqToVictoria(cfg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func reqToVictoria(cfg *cfgs.Config) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewBufferString(buildParams(cfg)))
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// maybe we can set timeout to cfg
	cli := &http.Client{Timeout: 30 * time.Second}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func buildParams(cfg *cfgs.Config) string {
	// build query params
	var query string
	query += "_time:" + cfg.LastDuration
	query += " " + cfg.Query
	query += " topic:" + cfg.Topic
	query += " caller:" + cfg.Caller
	query += " _stream:" + "{service=" + `"` + cfg.Stream.Service + `"}`
	query += " level:" + cfg.Level

	query += " | fields " + strings.Join(cfg.Fileds, ",")
	query += " | sort by (_time) desc"

	fmt.Println("Query:", query)
	// url encode
	query = strings.ReplaceAll(url.QueryEscape(query), "+", "%20")

	// time format to ts
	start, _ := time.Parse(time.RFC3339, cfg.Start)
	end, _ := time.Parse(time.RFC3339, cfg.End)

	// manually construct the form
	return fmt.Sprintf("query=%s&limit=%d&start=%d&end=%d", query, cfg.Limit, start.Unix(), end.Unix())
}
