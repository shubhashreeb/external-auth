package redis

import (
	"time"
)

// Interface for the cache
type Cache interface {
	Get(k string) ([]byte, error)
	Set(k string, v []byte, d time.Duration) error
	Delete(k string) error
}
