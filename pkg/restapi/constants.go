package restapi

import (
	"time"
)

const (
	API_PATH_LICENSE = "/vmturbo/rest/license"
	API_PATH_TARGET  = "/vmturbo/rest/targets"

	defaultTimeOut = time.Duration(60 * time.Second)
)

