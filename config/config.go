package config

import (
	"fmt"
	"log"
	"os"
)

// Services holds the mapping of image names to service IDs
var Services map[string]string

// LoadServices loads services from environment variables into the Services map
func LoadServices() error {
	Services = make(map[string]string)

	// Load services from environment variables
	for i := 1; ; i++ {
		serviceID := os.Getenv(fmt.Sprintf("SERVICE_%d_ID", i))
		imageName := os.Getenv(fmt.Sprintf("SERVICE_%d_IMAGE", i))

		if serviceID == "" || imageName == "" {
			// No more services found in environment variables
			break
		}

		Services[imageName] = serviceID
	}

	if len(Services) == 0 {
		return fmt.Errorf(
			"no services were loaded from environment variables",
		)
	}

	log.Println("Services map:")
	for imageName, serviceID := range Services {
		log.Printf("  %s -> %s", imageName, serviceID)
	}

	return nil
}
