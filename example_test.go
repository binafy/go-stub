package stub_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
	defer func() { _ = os.RemoveAll(filepath.Dir(dst)) }()

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

func ExampleRenderContent() {
	// The stub is already in memory — no file needed.
	out, err := stub.RenderContent("type {{ NAME }} struct{}",
		stub.WithReplace("NAME", "User"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: type User struct{}
}

func ExampleRenderFS() {
	// embedded is an embed.FS declared in the test package.
	out, err := stub.RenderFS(embedded, "testdata/model.stub",
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "Account"}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(out[:14])
	// Output: package models
}

func ExampleGenerateDir() {
	dstDir := filepath.Join(os.TempDir(), "go-stub-example-dir")
	defer func() { _ = os.RemoveAll(dstDir) }()

	err := stub.GenerateDir("testdata/scaffold", dstDir,
		stub.WithReplaces(map[string]any{"PACKAGE": "user", "NAME": "User"}),
		stub.WithTrimSuffix(".stub"),
	)
	if err != nil {
		panic(err)
	}

	// Print the generated tree, using forward slashes so the output is the
	// same on every OS.
	var files []string
	_ = filepath.Walk(dstDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			rel, _ := filepath.Rel(dstDir, path)
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	sort.Strings(files)
	for _, f := range files {
		fmt.Println(f)
	}
	// Output:
	// User.go
	// inner/note.txt
}

func ExampleGenerateJobs() {
	dir := filepath.Join(os.TempDir(), "go-stub-example-jobs")
	defer func() { _ = os.RemoveAll(dir) }()

	err := stub.GenerateJobs([]stub.Job{
		{Src: "testdata/model.stub", Dst: filepath.Join(dir, "user.go"),
			Opts: []stub.Option{stub.WithReplace("NAME", "User")}},
		{Src: "testdata/model.stub", Dst: filepath.Join(dir, "account.go"),
			Opts: []stub.Option{stub.WithReplace("NAME", "Account")}},
	}, stub.WithReplace("PACKAGE", "models")) // shared by every job
	if err != nil {
		panic(err)
	}

	fmt.Println("generated 2 files")
	// Output: generated 2 files
}

func ExampleWithStrict() {
	// The stub needs NAME, but it is not provided.
	_, err := stub.Render("testdata/model.stub",
		stub.WithReplace("PACKAGE", "models"),
		stub.WithStrict(),
	)

	var mk *stub.MissingKeysError
	if errors.As(err, &mk) {
		fmt.Println("missing:", mk.Keys)
	}
	fmt.Println("is ErrMissingKeys:", errors.Is(err, stub.ErrMissingKeys))
	// Output:
	// missing: [NAME]
	// is ErrMissingKeys: true
}
