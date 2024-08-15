package tests

import (
	"os"
	"strings"
)

// RunningInTest returns true if the current executing is within the go test framework.
// The function first checks the `CI` env variable defined by various CI environments,
// then tests whether the executable ends with the `.test` suffix generated by `go test`.
func RunningInTest() bool {
	// Partially taken from: https://stackoverflow.com/a/45913089 @ 2021-06-02T14:55:01+00:00
	return len(os.Getenv("CI")) > 0 || strings.HasSuffix(os.Args[0], ".test")
}
