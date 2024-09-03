package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
	"unsafe"

	"github.com/VictoriaMetrics-Community/victorialogs-tool/cfgs"
)

// QueryLogs is a function that queries logs from the victoriametrics database
func QueryLogs() ([]string, error) {
	cfg := getCfgByToml()

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
