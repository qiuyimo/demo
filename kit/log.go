package kit

import (
	"encoding/json"
	"log"
)

func J(v any) string {
	res, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("json.Marshal(ifs) err: %v", err)
	}
	return string(res)
}
