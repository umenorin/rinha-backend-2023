package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

)


func main() {
	

	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")

	} else if err != nil {
		fmt.Printf("error startin server: %s\n", err)
		os.Exit(1)
	}
}
