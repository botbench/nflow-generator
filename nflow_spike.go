package main

import (
	"log"
	"os"
)

//import "fmt"

//Generate a netflow packet w/ user-defined record count
func GenerateSpike(spike_proto Proto) Netflow {
	data := new(Netflow)
	header := CreateNFlowHeader(1)
	records := []NetflowPayload{}
	records = spikeFlowPayload(spike_proto)
	data.Header = header
	data.Records = records
	return *data
}

func spikeFlowPayload(spike_proto Proto) []NetflowPayload {
	payload := make([]NetflowPayload, 1)
	switch spike_proto {
	case PROTO_SSH:
		payload[0] = CreateSshFlow()
	case PROTO_FTP:
		payload[0] = CreateFTPFlow()
	case PROTO_HTTP:
		payload[0] = CreateHttpFlow()
	case PROTO_HTTPS:
		payload[0] = CreateHttpsFlow()
	case PROTO_NTP:
		payload[0] = CreateNtpFlow()
	case PROTO_SNMP:
		payload[0] = CreateSnmpFlow()
	case PROTO_IMAPS:
		payload[0] = CreateImapsFlow()
	case PROTO_MYSQL:
		payload[0] = CreateMySqlFlow()
	case PROTO_HTTPS_ALT:
		payload[0] = CreateHttpAltFlow()
	case PROTO_P2P:
		payload[0] = CreateP2pFlow()
	case PROTO_BITTORRENT:
		payload[0] = CreateBitorrentFlow()
	default:
		log.Printf("Invalid protocol option for spike")
		os.Exit(1)
	}
	return payload
}
