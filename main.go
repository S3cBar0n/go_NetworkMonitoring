package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/go-ping/ping"
)

func netWorker(host string, receive chan string) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Println("Failed to launch the pinger service.")
		return
	}
	pinger.Timeout = 10000000 * 1000
	pinger.OnFinish = func(stats *ping.Statistics) {
		//fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v packet loss\n",
		//	stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
		if stats.PacketLoss >= 35.0 {
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

func emailSend(offlineHosts []string) {
	fmt.Println(offlineHosts)
	authentication := smtp.PlainAuth("", "USERNAME", "SECURELY STORED PASSWORD", "MAILSERVER WITHOUT PORT")

	to := []string{"ADDRESS TO SEND THE MAIL TO"}
	msg := []byte(
		"To: ADDRESS TO SEND THE MAIL TO\r\n" +
			"Subject: Network Down\r\n" +
			"MIME-version: 1.0;\nContent-Type: text/html;\r\n" +
			"\r\n" +
			"<p>The <strong>following</strong> IP's have gone offline:\r\n")

	for item := range offlineHosts {
		item_msg := fmt.Sprintf("<p>%s</p>", offlineHosts[item])
		msg = append(msg, item_msg...)
	}
	msg = append(msg, "<p>Report any strange behavior to the Network Administrator.<br><br>Thanks,<br>NetworkMonitoring</p>"...)

	err := smtp.SendMail("SMTP_SERVER:PORT", authentication, "FROM_ADDRESS", to, msg)
	if err != nil {
		log.Fatalf("%s - Failed to send email", err)
	}
	fmt.Print("email sent")
}

func main() {
	var hosts []string
	var status []string
	var offlineHosts []string
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
		if strings.Contains(status[i], "Offline") {
			offlineHosts = append(offlineHosts, status[i])
		} else {
			continue
		}
	}

	if len(offlineHosts) >= 1 {
		emailSend(offlineHosts)
	}
}
