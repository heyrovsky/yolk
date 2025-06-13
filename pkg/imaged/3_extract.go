package imaged

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

func (q *Qcow2Imaged) Extract() error {
	q.mu.Lock()
	q.Stage = "Extracting"
	q.Percentage = 0
	q.mu.Unlock()
	fmt.Println(q.Get7zFileLocation())
	r, err := sevenzip.OpenReader(q.Get7zFileLocation())
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		if !strings.HasSuffix(file.Name, ".vmdk") {
			continue
		}

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file inside archive: %w", err)
		}
		defer src.Close()

		dstPath := q.GetVmdkFileLocation()
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}

		dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dst.Close()

		totalSize := int64(file.UncompressedSize)
		var written int64
		buf := make([]byte, 4068)

		for {
			n, err := src.Read(buf)
			if n > 0 {
				written += int64(n)
				if _, err := dst.Write(buf[:n]); err != nil {
					return fmt.Errorf("failed to write to destination file: %w", err)
				}
				if totalSize > 0 {
					q.mu.Lock()
					q.Percentage = (float64(written) / float64(totalSize)) * 100
					q.mu.Unlock()
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error during extraction read: %w", err)
			}
		}

		break // only extract one .vdi file
	}

	q.mu.Lock()
	q.Percentage = 100
	q.mu.Unlock()

	return nil
}
