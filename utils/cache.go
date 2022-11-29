package utils

import (
	"net/url"

	"github.com/jellydator/ttlcache/v3"
)

var Cache *ttlcache.Cache[string, *url.URL]
