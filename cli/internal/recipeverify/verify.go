package recipeverify

import "github.com/nkapatos/sbx-kit/cli/internal/experimental"

// Describe returns a short status blurb for help text.
func Describe() string {
	return "planned: recipe manifest validation and sbx kit verify delegation (see sbx-kit experimental verify)"
}

// VerifyRecipe checks a recipe manifest and references (stub).
func VerifyRecipe(catalogRoot, id string) error {
	_ = catalogRoot
	_ = id
	return experimental.ErrNotReady{Feature: "recipe verify", Track: "verify recipe"}
}

// VerifyKits checks kits in a recipe directory via sbx (stub).
func VerifyKits(catalogRoot, dir string) error {
	_ = catalogRoot
	_ = dir
	return experimental.ErrNotReady{Feature: "kit verify", Track: "verify kit"}
}
