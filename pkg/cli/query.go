package cli

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

type QueryOptions struct {
	RecordType string
	Verbose    bool
	UseDoH     bool
	Reverse    bool
}

var opts QueryOptions

var queryCmd = &cobra.Command{
	Use:   "query [server] [domain or IP]",
	Short: "Query a DNS server for a domain or perform a reverse lookup",
	Args:  cobra.ExactArgs(2),
	RunE:  runQuery,
}

func init() {
	queryCmd.Flags().StringVarP(&opts.RecordType, "type", "t", "A", "The DNS record type to query (A, AAAA, MX, CNAME, PTR, etc.)")
	queryCmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Enable verbose output")
	queryCmd.Flags().BoolVarP(&opts.UseDoH, "doh", "d", false, "Use DNS over HTTPS (DoH)")
	queryCmd.Flags().BoolVarP(&opts.Reverse, "reverse", "r", false, "Perform a reverse DNS lookup")
	rootCmd.AddCommand(queryCmd)
}

func runQuery(cmd *cobra.Command, args []string) error {
	server := args[0]
	target := args[1]

	if opts.Reverse {
		ip := net.ParseIP(target)
		if ip == nil {
			return fmt.Errorf("invalid IP address for reverse lookup: %s", target)
		}
		target = reverseIP(ip.String()) + ".in-addr.arpa"
		opts.RecordType = "PTR"
	}

	var err error
	if opts.UseDoH {
		err = QueryDoH(server, target, opts)
	} else {
		err = QueryDNS(server, target, opts)
	}

	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	return nil
}

func reverseIP(ip string) string {
	parts := strings.Split(ip, ".")
	for i := 0; i < len(parts)/2; i++ {
		j := len(parts) - 1 - i
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}
