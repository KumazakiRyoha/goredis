package tcp

import (
	"context"
	"github.com/hdt3213/godis/lib/logger"
	"goredis/interface/tcp"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	// os.Signal 底层操作系统级别给进程或者线程的信号
	sigChan := make(chan os.Signal)
	// 操作系统收到SIGHUP\SIGQUIT等信号时会转发给sigChan
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		signal := <-sigChan
		switch signal {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	ListenAndServe(listener, handler, closeChan)
	return nil

}

func ListenAndServe(listener net.Listener, handle tcp.Handler,
	closeChan <-chan struct{}) {
	ctx := context.Background()
	var waitDone sync.WaitGroup

	go func() {
		<-closeChan
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handle.Close()
	}()

	defer func() {
		_ = listener.Close()
		_ = handle.Close()
	}()
	for true {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		// 每为一个新链接服务就加一
		waitDone.Add(1)
		go func() {
			defer func() {
				// 服务结束就减一
				waitDone.Done()
			}()
			handle.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
