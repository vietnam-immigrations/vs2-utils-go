package text_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vietnam-immigrations/vs2-utils-go/pkg/text"
)

func TestRemoveNonASCII(t *testing.T) {
	assert.Equal(t, "L Nam Trng.pdf", text.RemoveNonASCII("Lê Nam Trường.pdf"))
}

func TestSanitizeCVFileName(t *testing.T) {
	sanitized := text.SanitizeCVFileName("Lê. Nam Trường.pdf")
	fmt.Println(sanitized)
	assert.True(t, strings.HasPrefix(sanitized, "L  Nam Trng"))
	assert.True(t, strings.HasSuffix(sanitized, ".pdf"))
}

func TestToASCII(t *testing.T) {
	r, err := text.ToASCII("Mã: E220915AREFZYP9897682")
	assert.NoError(t, err)
	assert.Equal(t, "Ma: E220915AREFZYP9897682", r)

	r, err = text.ToASCII("SÓ HO CHIÉU: FZYP98976")
	assert.NoError(t, err)
	assert.Equal(t, "SO HO CHIEU: FZYP98976", r)

	r, err = text.ToASCII("Lê Nam Trường")
	assert.NoError(t, err)
	assert.Equal(t, "Le Nam Truong", r)
}
