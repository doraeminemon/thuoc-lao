package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

func main() {
	dbName := "local.db"
	primaryUrl := "https://thuoclaodb-doraeminemon.turso.io"
	authToken := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDAxMzYxMzEsImlhdCI6MTc0MDEzMjUzMSwiaWQiOiI5NDg4YzE0NS1kMDY3LTQwYWEtYThkOC0yMTAyYzY2ZGQ0OWIifQ.ubE-kcIlFEh-oLjxGtX6svx3BGSZwIyL3D9BVr3iKbdPRo3Y1r0tLZrpTitHNKQcVmZnSwA0AkVjMcZvoL_cAQ"

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	queryProducts(db)
}

type Product struct {
	ID   int
	Name string
}

func queryProducts(db *sql.DB) {
	rows, err := db.Query("SELECT id, name FROM products")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v", err)
		os.Exit(1)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product

		if err := rows.Scan(&product.ID, &product.Name); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		products = append(products, product)
		fmt.Println(product.ID, product.Name)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
}
