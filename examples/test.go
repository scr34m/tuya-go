package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	tuya "github.com/scr34m/tuya-go"
)

type Config struct {
	CertSign  string
	Secret    string
	AppSecret string
	DeviceId  string
	ClientId  string
	Email     string
	Password  string
}

func main() {
	var config Config
	if _, err := toml.DecodeFile("tuya.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	t, _ := tuya.New(tuya.Config{
		CertSign:  config.CertSign,
		Secret:    config.Secret,
		AppSecret: config.AppSecret,
	}, config.DeviceId, config.ClientId)
	//t.Call("tuya.m.country.list", "2.0", nil, false)

	if len(os.Args) != 2 {
		t.Login(config.Email, config.Password, 36)
		return
	}

	t.SessionIdSet(os.Args[1])
	// b.m.device.register
	// tuya.m.notice.new.get
	/*
		res := t.Call("tuya.m.location.list", "2.1", nil, true)
		if !res.Success {
			fmt.Println(res.ErrorMessage)
			return
		}
		locations := res.GetLocationList()
		t.GroupIdSet(locations[0].GroupID)
	*/

	t.GroupIdSet("17957495")

	res := t.Call("tuya.m.my.group.device.list", "1.0", nil, true)
	if !res.Success {
		fmt.Println(res.ErrorMessage)
		return
	}

	devices := res.GetDeviceList()
	for _, device := range devices {
		fmt.Printf("Name: %s\n", device.Name)
		fmt.Printf("DevID: %s\n", device.DevID)
		fmt.Printf("LocalKey: %s\n", device.LocalKey)
		fmt.Printf("DPS: %v\n", device.Dps)
		//fmt.Printf("%#v\n", device)
		fmt.Printf("\n")
	}
}
