package exec

import (
	"context"

	"github.com/go-logr/logr"

	"github.com/AyakuraYuki/go-toolkits/cmd/ffmpeg-make/internal/args"
)

var _ Execute = (*execute)(nil)

type Execute interface {
	MakeM4r(ctx context.Context, args args.M4RArgs) (err error)
}

type Options struct {
	Logger logr.Logger
}

type execute struct {
	logger logr.Logger
}

func NewExecute(opts Options) Execute {
	return &execute{logger: opts.Logger}
}
