package merkletree

import (
	"fmt"
	"math"
	"sort"
)

type MultiProof struct {
	Leaves     []HexString // Hash delle foglie incluse nella proof
	Proof      []HexString // Lista dei nodi necessari per il calcolo della root
	ProofFlags []bool      // Indica quali nodi devono essere combinati
}

// IsTreeNode verifica se l'indice `i` è un nodo valido nell'albero
func IsTreeNode(tree []BytesLike, i int) bool {
	return i >= 0 && i < len(tree)
}

// IsInternalNode verifica se l'indice `i` è un nodo interno dell'albero di Merkle
func IsInternalNode(tree []BytesLike, i int) bool {
	return IsTreeNode(tree, LeftChildIndex(i))
}

// IsLeafNode verifica se un indice `i` è una foglia nell'albero di Merkle
func IsLeafNode(tree []BytesLike, i int) bool {
	return IsTreeNode(tree, i) && !IsInternalNode(tree, i)
}

// CheckLeafNode verifica se un indice è una foglia nell'albero di Merkle
func CheckLeafNode(tree []BytesLike, i int) {
	if !IsLeafNode(tree, i) {
		panic("Index is not a leaf")
	}
}

func IsValidMerkleNode(node BytesLike) bool {
	bytes, err := ToBytes(node)
	if err != nil {
		fmt.Errorf("Not valide Merkle Node: ", err)
	}
	return len(bytes) == 32
}

func CheckValidMerkleNode(node BytesLike) {
	if !IsValidMerkleNode(node) {
		panic("Merkle tree nodes must be Uint8Array of length 32")
	}
}

// MakeMerkleTree costruisce un albero di Merkle a partire da una lista di hash delle foglie
func MakeMerkleTree(hashes []BytesLike, nodeHash NodeHash) []HexString {
	if len(hashes) == 0 {
		panic("Errore: impossibile costruire un albero di Merkle con 0 elementi")
	}
	// Converte tutti gli hash in BytesLike
	leaves := make([]HexString, len(hashes))
	for i, h := range hashes {
		leaf, err := ToHex(h)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		leaves[i] = leaf
	}
	// Costruisce l'albero di Merkle
	tree := make([]HexString, 2*len(leaves)-1)
	copy(tree[len(tree)-len(leaves):], leaves)

	// Generazione dei nodi interni
	for i := len(tree) - len(leaves) - 1; i >= 0; i-- {
		leftChild := tree[LeftChildIndex(i)]
		rightChild := tree[RightChildIndex(i)]
		tree[i] = nodeHash(leftChild, rightChild)
	}

	return tree
}

// GetProof restituisce la proof di Merkle per un nodo specifico
func GetProof(tree []BytesLike, index int) []HexString {
	CheckLeafNode(tree, index)
	var proof []HexString
	for index > 0 {
		siblingIdx := SiblingIndex(index)
		value, err := ToHex(tree[siblingIdx])
		if err != nil {
			fmt.Errorf("Errore in GetProof: ", err)
		}
		proof = append(proof, value)
		index = ParentIndex(index)
	}
	return proof
}

// ProcessProof verifica la proof data e calcola la root risultante
func ProcessProof(leaf BytesLike, proof []BytesLike, nodeHash NodeHash) HexString {
	// Verifica che il nodo foglia sia valido
	CheckValidMerkleNode(leaf)

	// Verifica che tutti gli elementi della proof siano nodi validi
	for _, node := range proof {
		CheckValidMerkleNode(node)
	}

	// Applica la funzione di hash riducendo la proof a un singolo valore
	result, err := ToHex(leaf)
	if err != nil {
		fmt.Errorf("Error in ProcessProof: ", err)
	}
	for _, sibling := range proof {
		sibling, err := ToHex(sibling)
		if err != nil {
			fmt.Errorf("Error in ProcessProof: ", err)
		}
		result = nodeHash(result, sibling)
	}
	resultHex, err := ToHex(result)
	if err != nil {
		fmt.Printf("Error in ProcessProof: ", err)
	}
	return resultHex
}

// GetMultiProof genera una proof multipla per un insieme di foglie
func GetMultiProof(tree []BytesLike, indices []int) MultiProof {
	if len(indices) == 0 {
		panic("Errore: impossibile generare una proof multipla per 0 elementi")
	}

	var proof []HexString
	var proofFlags []bool
	stack := indices

	for len(stack) > 0 && stack[0] > 0 {
		j := stack[0]
		stack = stack[1:]

		s := SiblingIndex(j)
		p := ParentIndex(j)

		if len(stack) > 0 && s == stack[0] {
			proofFlags = append(proofFlags, true)
			stack = stack[1:]
		} else {
			proofFlags = append(proofFlags, false)
			proofVal, err := ToHex(tree[s])
			if err != nil {
				fmt.Errorf("Error: ", err)
			}
			proof = append(proof, proofVal)
		}

		stack = append(stack, p)
	}

	leavesHex := make([]HexString, len(indices))
	for i, idx := range indices {
		index, err := ToHex(idx)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		leavesHex[i] = index // Converte l'indice in formato esadecimale
	}
	return MultiProof{
		Leaves:     leavesHex,
		Proof:      proof,
		ProofFlags: proofFlags,
	}
}

