package merkletree

import "fmt"

// StandardMerkleTree rappresenta un Merkle Tree con encoding standard
type StandardMerkleTree[T any] struct {
	MerkleTreeImpl[T]
}

// NewStandardMerkleTree crea un nuovo StandardMerkleTree con i valori dati
func NewStandardMerkleTree[T any](values []T, options MerkleTreeOptions) *StandardMerkleTree[T] {
	options = NewMerkleTreeOptions(&options) // Usa le opzioni predefinite se non specificate

	tree, indexedValues := PrepareMerkleTree(values, options, StandardLeafHash[T], StandardNodeHash)

	return &StandardMerkleTree[T]{
		MerkleTreeImpl: MerkleTreeImpl[T]{
			Tree:       tree,
			Values:     indexedValues,
			LeafHash:   StandardLeafHash[T],
			NodeHash:   StandardNodeHash,
			HashLookup: make(map[HexString]int),
		},
	}
}

// Verify verifica una proof di Merkle per un valore specifico
func VerifyStandardMerkleTree[T any](root BytesLike, leaf T, proof []BytesLike) bool {
	leafHash := StandardLeafHash(leaf)
	leafHashVal, err := ToHex(leafHash)
	if err != nil {
		fmt.Errorf("Error: ", err)
	}
	// Debug
	fmt.Println("📌 DEBUG VerifyStandardMerkleTree: Hash calcolato per la leaf:", leafHashVal)

	// Calcola la root derivata dalla proof
	computedRoot := ProcessProof(leafHash, proof, StandardNodeHash)
	computedRootVal, err := ToHex(computedRoot)
	rootVal, err := ToHex(root)
	// Debug
	fmt.Println("📌 DEBUG VerifyStandardMerkleTree: Root derivata:", computedRootVal)
	fmt.Println("📌 DEBUG VerifyStandardMerkleTree: Root attesa:", rootVal)

	// Confronto tra root derivata e attesa
	return computedRootVal == rootVal
}

// StandardMerkleTreeData rappresenta i dati esportabili di un Standard Merkle Tree
type StandardMerkleTreeData[T any] struct {
	Format string
	Tree   []HexString
	Values []struct {
		Value     T
		TreeIndex int
	}
}

// Dump esporta i dati dell'albero per debugging o archiviazione
func (m *StandardMerkleTree[T]) Dump() StandardMerkleTreeData[T] {
	return StandardMerkleTreeData[T]{
		Format: "standard-v1",
		Tree:   m.Tree,
		Values: m.Values,
	}
}
