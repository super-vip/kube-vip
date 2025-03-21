package equinixmetal

import (
	"encoding/json"
	"fmt"
	"os"

	log "log/slog"

	"github.com/packethost/packngo"
)

func findProject(project string, c *packngo.Client) *packngo.Project {
	l := &packngo.ListOptions{Includes: []string{project}}
	ps, _, err := c.Projects.List(l)
	if err != nil {
		log.Error(err.Error())
	}
	for _, p := range ps {

		// Find our project
		if p.Name == project {
			return &p
		}
	}
	return nil
}

func findSelf(c *packngo.Client, projectID string) *packngo.Device {
	// Go through devices
	dev, _, _ := c.Devices.List(projectID, &packngo.ListOptions{})
	for _, d := range dev {
		// TODO do we need to replace os.Hostname with config.NodeName here?
		me, _ := os.Hostname()
		if me == d.Hostname {
			return &d
		}
	}
	return nil
}

// GetPacketConfig will lookup the configuration from a file path
func GetPacketConfig(providerConfig string) (string, string, error) {
	var config struct {
		AuthToken string `json:"apiKey"`
		ProjectID string `json:"projectId"`
	}
	// get our token and project
	if providerConfig != "" {
		configBytes, err := os.ReadFile(providerConfig)
		if err != nil {
			return "", "", fmt.Errorf("failed to get read configuration file at path %s: %v", providerConfig, err)
		}
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			return "", "", fmt.Errorf("failed to process json of configuration file at path %s: %v", providerConfig, err)
		}
	}
	return config.AuthToken, config.ProjectID, nil
}
