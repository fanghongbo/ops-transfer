package g

import "context"

func InitAll() {
	InitConfig()
	InitLog()
	InitRuntime()
}

func Shutdown(ctx context.Context) error {
	defer ctx.Done()

	return nil
}
