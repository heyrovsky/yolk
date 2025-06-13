package imaged

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/heyrovsky/yolk/common/config"
)

func (q *Qcow2Imaged) Status() ([]byte, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return json.Marshal(q)
}

func (q *Qcow2Imaged) GetBaseDirectory() string {
	return filepath.Join(config.BASE_ARCHIVE_LOCATION, q.name)
}

func (q *Qcow2Imaged) Get7zFileLocation() string {
	return filepath.Join(config.BASE_ARCHIVE_LOCATION, q.name, fmt.Sprintf("%s.7z", q.version))
}

func (q *Qcow2Imaged) GetVmdkFileLocation() string {
	return filepath.Join(config.BASE_ARCHIVE_LOCATION, q.name, fmt.Sprintf("%s.vmdk", q.version))
}

func (q *Qcow2Imaged) GetQcow2FileLocation() string {
	return filepath.Join(config.BASE_IMAGE_LOCATION, q.name, fmt.Sprintf("%s.qcow2", q.version))
}
