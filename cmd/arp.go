package cmd

import (
	"bytes"
	"demo/kit"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/mdlayher/arp"
	"github.com/mdlayher/ethernet"
	"github.com/spf13/cobra"
)

var (
	vsIP   string
	vsPort uint
	ifName string
	ifIP   string
	dnat   bool
)

func init() {
	rootCmd.AddCommand(arpCmd)
	arpCmd.Flags().StringVarP(&vsIP, "ip", "i", "192.168.7.44", "ip string")
	arpCmd.Flags().UintVarP(&vsPort, "port", "p", 81, "port")
	arpCmd.Flags().StringVarP(&ifName, "ifname", "n", "ens224", "interface name")
	arpCmd.Flags().StringVarP(&ifIP, "ifip", "", "192.168.7.30", "interface ip addr")
	arpCmd.Flags().BoolVarP(&dnat, "dnat", "d", true, "is generate DNAT rule")
}

var arpCmd = &cobra.Command{
	Use:   "arp",
	Short: "arp server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("vsIP:", vsIP, "vsPort:", vsPort)
		fmt.Println("ifName:", ifName, "ifIP:", ifIP)
		fmt.Println("dnat:", dnat)

		if dnat {
			if err := iptablesHandle(vsIP, vsPort, ifIP); err != nil {
				return err
			}
		}

		ifs, err := net.Interfaces()
		if err != nil {
			return fmt.Errorf("net.Interfaces() err: %w", err)
		}
		m := make(map[string]*net.Interface)
		for _, v := range ifs {
			v := v
			m[v.Name] = &v
		}

		ifData, ok := m[ifName]
		if !ok {
			return fmt.Errorf("%v is not existed", ifName)
		}

		fmt.Printf("ifData: %v\n", kit.J(ifData))

		client, err := arp.Dial(ifData)
		if err != nil {
			return fmt.Errorf("arp.Dial(ifData): %w", err)
		}

		for {
			pkt, eth, err := client.Read()
			if err != nil {
				return fmt.Errorf("client.Read err: %w", err)
			}

			if pkt.TargetIP.String() != vsIP {
				continue
			}

			if pkt.Operation != arp.OperationRequest {
				log.Println("not arp request")
				continue
			}

			if !bytes.Equal(eth.Destination, ethernet.Broadcast) && !bytes.Equal(eth.Destination, ifData.HardwareAddr) {
				return fmt.Errorf("mac addr not match")
			}

			if err := client.Reply(pkt, ifData.HardwareAddr, pkt.TargetIP); err != nil {
				return fmt.Errorf("client.Reply err: %w", err)
			}
			log.Printf("apr ok, sourceIP: %v, targetIP: %v, sourceMAC: %v, targetMAC: %v", pkt.SenderIP, pkt.TargetIP, pkt.SenderHardwareAddr, ifData.HardwareAddr)
		}
	},
}

func iptablesHandle(vsIP string, vsPort uint, ifIP string) error {
	// iptables -t nat -A PREROUTING -d 192.168.7.44 -p tcp --dport 81 -j DNAT --to-destination 192.168.7.30
	// iptables -t nat -A PREROUTING -d 192.168.7.44 -j DNAT --to-destination 192.168.7.30
	tables, err := iptables.New()
	if err != nil {
		return err
	}
	// ruleSpec := []string{"-d", vsIP, "-p", "tcp", "--dport", gconv.String(vsPort), "-j", "DNAT", "--to-destination", ifIP}
	ruleSpec := []string{"-d", vsIP, "-j", "DNAT", "--to-destination", ifIP}

	ok, err := tables.Exists("nat", "PREROUTING", ruleSpec...)
	if err != nil {
		return err
	}
	if ok {
		fmt.Println("nat rule is existed")
		return nil
	}
	err = tables.Append("nat", "PREROUTING", ruleSpec...)
	if err != nil {
		return err
	}
	fmt.Printf("add iptables rule: `iptables -t nat -L PREROUTING %v`\n", strings.Join(ruleSpec, " "))

	return nil
}
