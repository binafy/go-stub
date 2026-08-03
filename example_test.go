package stub_test

import (
	"fmt"
	"os"
	"path/filepath"

	stub "github.com/binafy/go-stub"
)

func ExampleRender() {
	out, err := stub.Render("testdata/model.stub",
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
	// Output:
	// package models
	//
	// type User struct {
	// 	ID   int
	// 	Name string
	// }
	//
	// func NewUser() *User {
	// 	return &User{}
	// }
}

func ExampleGenerate() {
	dst := filepath.Join(os.TempDir(), "go-stub-example", "user.go")
	defer os.RemoveAll(filepath.Dir(dst))

	err := stub.Generate("testdata/model.stub", dst,
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
		stub.WithForce(),
		stub.WithFormat(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("generated:", filepath.Base(dst))
	// Output: generated: user.go
}

func ExampleNew() {
	out, err := stub.New().
		From("testdata/model.stub").
		Replaces(map[string]any{"PACKAGE": "models", "NAME": "Account"}).
		Render()
	if err != nil {
		panic(err)
	}
	fmt.Println(out[:14])
	// Output: package models
}

func ExampleToSnake() {
	fmt.Println(stub.ToSnake("HTTPServer"))
	fmt.Println(stub.ToPascal("user_name"))
	// Output:
	// http_server
	// UserName
}
