package vdi2qcow2

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VditoQcow2JobStruct struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	SourceUrl    string    `json:"url"`
	CurrentState string    `json:"state"`
	Progress     int       `json:"progress"`
	IsCompleted  bool      `json:"completed"`
	ErrorMsg     string    `json:"error,omitempty"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`

	mutex  sync.RWMutex
	logger *zap.Logger
}

func NewVditoQcow2Job(name, version, sourceUrl string, logger *zap.Logger) *VditoQcow2JobStruct {
	return &VditoQcow2JobStruct{
		Id:           uuid.New(),
		Name:         name,
		Version:      version,
		SourceUrl:    sourceUrl,
		CurrentState: "initialized",
		logger:       logger,
	}
}
