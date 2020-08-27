package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HookSignals(app App) {
	sigChan := make(chan os.Signal)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	go func() {
		var sig os.Signal
		for {
			sig = <-sigChan
			switch sig {
			case syscall.SIGQUIT:
				_ = app.Stop() // graceful stop
			case syscall.SIGHUP:
				_ = app.GracefulStop(context.TODO())
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
				_ = app.Stop() // terminate now
			}
			time.Sleep(time.Second * 3)
		}
	}()
}
