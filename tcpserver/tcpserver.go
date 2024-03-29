// go run tcpserver.go
// nc localhost 9000

package main

import (
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

// Node *
type Node struct {
	ID      string
	country string
	region  string
	city    string
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

		//fmt.Println(nodesCountries)
		//fmt.Println(nodesRegions)
		//fmt.Println(nodesCities)
		/*
			io.WriteString(conn, fmt.Sprintln(nodesIds))
			io.WriteString(conn, fmt.Sprintln(nodesCountries))
			io.WriteString(conn, fmt.Sprintln(nodesRegions))
			io.WriteString(conn, fmt.Sprintln(nodesCities))
		*/
		go recommend(conn, node.ID, nodesIds, nodesCountries, nodesRegions, nodesCities)

	}
}

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

// NewNode *
func NewNode(conn net.Conn, IDS []string) *Node {
	node := &Node{
		ID:      calculateID(IDS),
		country: findCou(),
		region:  findReg(),
		city:    findCity(),
	}

	return node
}

func addID(ID string, IDS []string) []string {
	IDS = append(IDS, ID)
	return IDS
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

func calculateID(IDS []string) string {

	nodeID := shortuuid.New()
	if !isIDvalid(nodeID, IDS) {
		go calculateID(IDS)
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

/*
func findLoc() (string, string, string) {
	// Make HTTP GET request
	response, err := http.Get("https://mylocation.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	// Find a substr
	textStartIndex := strings.Index(pageContent, "<table>")
	if textStartIndex == -1 {
		fmt.Println("No text element found")
		os.Exit(0)
	}
	// The start index of the text is the index of the first
	// character, the < symbol. We don't want to include
	// <table> as part of the final value, so let's offset
	// the index by the number of characers in <table>
	textStartIndex += 7

	// Find the index of the closing tag
	textEndIndex := strings.Index(pageContent, "</table>")
	if textEndIndex == -1 {
		fmt.Println("No closing tag for text found.")
		os.Exit(0)
	}

	// (Optional)
	// Copy the substring in to a separate variable so the
	// variables with the full document data can be garbage collected
	pagetext := []byte(pageContent[textStartIndex:textEndIndex])
	s := strings.Split(string(pagetext), "<td>")

	var (
		country = "Unknown"
		region  = "Unknown"
		city    = "Unknown"
	)

	for i := 0; i < len(s); i++ {

		if strings.Contains(s[i], "City") {

			c := strings.Split(string(s[i+1]), "</td>")
			city = c[0]
			break
		}

		if strings.Contains(s[i], "Region") {

			c := strings.Split(string(s[i+1]), "</td>")
			region = c[0]
			continue
		}

		if strings.Contains(s[i], "Country") {

			c := strings.Split(string(s[i+1]), "</td>")
			country = c[0]
			continue
		}

	}

	return country, region, city
}
*/
