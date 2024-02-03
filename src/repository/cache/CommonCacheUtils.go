package cache

import (
	"fmt"
	"github.com/strikersk/user-auth/config"
	"log"
)

const cachePrefix = "user:"

func cacheUrlResolver(configuration config.CacheConfiguration) (address string) {
	if configuration.URL != "" {
		address = configuration.URL
	} else if configuration.Host != "" && configuration.Port != "" {
		address = fmt.Sprintf("%s:%s", configuration.Host, configuration.Port)
	} else {
		log.Fatal("cache address not provide")
	}
	return
}
