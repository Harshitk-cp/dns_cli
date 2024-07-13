package cli

import (
	"github.com/spf13/cobra"
)

var recordType string
var verbose bool
var useDoH bool

var queryCmd = &cobra.Command{
	Use:   "query [server] [domain]",
	Short: "Query a DNS server for a domain",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		domain := args[1]
		if useDoH {
			QueryDoH(server, domain, recordType, verbose)
		} else {
			QueryDNS(server, domain, recordType, verbose)
		}
	},
}

func init() {
	queryCmd.Flags().StringVarP(&recordType, "type", "t", "A", "The DNS record type to query (A, AAAA, MX, CNAME, etc.)")
	queryCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	queryCmd.Flags().BoolVarP(&useDoH, "doh", "d", false, "Use DNS over HTTPS (DoH)")
	rootCmd.AddCommand(queryCmd)
}
