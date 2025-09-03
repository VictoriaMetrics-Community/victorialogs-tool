package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// TailLogs is a function that tails logs from the victoriametrics database
func TailLogs() error {
	cfg := getCfgByToml()

	// replace url from xxx/select/logsql/query to xxx/select/logsql/tail
	cfg.URL = strings.Replace(cfg.URL, "query", "tail", 1)

	// send request to victoria
	param, err := buildParams(cfg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewBufferString(param))
	if err != nil {
		return err
	}

	// set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// maybe we can set timeout to cfg
	cli := &http.Client{Timeout: 20 * time.Minute}

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// copy the response body to stdout
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	return nil
}
