package discovery

import (
	"fmt"
	"os"
	"os/exec"
)

type DiscoveryFunc func() (string, error)

type Discovery struct{}

func (d *Discovery) LookPath() (string, error) {
	osPath := os.Getenv("PATH")
	path, err := exec.LookPath(executableName)
	if err != nil {
		return "", fmt.Errorf("unable to find %s (PATH = %q): %s",
			executableName, osPath, err)
	}
	return path, nil
}
