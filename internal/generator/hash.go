package generator

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

// HashSlots provides deterministic values derived from an input string.
type HashSlots struct {
	values []int
}

// NewHashSlots creates the fixed sequence used to select image assets.
func NewHashSlots(input string, count int, ignoreExtension bool) (*HashSlots, error) {
	if count < 1 || count > 30 {
		return nil, fmt.Errorf("slots must be between 1 and 30")
	}

	if ignoreExtension {
		input = normalizeInput(input)
	}

	digest := sha512.Sum512([]byte(input))
	hexDigest := hex.EncodeToString(digest[:])
	blockSize := len(hexDigest) / count
	values := make([]int, 0, count*2)

	for i := 0; i < count; i++ {
		// The historical default of 11 slots creates odd-length hexadecimal
		// chunks. DecodeString returns the usable prefix alongside an error;
		// retain that prefix to preserve existing avatar selections.
		block, _ := hex.DecodeString(hexDigest[i*blockSize : (i+1)*blockSize])
		if len(block) == 0 {
			return nil, fmt.Errorf("decode hash slot %d", i)
		}
		values = append(values, int(block[0]))
	}

	// Preserve the original RoboHash slot behavior when more asset groups are
	// added than the requested number of hash slots.
	values = append(values, values...)
	return &HashSlots{values: values}, nil
}

// At returns a stable value for the requested slot.
func (h *HashSlots) At(index int) int {
	return h.values[index%len(h.values)]
}

func normalizeInput(input string) string {
	input = strings.ToLower(input)
	if strings.HasSuffix(input, ".jpg") {
		return strings.TrimSuffix(input, ".jpg") + ".jpeg"
	}
	return input
}
