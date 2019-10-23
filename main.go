package main

import (
	"./packs/gin"
	"./packs/util"
	"./packs/config"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	listener net.Listener
	err error
	// formating the hotfix cmd name: hf, and init the bool value to false
	graceful = flag.Bool("hf", false, "Listen on fd open 3 [internal user only]")
)

// some init before the main func
func init() {

	log.Println("main.init")
	// log path setting
	RootPath, err := os.Getwd()
	util.Assert(err)
	logPath := RootPath + config.ConstConfigLogPath
	logFile, _ := os.Create(logPath)

	// request log out put, file and terminate stdout
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
}

func main() {
	port := config.DynamicSettingInstance.GetString("HttpPort", config.ConstConfigDefaultValue)
	flag.Parse()
	fmt.Println("start-up at " , time.Now(), *graceful)
	if *graceful {
		f := os.NewFile(3, "")
		listener, err = net.FileListener(f)
		fmt.Printf( "graceful-reborn  %v %v  %#v \n", f.Fd(), f.Name(), listener)
	}else{
		listener, err = net.Listen("tcp", port)
		tcp,_ := listener.(*net.TCPListener)
		fd,_ := tcp.File()
		fmt.Printf( "first-boot  %v  %#v \n ", fd.Fd(), listener)
	}


	server := &http.Server{
		Addr:			port,
		Handler:		LoadRouters(),
		ReadTimeout:	config.HTTP_READTIMEOUT,
		WriteTimeout:	config.HTTP_WRETETIMEOUT,
		MaxHeaderBytes:	config.HTTP_MAXHEADERSBYTES,
	}

	log.Printf("Actual pid is %d\n", syscall.Getpid())
	if err != nil {
		println(err)
		return
	}
	log.Printf(" listener: %v\n",   listener)
	go func(){//不要阻塞主进程
		err := server.Serve(listener)
		if err != nil {
			log.Println(err)
		}
	}()

	//signals
	func(){
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM)
		for{//阻塞主进程， 不停的监听系统信号
			sig := <- ch
			log.Printf("signal: %v", sig)
			ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
			switch sig {
			case syscall.SIGTERM, syscall.SIGHUP:
				println("signal cause reloading")
				signal.Stop(ch)
				{//fork new child process
					tl, ok := listener.(*net.TCPListener)
					if !ok {
						fmt.Println("listener is not tcp listener")
						return
					}
					currentFD, err := tl.File()
					if err != nil {
						fmt.Println("acquiring listener file failed")
						return
					}
					cmd := exec.Command(os.Args[0], "-hf")
					cmd.ExtraFiles, cmd.Stdout,cmd.Stderr = []*os.File{currentFD} ,os.Stdout, os.Stderr
					err = cmd.Start()

					if err != nil {
						fmt.Println("cmd.Start fail: ", err)
						return
					}
					fmt.Println("forked new pid : ",cmd.Process.Pid)
				}

				server.Shutdown(ctx)
				fmt.Println("graceful shutdown at ", time.Now())
			}

		}
	}()
}
