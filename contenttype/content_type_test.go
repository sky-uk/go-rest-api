package contenttype

import (
	//"github.com/sky-uk/go-rest-api/contenttype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContentType(t *testing.T) {
	assert.Equal(t, "html", GetType("Content-Type: text/html"))
	assert.Equal(t, "plain", GetType("Content-Type: text/plain"))
	assert.Equal(t, "x-www-form-urlencoded", GetType("Content-Type: application/x-www-form-urlencoded"))
}
