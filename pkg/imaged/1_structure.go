package imaged

import (
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type Qcow2Imaged struct {
	mu         sync.RWMutex
	Jobid      uuid.UUID `json:"id"`
	link       string    `json:"-"`
	name       string    `json:"-"`
	version    string    `json:"-"`
	Percentage float64   `json:"percentage"`
	Stage      string    `json:"stage"`
}

func NewQcow2ImageDaemon(name, version, link string) *Qcow2Imaged {
	return &Qcow2Imaged{
		Jobid:   uuid.New(),
		link:    link,
		name:    name,
		version: version,
	}
}

func (q *Qcow2Imaged) Exec() error {
	if err := os.MkdirAll(q.GetBaseDirectory(), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := q.Download(); err != nil {
		return err
	}

	if err := q.Extract(); err != nil {
		return err
	}

	if err := q.Convert(); err != nil {
		return err
	}

	return nil
}
