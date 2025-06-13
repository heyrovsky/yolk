package config

import (
	"os"
)

// qcow2imaged variables
var (
	BASE_ARCHIVE_LOCATION = os.TempDir()
	BASE_IMAGE_LOCATION   = "/opt"
	BASE_CONFIG_APP       = "qcow2imaged"
)
