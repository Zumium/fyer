package bin

import (
	"fmt"
	"os"
)

//ReportErr prints error message to stderr
func ReportErr(err error) {
	fmt.Fprintf(os.Stderr, "error happened: %s\n", err)
}
