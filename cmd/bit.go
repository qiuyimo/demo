package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Flags uint

const (
	FlagUp           Flags = 1 << iota // interface is up
	FlagBroadcast                      // interface supports broadcast access capability
	FlagLoopback                       // interface is a loopback interface
	FlagPointToPoint                   // interface belongs to a point-to-point link
	FlagMulticast                      // interface supports multicast access capability
)

var flagNames = []string{
	"up",
	"broadcast",
	"loopback",
	"pointtopoint",
	"multicast",
}

func (f Flags) String() string {
	s := ""
	for i, name := range flagNames {
		if f&(1<<uint(i)) != 0 {
			if s != "" {
				s += "|"
			}
			s += name
		}
	}
	if s == "" {
		s = "0"
	}
	return s
}

// bitCmd represents the bit command
var bitCmd = &cobra.Command{
	Use:   "bit",
	Short: "如何使用 bit 来作为多选使用",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("all: %v\n", flagNames)

		var f Flags

		f = 0b00001
		fmt.Println(f) // up

		f = 0b00011
		fmt.Println(f) // up|broadcast

		f = 0b10010
		fmt.Println(f) // broadcast|multicast
	},
}

func init() {
	rootCmd.AddCommand(bitCmd)
}
