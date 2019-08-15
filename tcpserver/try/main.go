/*
> cd .../try
> make deps
> go build
> ./try
*/

package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func addID(ID string, IDS []string) []string {
	IDS = append(IDS, ID)
	return IDS
}

func main() {
	//Listen for incoming connections
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	//Close the listener when application closes
	defer l.Close()

	var nodesIds []string
	nodesCountries := make(map[string]string)
	nodesRegions := make(map[string]string)
	nodesCities := make(map[string]string)

	fmt.Println("Listening on :9000")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		node := NewNode(conn, nodesIds)

		io.WriteString(conn, fmt.Sprintf("Your node ID is: %s", node.ID))
		io.WriteString(conn, fmt.Sprintf("\nCountry: %s\nRegion: %s\nCity: %s", node.country, node.region, node.city))

		nodesIds = addID(node.ID,nodesIds)
		nodesCountries[node.ID] = node.country
		nodesRegions[node.ID] = node.region
		nodesCities[node.ID] = node.city

		fmt.Println(nodesCountries)
		fmt.Println(nodesRegions)
		fmt.Println(nodesCities)

		io.WriteString(conn, fmt.Sprintln(nodesIds))
		io.WriteString(conn, fmt.Sprintln(nodesCountries))
		io.WriteString(conn, fmt.Sprintln(nodesRegions))
		io.WriteString(conn, fmt.Sprintln(nodesCities))

		go recommend(conn, node.ID, nodesIds, nodesCountries, nodesRegions, nodesCities)

	}
}