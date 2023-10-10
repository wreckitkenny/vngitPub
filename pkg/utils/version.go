package utils

import "os"


// GetVersion returns version from file
func GetVersion() string {
	logger := ConfigZap()

	version, err := os.ReadFile("VERSION")
	if err != nil {
		logger.Errorf("Loading version...FAILED: %s", err)
	} else {
		logger.Infof("Loading version...%s", version)
	}

	return string(version)
}