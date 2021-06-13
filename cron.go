package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)

func main() {

	c := cron.New()
	c.AddFunc("*/1 * * * *", ReadRouter())
	c.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
