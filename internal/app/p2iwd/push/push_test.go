package push

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsoluteURL(t *testing.T) {
	m := make(map[string]bool)
	m["http://localhost/v2/library/n3dr/blobs/uploads/530c9d69-4b4a-4809-9c66-83d98d664e56?_state=zxExkA9bBHlHKTLQpmrIP2e9SqRDVY7UB0N_uqSqWMd7Ik5hbWUiOiJsaWJyYXJ5L24zZHIiLCJVVUlEIjoiNTMwYzlkNjktNGI0YS00ODA5LTljNjYtODNkOThkNjY0ZTU2IiwiT2Zmc2V0IjowLCJTdGFydGVkQXQiOiIyMDIyLTA0LTE5VDA1OjM5OjE5Ljg4NTg5NDk2WiJ9"] = true
	m["http://localhost"] = true
	m["https://localhost"] = true
	m["localhost"] = false

	for k, v := range m {
		b, _ := absoluteURL(k)
		assert.Equal(t, v, b)
	}
}
