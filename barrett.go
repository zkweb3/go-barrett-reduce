package barrett

import (
	"errors"
	"math"
	"math/big"
)

type Barrett struct {
	_p *big.Int
	_u *big.Int
	_k uint
}

/*

INPUT:
 - Non-negative integers `z` and modulus `p`.
 - Radix `b`, the base of `z` and `p` representation.
 - Integer `k` such that `k = ⌊log_b(p)⌋ + 1`.
 - Integer `z` such that `0 ≤ z < b^(2k)`.
 - Precomputed `µ` as `µ = ⌊b^(2k) / p⌋`.

OUTPUT: `z mod p`.

 1. Compute `q̄` as `⌊⌊z / b^(k-1)⌋ * µ / b^(k+1)⌋`.
 2. Compute `r` as `(z mod b^(k+1)) - (q̄ * p mod b^(k+1))`.
 3. If `r < 0` then `r <- r + b^(k+1)`.
 4. While `r ≥ p` do `r <- r - p`.
 5. Return `r`.

*/

func NewBarrett(p *big.Int) *Barrett {
	bf, _, err := big.ParseFloat(p.Text(10), 10, 256, big.ToNearestEven)
	if err != nil {
		return nil
	}
	f, _ := bf.Float64()
	bit := math.Log2(f)
	var k uint
	switch {
	case bit <= 16:
		k = 16
	case bit <= 32:
		k = 32
	case bit <= 64:
		k = 64
	case bit <= 128:
		k = 128
	case bit <= 256:
		k = 256
	}
	u := big.NewInt(1)
	u.Lsh(u, 2*k) // (1ULL << 2k) / p => 2^2k / p
	u.Div(u, p)   // once
	return &Barrett{
		_p: p,
		_u: u,
		_k: k,
	}
}

func (b *Barrett) Reduce(z *big.Int) (*big.Int, error) {
	x := new(big.Int).Rsh(z, 2*b._k)
	if x.Cmp(big.NewInt(1)) >= 0 {
		return nil, errors.New("z < b^(2k)")
	}
	q1 := new(big.Int).Rsh(z, b._k-1)                                     // q1 = z / b^(k-1)
	q2 := new(big.Int).Mul(q1, b._u)                                      // q2 = q1 * u
	q3 := new(big.Int).Rsh(q2, b._k+1)                                    // q3 = q2 / b^(k+1)
	k1 := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(b._k+1)), nil) // k1 = b^(k+1)
	r1 := new(big.Int).Mod(z, k1)                                         // r1 = z % k1
	r2 := new(big.Int).Mod(new(big.Int).Mul(q3, b._p), k1)                // r2 = (q3 * p) % k1
	r := new(big.Int).Sub(r1, r2)                                         // r = r1 - r2
	if r.Cmp(big.NewInt(0)) < 0 {                                         // if r < 0, then r = r + k1
		r.Add(r, k1)
	}
	for r.Cmp(b._p) >= 0 { // while r >= mod, do r = r - p
		r.Sub(r, b._p)
	}
	return r, nil
}
