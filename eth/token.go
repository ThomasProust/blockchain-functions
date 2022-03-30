package eth

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/SoteriaTech/blockchain-functions/env"
	"golang.org/x/crypto/sha3"
)

func parseTokenTransfer(tx *Transaction, c *env.CurrencyConfig) (*Transaction, error) {
	receipt, err := ethService.api.GetReceipt(tx.Hash)
	if err != nil {
		return nil, err
	}
	for idx, l := range receipt.Logs {
		if l.Topics[0].String() == c.TransferSignature {
			r, _ := hex.DecodeString(l.Topics[2].String()[2:])
			tx.Receiver = parseDataAddr(r)
			tx.Value = parseDataValue(l.Data.String())
			tx.LogIdx = strconv.Itoa(idx)
		}
	}
	tx.Currency = c.Name
	return tx, nil
}

func parseDataValue(hex string) *big.Int {
	// v, _ := strconv.ParseUint(hex[2:], 16, 64)
	// log.Fatal(v)
	// return new(big.Int).SetUint64(v)
	v, err := decodeBig(hex)
	if err != nil {
		log.Fatal("err", err)
	}
	return v
}

// from ethereum-go library  https://github.com/ethereum/go-ethereum/blob/991384a7f6719e1125ca0be7fb27d0c4d1c5d2d3/common/types.go#L235
// That's how they transform the hex encoded data in the topics into an EIP55 compliant address
func parseDataAddr(b []byte) string {
	var a [20]byte
	copy(a[:], b[len(b)-20:])
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return string(buf[:])
}

/**
* code coming from https://github.com/ethereum/go-ethereum/blob/991384a7f6719e1125ca0be7fb27d0c4d1c5d2d3/common/hexutil/hexutil.go#L139
* allows to handle cases where value of log.Data exceeds a regular Uint64 by directly dealing with big.Int
* logic is modified to bypass verification of trailing 0s
 */
const badNibble = ^uint64(0)

func decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return badNibble
	}
}

var bigWordNibbles int

// DecodeBig decodes a hex string with 0x prefix as a quantity.
// Numbers larger than 256 bits are not accepted.
func decodeBig(input string) (*big.Int, error) {

	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		bigWordNibbles = 16
	case 2:
		bigWordNibbles = 8
	default:
		panic("weird big.Word size")
	}
	raw, err := checkNumber(input)
	if err != nil {
		return nil, err
	}
	if len(raw) > 64 {
		return nil, ErrBig256Range
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == badNibble {
				return nil, ErrSyntax
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	dec := new(big.Int).SetBits(words)
	return dec, nil
}

// Errors
type decError struct{ msg string }

func (err decError) Error() string { return err.msg }

const uintBits = 32 << (uint64(^uint(0)) >> 63)

var (
	ErrEmptyString   = &decError{"empty hex string"}
	ErrSyntax        = &decError{"invalid hex string"}
	ErrMissingPrefix = &decError{"hex string without 0x prefix"}
	ErrOddLength     = &decError{"hex string of odd length"}
	ErrEmptyNumber   = &decError{"hex string \"0x\""}
	ErrLeadingZero   = &decError{"hex number with leading zero digits"}
	ErrUint64Range   = &decError{"hex number > 64 bits"}
	ErrUintRange     = &decError{fmt.Sprintf("hex number > %d bits", uintBits)}
	ErrBig256Range   = &decError{"hex number > 256 bits"}
)

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func checkNumber(input string) (raw string, err error) {
	if len(input) == 0 {
		return "", ErrEmptyString
	}
	if !has0xPrefix(input) {
		return "", ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return "", ErrEmptyNumber
	}
	// if len(input) > 1 && input[0] == '0' {
	// 	return "", ErrLeadingZero
	// }
	return input, nil
}
