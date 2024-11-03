//go:build !darwin
// +build !darwin

/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-01 15:11:36
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-09-05 16:28:46
 * @FilePath: /DappsterOS/route/periodical.go
 * @Description:
 * @Website: https://www.dappsteros.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package route

import (
	"strings"
	"time"
	"unsafe"

	"github.com/dappster-io/DappsterOS/model"
	"github.com/dappster-io/DappsterOS/service"
)

func SendAllHardwareStatusBySocket() {
	netList := service.MyService.System().GetNetInfo()
	newNet := []model.IOCountersStat{}
	nets := service.MyService.System().GetNet(true)
	for _, n := range netList {
		for _, netCardName := range nets {
			if n.Name == netCardName {
				item := *(*model.IOCountersStat)(unsafe.Pointer(&n))
				item.State = strings.TrimSpace(service.MyService.System().GetNetState(n.Name))
				item.Time = time.Now().Unix()
				newNet = append(newNet, item)
				break
			}
		}
	}
	cpuPercents := service.MyService.System().GetCpuPercents()
	cpuPercent := service.MyService.System().GetCpuPercent()

	var cpuModel = "arm"
	var cpuModelName = ""
	if cpus := service.MyService.System().GetCpuInfo(); len(cpus) > 0 {
		if strings.Count(strings.ToLower(strings.TrimSpace(cpus[0].ModelName)), "intel") > 0 {
			cpuModel = "intel"
		} else if strings.Count(strings.ToLower(strings.TrimSpace(cpus[0].ModelName)), "amd") > 0 {
			cpuModel = "amd"
		}
		cpuModelName = strings.TrimSpace(cpus[0].ModelName)
	}

	num := service.MyService.System().GetCpuCoreNum()
	cpuData := make(map[string]interface{})
	cpuData["percents"] = cpuPercents
	cpuData["percent"] = cpuPercent
	cpuData["num"] = num
	cpuData["temperature"] = service.MyService.System().GetCPUTemperature()
	cpuData["power"] = service.MyService.System().GetCPUPower()
	cpuData["model"] = cpuModel
	cpuData["modelName"] = cpuModelName

	memInfo := service.MyService.System().GetMemInfo()

	body := make(map[string]interface{})

	body["sys_mem"] = memInfo

	body["sys_cpu"] = cpuData

	body["sys_net"] = newNet
	systemTempMap := service.MyService.Notify().GetSystemTempMap()
	systemTempMap.Range(func(key, value interface{}) bool {
		body[key.(string)] = value
		return true
	})

	body["hardware"] = service.MyService.System().GetDeviceInfo()
	service.MyService.Notify().SendNotify("dappsteros:system:utilization", body)
}

// func MonitoryUSB() {
// 	var matcher netlink.Matcher

// 	conn := new(netlink.UEventConn)
// 	if err := conn.Connect(netlink.UdevEvent); err != nil {
// 		logger.Error("udev err", zap.Any("Unable to connect to Netlink Kobject UEvent socket", err))
// 	}
// 	defer conn.Close()

// 	queue := make(chan netlink.UEvent)
// 	errors := make(chan error)
// 	quit := conn.Monitor(queue, errors, matcher)

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
// 	go func() {
// 		<-signals
// 		close(quit)
// 		os.Exit(0)
// 	}()

// 	for {
// 		select {
// 		case uevent := <-queue:
// 			if uevent.Env["DEVTYPE"] == "disk" {
// 				time.Sleep(time.Microsecond * 500)
// 				SendUSBBySocket()
// 				continue
// 			}
// 		case err := <-errors:
// 			logger.Error("udev err", zap.Any("err", err))
// 		}
// 	}

// }
