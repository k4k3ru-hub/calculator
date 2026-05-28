//
// calculator.go
//
package calculator

import (
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
//   - decimals=0 -> step="1",      num=1, den=1
//   - decimals=1 -> step="0.1",    num=1, den=10
//   - decimals=2 -> step="0.01",   num=1, den=100
//   - decimals=4 -> step="0.0001", num=1, den=10000
//
// Returns:
//   - step: Decimal step string.
//   - num:  Numerator of the step fraction.
//   - den:  Denominator of the step fraction.
//
// Version:
//   - 2026-04-15: Added.
//
func BuildStepFromDecimals(decimals uint8) (step string, num uint64, den uint64, err error) {
    // Calculate denominator: 10^decimals.
    den = 1
    for i := uint8(0); i < decimals; i++ {
        if den > ^uint64(0)/10 {
            return "", 0, 0, fmt.Errorf("failed to build step from decimals: decimals overflow: decimals=%d", decimals)
        }
        den *= 10
    }

    // Integer step.
    if decimals == 0 {
        return "1", 1, 1, nil
    }

    // Decimal step string.
    // Example:
    //   - decimals=1 -> "0.1"
    //   - decimals=4 -> "0.0001"
    step = "0." + strings.Repeat("0", int(decimals)-1) + "1"

    return step, 1, den, nil
}


//
// Counts the number of digits in the fractional part of a decimal string.
//
// Return:
//   - The number of fraction digits in s.
//   - The number of effective fraction digits in s, excluding trailing zeros.
//   - An error if s is empty or not a valid decimal number.
//
// Example:
//   - "123":      0, 0, nil
//   - "123.45":   2, 2, nil
//   - "123.4500": 4, 2, nil
//   - "123.000":  3, 0, nil
//   - "123.":     0, 0, nil
//   - ".45":      2, 2, nil
//   - ".4500":    4, 2, nil
//   - "-12.300":  3, 1, nil
//   - "+12":      0, 0, nil
//   - ".":        0, 0, error
//
// Version:
//   - 2026-04-15: Added.
//
func CountFractionDigits(s string) (int, int, error) {
    if s == "" {
        return 0, 0, fmt.Errorf("failed to count fraction digits: variable=empty")
    }

    // Check leading +/-.
    start := 0
    if s[0] == '+' || s[0] == '-' {
        start = 1
        if len(s) == 1 {
            return 0, 0, fmt.Errorf("failed to count fraction digits: only sign: variable=%q", s)
        }
    }

    dotCount := 0
	seenDot := false
	digitCount := 0

	fracDigits := 0
	trailingZeroCount := 0

    for i := start; i < len(s); i++ {
        c := s[i]

        switch {
        case c >= '0' && c <= '9':
            digitCount++

            if seenDot {
				fracDigits++
				if c == '0' {
					trailingZeroCount++
				} else {
					trailingZeroCount = 0
				}
			}
        case c == '.':
            dotCount++
            if dotCount > 1 {
                return 0, 0, fmt.Errorf("failed to count fraction digits: multiple dots: variable=%q", s)
            }
            seenDot = true
        default:
            return 0, 0, fmt.Errorf("failed to count fraction digits: invalid character: variable=%q", s)
        }
    }

    if digitCount == 0 {
        return 0, 0, fmt.Errorf("failed to count fraction digits: no digits: variable=%q", s)
    }

    effectiveFracDigits := fracDigits - trailingZeroCount
	if effectiveFracDigits < 0 {
		return 0, 0, fmt.Errorf("failed to count fraction digits: invalid value: variable=%q", s)
	}

	return fracDigits, effectiveFracDigits, nil
}


//
// Counts the number of digits in the integer part of a decimal string.
//
// Note:
//   - A leading '+' or '-' is allowed.
//   - Scientific notation is not supported.
//   - Leading zeros in the integer part are ignored.
//   - If the integer part is zero, returns 1.
//
// Example:
//   - "2316.7"   -> 4
//   - "10000"    -> 5
//   - "00012.34" -> 2
//   - "0.001"    -> 1
//   - ".123"     -> 1
//   - "-12.34"   -> 2
//
// Version:
//   - 2026-04-15: Added.
//
func CountIntegerDigits(s string) (int, error) {
    if s == "" {
        return 0, fmt.Errorf("failed to count integer digits: variable=empty")
    }

    // Check leading +/-.
    start := 0
    if s[0] == '+' || s[0] == '-' {
        start = 1
        if len(s) == 1 {
            return 0, fmt.Errorf("failed to count integer digits: only sign: variable=%q", s)
        }
    }

    dotCount := 0
    seenDot := false
    digitCount := 0
    intDigits := 0
    nonZeroSeen := false

    for i := start; i < len(s); i++ {
		c := s[i]

		switch {
		case c >= '0' && c <= '9':
			digitCount++

            if !seenDot {
                if c != '0' {
					nonZeroSeen = true
				}
				if nonZeroSeen {
					intDigits++
				}
			}
		case c == '.':
			dotCount++
			if dotCount > 1 {
				return 0, fmt.Errorf("failed to count integer digits: multiple dots: variable=%q", s)
			}
			seenDot = true
		default:
			return 0, fmt.Errorf("failed to count integer digits: invalid character: variable=%q", s)
		}
	}

	if digitCount == 0 {
		return 0, fmt.Errorf("failed to count integer digits: no digits: variable=%q", s)
	}

    // Integer part is zero, such as ".123", "0.123", or "000.1".
    if intDigits == 0 {
		return 1, nil
	}

    return intDigits, nil
}
