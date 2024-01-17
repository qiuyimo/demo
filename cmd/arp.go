package cmd

import (
	"bytes"
	"demo/kit"
	"fmt"
	"github.com/mdlayher/arp"
	"github.com/mdlayher/ethernet"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var ipRange string
var ifName string

func init() {
	rootCmd.AddCommand(arpCmd)
	arpCmd.Flags().StringVarP(&ipRange, "range", "r", "172.16.98.150-172.16.98.160", "，example: 172.16.98.150-172.16.98.160")
	arpCmd.Flags().StringVarP(&ifName, "ifname", "n", "ens160", "interface name")
}

// arpCmd represents the arp command
var arpCmd = &cobra.Command{
	Use:   "arp",
	Short: "创建一个 arp 服务端，会根据配置的 ip 范围，响应 arp 请求",
	Run: func(cmd *cobra.Command, args []string) {
		ifs, err := net.Interfaces()
		if err != nil {
			log.Fatalf("net.Interfaces() err: %v", err)
		}
		m := make(map[string]*net.Interface)
		for _, v := range ifs {
			v := v
			m[v.Name] = &v
		}

		ifData, ok := m[ifName]
		if !ok {
			log.Fatalf("%v is not existed", ifName)
		}

		fmt.Printf("ifData: %v\n", kit.J(ifData))

		client, err := arp.Dial(ifData)
		if err != nil {
			log.Fatalf("arp.Dial(ifData): %v", err)
		}

		var readCount uint

		for {
			pkt, eth, err := client.Read()
			if err != nil {
				log.Fatalf("client.Read err: %v\n", err)
			}
			readCount++
			if readCount%10 == 0 {
				fmt.Println("accept arp request count:", readCount)
			}

			if !kit.IsContainIPWithRange(pkt.TargetIP.String(), ipRange) {
				continue
			}

			if pkt.Operation != arp.OperationRequest {
				log.Println("not arp request")
				continue
			}

			if !bytes.Equal(eth.Destination, ethernet.Broadcast) && !bytes.Equal(eth.Destination, ifData.HardwareAddr) {
				log.Fatal("mac addr not match")
			}

			if err := client.Reply(pkt, ifData.HardwareAddr, pkt.TargetIP); err != nil {
				log.Fatalf("client.Reply err: %v\n", err)
			}
			log.Printf("apr ok, sourceIP: %v, targetIP: %v, sourceMAC: %v, targetMAC: %v", pkt.SenderIP, pkt.TargetIP, pkt.SenderHardwareAddr, ifData.HardwareAddr)
		}
	},
}
