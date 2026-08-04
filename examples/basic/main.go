// Command basic demonstrates go-stub: it embeds a stub file, renders it with
// both the functional and fluent APIs, and prints the result.
//
// Run it from the repository root:
//
//	go run ./examples/basic
package main

import (
	"embed"
	"fmt"
	"log"

	stub "github.com/binafy/go-stub"
)

//go:embed model.stub
var stubs embed.FS

func main() {
	name := "User"

	// Functional API, reading the stub from the embedded filesystem.
	out, err := stub.RenderFS(stubs, "model.stub",
		stub.WithReplaces(map[string]any{
			"PACKAGE": stub.ToSnake(name) + "s", // "users"
			"NAME":    stub.ToPascal(name),      // "User"
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- functional (RenderFS) ---")
	fmt.Println(out)

	// Fluent API, same source and values.
	out, err = stub.New().
		FromFS(stubs, "model.stub").
		Replace("PACKAGE", "users").
		Replace("NAME", "User").
		Render()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- fluent (New) ---")
	fmt.Println(out)
}
