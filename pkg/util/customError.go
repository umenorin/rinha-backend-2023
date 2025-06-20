package util

import (
	"fmt"
	"os"
)

func CheckErrorQuery(err error) {

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}
