package main

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

const (
	URLS  = "url"
	PATH  = "path"
	TITLE = "title"
	START = "start"
)

var (
	g_stURLs []string
	g_path   string
	g_title  string
	g_start  int
)

func loadCfg() {
	cfg := viper.New()
	p, _ := os.Getwd()
	cfg.AddConfigPath(p)
	cfg.SetConfigName("cfg")
	cfg.SetConfigType("yaml")

	var lock sync.Mutex
	err := cfg.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal("ReadInConfig ", err)
	}
	lock.Lock()
	defer lock.Unlock()

	g_stURLs = cfg.GetStringSlice(URLS)
	g_path = cfg.GetString(PATH)
	g_title = cfg.GetString(TITLE)
	g_start = cfg.GetInt(START)
}

func myPanic(v any) {
	viper.Set(START, g_count)
	panic(v)
}
