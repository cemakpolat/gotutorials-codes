// main.go
package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

type Record struct {
	Name  string
	Type  uint16
	Value string
}

var dnsRecords = map[string][]Record{
	"example.com.": {
		Record{
			Name:  "example.com.",
			Type:  dns.TypeA,
			Value: "192.168.1.100",
		},
		Record{
			Name:  "example.com.",
			Type:  dns.TypeAAAA,
			Value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
	},
	"www.example.com.": {
		Record{
			Name:  "www.example.com.",
			Type:  dns.TypeCNAME,
			Value: "example.com.",
		},
	},
	"mail.example.com.": {
		Record{
			Name:  "mail.example.com.",
			Type:  dns.TypeMX,
			Value: "10 mail.example.com.", // using 10 as priority for MX records
		},
	},
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	for _, question := range r.Question {
		fmt.Printf("Query for %s Type %d\n", question.Name, question.Qtype)
		if records, ok := dnsRecords[question.Name]; ok {
			for _, record := range records {
				if record.Type == question.Qtype {
					switch question.Qtype {
					case dns.TypeA:
						rr, err := dns.NewRR(fmt.Sprintf("%s A %s", question.Name, record.Value))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
					case dns.TypeAAAA:
						rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", question.Name, record.Value))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
					case dns.TypeCNAME:
						rr, err := dns.NewRR(fmt.Sprintf("%s CNAME %s", question.Name, record.Value))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
					case dns.TypeMX:
						rr, err := dns.NewRR(fmt.Sprintf("%s MX %s", question.Name, record.Value))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
					}
				}
			}
		}
	}

	err := w.WriteMsg(m)
	if err != nil {
		log.Println("Error responding to dns query:", err)
	}
}

func main() {
	addr := ":53"
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	defer pc.Close()
	fmt.Println("DNS Server is running at port 53")

	dns.HandleFunc(".", handleDNSRequest)

	err = dns.ActivateAndServe(nil, pc, nil)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
