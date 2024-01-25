package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/sitenet/svclib/logger"
)

func TestCache(t *testing.T) {
	log, _ := logger.NewLogger()

	config := RedisConfig{Addrs: []string{"192.168.86.211:32379"}}
	c, err := NewRedisCache(&config)

	if err != nil {
		log.Fatal("Error in connecting server")
	}

	_, err = c.Get("Aaron")
	assert.NotNil(t, err)

	err = c.Set("Aaron", []byte("4"), 2*time.Second)
	assert.Nil(t, err)

	by, err := c.Get("Aaron")
	assert.Nil(t, err)

	log.Info("Here is the value fetched", string(by))

	time.Sleep(2 * time.Second)

	by, err = c.Get("Aaron")
	assert.NotNil(t, err)

}
