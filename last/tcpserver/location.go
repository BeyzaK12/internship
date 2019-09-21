package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

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
