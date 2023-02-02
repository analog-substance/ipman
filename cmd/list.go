package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/analog-substance/ipman/internal/fileutil"
	"github.com/analog-substance/ipman/pkg/ip"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [ip|CIDR]",
	Short: "Sort and list IPs",
	Run: func(cmd *cobra.Command, args []string) {
		ipSet := ip.NewSet()

		ips := append([]string{}, args...)
		files, _ := cmd.Flags().GetStringSlice("file")
		for _, f := range files {
			lines, err := fileutil.ReadLines(f)
			checkErr(err)
			ips = append(ips, lines...)
		}

		for _, currentIP := range ips {
			if strings.Contains(currentIP, "/") {
				_, network, err := net.ParseCIDR(currentIP)
				checkErr(err)

				ipSet.AddNetwork(network)
			} else {
				parsed := net.ParseIP(currentIP)
				if parsed != nil {
					ipSet.Add(parsed)
				}
			}
		}

		allIPs := ipSet.SortedSlice()
		for _, ip := range allIPs {
			fmt.Println(ip)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceP("file", "f", []string{}, "File(s) containing IPs to list.")
}
