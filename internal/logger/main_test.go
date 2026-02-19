package logger

import(
	"testing"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Run("Test adapter", func(t *testing.T) {
		var b bytes.Buffer
		f := func(m ...any) {
			fmt.Fprint(&b, m...)
		}
		la := LoggerAdapter(f)
		la.Log("HI")
		assert.Equal(t, "HI", b.String())
	})
}
