// go run tcpserver.go
// nc localhost 9000

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/lithammer/shortuuid"
)

var nodesIds []string

// Node *
type Node struct {
	ID      string
	country string
	region  string
	city    string
}

// NewNode *
func NewNode() *Node {
	node := &Node{
		ID:      calculateID(),
		country: findCou(),
		region:  findReg(),
		city:    findCity(),
	}

	return node
}
func calculateID() string {

	nodeID := shortuuid.New()
	if !isIDvalid(nodeID, nodesIds) {
		go calculateID()
	}
	return nodeID
}

func isIDvalid(ID string, IDS []string) bool {
	i := 0
	for _, id := range IDS {
		if ID != id {
			if i == len(IDS)-1 {
				return true
			}
			continue
		}
		i++
		break
	}
	return false
}

// handleConnection *
func handleConnection(c net.Conn, countries map[string]string, regions map[string]string, cities map[string]string) {

	c.Write([]byte("Please enter the key >>> "))

	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	key := strings.TrimSpace(string(netData))
	if key != "hello" {
		c.Write([]byte("This is not the key!\nPress enter"))
		c.Close()
		return
	}

	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	node := NewNode()

	c.Write([]byte("Your node ID is: " + node.ID))
	c.Write([]byte("\nCountry: " + node.country + "\nRegion: " + node.region + "\nCity: " + node.city))

	nodesIds = append(nodesIds, node.ID)
	countries[node.ID] = node.country
	regions[node.ID] = node.region
	cities[node.ID] = node.city

	go recommend(c, node.ID, countries, regions, cities)

	for {

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}

	}
	c.Close()
}

func recommend(c net.Conn, ID string, countries map[string]string, regions map[string]string, cities map[string]string) {

	if len(nodesIds) > 5 {

		io.WriteString(c, fmt.Sprintln(nodesIds))
		io.WriteString(c, fmt.Sprintln(countries))
		io.WriteString(c, fmt.Sprintln(regions))
		io.WriteString(c, fmt.Sprintln(cities))

		for _, opp := range nodesIds {

			if cities[ID] == cities[opp] && ID != opp {
				c.Write([]byte("\nYou are soo close to " + opp))
				continue
			} else if regions[ID] == regions[opp] && ID != opp {
				c.Write([]byte("\nYou are near to " + opp))
				continue
			} else if countries[ID] == countries[opp] && ID != opp {
				c.Write([]byte("\nAt least you are in the same country with " + opp))
				continue
			}
		}
	}
}

func main() {

	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	nodesCountries := make(map[string]string)
	nodesRegions := make(map[string]string)
	nodesCities := make(map[string]string)

	fmt.Println("Listening on :9000")
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c, nodesCountries, nodesRegions, nodesCities)
	}
}

func findCou() string {
	response, err := http.Get("https://mylocation.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	textStartIndex := strings.Index(pageContent, "<table>")
	if textStartIndex == -1 {
		fmt.Println("No text element found")
		os.Exit(0)
	}

	textStartIndex += 7

	textEndIndex := strings.Index(pageContent, "</table>")
	if textEndIndex == -1 {
		fmt.Println("No closing tag for text found.")
		os.Exit(0)
	}

	pagetext := []byte(pageContent[textStartIndex:textEndIndex])
	s := strings.Split(string(pagetext), "<td>")

	var country = "Unknown"

	for i := 0; i < len(s); i++ {

		if strings.Contains(s[i], "Country") {

			c := strings.Split(string(s[i+1]), "</td>")
			country = c[0]
			break
		}

	}

	return country
}

func findReg() string {
	response, err := http.Get("https://mylocation.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	textStartIndex := strings.Index(pageContent, "<table>")
	if textStartIndex == -1 {
		fmt.Println("No text element found")
		os.Exit(0)
	}

	textStartIndex += 7

	textEndIndex := strings.Index(pageContent, "</table>")
	if textEndIndex == -1 {
		fmt.Println("No closing tag for text found.")
		os.Exit(0)
	}

	pagetext := []byte(pageContent[textStartIndex:textEndIndex])
	s := strings.Split(string(pagetext), "<td>")

	var region = "Unknown"

	for i := 0; i < len(s); i++ {

		if strings.Contains(s[i], "Region") {

			c := strings.Split(string(s[i+1]), "</td>")
			region = c[0]
			break
		}

	}

	return region
}

func findCity() string {
	response, err := http.Get("https://mylocation.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	textStartIndex := strings.Index(pageContent, "<table>")
	if textStartIndex == -1 {
		fmt.Println("No text element found")
		os.Exit(0)
	}

	textStartIndex += 7

	textEndIndex := strings.Index(pageContent, "</table>")
	if textEndIndex == -1 {
		fmt.Println("No closing tag for text found.")
		os.Exit(0)
	}

	pagetext := []byte(pageContent[textStartIndex:textEndIndex])
	s := strings.Split(string(pagetext), "<td>")

	var city = "Unknown"

	for i := 0; i < len(s); i++ {

		if strings.Contains(s[i], "City") {

			c := strings.Split(string(s[i+1]), "</td>")
			city = c[0]
			break
		}

	}

	return city
}
