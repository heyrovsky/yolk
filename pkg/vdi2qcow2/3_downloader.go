package vdi2qcow2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/heyrovsky/yolk/common/utils"
)

func (v *VditoQcow2JobStruct) VdiDownloader(ctx context.Context, destPath string, isChecksum bool) error {
	v.updateState("downloader:resolving")
	v.StartTime = time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", func() string {
		if isChecksum {
			return v.ChecksumUrl
		}
		return v.SourceUrl
	}(), nil) // dont ask me why but i just wanted ternary operator
	if err != nil {
		return v.fail("downloader:error", fmt.Errorf("create request: %w", err))
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
		Timeout: 0,
	}

	resp, err := client.Do(req)
	if err != nil {
		return v.fail("downloader:error", fmt.Errorf("request failed: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return v.fail("downloader:error", fmt.Errorf("unexpected status: %s", resp.Status))
	}
	if err := utils.DeleteFileIfExists(destPath); err != nil {
		return v.fail("downloader:error", fmt.Errorf("delete file: %w", err))
	}
	out, err := utils.CreateFileWithDirs(destPath)
	if err != nil {
		return v.fail("downloader:error", fmt.Errorf("create file: %w", err))
	}
	defer out.Close()

	v.updateState("downloader:downloading")

	var (
		totalSize  int64
		downloaded int64
		lastUpdate = time.Now()
		buf        = make([]byte, 64*1024)
	)
	if sizeStr := resp.Header.Get("Content-Length"); sizeStr != "" {
		totalSize, _ = strconv.ParseInt(sizeStr, 10, 64)
	}

	for {
		select {
		case <-ctx.Done():
			return v.fail("downloader:cancelled", ctx.Err())
		default:
			nr, er := resp.Body.Read(buf)
			if nr > 0 {
				nw, ew := out.Write(buf[:nr])
				if ew != nil || nw != nr {
					return v.fail("downloader:error", fmt.Errorf("write failed: %w", ew))
				}
				downloaded += int64(nw)
				if totalSize > 0 && time.Since(lastUpdate) > 500*time.Millisecond {
					v.setProgress(int(downloaded * 100 / totalSize))
					lastUpdate = time.Now()
				}
			}
			if er != nil {
				if er == io.EOF {
					break
				}
				return v.fail("downloader:error", fmt.Errorf("read error: %w", er))
			}
		}
	}

}
