package main

import (
	"fmt"
	"log"
	"net/http"
	"sysmonitor/monitor"
)

func main() {
	http.HandleFunc("/metrics", monitor.MetricsHandler)
	fmt.Println("Server is starting at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
