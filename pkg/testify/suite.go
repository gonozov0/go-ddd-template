package testify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Suite interface {
	Require() *require.Assertions
	T() *testing.T
	Run(name string, subtest func()) bool
}
