package imaged

import (
	"io"
	"net/http"
	"os"
	"strconv"
)

func (q *Qcow2Imaged) Download() error {
	q.mu.Lock()
	q.Stage = "Downloading"
	q.mu.Unlock()
	request, err := http.NewRequest("GET", q.link, nil)
	if err != nil {
		return err
	}
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	length := int64(-1)
	if lenStr := resp.Header.Get("Content-Length"); lenStr != "" {
		length, _ = strconv.ParseInt(lenStr, 10, 32)
	}
	rangeEffective := resp.Header.Get("Content-Range") != ""

	var out io.Writer
	var outFile *os.File
	flags := os.O_WRONLY | os.O_CREATE
	if rangeEffective {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	outFile, err = os.OpenFile(q.Get7zFileLocation(), flags, 0660)
	if err != nil {
		return err
	}
	defer outFile.Close()
	out = outFile

	buf := make([]byte, 4068)
	tot := int64(0)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			tot += int64(n)
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
			if length > 0 {
				q.mu.Lock()
				q.Percentage = float64(100*tot) / float64(length)
				q.mu.Unlock()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	q.mu.Lock()
	q.Percentage = 100
	q.mu.Unlock()

	return nil
}
