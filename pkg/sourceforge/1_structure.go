package sourceforge

import "go.uber.org/zap"

type SourceForgeVersionInfo struct {
	Version  string
	Archive  string // 64bit.7z
	Checksum string // SHA256.txt
}

type SourceForgeReader struct {
	BaseURL  string
	Path     string
	Logger   *zap.Logger
	Versions map[string]*SourceForgeVersionInfo
}

// #########################
// # Constructor Functions #
// #########################

// NewSourceForgeReader creates a SourceForgeReader, initializes version data by parsing the RSS feed,
// and returns the reader along with any error encountered during initialization.
func NewSourceForgeReader(baseURL, path string, logger *zap.Logger) (*SourceForgeReader, error) {
	r := &SourceForgeReader{
		BaseURL:  baseURL,
		Path:     path,
		Logger:   logger,
		Versions: make(map[string]*SourceForgeVersionInfo),
	}

	if err := r.ReadVersions(); err != nil {
		return nil, err
	}
	return r, nil
}
