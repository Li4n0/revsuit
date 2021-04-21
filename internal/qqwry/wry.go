package qqwry

import (
	"os"
	"sync"
	"time"

	"github.com/sinlov/qqwry-golang/qqwry"
	log "unknwon.dev/clog/v2"
)

var wry *qqwry.QQwry
var once sync.Once

func init() {
	info, err := os.Stat("qqwry.dat")
	if err != nil {
		if os.IsNotExist(err) {
			err := download()
			if err != nil {
				log.Error("Download qqwry.dat failed, caused by:%v, recommend to download it by yourself otherwise the `IpArea` will be null", err.Error())
			}
		}
	} else if time.Until(info.ModTime()) > 5*24*time.Hour {
		log.Info("Updating qqwry.dat...")
		err := download()
		if err != nil {
			log.Warn("Update qqwry.dat failed, please download qqwry.dat by yourself")
		}
	}
}

func GetQQWry() *qqwry.QQwry {
	once.Do(func() {
		qqwry.DatData.FilePath = "qqwry.dat"
		init := qqwry.DatData.InitDatFile()
		if v, ok := init.(error); ok {
			if v != nil {
				log.Error("qqwry init failed")
				wry = nil
			}
		}
		wry = qqwry.NewQQwry()
	})

	return wry
}
