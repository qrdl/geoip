package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

type geoRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
}

func usage(name string) {
	fmt.Fprintf(os.Stderr, "Usage: %s <geoip2 DB file> <column number with IP address>\n", name)
}

func main() {
	if len(os.Args) < 3 {
		usage(os.Args[0])
		os.Exit(1)
	}

	colnum, err := strconv.Atoi(os.Args[2])
	if err != nil {
		usage(os.Args[0])
		os.Exit(1)
	}

	db, err := maxminddb.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open GeoIP2 db: %v", err)
		usage(os.Args[0])
		os.Exit(1)
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < colnum {
			continue // skip malformed line
		}
		ip := net.ParseIP(fields[colnum-1])
		if ip == nil {
			continue // skip invalid IP address
		}

		var rec geoRecord
		var city, country string
		err = db.Lookup(ip, &rec)
		if err != nil {
			city = "-"
			country = "-"
		} else {
			country = rec.Country.ISOCode
			if len(country) == 0 {
				country = "-"
			} else {
				country = strings.ReplaceAll(country, " ", "_")
			}
			city = rec.City.Names["en"]
			if len(city) == 0 {
				city = "-"
			} else {
				city = strings.ReplaceAll(city, " ", "_")
			}
		}
		fmt.Printf("%s %s %s\n", country, city, line)
	}

	if err = scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v", err)
	}
}
