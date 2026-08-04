package stub_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	stub "github.com/binafy/go-stub"
)

func TestRenderContent(t *testing.T) {
	got, err := stub.RenderContent("Hello {{ NAME }}!",
		stub.WithReplace("NAME", "World"),
	)
	if err != nil {
		t.Fatalf("RenderContent() error = %v", err)
	}
	if got != "Hello World!" {
		t.Errorf("RenderContent() = %q", got)
	}
}

func TestRenderContentStrict(t *testing.T) {
	_, err := stub.RenderContent("{{ A }} {{ B }}",
		stub.WithReplace("A", "x"),
		stub.WithStrict(),
	)
	if !errors.Is(err, stub.ErrMissingKeys) {
		t.Fatalf("RenderContent() error = %v, want ErrMissingKeys", err)
	}
}

func TestBuilderContentRender(t *testing.T) {
	got, err := stub.New().
		Content("type {{ NAME }} struct{}").
		Replace("NAME", "User").
		Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got != "type User struct{}" {
		t.Errorf("Render() = %q", got)
	}
}

func TestBuilderContentGenerate(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out", "user.go")

	err := stub.New().
		Content("package {{ PKG }}\n").
		To(dst).
		Replace("PKG", "models").
		Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "package models\n" {
		t.Errorf("generated = %q", string(data))
	}
}

func TestBuilderContentGenerateStrict(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "out.go")

	err := stub.New().
		Content("{{ MISSING }}").
		To(dst).
		Strict().
		Generate()
	if !errors.Is(err, stub.ErrMissingKeys) {
		t.Fatalf("Generate() error = %v, want ErrMissingKeys", err)
	}
	if _, statErr := os.Stat(dst); statErr == nil {
		t.Error("file should not be written on strict failure")
	}
}

func TestBuilderContentNeedsDestination(t *testing.T) {
	if err := stub.New().Content("x").Generate(); err == nil {
		t.Error("Generate() expected errNoDestination")
	}
}

func TestBuilderContentOverridesFrom(t *testing.T) {
	// From then Content: Content wins (last source set).
	got, err := stub.New().
		From("testdata/model.stub").
		Content("just {{ X }}").
		Replace("X", "content").
		Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got != "just content" {
		t.Errorf("Render() = %q", got)
	}
}
