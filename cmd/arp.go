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
	arpCmd.Flags().StringVarP(&ipRange, "range", "r", "172.16.98.150-172.16.98.160", "example: 172.16.98.150-172.16.98.160")
	arpCmd.Flags().StringVarP(&ifName, "ifname", "n", "ens160", "interface name")
}

// arpCmd represents the arp command
var arpCmd = &cobra.Command{
	Use:   "arp",
	Short: "创建一个 arp 服务端，会根据配置的 ip 范围，响应 arp 请求",
	Long: `通过此服务，指定 ip 范围，这些 ips 要求是没有使用的。
再给加上 DNAT 规则，即可使用这些 ip 访问了。就像服务器已经设置了这个 ip 一样。

举例：
有一个服务，网卡是：
	root@sag-30:~# ip a show ens224
	4: ens224: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
		link/ether 00:0c:29:65:e7:44 brd ff:ff:ff:ff:ff:ff
		altname enp19s0
		inet 192.168.7.30/24 brd 192.168.7.255 scope global ens224
		   valid_lft forever preferred_lft forever
		inet6 fe80::20c:29ff:fe65:e744/64 scope link
		   valid_lft forever preferred_lft forever
	root@sag-30:~#

这个服务器上部署了 nginx，监听了 0.0.0.0:80，会响应 hello world。
	root@sag-29:~# curl --interface ens224 192.168.7.30
	hello world
	root@sag-29:~#

启动此服务，demo arp --range 192.168.7.40-192.168.7.50 --ifname ens224
这样，此服务就会把 ens224 网卡的 mac 地址，作为 arp 请求的 dst ip 是 192.168.7.40-192.168.7.50 的响应。

再加上一个 iptables DNAT 规则：

iptables -t nat -A PREROUTING -d 192.168.7.44 -p tcp --dport 80 -j DNAT --to-destination 192.168.7.30:80

即可直接使用 192.168.7.44 访问了。
	root@sag-29:~# curl --interface ens224 192.168.7.44
	hello world
	root@sag-29:~#

`,
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

		for {
			pkt, eth, err := client.Read()
			if err != nil {
				log.Fatalf("client.Read err: %v\n", err)
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
