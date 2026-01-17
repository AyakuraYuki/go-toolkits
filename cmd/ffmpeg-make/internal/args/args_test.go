package args_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AyakuraYuki/go-toolkits/cmd/ffmpeg-make/internal/args"
)

func TestM4RArgs_GetSS_GetTo(t *testing.T) {
	tests := []struct {
		args args.M4RArgs
		ss   time.Duration
		to   time.Duration
	}{
		{
			args: args.M4RArgs{Start: ""},
			ss:   time.Duration(0),
			to:   time.Duration(0),
		},
		{
			args: args.M4RArgs{Start: "00:00:00.000"},
			ss:   time.Duration(0),
			to:   time.Duration(0),
		},
		{
			args: args.M4RArgs{Start: "00:00:00.000", End: "00:01:10.123"},
			ss:   time.Duration(0),
			to:   time.Duration(70.123 * float64(time.Second)),
		},
	}

	for _, tt := range tests {
		assert.EqualValues(t, tt.args.GetSS(), tt.ss)
		assert.EqualValues(t, tt.args.GetTo(), tt.to)
	}
}
