package adblock

import (
	"bufio"
	"log"
	"net/http"
	"strings"
)

func FetchAdDomains(url string) map[string]struct{} {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch ad domain list from %s: %s\n", url, err.Error())
		return nil
	}
	defer resp.Body.Close()

	adDomainList := make(map[string]struct{})
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "0.0.0.0 ") {
			domain := strings.Fields(line)[1]
			adDomainList[domain] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading response body from %s: %s\n", url, err.Error())
		return nil
	}

	return adDomainList
}
