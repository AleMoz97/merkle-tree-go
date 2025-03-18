package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/AleMoz97/merkle-tree-go/merkletree"
)

func main() {
	fmt.Println("🚀 Inizio test per SimpleMerkleTree")

	// 1️⃣ Creiamo un array di dati da includere nell'albero
	values := []merkletree.BytesLike{
		"ciao",
		"0x2222",
		"0x3333",
		"0x4444",
	}

	// 2️⃣ Creiamo l'albero di Merkle
	tree := merkletree.NewSimpleMerkleTree(values, merkletree.SimpleMerkleTreeOptions{})

	// 3️⃣ Stampiamo la root dell'albero
	fmt.Println("✅ Merkle Root:", tree.Root())

	// 5️⃣ Selezioniamo un valore dall'albero per testare la proof
	testLeaf := values[2] // "0x3333"

	// Generiamo la proof per il valore selezionato
	proof := tree.GetProof(testLeaf)

	// 6️⃣ Stampiamo la proof generata
	fmt.Println("\n🔍 DEBUG Proof Generata:")
	for i, p := range proof {
		fmt.Printf("  Step %d: %s\n", i+1, p)
	}

	// 7️⃣ Convertiamo la proof in BytesLike
	proofBytes := make([]merkletree.BytesLike, len(proof))
	for i, p := range proof {
		proofVal, err := merkletree.ToBytes(p)
		if err != nil {
			fmt.Errorf("Error: ", err)
		}
		proofBytes[i] = proofVal
	}

	// 8️⃣ Verifichiamo se la proof è valida
	isValid := merkletree.VerifySimpleMerkleTree(tree.Root(), testLeaf, proofBytes, nil)
	fmt.Println("\n✅ Proof valida?", isValid)

	// 9️⃣ Testiamo il dump dell'albero

	treeData := tree.Dump()
	jsonData, err := json.MarshalIndent(treeData, "", "  ")
	if err != nil {
		log.Fatalf("Errore nella serializzazione JSON: %v", err)
	}
	fmt.Println("📋 JSON dell'albero di Merkle:\n", string(jsonData))
	filename := "tmp/jsonMerkle.json"
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("❌ Errore nella scrittura del file: %v", err)
	}

	fmt.Printf("✅ Albero di Merkle salvato con successo in %s\n", filename)
}
