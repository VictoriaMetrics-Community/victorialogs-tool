package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/BurntSushi/toml"
	"github.com/here-Leslie-Lau/victorialogs-tool/cfgs"
)

// QueryLogs is a function that queries logs from the victoriametrics database
func QueryLogs() ([]string, error) {
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
	// batch send request by start and end
	configs := splitRequest(cfg)
	wg := sync.WaitGroup{}
	var list []string

	for _, config := range configs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := reqToVictoria(config)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			list = append(list, res)
		}()
	}
	wg.Wait()

	// order the result by num
	if cfg.Sort == cfgs.SortTypeDesc {
		// desc
		sort.Slice(list, func(i, j int) bool {
			return configs[i].Num > configs[j].Num
		})
	} else {
		// default is asc
		sort.Slice(list, func(i, j int) bool {
			return configs[i].Num < configs[j].Num
		})

	}

	return list, nil
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

func reqToVictoria(cfg *cfgs.Config) (string, error) {
	req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewBufferString(buildParams(cfg)))
	if err != nil {
		return "", err
	}

	// set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// maybe we can set timeout to cfg
	cli := &http.Client{Timeout: 30 * time.Second}

	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	str := *(*string)(unsafe.Pointer(&b))
	return str, nil
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

	query += " | fields " + strings.Join(cfg.Fileds, ",")
	if cfg.Sort == cfgs.SortTypeDesc {
		query += " | sort by (_time) desc"
	} else {
		query += " | sort by (_time) asc"
	}

	fmt.Println("Query:", query)
	// url encode
	query = strings.ReplaceAll(url.QueryEscape(query), "+", "%20")

	// time format to ts
	start, _ := time.Parse(time.RFC3339, cfg.Start)
	end, _ := time.Parse(time.RFC3339, cfg.End)

	// manually construct the form
	return fmt.Sprintf("query=%s&limit=%d&start=%d&end=%d", query, cfg.Limit, start.Unix(), end.Unix())
}
