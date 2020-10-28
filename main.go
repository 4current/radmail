package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}

func main() {

	v := viper.New()
	cfgDir := "$HOME/Dropbox/ham/Winlink/pat_mail/.wl2k"
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(cfgDir)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("A File Event happened: ", e.Name, e.Op)
	})
	v.WatchConfig()

	var settings map[string]interface{}
	settings = v.AllSettings()
	viper.Unmarshal(&settings)

	fmt.Println(reflect.TypeOf(settings))

	for k, v := range settings {
		fmt.Println("key: ", k, "    value: ", reflect.TypeOf(v))
	}

	fmt.Printf("Press Ctrl+C to end\n")
	waitForCtrlC()
	fmt.Printf("\n")

}
