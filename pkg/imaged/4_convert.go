package imaged

import (
	"fmt"
	"os/exec"
	"regexp"
)

func (q *Qcow2Imaged) Convert() error {
	q.mu.Lock()
	q.Stage = "Converting"
	q.Percentage = 0
	q.mu.Unlock()

	if _, err := exec.LookPath("qemu-img"); err != nil {
		return fmt.Errorf("qemu-img not found in PATH")
	}
	cmd := exec.Command("qemu-img", "convert", "-p", "-f", "vmdk", "-O", "qcow2", q.GetVmdkFileLocation(), q.GetQcow2FileLocation())
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 1024)
		re := regexp.MustCompile(`\((\d+\.\d+)/100%\)`)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				output := string(buf[:n])
				matches := re.FindStringSubmatch(output)
				if len(matches) > 1 {
					var percent float64
					if _, err := fmt.Sscanf(matches[1], "%f", &percent); err == nil {
						q.mu.Lock()
						q.Percentage = percent
						q.mu.Unlock()
					}
				}
			}
			if err != nil {
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		q.mu.Lock()
		q.Stage = "failed"
		q.mu.Unlock()
		return fmt.Errorf("command execution failed: %w", err)
	}
	<-done
	q.mu.Lock()
	q.Percentage = 100
	q.Stage = "done"
	q.mu.Unlock()

	return nil
}
