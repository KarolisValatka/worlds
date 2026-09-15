package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAndValidate(t *testing.T) {
	dir := t.TempDir()
	v := filepath.Join(dir, "v")
	if err := Init(v); err != nil {
		t.Fatal(err)
	}
	if err := Validate(v); err != nil {
		t.Fatal(err)
	}
}

func TestValidateExample(t *testing.T) {
	// run from module root in CI; skip if missing
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	ex := filepath.Join(root, "example-vault")
	if _, err := os.Stat(ex); err != nil {
		t.Skip("example-vault not present")
	}
	if err := Validate(ex); err != nil {
		t.Fatal(err)
	}
}
