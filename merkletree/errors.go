package merkletree

import (
	"fmt"
	"log"
)

// Invariant verifica una condizione e causa un panic se la condizione è falsa
func Invariant(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("InvariantError: %s", message))
	}
}

// InvariantWithDebug permette di fornire messaggi più dettagliati per il debug
func InvariantWithDebug(condition bool, message string, debugInfo interface{}) {
	if !condition {
		log.Fatalf("InvariantError: %s | Debug Info: %v", message, debugInfo)
		panic(fmt.Sprintf("InvariantError: %s", message))
	}
}

// SafePanic permette di catturare gli errori e fornire messaggi dettagliati senza crash immediato
func SafePanic(message string) {
	log.Println("❌ ERRORE CRITICO:", message)
	panic(fmt.Sprintf("FatalError: %s", message))
}

// ValidateArgument verifica una condizione e causa un panic se la condizione è falsa
func ValidateArgument(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("InvalidArgumentError: %s", message))
	}
}

// Assert verifica una condizione e fornisce un log di errore senza crash immediato
func Assert(condition bool, message string) {
	if !condition {
		log.Println("⚠️ AssertError:", message)
	}
}
