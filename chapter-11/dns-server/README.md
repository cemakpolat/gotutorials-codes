## Implementing a simple DNS Server
Go's networking capabilities make it well-suited for building custom DNS servers, allowing for control over domian name resolution. This project shows the implementation of a DNS server using the library github.com/miekg/dns that listens on UDP port 53. It defines the DNS records in a map, that represents a `Record` struct (Name, Type, Value). The handleDNSRequest function parses the query and responds with the DNS records if they are available in the map. The main function creates a UDP listener and registers the handleDNSRequest to be called when a new request arrives.

Core Functionality:
- Implement a basic DNS server that listens on UDP.
- Parse DNS queries.
- Respond with configured DNS records.

## How to test

Open a terminal, navigate to the directory containing `main.go`, and run:

`go run .`

You can use a command line tool such as dig, assuming that it is installed, to test the local dns server.

`dig @localhost example.com`
