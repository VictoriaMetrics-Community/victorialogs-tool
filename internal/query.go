package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/here-Leslie-Lau/victorialogs-tool/cfgs"
)

const baseJSON = "../cfgs/base.json"

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
	req, err := http.NewRequest("POST", cfg.URL, bytes.NewBuffer(buildParams(cfg)))
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cli := &http.Client{Timeout: 30 * time.Second}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func buildParams(cfg *cfgs.Config) []byte {
	streamStr, _ := json.Marshal(cfg.Stream)

	params := cfgs.RequestParams{
		Limit: cfg.Limit,
		Start: cfg.Start,
		End:   cfg.End,
	}
	// build query params
	params.Query += "_time:" + cfg.LastDuration
	params.Query += " " + cfg.Query
	params.Query += " topic:" + cfg.Topic
	params.Query += " caller:" + cfg.Caller
	params.Query += " _stream:" + string(streamStr)
	params.Query += " level:" + cfg.Level

	params.Query += " | fields " + strings.Join(cfg.Fileds, ",")
	params.Query += " | sort by (_time) desc"

	fmt.Println("Query: ", params.Query)

	byt, _ := json.Marshal(params)
	return byt
}
