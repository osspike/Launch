package color

import (
	"reflect"
	"testing"
)

func TestPrefix(t *testing.T) {
	assert := func(want, got []byte) {
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v got %v", want, got)
		}
	}

	assert([]byte("\033[32m"), prefix(Color{Value: GREEN}))
	assert([]byte("\033[33m"), prefix(Color{Value: YELLOW}))
}
