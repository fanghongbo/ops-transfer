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

	// 刷新日志缓存
	dlog.Close()

	return nil
}
