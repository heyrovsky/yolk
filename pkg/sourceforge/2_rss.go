package sourceforge

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

// ReadVersions fetches and parses the SourceForge RSS feed to populate version info including archives and checksums.
func (s *SourceForgeReader) ReadVersions() error {
	feed, err := gofeed.NewParser().ParseURL(fmt.Sprintf("%s/rss?path=%s", strings.TrimSuffix(s.BaseURL, "/"), s.Path))
	if err != nil {
		s.Logger.Error("failed to parse RSS feed", zap.Error(err))
		return err
	}

	s.Versions = make(map[string]*SourceForgeVersionInfo)
	for _, item := range feed.Items {
		if item == nil || item.Link == "" {
			continue
		}

		version, fileType := extractVersionAndType(item.Link, s.Path)
		if version == "" || fileType == "" {
			continue
		}

		if _, ok := s.Versions[version]; !ok {
			s.Versions[version] = &SourceForgeVersionInfo{Version: version}
		}

		switch fileType {
		case "archive":
			s.Versions[version].Archive = item.Link
		case "checksum":
			s.Versions[version].Checksum = item.Link
		}
	}

	return nil
}

// VersionsList returns a slice of all parsed version strings.
func (s *SourceForgeReader) VersionsList() []string {
	keys := make([]string, 0, len(s.Versions))
	for k := range s.Versions {
		keys = append(keys, k)
	}
	return keys
}

// VersionDetails returns the detailed information for a specific version.
// Returns an error if the version does not exist.
func (s *SourceForgeReader) VersionDetails(version string) (*SourceForgeVersionInfo, error) {
	v, ok := s.Versions[version]
	if !ok {
		return nil, fmt.Errorf("version not found: %s", version)
	}
	return v, nil
}

// extractVersionAndType extracts version string and file type ("archive" or "checksum") from a given URL.
func extractVersionAndType(link string, basePath string) (string, string) {

	u, err := url.Parse(link)
	if err != nil {
		return "", ""
	}
	cleanPath := path.Clean(u.Path)

	bp := path.Clean(basePath)
	idx := strings.Index(cleanPath, bp)
	if idx == -1 {
		return "", ""
	}
	subPath := cleanPath[idx+len(bp):]
	parts := strings.Split(strings.TrimPrefix(subPath, "/"), "/")
	parts = parts[len(parts)-3:] // simple fix : lol
	version := parts[0]
	filename := parts[1]
	if strings.HasSuffix(filename, ".7z") {
		return version, "archive"
	} else if strings.HasSuffix(filename, "SHA256.txt") {
		return version, "checksum"
	}

	return "", ""
}
