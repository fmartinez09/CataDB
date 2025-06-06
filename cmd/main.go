package main

import (
	"CataDB/internal/engine"
	"CataDB/internal/tx"
	"fmt"
)

func main() {
	// Inicializa la store MVCC
	mvcc := engine.NewMVCCStore()

	// Transacción 1: escribe un valor con timestamp 100
	txn1 := tx.NewTransaction(mvcc, 0)
	txn1.Write("blob-123", []byte("primer valor"))
	err := txn1.Commit(100)
	if err != nil {
		panic(err)
	}

	// Transacción 2: lee antes de commit (snapshot 50)
	txn2 := tx.NewTransaction(mvcc, 50)
	val2, err := txn2.Read("blob-123")
	if err != nil {
		fmt.Println("TXN2: No se encontró valor (snapshot 50)")
	} else {
		fmt.Printf("TXN2: Valor en snapshot 50 = %s\n", val2)
	}

	// Transacción 3: lee después de commit (snapshot 150)
	txn3 := tx.NewTransaction(mvcc, 150)
	val3, err := txn3.Read("blob-123")
	if err != nil {
		fmt.Println("TXN3: No se encontró valor (snapshot 150)")
	} else {
		fmt.Printf("TXN3: Valor en snapshot 150 = %s\n", val3)
	}
}
