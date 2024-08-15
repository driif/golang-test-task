package tests_test

import (
	"testing"

	"github.com/driif/golang-test-task/pkg/tests"
	"github.com/stretchr/testify/require"
)

func TestRunningInTest(t *testing.T) {
	require.True(t, tests.RunningInTest(), "Should be true while we are running in the test env/context")
}
