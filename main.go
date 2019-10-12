package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"./packs/gin"
	"./packs/util"
	"time"
)

// some init before the main func
func init() {
	log.Println("main.init")
	// log path setting
	RootPath, err := os.Getwd()
	util.CheckErr(err)
	logPath := RootPath + util.ConstConfigLogPath
	logFile, _ := os.Create(logPath)

	// request log out put, file and terminate stdout
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
}

func main() {
	if "dev" == util.ConstConfigEnv {
		gin.SetMode(gin.DebugMode)
	}else{
		gin.SetMode(gin.ReleaseMode)
	}

	r := LoadRouters()

	r.Delims("{{", "}}")

	r.SetFuncMap(template.FuncMap{
		"echo"		:	fmt.Sprintf,
		"date"		:	util.DateFormat,
		"str2html"	:	util.Str2html,
		"unix2time"	:	util.Unix2time,
		"unix2date"	:	util.Unix2date,
		"date2unix"	:	util.Date2unix,
	})

	r.LoadHTMLGlob("views/**/**/*")

	r.Static("/static", "static")

	port :=  util.DynamicSettingInstance.GetString("HttpPort", util.ConstConfigDefaultValue)

	s := &http.Server{
		Addr:           port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s \n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<- quit

	log.Println("Shutdown  Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
