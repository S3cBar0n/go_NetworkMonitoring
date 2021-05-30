package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-ping/ping"
)

func netWorker(host string, receive chan string) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Println("Failed to launch the pinger service.")
		return
	}
	pinger.SetPrivileged(true)
	pinger.Timeout = 10000000 * 1000
	pinger.OnFinish = func(stats *ping.Statistics) {
		//fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v packet loss\n",
		//	stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
		if stats.PacketLoss >= 50.0 {
			host = fmt.Sprintf("%s - %s with %v", "Offline", host, stats.AvgRtt)
			receive <- host
			return
		}
		host = fmt.Sprintf("%s - %s, with %v", "Online", host, stats.AvgRtt)
		receive <- host
	}
	err = pinger.Run()
	if err != nil {
		host = fmt.Sprintf("%s - %s", "Could not ping", host)
		receive <- host
	}
}

func main() {
	var hosts []string
	var status []string
	receive := make(chan string)

	hostsfile, err := os.Open("hosts.csv")
	if err != nil {
		log.Fatalln("Could not read the csv file provided.")
	}
	read_csv := csv.NewReader(hostsfile)
	for {
		item, err := read_csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		hosts = append(hosts, item[0])
	}

	for i := range hosts {
		go netWorker(hosts[i], receive)
	}

	for range hosts {
		processedHosts := <-receive
		status = append(status, processedHosts)
	}

	for i := range status {
		fmt.Println(status[i])
	}

}
