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

type beacon struct {
	every       int    `json:"every"`
	message     string `json:"message"`
	destination string `json:"destination"`
}

type ax25 struct {
	port   int    `json:"port"`
	beacon int    `json:"beacon"`
	rig    string `json:"rig"`
}

type hamlibRig struct {
	Address string `json:"address"`
	Network string `json:"network"`
}

type serialTNC struct {
	Path     string `json:"path"`
	Baudrate int    `json:"baudrate"`
	Type     string `json:"type"`
}

type winmor struct {
	Addr             string `json:"addr"`
	InboundBandwidth int    `json:"inbound_bandwidth"`
	DriveLevel       int    `json:"drive_level"`
	Rig              string `json:"rig"`
	PttCtrl          bool   `json:"ptt_ctrl"`
}
type pactor struct {
	Path             string `json:"path"`
	BaudRate         int    `json:"baud_rate"`
	Rig              string `json:"rig"`
	CustomInitScript bool   `json:"custom_init_script"`
}
type telnet struct {
	ListenAddr string `json:"listen_addr"`
	Password   int    `json:"password"`
}

type gpsd struct {
	Addr          string `json:"addr"`
	EnableHTTP    bool   `json:"enable_http"`
	UseServerTime bool   `json:"use_server_time"`
}

type schedule map[string]string

type config struct {
	MyCall                   string               `json:"mycall"`
	FormsPath                string               `json:"forms_path"`
	SecureLoginPassword      string               `json:"secure_login_password"`
	AuxiliaryAddresses       []string             `json:"auxiliary_addresses"`
	Locator                  string               `json:"locator"`
	ServiceCodes             []string             `json:"service_codes"`
	HTTPAddr                 string               `json:"http_addr"`
	MOTD                     []string             `json:"motd"`
	ConnectAliases           []string             `json:"connect_aliases"`
	Listen                   []int                `json:"listen"`
	HamlibRigs               map[string]hamlibRig `json:"hamlib_rigs"`
	AX25                     ax25                 `json:"ax25"`
	SerialTNC                serialTNC            `json:"serial-tnc"`
	WINMOR                   winmor               `json:"winmor"`
	Pactor                   pactor               `json:"pactor"`
	Telnet                   telnet               `json:"telnet"`
	GPSD                     gpsd                 `json:"gpsd"`
	Schedule                 schedule             `json:"schedule"`
	VersionReportingDisabled bool                 `json:"version_reporting_disabled"`
}

var myCfg config

func showIt(x map[string]interface{}) {
	for k, v := range x {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			fmt.Println(k, "     ", v)
		case reflect.Bool:
			fmt.Println(k, "     ", v)
		case reflect.Slice:
			fmt.Println(k, "     slice - ", reflect.TypeOf(v))
		case reflect.Map:
			fmt.Println(k, "     map - ", reflect.TypeOf(v))
			fmt.Println(k, "     map - ", reflect.TypeOf(v).Kind())
		default:
			fmt.Println(k, "     ", reflect.TypeOf(v))
			fmt.Println(k, "     ", reflect.TypeOf(v).Kind())
		}
	}
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
	showIt(settings)

	var serialTNCSettings map[string]string
	serialTNCSettings = v.GetStringMapString("serial-tnc")
	viper.Unmarshal(&serialTNCSettings)
	fmt.Println("serial-tnc")
	for k, v := range serialTNCSettings {
		fmt.Println("  ", k, "   ", v)
	}

	myCfg.myCall = v.GetString("mycall")
	myCfg.formsPath = v.GetString("forms_path")

	fmt.Println(myCfg)

	fmt.Printf("Press Ctrl+C to end\n")
	waitForCtrlC()
	fmt.Printf("\n")

}
