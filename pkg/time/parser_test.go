package time_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vietnam-immigrations/vs2-utils-go/pkg/time"
)

func TestFromWooDateString(t *testing.T) {
	result, err := time.FromWooDateString("03/05/1988")
	assert.NoError(t, err)
	timezone, _ := result.Zone()
	assert.Equal(t, "+07", timezone)
	assert.Equal(t, 1988, result.Year())
	assert.Equal(t, 5, int(result.Month()))
	assert.Equal(t, 3, result.Day())
	assert.Equal(t, 0, result.Hour())
	assert.Equal(t, 0, result.Minute())
	assert.Equal(t, 0, result.Second())
	assert.Equal(t, 0, result.Nanosecond())
}
