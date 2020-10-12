package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitSignal() {
	ch := make(chan os.Signal)
	//监听所有信号
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	s := <-ch
	fmt.Println("退出信号", s)
}
