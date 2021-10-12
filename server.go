package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/robfig/cron/v3"
	"iot-manager/routes"
	"iot-manager/controller"
	"log"
	"net/http"
	"time"
)

func pingPhone() {
	pinger, err := ping.NewPinger("192.168.1.74")
	if err != nil {
		panic(err)
	}

	pinger.Count = 10
	pinger.Timeout = time.Second * 13

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		if stats.PacketsRecv == 0 && stats.PacketsSent > 5 {
			fmt.Printf("Can't ping device, turning everything off\n")
			controller.TurnOffAllDevices()
		}
	}
	err = pinger.Run()
	if err != nil {
		panic(err)
	}

}
func setUpCron() {
	c := cron.New()
	c.AddFunc("*/30 * * * *", func() { pingPhone() })
	c.Start()
}

func main() {
	setUpCron()
	router := routes.IOTRoutes()
	http.Handle("/iotapi", router)
	log.Fatal(http.ListenAndServe(":8081", router))
}
