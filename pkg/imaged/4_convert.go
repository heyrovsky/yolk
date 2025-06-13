package imaged

import (
	"bufio"
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
	cmd := exec.Command("qemu-img", "convert", "-p", "-f", "vmdk", "-O", "qcow2", q.Get7zFileLocation(), q.GetQcow2FileLocation())

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
		scanner := bufio.NewScanner(stderr)
		re := regexp.MustCompile(`(\d+)%`)
		for scanner.Scan() {
			line := scanner.Text()
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				var percent float64
				if _, err := fmt.Sscanf(matches[1], "%f", &percent); err == nil {
					q.mu.Lock()
					q.Percentage = percent
					q.mu.Unlock()
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Scanner error:", err)
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
