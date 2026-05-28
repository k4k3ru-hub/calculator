//
// parse.go
//
package decimals

import (
    "fmt"
    "math/big"
    "strings"
)

//
// Decimal parts.
//
type DecimalParts struct {
    Sign     int
    IntPart  string
    FracPart string
}


//
// Parse decimal string to integer.
//
// Note:
//   - Leading / trailing spaces are trimmed.
//   - Leading '+' or '-' is allowed.
//   - Scientific notation is not supported.
//   - Empty integer part is not allowed.
//   - Empty fractional part after dot is not allowed.
//   - Trailing zeros in the fractional part are ignored for decimals validation.
//   - Zero is allowed.
//
// Example:
//   - decimalString="1.23",    decimals=6 -> 1230000
//   - decimalString="1.23",    decimals=2 -> 123
//   - decimalString="1.2300",  decimals=2 -> 123
//   - decimalString="-1.2300", decimals=2 -> -123
//   - decimalString="+1.0000", decimals=0 -> 1
//   - decimalString="0",       decimals=6 -> 0
//   - decimalString="1.2301",  decimals=2 -> error
//   - decimalString="123.",    decimals=2 -> error
//   - decimalString=".45",     decimals=2 -> error
//   - decimalString=".4500",   decimals=2 -> error
//   - decimalString=".123",    decimals=3 -> error
//   - decimalString="-12.34",  decimals=2 -> -1234
//
// Version:
//   - 2026-05-28: Added.
//
func ParseToInteger(decimalString string, decimals uint8) (*big.Int, error) {
    parts, err := ParseDecimalParts(decimalString)
    if err != nil {
        return nil, fmt.Errorf("failed to parse decimal to integer: %w", err)
    }

    fracPart := strings.TrimRight(parts.FracPart, "0")
    if len(fracPart) > int(decimals) {
        return nil, fmt.Errorf("failed to parse decimal to integer: decimal places exceed decimals: decimal=%q decimals=%d", decimalString, decimals)
    }

    fracPart += strings.Repeat("0", int(decimals)-len(fracPart))

    integerString := strings.TrimLeft(parts.IntPart+fracPart, "0")
    if integerString == "" {
        return big.NewInt(0), nil
    }

    n, ok := new(big.Int).SetString(integerString, 10)
    if !ok {
        return nil, fmt.Errorf("failed to parse decimal to integer: decimal=%q", decimalString)
    }

    if parts.Sign < 0 {
        n.Neg(n)
    }

    return n, nil
}


//
// Parse decimal parts.
//
// Note:
//   - Leading / trailing spaces are trimmed.
//   - Leading '+' or '-' is allowed.
//   - Scientific notation is not supported.
//   - Empty integer part is not allowed.
//   - Empty fractional part after dot is not allowed.
//
// Example:
//   - "123"     -> Sign=1,  IntPart="123", FracPart=""
//   - "+123.45" -> Sign=1,  IntPart="123", FracPart="45"
//   - "-0.001"  -> Sign=-1, IntPart="0",   FracPart="001"
//   - ".45"     -> error
//   - "123."    -> error
//   - "1e6"     -> error
//
// Version:
//   - 2026-05-28: Added.
//
func ParseDecimalParts(decimalString string) (*DecimalParts, error) {
    s := strings.TrimSpace(decimalString)
    if s == "" {
        return nil, fmt.Errorf("failed to parse decimal parts: variable=empty")
    }

    sign := 1
    if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
        if s[0] == '-' {
            sign = -1
        }
        s = s[1:]
        if s == "" {
            return nil, fmt.Errorf("failed to parse decimal parts: only sign: variable=%q", decimalString)
        }
    }

    if strings.ContainsAny(s, "eE") {
        return nil, fmt.Errorf("failed to parse decimal parts: scientific notation is not supported: variable=%q", decimalString)
    }

    parts := strings.Split(s, ".")
    if len(parts) > 2 {
        return nil, fmt.Errorf("failed to parse decimal parts: multiple dots: variable=%q", decimalString)
    }

    intPart := parts[0]
    fracPart := ""

    if intPart == "" {
        return nil, fmt.Errorf("failed to parse decimal parts: empty integer part: variable=%q", decimalString)
    }
    if !isDigits(intPart) {
        return nil, fmt.Errorf("failed to parse decimal parts: invalid integer part: variable=%q", decimalString)
    }

    if len(parts) == 2 {
        fracPart = parts[1]
        if fracPart == "" {
            return nil, fmt.Errorf("failed to parse decimal parts: empty fractional part: variable=%q", decimalString)
        }
        if !isDigits(fracPart) {
            return nil, fmt.Errorf("failed to parse decimal parts: invalid fractional part: variable=%q", decimalString)
        }
    }

    return &DecimalParts{
        Sign:     sign,
        IntPart:  intPart,
        FracPart: fracPart,
    }, nil
}



func isDigits(s string) bool {
    if s == "" {
        return false
    }

    for _, ch := range s {
        if ch < '0' || ch > '9' {
            return false
        }
    }

    return true
}






