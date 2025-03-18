package merkletree

import "fmt"

// MerkleTreeOptions definisce le opzioni di configurazione per la costruzione dell'albero di Merkle.
type MerkleTreeOptions struct {
	SortLeaves bool `json:"sortLeaves"` // Se true, le foglie vengono ordinate per facilitare le multiproof
}

// DefaultOptions rappresenta la configurazione predefinita per un Merkle Tree
var DefaultOptions = MerkleTreeOptions{
	SortLeaves: true, // Ordinamento delle foglie abilitato di default per multiproof più efficienti
}

// NewMerkleTreeOptions crea un oggetto `MerkleTreeOptions` con valori predefiniti se non specificati
func NewMerkleTreeOptions(options *MerkleTreeOptions) MerkleTreeOptions {
	if options == nil {
		fmt.Println(DefaultOptions)
		return DefaultOptions
	}
	// sto ritornando sempre DefaultOptions perchè se non metto nulla prende che ho messo false
	return DefaultOptions
}
