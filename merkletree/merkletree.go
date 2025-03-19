package merkletree

import (
	"fmt"
)

// MerkleTreeImpl è la struttura base del Merkle Tree
type MerkleTreeImpl[T any] struct {
	Tree   []HexString
	Values []struct {
		Value     T
		TreeIndex int
	}
	LeafHash   func(T) HexString
	NodeHash   NodeHash
	HashLookup map[HexString]int
}

// Root restituisce la root dell'albero di Merkle
func (m *MerkleTreeImpl[T]) Root() HexString {
	return m.Tree[0]
}

// getLeafIndex restituisce l'indice di un valore nel Merkle Tree
func (m *MerkleTreeImpl[T]) getLeafIndex(leaf interface{}) int {
	switch v := leaf.(type) {
	case int:
		if v < 0 || v >= len(m.Values) {
			panic(fmt.Sprintf("❌ ERRORE: Indice foglia %d fuori dai limiti!", v))
		}
		return v
	default:
		hashedLeaf := m.LeafHash(v.(T))
		if index, found := m.HashLookup[hashedLeaf]; found {
			return index
		}
		panic("❌ ERRORE: Il valore richiesto non esiste nel Merkle Tree")
	}
}

// validateValueAt verifica che il valore sia valido nel Merkle Tree
func (m *MerkleTreeImpl[T]) validateValueAt(index int) {
	if index < 0 || index >= len(m.Values) {
		panic(fmt.Sprintf("❌ ERRORE: Indice %d fuori dai limiti!", index))
	}

	expectedHash := m.LeafHash(m.Values[index].Value)
	actualHash := m.Tree[m.Values[index].TreeIndex]

	if expectedHash != actualHash {
		panic(fmt.Sprintf("❌ ERRORE: Valore atteso %s, ma trovato %s", expectedHash, actualHash))
	}
}

// IsValidMerkleTree verifica se un Merkle Tree è valido
func IsValidMerkleTree(tree []HexString, nodeHash NodeHash) bool {
	if len(tree) == 0 {
		return false
	}

	// Controlliamo ogni nodo per assicurarci che i figli producano il valore corretto
	for i, node := range tree {
		left := LeftChildIndex(i)
		right := RightChildIndex(i)

		if right < len(tree) {
			expected := nodeHash(tree[left], tree[right])
			if expected != node {
				return false
			}
		}
	}
	return true
}

// LeafHashFromInput calcola l'hash della foglia, assicurando coerenza con la costruzione
func (m *MerkleTreeImpl[T]) LeafHashFromInput(leaf interface{}) HexString {
	switch v := leaf.(type) {
	case int:
		if v < 0 || v >= len(m.Values) {
			panic(fmt.Sprintf("❌ ERRORE: Indice foglia %d fuori dai limiti!", v))
		}
		hashed := m.LeafHash(m.Values[v].Value)
		return hashed

	default:
		hashed := m.LeafHash(v.(T))
		return hashed
	}
}

// GetProof genera una proof per un valore specifico
func (m *MerkleTreeImpl[T]) GetProof(leaf interface{}) []HexString {
	valueIndex := m.getLeafIndex(leaf)
	m.validateValueAt(valueIndex)

	treeIndex := m.Values[valueIndex].TreeIndex
	bytesTree := make([]BytesLike, len(m.Tree))
	for i, hexStr := range m.Tree {
		hexStrVal, err := ToBytes(hexStr)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		bytesTree[i] = hexStrVal
	}

	proof := GetProof(bytesTree, treeIndex)

	if len(proof) == 0 {
		panic("❌ ERRORE: Proof generata è vuota!")
	}

	return proof
}

// Verify verifica se una proof è valida
func (m *MerkleTreeImpl[T]) Verify(leaf interface{}, proof []HexString) bool {
	bytesProof := make([]BytesLike, len(proof))
	for i, hexStr := range proof {
		proofVal, err := ToBytes(hexStr)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		bytesProof[i] = proofVal
	}

	leafHash := m.LeafHashFromInput(leaf)
	hashFunc := m.NodeHash
	if hashFunc == nil {
		hashFunc = StandardNodeHash
	}

	computedRoot := ProcessProof(leafHash, bytesProof, hashFunc)

	if computedRoot != m.Root() {
		return false
	}
	return true
}

// Validate verifica se l'albero è strutturalmente valido
func (m *MerkleTreeImpl[T]) Validate() {
	for i := range m.Values {
		m.validateValueAt(i)
	}

	if !IsValidMerkleTree(m.Tree, m.NodeHash) {
		panic("❌ ERRORE: L'albero di Merkle non è valido!")
	}

	fmt.Println("✅ Albero di Merkle validato con successo!")
}
