package redis

import (
	"fmt"
	"testing"

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

	bytes, err := c.Get("eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwdXZPdzZoMzNQYVNETFQzSVprZnRGdXpGRzBrOTBHU3hWZzhfYmt4c2t3In0.eyJleHAiOjE3MDYzMDU3NzIsImlhdCI6MTcwNjMwNTQ3MiwianRpIjoiMDZjODg4NzEtYmRlNC00Mjg2LTg1ZTQtZjMxZGU1NmU2ZWU5IiwiaXNzIjoiaHR0cDovLzE5Mi4xNjguODYuMjExOjMyMDg4L3JlYWxtcy9uc2h1YiIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJhMTJmZmZlOC1iZTIxLTQ5NWYtOGZhZS1hYWYxYTc5MmJhMGIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhdXRoLXN2YyIsInNlc3Npb25fc3RhdGUiOiI4NjM1OWIzOS00ZjU1LTRmOWItODQwMi04ZWVlNjhlY2YzMTIiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIi8qIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLW5zaHViIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSIsInNpZCI6Ijg2MzU5YjM5LTRmNTUtNGY5Yi04NDAyLThlZWU2OGVjZjMxMiIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYWFyb24ifQ.WYFC5e1V0U6y1L1lDE8JCrRNvbYZMbqstTTucAZHtL0aNAgkRxTnaKLSv7qfz8yFIB6isv6K3R3s9aG4G96BmdSv7nLj2TOVhOg5Av7dKjChtxwKxtILnnQfrNX-NivI0fl83ZNVMaNJSkyOUNcBmL11VWICi8KQ9cuuJ9zIL4ja23Izo1nPtEM9OgHJ8QB5Ggy1aWWivuCR7kHikp3_yRovSLdqnKE9XqSNswSf76aIFKIMOoWjY4vJw40vGLzCcCzljr-KQyRV6r7QJIaMMr84fQg0_HPpQyMtLs1SQ0uHKJH8tglJsnK11-iHl7D727wGDLnAOU3n-Phy5JIhUA")
	assert.Nil(t, err)
	assert.Equal(t, string(bytes), "")
	fmt.Println(bytes)
	/*
		err = c.Set("Aaron", []byte("4"), 2*time.Second)
		assert.Nil(t, err)

		by, err := c.Get("Aaron")
		assert.Nil(t, err)

		log.Info("Here is the value fetched", string(by))

		time.Sleep(2 * time.Second)

		by, err = c.Get("Aaron")
		assert.NotNil(t, err)
	*/

}
