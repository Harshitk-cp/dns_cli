package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dns-query-cli",
	Short: "A simple DNS query tool",
	Long:  `A simple DNS query tool built using Cobra CLI in Go.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var recordType string

var queryCmd = &cobra.Command{
	Use:   "query [server] [domain]",
	Short: "Query a DNS server for a domain",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		domain := args[1]
		queryDNS(server, domain, recordType)
	},
}

var verbose bool

func init() {
	queryCmd.Flags().StringVarP(&recordType, "type", "t", "A", "The DNS record type to query (A, AAAA, MX, CNAME, etc.)")
	queryCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.AddCommand(queryCmd)
}

func queryDNS(server, domain, recordType string) {
	c := new(dns.Client)
	m := new(dns.Msg)
	rt, ok := dns.StringToType[recordType]
	if !ok {
		log.Fatalf("Unknown record type: %s\n", recordType)
	}
	m.SetQuestion(dns.Fqdn(domain), rt)
	r, _, err := c.Exchange(m, server+":53")
	if err != nil {
		log.Fatalf("Failed to query DNS: %v\n", err)
	}

	if verbose {
		fmt.Println(r.String())
	} else {
		if len(r.Answer) == 0 {
			fmt.Println("No records found")
			return
		}

		for _, ans := range r.Answer {
			fmt.Println(ans.String())
		}
	}
}
