package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/fatih/color"
)

func main() {
	// Get URL base from environment variable or use default value
	color.Red("\nWELCOME TO WEB-STRESSER")
	urlBase := os.Getenv("URL_BASE")
	if urlBase == "" {
		fmt.Println("Please introduce the web to test as ENV like URL_BASE=www.example.com")
	}

	// Get interval from environment variable or use default value
	interval := os.Getenv("INTERVAL")
	if interval == "" {
		interval = "5000"
	}

	color.Yellow("\nAttacking %v every %v milliseconds", urlBase, interval)
	// Convert interval to an integer value
	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		fmt.Println("Error converting interval to integer:", err)
		return
	}

	// Perform an HTTP GET request to the URL base
	response, err := http.Get(urlBase)
	if err != nil {
		fmt.Println("Error performing HTTP request:", err)
		return
	}

	// Parse the HTML content of the webpage
	page, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println("Error parsing HTML content:", err)
		return
	}

	// Extract all links from the webpage
	var links []string
	links = extractLinks(page, links)

	// Filter out links that point to non-endpoints
	var endpoints []string
	for _, link := range links {
		if strings.HasPrefix(link, urlBase) {
			endpoints = append(endpoints, link)
		}
	}
	if len(endpoints) == 0 {
		fmt.Printf("No endpoints found on %v to atack", urlBase)
		return
	}
	// Generate a list of IP addresses to use for the "X-Forwarded-For" header
	var ipAddresses []string
	for i := 0; i < 10; i++ {
		ipAddress := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
		ipAddresses = append(ipAddresses, ipAddress)
	}

	for {
		// Get a random IP address from the list
		ipAddress := ipAddresses[rand.Intn(len(ipAddresses))]

		// Perform an HTTP GET request to the URL base with the custom header
		req, err := http.NewRequest("GET", urlBase, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request:", err)
			continue
		}
		req.Header.Set("X-Forwarded-For", ipAddress)
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Error performing HTTP request:", err)
		} else {
			// Print the HTTP response status code
			fmt.Println("HTTP response status code:", response.StatusCode)
		}

		// Wait for the specified interval
		time.Sleep(time.Duration(intervalInt) * time.Millisecond)
	}
}

// Function that extracts all links from an HTML page
func extractLinks(node *html.Node, links []string) []string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		links = extractLinks(child, links)
	}
	return links
}
