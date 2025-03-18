package merkletree

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// LeafHash rappresenta una funzione che calcola l'hash di una foglia
type LeafHash[T any] func(leaf T) HexString

// NodeHash rappresenta una funzione che calcola l'hash di un nodo
type NodeHash func(left BytesLike, right BytesLike) HexString

// Keccak256 calcola l'hash Keccak-256 di un input
func Keccak256(input BytesLike) HexString {
	inputB, err := ToBytes(input)
	if err != nil {
		fmt.Errorf("Error: ", err)
	}
	hash := crypto.Keccak256(inputB)
	hashHex, _ := ToHex(hash)

	return hashHex
}

// StandardLeafHash calcola l'hash standard di una foglia, utilizzando l'encoding ABI come su Ethereum
func StandardLeafHash[T any](value T) HexString {
	// Convertiamo il valore in bytes
	/*valArr, err := ToBytes(value)
	if err != nil {
		panic("Errore nella conversione in bytes: " + err.Error())
	}

	// ✅ **Creiamo un array di 32 byte (non slice!)**
	var fixedArray [32]byte
	copy(fixedArray[:], valArr)

	// 📌 DEBUG: Stampiamo il valore convertito
	fixedArrayHex, _ := ToHex(fixedArray[:])
	fmt.Println("📌 DEBUG StandardLeafHash - Valore convertito:", fixedArrayHex)*/

	// Codifica ABI dei dati come fa lo smart contract di OpenZeppelin
	/*abiType, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		panic("Errore nella creazione del tipo ABI: " + err.Error())
	}*/

	/*encoded, err := abi.Arguments{
		{Type: abiType},
	}.Pack(fixedArray) // ✅ Ora stiamo passando un array di 32 byte!
	if err != nil {
		panic("Errore nella codifica ABI: " + err.Error())
	}*/
	encodedPacked, err := keccak256HashedData(value)
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}
	encodedPackedHex, err := ToHex(encodedPacked)
	return encodedPackedHex
}

// StandardNodeHash calcola l'hash standard di due nodi
func StandardNodeHash(a BytesLike, b BytesLike) HexString {
	// Ordiniamo i due nodi per garantire consistenza
	nodes := []BytesLike{a, b}
	sort.Slice(nodes, func(i, j int) bool {
		result, err := Compare(nodes[i], nodes[j])
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		return result < 0
	})
	concatenated, err := Concat(nodes[0], nodes[1])
	if err != nil {
		fmt.Errorf("Error: ", err)
	}
	hashed, _ := keccak256HashedData(concatenated)
	hashedHex, _ := ToHex(hashed)

	return hashedHex
}

func abiEncodePacked(args ...interface{}) ([]byte, error) {
	var buf bytes.Buffer

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			buf.Write([]byte(v)) // Converti stringa in byte senza padding
		case []byte:
			buf.Write(v) // Scrivi direttamente i byte
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64:
			buf.Write(uintToBytes(v)) // Converte gli interi in byte
		default:
			return nil, fmt.Errorf("tipo non supportato: %T", v)
		}
	}

	return buf.Bytes(), nil
}

// Converte interi in byte senza padding extra
func uintToBytes(num interface{}) []byte {
	switch v := num.(type) {
	case uint8:
		return []byte{v}
	case uint16:
		return []byte{byte(v >> 8), byte(v)}
	case uint32:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	case uint64:
		return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32),
			byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	default:
		return nil
	}
}

// Funzione per calcolare il Keccak256 dei dati codificati
func keccak256HashedData(args ...interface{}) ([]byte, error) {
	encodedData, err := abiEncodePacked(args...)
	if err != nil {
		return nil, err
	}

	// Calcola Keccak256 (SHA3 con Ethereum specifica)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(encodedData)
	return hash.Sum(nil), nil
}
func doubleKeccak256HashedData(args ...interface{}) ([]byte, error) {
	// Primo hash
	firstHash, err := keccak256HashedData(args...)
	if err != nil {
		return nil, err
	}

	// Secondo hash
	secondHash := sha3.NewLegacyKeccak256()
	secondHash.Write(firstHash)
	return secondHash.Sum(nil), nil
}
