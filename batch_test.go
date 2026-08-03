package stub_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	stub "github.com/binafy/go-stub"
)

func TestGenerateJobs(t *testing.T) {
	dir := t.TempDir()
	userDst := filepath.Join(dir, "user.go")
	acctDst := filepath.Join(dir, "acct.go")

	jobs := []stub.Job{
		{Src: "testdata/model.stub", Dst: userDst, Opts: []stub.Option{stub.WithReplace("NAME", "User")}},
		{FS: embedded, Src: "testdata/model.stub", Dst: acctDst, Opts: []stub.Option{stub.WithReplace("NAME", "Account")}},
	}

	err := stub.GenerateJobs(jobs, stub.WithReplace("PACKAGE", "models"))
	if err != nil {
		t.Fatalf("GenerateJobs() error = %v", err)
	}

	if got := readFile(t, userDst); !contains(got, "package models") || !contains(got, "func NewUser() *User") {
		t.Errorf("user file = %q", got)
	}
	if got := readFile(t, acctDst); !contains(got, "func NewAccount() *Account") {
		t.Errorf("acct file = %q", got)
	}
}

func TestGenerateJobsPerJobOverride(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "out.go")
	if err := os.WriteFile(dst, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Shared options refuse overwrite; the job opts add Force to override.
	jobs := []stub.Job{
		{Src: "testdata/model.stub", Dst: dst, Opts: []stub.Option{stub.WithReplace("NAME", "X"), stub.WithForce()}},
	}
	if err := stub.GenerateJobs(jobs, stub.WithReplace("PACKAGE", "m")); err != nil {
		t.Fatalf("GenerateJobs() error = %v", err)
	}
	if got := readFile(t, dst); got == "old" {
		t.Error("per-job Force did not overwrite")
	}
}

func TestGenerateJobsStopsOnError(t *testing.T) {
	dir := t.TempDir()
	jobs := []stub.Job{
		{Src: "testdata/missing.stub", Dst: filepath.Join(dir, "a.go")},
	}
	err := stub.GenerateJobs(jobs)
	if err == nil {
		t.Fatal("GenerateJobs() expected error")
	}
	if !contains(err.Error(), "job 0") {
		t.Errorf("error should identify the failing job index: %v", err)
	}
}

func TestGenerateJobsExistsError(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "out.go")
	if err := os.WriteFile(dst, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	jobs := []stub.Job{
		{Src: "testdata/model.stub", Dst: dst, Opts: []stub.Option{stub.WithReplaces(map[string]any{"PACKAGE": "m", "NAME": "X"})}},
	}
	err := stub.GenerateJobs(jobs)
	if !errors.Is(err, stub.ErrExists) {
		t.Fatalf("GenerateJobs() error = %v, want ErrExists", err)
	}
}
