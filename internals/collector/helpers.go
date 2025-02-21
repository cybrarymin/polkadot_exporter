package collector

import (
	"bytes"
	"fmt"

	"github.com/polkadot-go/api/v4/scale"
)

func encodeU32(u uint32) ([]byte, error) {
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)

	// Encode your eraIndex (or any uint32) into the buffer
	if err := encoder.Encode(u); err != nil {
		return nil, fmt.Errorf("failed to SCALE-encode uint32: %w", err)
	}

	// Return the raw bytes
	return buf.Bytes(), nil
}