// ProcessMultiProof verifica una proof multipla e calcola la root risultante
func ProcessMultiProof(multiproof MultiProof, nodeHash NodeHash) HexString {
	stack := multiproof.Leaves
	proof := multiproof.Proof

	for _, flag := range multiproof.ProofFlags {
		if len(stack) < 1 || (!flag && len(proof) < 1) {
			panic("Errore: multiproof non valida")
		}

		a := stack[0]
		stack = stack[1:]
		var b HexString
		if flag {
			b = stack[0]
			stack = stack[1:]
		} else {
			b = proof[0]
			proof = proof[1:]
		}
		leafA, err := ToHex(a)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		leafB, err := ToHex(b)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		stack = append(stack, nodeHash(leafA, leafB))
	}

	if len(stack)+len(proof) != 1 {
		panic("Errore: multiproof non valida")
	}

	return stack[0]
}

// Funzioni di supporto per gli indici degli alberi di Merkle
// ParentIndex restituisce l'indice del nodo genitore per un nodo dato
func ParentIndex(i int) int {
	if i > 0 {
		return int(math.Floor((float64(i) - 1) / 2))
	}
	panic("❌ ERRORE: La radice non ha un nodo genitore!")
}

// SiblingIndex restituisce l'indice del nodo fratello per un nodo dato
func SiblingIndex(i int) int {
	if i > 0 {
		val := i - int(math.Pow(-1, float64(i%2)))
		//fmt.Println("sibling", val) ora corretto!
		return val
	}
	panic("❌ ERRORE: La radice non ha fratelli!")
}

func LeftChildIndex(i int) int {
	return 2*i + 1
}

func RightChildIndex(i int) int {
	return 2*i + 2
}

// PrepareMerkleTree costruisce l'albero di Merkle e assegna gli indici corretti alle foglie
func PrepareMerkleTree[T any](values []T, options MerkleTreeOptions, leafHash func(T) HexString, nodeHash NodeHash) ([]HexString, []struct {
	Value     T
	TreeIndex int
}) {

	// Se `nodeHash` è nil, assegniamo la funzione standard
	if nodeHash == nil {
		nodeHash = StandardNodeHash
	}

	// Assicuriamoci che `leafHash` sia sempre `StandardLeafHash`
	leafHash = StandardLeafHash[T]

	// Creiamo una struttura per memorizzare i valori hashati
	hashedValues := make([]struct {
		Value      T
		ValueIndex int
		Hash       HexString
	}, len(values))

	// Applica la funzione di hash alle foglie
	for i, value := range values {
		hashedValues[i] = struct {
			Value      T
			ValueIndex int
			Hash       HexString
		}{
			Value:      value,
			ValueIndex: i,
			Hash:       leafHash(value),
		}

	}

	// Se l'opzione `sortLeaves` è attiva, ordiniamo le foglie
	if options.SortLeaves {
		sort.Slice(hashedValues, func(i, j int) bool {
			result, err := Compare(hashedValues[i].Hash, hashedValues[j].Hash)
			if err != nil {
				fmt.Errorf("Error: ", err)
			}
			return result < 0
		})
	}

	// Costruiamo l'albero di Merkle
	tree := MakeMerkleTree(
		func() []BytesLike {
			hashes := make([]BytesLike, len(hashedValues))
			for i, v := range hashedValues {
				hashes[i] = v.Hash
			}
			return hashes
		}(),
		nodeHash,
	)
	// Assegniamo gli indici corretti alle foglie
	indexedValues := make([]struct {
		Value     T
		TreeIndex int
	}, len(values))

	for leafIndex, hv := range hashedValues {
		correctedIndex := len(tree) - len(hashedValues) + leafIndex
		indexedValues[hv.ValueIndex] = struct {
			Value     T
			TreeIndex int
		}{
			Value:     hv.Value,
			TreeIndex: correctedIndex,
		}
	}

	// Verifica che gli indici siano validi
	for _, v := range indexedValues {
		if v.TreeIndex < 0 || v.TreeIndex >= len(tree) {
			fmt.Printf("❌ ERRORE: TreeIndex %d è fuori dai limiti! (Max: %d)\n", v.TreeIndex, len(tree)-1)
			panic("TreeIndex fuori dai limiti!")
		}
	}

	return tree, indexedValues
}
