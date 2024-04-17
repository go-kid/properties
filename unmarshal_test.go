package properties

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	var m = map[string]any{}
	err := Unmarshal([]byte(`
a.b=1
a.c=2
b.c=3
b.d=4
`), &m)
	assert.NoError(t, err)
	fmt.Println(m)
}
