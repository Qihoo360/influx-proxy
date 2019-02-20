// Copyright 2016 Eleme. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/chengshiwen/influx-proxy/backend"
)

var (
	ErrConfig  = errors.New("config parse error")
	ConfigFile string
	NodeName   string
	LogPath    string
	StoreDir   string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	flag.StringVar(&ConfigFile, "config", "proxy.json", "proxy config file")
	flag.StringVar(&NodeName, "node", "l1", "node name")
	flag.StringVar(&LogPath, "log-path", "bin/influx-proxy.log", "log file path")
	flag.StringVar(&StoreDir, "data-dir", "bin/place", "dir to store .dat .rec")
	flag.Parse()
}

// initLog log初始化
func initLog() {
	if LogPath == "" {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LogPath,
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     7,
		})
	}
}

// PathExists 检查目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	initLog()
	exist, err := PathExists(StoreDir)
	if err != nil {
		log.Println("check dir error!")
		return
	}
	if !exist {
		err = os.MkdirAll(StoreDir, os.ModePerm)
		if err != nil {
			log.Println("create dir error!")
			return
		}
	}

	fcs := backend.NewFileConfigSource(ConfigFile, NodeName)
	var nodecfg backend.NodeConfig
	nodecfg, err = fcs.LoadNode()
	if err != nil {
		log.Printf("config source load failed.")
		return
	}

	ic := backend.NewInfluxCluster(fcs, &nodecfg, StoreDir)
	ic.LoadConfig()

	mux := http.NewServeMux()
	NewHttpService(ic, nodecfg.DB).Register(mux)
	log.Printf("http service start.")
	server := &http.Server{
		Addr:        nodecfg.ListenAddr,
		Handler:     mux,
		IdleTimeout: time.Duration(nodecfg.IdleTimeout) * time.Second,
	}

	if nodecfg.IdleTimeout <= 0 {
		server.IdleTimeout = 10 * time.Second
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
