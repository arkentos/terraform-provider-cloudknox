package common

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

/* Private Variables */

var configuration Configuration
var configurationError error
var configOnce sync.Once

/* Public Variables */

const (
	NewPolicy string = "cloudknox_policy"
)

/* Public Functions */

func SetConfiguration(resource_path string) error {
	configOnce.Do(
		func() {
			logger := GetLogger()
			logger.Debug("msg", "Setting Constants")

			yamlFile, err := ioutil.ReadFile(resource_path)
			if err != nil {
				logger.Debug("msg", "Error Reading Configuration File", "file_read_error", err)
				configurationError = err
			}
			err = yaml.Unmarshal(yamlFile, &configuration)
			if err != nil {
				logger.Debug("msg", "Unable to Decode Into Struct", "yaml_decode_error", err)
				configurationError = err
			}

		},
	)
	return configurationError
}

func GetConfiguration() Configuration {
	return configuration
}
