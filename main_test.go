package tavern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	assert := assert.New(t)
	err := Add("123456").Check()
	assert.NoError(err)
	err = Add("").Check()
	assert.NoError(err)
	err = Add("0").Check()
	assert.NoError(err)

	err = Add("123456").Required().Check()
	assert.NoError(err)
	err = Add("").Required().Check()
	assert.Error(err)
	err = Add(0).Required().Check()
	assert.NoError(err)

	var s interface{}
	s = ""
	err = Add(s).Required().Check()
	assert.Error(err)
	s = nil
	err = Add(s).Required().Check()
	assert.Error(err)
	s = nil
	err = Add(s).Check()
	assert.NoError(err)
	s = 0
	err = Add(s).Required().Check()
	assert.NoError(err)
}

func TestLength(t *testing.T) {
	assert := assert.New(t)
	err := Add("123456").Length(1, 3).Check()
	assert.Error(err)
	err = Add("123456").Length(1, 6).Check()
	assert.NoError(err)
	err = Add("123456").Length(0, 1).Check()
	assert.Error(err)
	err = Add("123456").Length(0, 0).Check()
	assert.Error(err)
	err = Add("").Length(0, 0).Check()
	assert.NoError(err)
	err = Add("").Length(1, 2).Check()
	assert.Error(err)

	err = Add(123456).Length(1, 3).Check()
	assert.Error(err)
	err = Add(123456).Length(1, 6).Check()
	assert.NoError(err)
	err = Add(123456).Length(0, 1).Check()
	assert.Error(err)
	err = Add(123456).Length(0, 0).Check()
	assert.Error(err)
}
