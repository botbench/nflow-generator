# Usage - nflow-generator

This program generates mock netflow (v5) data that can be used to test netflow collector programs. 
The program simulates a router that is exporting flow records to the collector.
It is useful for determining whether the netflow collector is operating and/or receiving netflow datagrams.

nflow-generator generates several netflow datagrams per second, each with 8 or 16 records for varying kinds of traffic (HTTP, SSH, SNMP, DNS, MySQL, and many others.)

### Build

Install [Go](http://golang.org/doc/install), then:

	git clone https://github.com/botbench/nflow-generator.git 
	cd <dir>
	go build

Go build will leave a binary in the root directory that can be run.
	
### Run

Feed it the target collector and port. The use of environment variables allows this to be run as a container.

  COLLECTOR_IP=127.0.0.1 COLLECTOR_PORT=9001 go run .

### REST interface

The nflow_generator can be controlled via a simple REST interface over HTTP

	/api/v1/generate/start - start the netflow packet generation
	/api/v1/generate/stop  - stop the netflow generation
	/api/v1/status         - check if the packet generation is running or not

### Run a Test Collection

You can run a simple test collection using nfcapd from the nfdump package with the following.

- Start a netflow collector

```
sudo apt-get install nfdump
mkdir /tmp/nfcap-test
nfcapd -E  -p 9001 -l /tmp/nfcap-test
```

In a seperate console, run the netflow-generator pointing at an IP on the host the collector is running on (in this case the VM has an IP of 192.168.1.113).

```
./nflow_generator -t 192.168.1.113 -p 9001
```

- You should start seeing records displayed to the output of the screen running nfcapd like the following.

```
$> nfcapd -E  -p 9001 -l /tmp/nfcap-test
Add extension: 2 byte input/output interface index
Add extension: 4 byte input/output interface index
Add extension: 2 byte src/dst AS number
Add extension: 4 byte src/dst AS number
Bound to IPv4 host/IP: any, Port: 9001
Startup.
Init IPFIX: Max number of IPFIX tags: 62

Flow Record:
  Flags        =              0x00 FLOW, Unsampled
  export sysid =                 1
  size         =                56
  first        =        1552592037 [2019-03-14 15:33:57]
  last         =        1552592038 [2019-03-14 15:33:58]
  msec_first   =               973
  msec_last    =               414
  src addr     =      112.10.20.10
  dst addr     =     172.30.190.10
  src port     =                40
  dst port     =                80
  fwd status   =                 0
  tcp flags    =              0x00 ......
  proto        =                 6 TCP
  (src)tos     =                 0
  (in)packets  =               792
  (in)bytes    =                23
  input        =                 0
  output       =                 0
  src as       =             48730
  dst as       =             15401


Flow Record:
  Flags        =              0x00 FLOW, Unsampled
  export sysid =                 1
  size         =                56
  first        =        1552592038 [2019-03-14 15:33:58]
  last         =        1552592038 [2019-03-14 15:33:58]
  msec_first   =               229
  msec_last    =               379
  src addr     =     192.168.20.10
  dst addr     =     202.12.190.10
  src port     =                40
  dst port     =               443
  fwd status   =                 0
  tcp flags    =              0x00 ......
  proto        =                 6 TCP
  (src)tos     =                 0
  (in)packets  =               599
  (in)bytes    =               602
  input        =                 0
  output       =                 0
  src as       =              1115
  dst as       =             50617

```
