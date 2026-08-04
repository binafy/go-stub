package stub_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	stub "github.com/binafy/go-stub"
)

func TestRenderStrictReportsMissingKeys(t *testing.T) {
	path := writeTemp(t, "{{ A }} {{ B }} {{ A }} {{ C }}")

	_, err := stub.Render(path,
		stub.WithReplace("A", "x"),
		stub.WithStrict(),
	)
	if err == nil {
		t.Fatal("Render() expected strict error, got nil")
	}

	// errors.Is via the sentinel.
	if !errors.Is(err, stub.ErrMissingKeys) {
		t.Errorf("errors.Is(err, ErrMissingKeys) = false, err = %v", err)
	}

	// errors.As to read the keys, in first-seen order without duplicates.
	var mk *stub.MissingKeysError
	if !errors.As(err, &mk) {
		t.Fatalf("errors.As(*MissingKeysError) = false, err = %v", err)
	}
	want := []string{"B", "C"}
	if len(mk.Keys) != len(want) {
		t.Fatalf("keys = %v, want %v", mk.Keys, want)
	}
	for i := range want {
		if mk.Keys[i] != want[i] {
			t.Errorf("keys = %v, want %v", mk.Keys, want)
		}
	}
}

func TestRenderStrictAllResolved(t *testing.T) {
	path := writeTemp(t, "{{ A }}-{{ B }}")

	got, err := stub.Render(path,
		stub.WithReplaces(map[string]any{"A": "1", "B": "2"}),
		stub.WithStrict(),
	)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got != "1-2" {
		t.Errorf("Render() = %q, want %q", got, "1-2")
	}
}

func TestRenderNonStrictLeavesUnknown(t *testing.T) {
	path := writeTemp(t, "{{ A }} {{ B }}")

	got, err := stub.Render(path, stub.WithReplace("A", "x")) // no strict
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got != "x {{ B }}" {
		t.Errorf("Render() = %q, want %q", got, "x {{ B }}")
	}
}

func TestGenerateStrictWritesNothing(t *testing.T) {
	src := writeTemp(t, "package p\nvar X = {{ MISSING }}\n")
	dst := filepath.Join(t.TempDir(), "out.go")

	err := stub.Generate(src, dst, stub.WithStrict())
	if !errors.Is(err, stub.ErrMissingKeys) {
		t.Fatalf("Generate() error = %v, want ErrMissingKeys", err)
	}
	if _, statErr := os.Stat(dst); statErr == nil {
		t.Error("destination should not be created when strict check fails")
	}
}

func TestBuilderStrict(t *testing.T) {
	path := writeTemp(t, "{{ NAME }}")

	_, err := stub.New().From(path).Strict().Render()
	if !errors.Is(err, stub.ErrMissingKeys) {
		t.Fatalf("Builder Render() error = %v, want ErrMissingKeys", err)
	}
}

func TestGenerateDirStrictOnFileName(t *testing.T) {
	// A placeholder in the file name with no replacement should fail in strict
	// mode and identify it as a file-name error.
	dstDir := t.TempDir()

	err := stub.GenerateDir("testdata/scaffold", dstDir,
		stub.WithReplace("PACKAGE", "p"), // NAME intentionally missing
		stub.WithStrict(),
	)
	if !errors.Is(err, stub.ErrMissingKeys) {
		t.Fatalf("GenerateDir() error = %v, want ErrMissingKeys", err)
	}
}

func TestMissingKeysErrorMessage(t *testing.T) {
	err := &stub.MissingKeysError{Keys: []string{"A", "B"}}
	if got := err.Error(); got != "stub: unresolved placeholders: A, B" {
		t.Errorf("Error() = %q", got)
	}
}
