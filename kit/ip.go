package kit

import (
	"strconv"
	"strings"
)

// IsContainIPWithRange Determine whether the IP is within the ip range
func IsContainIPWithRange(ip, ipScope string) bool {
	ipSlice := strings.Split(ipScope, `-`)
	if len(ipSlice) != 2 {
		return false
	}
	return IP2Int(ip) >= IP2Int(ipSlice[0]) && IP2Int(ip) <= IP2Int(ipSlice[1])
}

func IP2Int(ip string) int64 {
	if len(ip) == 0 {
		return 0
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return 0
	}
	b0 := string2Int(bits[0])
	b1 := string2Int(bits[1])
	b2 := string2Int(bits[2])
	b3 := string2Int(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func string2Int(in string) (out int) {
	out, _ = strconv.Atoi(in)
	return
}
