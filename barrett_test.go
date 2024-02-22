package barrett

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	N string = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
)

func TestBarrettReduce(t *testing.T) {
	n, ok := new(big.Int).SetString(N, 16)
	require.True(t, ok)
	x := big.NewInt(math.MaxInt64)
	x.Exp(x, big.NewInt(8), nil)

	b := NewBarrett(n)
	require.NotNil(t, b)
	result, err := b.Reduce(x)
	require.Nil(t, err)

	// verify result
	correct := big.NewInt(math.MaxInt64)
	correct.Exp(correct, big.NewInt(8), n)
	require.Equal(t, result.Cmp(correct), 0)

	// print result
	fmt.Println("result", result.Text(16))
}
