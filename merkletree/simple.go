package merkletree

import (
	"fmt"
)

// SimpleMerkleTree rappresenta un Merkle Tree con hashing standard
type SimpleMerkleTree struct {
	MerkleTreeImpl[BytesLike]
}

// SimpleMerkleTreeOptions rappresenta le opzioni per il Simple Merkle Tree
type SimpleMerkleTreeOptions struct {
	MerkleTreeOptions // Includiamo le opzioni base del Merkle Tree
	NodeHash          NodeHash
}

// SimpleMerkleTreeData rappresenta i dati di un Simple Merkle Tree
type SimpleMerkleTreeData struct {
	Format string
	Tree   []HexString
	Values []struct {
		Value     BytesLike
		TreeIndex int
	}
	Hash string
}

// FormatLeaf converte un valore in un formato hashato per l'inserimento nel Merkle Tree
func FormatLeaf(value BytesLike) HexString {
	return StandardLeafHash(value)
}

// NewSimpleMerkleTree crea un nuovo SimpleMerkleTree con i valori dati
func NewSimpleMerkleTree(values []BytesLike, options SimpleMerkleTreeOptions) *SimpleMerkleTree {
	options.MerkleTreeOptions = NewMerkleTreeOptions(&options.MerkleTreeOptions) // Usa opzioni predefinite se non specificate

	tree, indexedValues := PrepareMerkleTree(values, options.MerkleTreeOptions, FormatLeaf, options.NodeHash)

	hashLookup := make(map[HexString]int)
	for i, v := range indexedValues {
		hash := FormatLeaf(v.Value) // 🔹 Assicuriamoci che sia lo stesso metodo usato per l'hashing
		hashLookup[hash] = i

		// Debug
		//fmt.Printf("📌 DEBUG HashLookup: Inserito %s -> Index %d\n", hash, i)
	}

	// Restituiamo il nuovo Merkle Tree
	return &SimpleMerkleTree{
		MerkleTreeImpl[BytesLike]{
			Tree:       tree,
			Values:     indexedValues,
			LeafHash:   FormatLeaf,
			NodeHash:   options.NodeHash,
			HashLookup: hashLookup, // 🔹 Ora contiene tutti i valori correttamente
		},
	}
}

// VerifySimpleMerkleTree verifica una proof di Merkle per un valore specifico
func VerifySimpleMerkleTree(root BytesLike, leaf BytesLike, proof []BytesLike, nodeHash NodeHash) bool {
	leafHash := StandardLeafHash(leaf)

	// Se `nodeHash` è nil, assegniamo la funzione standard
	if nodeHash == nil {
		nodeHash = StandardNodeHash
	}

	// Calcola la root derivata dalla proof
	computedRoot := ProcessProof(leafHash, proof, nodeHash)

	// Confronto tra root derivata e attesa
	computedRootVal, err := ToHex(computedRoot)
	if err != nil {
		fmt.Errorf("Error: ", err)
	}
	rootVal, err := ToHex(root)
	if err != nil {
		fmt.Errorf("Error: ", err)
	}
	if computedRootVal != rootVal {
		fmt.Println("❌ ERRORE: Root derivata e root attesa non corrispondono!")
	}
	return computedRootVal == rootVal
}

// Dump esporta i dati dell'albero per debugging o archiviazione
func (m *SimpleMerkleTree) Dump() SimpleMerkleTreeData {
	return SimpleMerkleTreeData{
		Format: "simple-v1",
		Tree:   m.Tree,
		Values: m.Values,
		Hash:   "custom",
	}
}
