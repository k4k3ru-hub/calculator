//
// step.go
//
package decimals

import (
    "math/big"
    "strings"
)


//
// Builds a decimal step expression from the number of fractional digits.
//
// Note:
//   - The returned step always represents 10^-decimals.
//   - This is useful for building quantity / price step metadata without float64.
//
// Example:
//   - decimals=0 -> step="1",        num=1, den=1
//   - decimals=1 -> step="0.1",      num=1, den=10
//   - decimals=2 -> step="0.01",     num=1, den=100
//   - decimals=4 -> step="0.0001",   num=1, den=10000
//
// Returns:
//   - step: Decimal step string.
//   - num:  Numerator of the step fraction.
//   - den:  Denominator of the step fraction.
//
// Version:
//   - 2026-05-28: Added.
//
func BuildStepFromDecimals(decimals uint8) (step string, num *big.Int, den *big.Int, err error) {
	num = big.NewInt(1)
	den = new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	if decimals == 0 {
		return "1", num, den, nil
	}

    step = "0." + strings.Repeat("0", int(decimals)-1) + "1"

	return step, num, den, nil
}
