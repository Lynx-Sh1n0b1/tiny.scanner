package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
)

func scanPort(target string, port int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", target, port))
	if err == nil {
		conn.Close()
		results <- fmt.Sprintf("%s %d open", target, port)
	}
}

func printSimpleBanner(text string) {
	fmt.Println(strings.Repeat("*", len(text)+4))
	fmt.Printf("* %s *\n", text)
	fmt.Println(strings.Repeat("*", len(text)+4))
}

func printHelp() {
	fmt.Println("Usage: tiny.scanner -t <target> -p <ports> [-c <csv_file>]")
	flag.PrintDefaults()
}

func loadIPsFromFile(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error loading IPs from file:", err)
		return nil
	}
	return strings.Split(string(data), "\n")
}

func expandCIDR(cidr string) ([]string, error) {
	ips := []string{}
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {
	printSimpleBanner("tiny.scanner")

	var ipsArg, portsArg, csvFileArg string
	flag.StringVar(&ipsArg, "t", "", "Target IPs (comma-separated), CIDR notation, or path to a text file with IPs per line")
	flag.StringVar(&portsArg, "p", "", "Port(s) to scan (e.g., 80,443 or 1-1024)")
	flag.StringVar(&csvFileArg, "c", "", "CSV file to save results")
	flag.Parse()

	if ipsArg == "" || portsArg == "" {
		printHelp()
		return
	}

	ips := parseIPInput(ipsArg)
	ports := parsePortInput(portsArg)
	var wg sync.WaitGroup
	results := make(chan string)

	for _, ip := range ips {
		for _, port := range ports {
			wg.Add(1)
			go scanPort(ip, port, results, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var resultStrings []string
	for result := range results {
		resultStrings = append(resultStrings, result)
	}

	if csvFileArg != "" {
		saveResultsToCSV(csvFileArg, resultStrings)
	} else {
		for _, result := range resultStrings {
			fmt.Println(result)
		}
	}
}

func parseIPInput(input string) []string {
	if strings.Contains(input, ".txt") {
		return loadIPsFromFile(input)
	} else if strings.Contains(input, "/") {
		ips, err := expandCIDR(input)
		if err != nil {
			fmt.Println("Error expanding CIDR:", err)
			os.Exit(1)
		}
		return ips
	}
	parts := strings.Split(input, ",")
	return parts
}

func parsePortInput(input string) []int {
	var ports []int
	parts := strings.Split(input, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			start := parseInt(rangeParts[0])
			end := parseInt(rangeParts[1])
			for port := start; port <= end; port++ {
				ports = append(ports, port)
			}
		} else {
			port := parseInt(part)
			ports = append(ports, port)
		}
	}

	return ports
}

func parseInt(s string) int {
	var val int
	_, err := fmt.Sscanf(s, "%d", &val)
	if err != nil {
		fmt.Printf("Error parsing port: %s\n", s)
		os.Exit(1)
	}
	return val
}

func saveResultsToCSV(filename string, results []string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	for _, result := range results {
		_, err := file.WriteString(result + "\n")
		if err != nil {
			fmt.Println("Error writing to CSV file:", err)
			return
		}
	}

	fmt.Println("Results saved to", filename)
}
