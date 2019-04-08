// Copyright 2016 Eleme. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/wilhelmguo/influx-proxy/backend"
	"github.com/wilhelmguo/influx-proxy/logs"
)

var (
	ErrConfig  = errors.New("config parse error")
	ConfigFile string
	NodeName   string
	StoreDir   string
	RavenDSN   string
)

func init() {

	flag.StringVar(&ConfigFile, "config", "proxy.json", "proxy config file")
	flag.StringVar(&NodeName, "node", "l1", "node name")
	flag.StringVar(&RavenDSN, "raven-dsn", "", "the sentry dsn, leave it empty if you not use sentry.")
	flag.StringVar(&StoreDir, "data-dir", "data", "dir to store .dat .rec")
	flag.Parse()
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
	logs.InitLog(RavenDSN)

	exist, err := PathExists(StoreDir)
	if err != nil {
		logs.Error("check data dir error")
		return
	}
	if !exist {
		err = os.MkdirAll(StoreDir, os.ModePerm)
		if err != nil {
			logs.Error("check data dir error")
			return
		}
	}

	fcs := backend.NewFileConfigSource(ConfigFile, NodeName)
	nodecfg, err := fcs.LoadNode()
	if err != nil {
		logs.Errorf("config source load failed.")
		return
	}

	ic := backend.NewInfluxCluster(fcs, &nodecfg, StoreDir)
	ic.LoadConfig()

	mux := http.NewServeMux()
	NewHttpService(ic).Register(mux)
	logs.Info("http service start.")
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
		logs.Error(err)
		return
	}
}
