package main

import (
	"context"
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/http"
	"github.com/fanghongbo/ops-transfer/rpc"
	"github.com/fanghongbo/ops-transfer/sender"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	g.InitAll()

	go sender.Start()
	go rpc.Start()
	go http.Start()

	// 等待中断信号以优雅地关闭 transfer（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Transfer ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := g.Shutdown(ctx); err != nil {
		log.Fatal("Transfer Shutdown:", err)
	} else {
		log.Println("Transfer Exiting")
	}
}
