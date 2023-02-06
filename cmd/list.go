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
		v4Only, _ := cmd.Flags().GetBool("ipv4")
		v6Only, _ := cmd.Flags().GetBool("ipv6")
		private, _ := cmd.Flags().GetBool("private")
		public, _ := cmd.Flags().GetBool("public")

		ipSet := ip.NewSet()

		ips := append([]string{}, args...)
		files, _ := cmd.Flags().GetStringSlice("file")
		for _, f := range files {
			lines, err := fileutil.ReadLines(f)
			checkErr(err)
			ips = append(ips, lines...)
		}

		isValid := func(ip net.IP) bool {
			v4 := ip.To4()
			if v4Only && v4 == nil ||
				v6Only && v4 != nil {
				return false
			}

			isPrivate := ip.IsPrivate()
			if private && !isPrivate ||
				public && isPrivate {
				return false
			}

			return true
		}

		for _, currentIP := range ips {
			if strings.Contains(currentIP, "/") {
				_, network, err := net.ParseCIDR(currentIP)
				checkErr(err)

				ipSet.AddNetworkWithFilter(network, isValid)
			} else {
				parsed := net.ParseIP(currentIP)
				if parsed != nil {
					if !isValid(parsed) {
						continue
					}

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
	listCmd.Flags().BoolP("ipv6", "6", false, "List only IPv6 IPs")
	listCmd.Flags().BoolP("ipv4", "4", false, "List only IPv4 IPs")
	listCmd.Flags().BoolP("private", "p", false, "List only private IPs")
	listCmd.Flags().BoolP("public", "P", false, "List only public IPs")
}
