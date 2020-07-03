package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "oogo",
		Price: 12,
		SKU:   "abd-ddd-ddd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
