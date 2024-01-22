package cmd

import (
	"fmt"
	"time"

	ct "github.com/florianl/go-conntrack"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var sleepSec uint

func init() {
	ginCmd.Flags().UintVarP(&sleepSec, "sleepSec", "s", 0, "sleep seconds")
	rootCmd.AddCommand(ginCmd)
}

// ginCmd represents the gin command
var ginCmd = &cobra.Command{
	Use:   "gin",
	Short: "gin demo",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			nfct, err := ct.Open(&ct.Config{})
			if err != nil {
				fmt.Println("Could not create nfct:", err)
				return
			}
			defer nfct.Close()
			sessions, err := nfct.Dump(ct.Conntrack, ct.IPv4)
			if err != nil {
				fmt.Println("Could not dump sessions:", err)
				return
			}
			for _, session := range sessions {
				// fmt.Printf("[%2d] %s - %s\n", session.Origin.Proto.Number, session.Origin.Src, session.Origin.Dst)
				// session.Reply.Proto.SrcPort

				fmt.Printf("[req src %v:%v dst %v:%v], [resp src %v:%v dst %v:%v]\n",
					session.Origin.Src.String(), session.Origin.Proto.SrcPort,
					session.Origin.Dst.String(), session.Origin.Proto.DstPort,
					session.Reply.Src.String(), session.Reply.Proto.SrcPort,
					session.Reply.Dst.String(), session.Reply.Proto.DstPort,
				)

			}

			// Print out all expected sessions.
			// for _, session := range sessions {
			// 	fmt.Printf("%#v\n", session)
			// }

			if sleepSec > 0 {
				time.Sleep(time.Duration(sleepSec) * time.Second)
			}
			c.JSON(200, c.ClientIP())
		})

		return r.Run(":82")
	},
}

func t() {
	// client, err := conntrack.Dial(nil)
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Close()
	//
	// // 创建 conntrack 对象
	// ct, err := client.New()
	// if err != nil {
	// 	panic(err)
	// }
	// defer ct.Close()
}
