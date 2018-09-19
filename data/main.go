package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	// Cmd line flags for external data sources
	ddnsFile := flag.String("dyndns", "dynamic-dns.txt", "File location for the list of Dynamic DNS providers")
	majesticFile := flag.String("majestic", "majestic_million.csv", "File location for the Majestic Million CSV")

	// Parse in the cmd-line flags
	flag.Parse()

	// Open the Dynamic DNS File
	ddns, err := os.Open(*ddnsFile)
	Check(err)
	defer ddns.Close()

	// Open the majestic million
	majestic, err := os.Open(*majesticFile)
	Check(err)

	// Open an appendonly file handle to populate
	aofFile, err := os.Create("appendonly.aof")
	Check(err)

	// Commit the DynamicDNS providers to the appendonly file
	DynamicDNSToRedis(aofFile, ddns)

	// Commit the Majestic Million to the appendonly file
	MajesticToRedis(aofFile, majestic)

}

// Check is an error check wrapper
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// GenerateRedisProto produces the Redis appendonly
// file formatted using the Redis Protocol
func GenerateRedisProto(cmd []string) []byte {
	proto := fmt.Sprintf("*%d\r\n", len(cmd))
	for c := range cmd {
		proto += fmt.Sprintf("$%d\r\n", len(cmd[c]))
		proto += fmt.Sprintf("%s\r\n", cmd[c])
	}

	return []byte(proto)
}

// DynamicDNSToRedis takes a byte stream of DynamicDNS
// providers and writes out the aof file formatted
// to be read by Redis on startup
func DynamicDNSToRedis(aof *os.File, dyndns io.Reader) {

	scanner := bufio.NewScanner(dyndns)
	for scanner.Scan() {
		domain := scanner.Text()
		domain = strings.Split(domain, "#")[0]
		aof.Write(GenerateRedisProto([]string{"sadd", "dynamicdns", domain}))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// MajesticToRedis takes in the Majestic Million CSV
// and writes each domain out to the aof file formatted to
// be read by Redis
func MajesticToRedis(aof *os.File, majesticFile *os.File) {

	r := csv.NewReader(majesticFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		domain := record[2]
		aof.Write(GenerateRedisProto([]string{"sadd", "majestic", domain}))
	}
}
