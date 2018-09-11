package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)
var (
	device string = "en0"
	snapshot_len int32 = 1024
	promiscuous bool = false
	err error
	timeout time.Duration = 30 * time.Second
	handle * pcap.Handle
)

func main (){
	handle,err = pcap.OpenLive(device,snapshot_len,promiscuous,timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle,handle.LinkType())

	for packet := range packetSource.Packets(){
		fmt.Println(packet)
	}
}

