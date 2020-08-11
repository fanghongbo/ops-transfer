package g

import (
	"context"
	"github.com/fanghongbo/dlog"
)

func InitAll() {
	InitConfig()
	InitLog()
	InitRuntime()
}

func Shutdown(ctx context.Context) error {
	defer ctx.Done()

	dlog.Close()

	return nil
}
