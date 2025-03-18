package merkletree

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// BytesLike rappresenta i tipi di dati compatibili con le operazioni di hashing
type BytesLike interface{}

// HexString rappresenta una stringa esadecimale
type HexString string

// ToBytes converte un BytesLike in un array di byte (equivalente a hexToBytes in TypeScript)
func ToBytes(value BytesLike) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil
	case HexString: // Se ricevi un HexString, convertilo in string
		return ToBytes(string(v)) // Ricorsivamente chiami ToBytes con stringa normale
	case string:
		if strings.HasPrefix(v, "0x") {
			hexData := v[2:] // Rimuove "0x"
			decoded, err := hex.DecodeString(hexData)
			if err != nil {
				fmt.Println("Errore nella decodifica esadecimale:", err)
				return nil, errors.New("stringa esadecimale non valida")
			}
			return decoded, nil
		}
		return []byte(v), nil
	case []int:
		bytes := make([]byte, len(v))
		for i, num := range v {
			bytes[i] = byte(num)
		}
		return bytes, nil
	default:
		fmt.Println("Errore: tipo non supportato in ToBytes")
		return nil, errors.New("tipo non supportato in ToBytes")
	}
}

// ToHex converte un BytesLike in un HexString (equivalente a bytesToHex in TypeScript)
func ToHex(value BytesLike) (HexString, error) {
	switch v := value.(type) {
	case string, HexString: // Supporta anche HexString direttamente
		str := fmt.Sprintf("%v", v) // Converti in stringa standard
		_, err := hex.DecodeString(strings.TrimPrefix(str, "0x"))
		if err != nil {
			return "", errors.New("stringa esadecimale non valida")
		}
		return HexString("0x" + strings.TrimPrefix(str, "0x")), nil
	case []byte:
		return HexString("0x" + hex.EncodeToString(v)), nil
	case []int:
		bytes, err := ToBytes(v)
		if err != nil {
			return "", err
		}
		return HexString("0x" + hex.EncodeToString(bytes)), nil
	default:
		return "", errors.New("tipo non supportato in ToHex")
	}
}

// Concat concatena più BytesLike in un unico array di byte (equivalente a concatBytes)
func Concat(values ...BytesLike) ([]byte, error) {
	var result []byte
	for _, v := range values {
		bytes, err := ToBytes(v)
		if err != nil {
			return nil, err
		}
		result = append(result, bytes...)
	}
	return result, nil
}

// Compare confronta due BytesLike e restituisce -1, 0, 1 (equivalente alla funzione compare in TypeScript)
func Compare(a BytesLike, b BytesLike) (int, error) {
	aHex, err := ToHex(a)
	if err != nil {
		return 0, err
	}
	bHex, err := ToHex(b)
	if err != nil {
		return 0, err
	}

	aBigInt := new(big.Int)
	bBigInt := new(big.Int)

	aBigInt.SetString(string(aHex)[2:], 16) // Rimuove "0x" e converte in BigInt
	bBigInt.SetString(string(bHex)[2:], 16)

	return aBigInt.Cmp(bBigInt), nil
}
