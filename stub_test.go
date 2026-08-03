package stub_test

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	stub "github.com/binafy/go-stub"
)

//go:embed testdata/model.stub
var embedded embed.FS

func TestApplyViaRender(t *testing.T) {
	tests := []struct {
		name    string
		content string
		opts    []stub.Option
		want    string
	}{
		{
			name:    "single key",
			content: "Hello {{ NAME }}!",
			opts:    []stub.Option{stub.WithReplace("NAME", "World")},
			want:    "Hello World!",
		},
		{
			name:    "no inner spaces",
			content: "Hello {{NAME}}!",
			opts:    []stub.Option{stub.WithReplace("NAME", "World")},
			want:    "Hello World!",
		},
		{
			name:    "repeated key",
			content: "{{ X }}-{{ X }}",
			opts:    []stub.Option{stub.WithReplace("X", "1")},
			want:    "1-1",
		},
		{
			name:    "multiple keys via map",
			content: "{{ A }} and {{ B }}",
			opts:    []stub.Option{stub.WithReplaces(map[string]any{"A": "x", "B": "y"})},
			want:    "x and y",
		},
		{
			name:    "unknown key preserved",
			content: "{{ KNOWN }} {{ UNKNOWN }}",
			opts:    []stub.Option{stub.WithReplace("KNOWN", "ok")},
			want:    "ok {{ UNKNOWN }}",
		},
		{
			name:    "non-string value",
			content: "count={{ N }}",
			opts:    []stub.Option{stub.WithReplace("N", 42)},
			want:    "count=42",
		},
		{
			name:    "unterminated delimiter left untouched",
			content: "prefix {{ NAME",
			opts:    []stub.Option{stub.WithReplace("NAME", "x")},
			want:    "prefix {{ NAME",
		},
		{
			name:    "custom delimiters",
			content: "Hello <NAME>!",
			opts:    []stub.Option{stub.WithDelimiters("<", ">"), stub.WithReplace("NAME", "World")},
			want:    "Hello World!",
		},
		{
			name:    "no placeholders",
			content: "plain text",
			opts:    nil,
			want:    "plain text",
		},
		{
			name:    "empty content",
			content: "",
			opts:    []stub.Option{stub.WithReplace("NAME", "x")},
			want:    "",
		},
		{
			name:    "later replace overrides earlier",
			content: "{{ K }}",
			opts:    []stub.Option{stub.WithReplace("K", "first"), stub.WithReplace("K", "second")},
			want:    "second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTemp(t, tt.content)
			got, err := stub.Render(path, tt.opts...)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderRealFixture(t *testing.T) {
	got, err := stub.Render("testdata/model.stub",
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "User"}),
	)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want := "package models\n\n" +
		"type User struct {\n\tID   int\n\tName string\n}\n\n" +
		"func NewUser() *User {\n\treturn &User{}\n}\n"
	if got != want {
		t.Errorf("Render() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderMissingFile(t *testing.T) {
	if _, err := stub.Render("testdata/does-not-exist.stub"); err == nil {
		t.Fatal("Render() expected error for missing file, got nil")
	}
}

func TestRenderFS(t *testing.T) {
	got, err := stub.RenderFS(embedded, "testdata/model.stub",
		stub.WithReplaces(map[string]any{"PACKAGE": "models", "NAME": "Account"}),
	)
	if err != nil {
		t.Fatalf("RenderFS() error = %v", err)
	}
	if want := "package models"; got[:len(want)] != want {
		t.Errorf("RenderFS() prefix = %q, want %q", got[:len(want)], want)
	}
	if !contains(got, "func NewAccount() *Account") {
		t.Errorf("RenderFS() missing rendered constructor, got:\n%s", got)
	}
}

func TestRenderFSMissingFile(t *testing.T) {
	if _, err := stub.RenderFS(embedded, "testdata/missing.stub"); err == nil {
		t.Fatal("RenderFS() expected error for missing file, got nil")
	}
}

// writeTemp writes content to a temporary .stub file and returns its path.
func writeTemp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fixture.stub")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp fixture: %v", err)
	}
	return path
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}
