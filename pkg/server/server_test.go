package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKey_ok(t *testing.T) {
	config := NewDefaultConfig()
	c, err := newKey(config)
	assert.Nil(t, err)
	assert.Equal(t, config, c.config)
	// TODO assert.NotEmpty on other elements of server struct
	//assert.NotEmpty(t, c.storer)
}

func TestNewKey_err(t *testing.T) {
	badConfigs := map[string]*Config{
	// TODO add bad config instances
	}
	for desc, badConfig := range badConfigs {
		c, err := newKey(badConfig)
		assert.NotNil(t, err, desc)
		assert.Nil(t, c)
	}
}

// TODO add TestKey_ENDPOINT_(ok|err) for each ENDPOINT
