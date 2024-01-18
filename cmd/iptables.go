package cmd

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/spf13/cobra"
)

// iptablesCmd represents the iptables command
var iptablesCmd = &cobra.Command{
	Use:   "iptables",
	Short: "iptables handle",
	RunE: func(cmd *cobra.Command, args []string) error {

		tables, err := iptables.New()
		if err != nil {
			return err
		}
		rules, err := tables.List("nat", "PREROUTING")
		if err != nil {
			return err
		}
		fmt.Println("table: nat, chain: PREROUTING, rules count:", len(rules))
		for _, r := range rules {
			fmt.Println("rule:", r)
		}

		// chains, err := tables.ListChains("nat")
		// if err != nil {
		// 	return err
		// }
		// for _, chain := range chains {
		// 	fmt.Println("chain:", chain)
		// 	rules, err := tables.List("nat", chain)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	for _, r := range rules {
		// 		fmt.Println("rule:", r)
		// 	}
		// }

		err = tables.Insert("nat", "PREROUTING", 1, "-p", "udp", "--dport", "80", "-j", "REDIRECT", "--to-ports", "8080")
		if err != nil {
			return err
		}
		err = tables.Append("nat", "PREROUTING", "-p", "udp", "--dport", "81", "-j", "REDIRECT", "--to-ports", "8081")
		if err != nil {
			return err
		}
		ok, err := tables.Exists("nat", "PREROUTING", "-p", "udp", "--dport", "80", "-j", "REDIRECT", "--to-ports", "8080")
		if err != nil {
			return err
		}
		fmt.Println("exist:", ok)
		ok, err = tables.Exists("nat", "PREROUTING", "-p", "udp", "--dport", "90", "-j", "REDIRECT", "--to-ports", "8090")
		if err != nil {
			return err
		}
		fmt.Println("exist:", ok)
		err = tables.Delete("nat", "PREROUTING", "-p", "udp", "--dport", "80", "-j", "REDIRECT", "--to-ports", "8080")
		if err != nil {
			return err
		}
		ok, err = tables.Exists("nat", "PREROUTING", "-p", "udp", "--dport", "80", "-j", "REDIRECT", "--to-ports", "8080")
		if err != nil {
			return err
		}
		fmt.Println("exist:", ok)
		err = tables.ClearChain("nat", "PREROUTING")
		if err != nil {
			return err
		}
		ok, err = tables.Exists("nat", "PREROUTING", "-p", "udp", "--dport", "80", "-j", "REDIRECT", "--to-ports", "8008")
		if err != nil {
			return err
		}
		fmt.Println("exist:", ok)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(iptablesCmd)
}
