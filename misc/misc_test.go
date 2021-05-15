package misc

import (
	"testing"

	"github.com/matryer/is"
)

func TestSuccess(t *testing.T) {
	is := is.New(t)

	is.True(!Success(199))
	is.True(!Success(300))
	is.True(Success(200))
}
