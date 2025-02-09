Go's networking capabilities make it well-suited for building custom DNS servers, allowing for control over domain name resolution.

Core Functionality:
- Implement a basic DNS server that listens on UDP.
- Parse DNS queries.
- Respond with configured DNS records.

This program implements a DNS server using the library github.com/miekg/dns that listens on UDP port 53.
It defines the DNS records in a map, that represents a Record struct (Name, Type, Value).
The handleDNSRequest function parses the query and responds with the DNS records if they are available in the map.
The main function creates a UDP listener and registers the handleDNSRequest to be called when a new request arrives.

Open a terminal, navigate to the directory containing main.go, and run:

`go run .`

You can use a command line tool such as dig to test the dns server.

`dig @localhost example.com`
