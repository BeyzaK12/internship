package main

import (
	"fmt"
	"io"
	"net"
)

func recommend(conn net.Conn, ID string, IDS []string, countries map[string]string, regions map[string]string, cities map[string]string) {

	/*
		io.WriteString(conn, fmt.Sprintln(IDS))
		io.WriteString(conn, fmt.Sprintln(countries))
		io.WriteString(conn, fmt.Sprintln(regions))
		io.WriteString(conn, fmt.Sprintln(cities))
	*/

	if len(IDS) > 5 {

		for _, opp := range IDS {

			if cities[ID] == cities[opp] && ID != opp {
				io.WriteString(conn, fmt.Sprintf("\nYou are soo close to %s!", opp))
				continue
			} else if regions[ID] == regions[opp] && ID != opp {
				io.WriteString(conn, fmt.Sprintf("\nYou are near to %s.", opp))
				continue
			} else if countries[ID] == countries[opp] && ID != opp {
				io.WriteString(conn, fmt.Sprintf("\nAt least you are in the same country with %s.", opp))
				continue
			}
		}
	}
}
