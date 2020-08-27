package signal

import (
	"context"
	"fmt"
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
		syscall.SIGSTOP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGKILL,
	)

	go func() {
		var sig os.Signal
		for {
			sig = <-sigChan
			fmt.Println()
			switch sig {
			case syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGUSR1:
				_ = app.Stop() // graceful stop
			case syscall.SIGHUP:
				_ = app.GracefulStop(context.TODO())
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGUSR2, syscall.SIGTERM:
				_ = app.Stop() // terminate now
			}
			time.Sleep(time.Second * 3)
		}
	}()
}
