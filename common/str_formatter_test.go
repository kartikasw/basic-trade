package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type strFormatterTest struct {
	Keyword string
}

func TestFormatStrForFullTextSearch(t *testing.T) {
	arg1 := strFormatterTest{}
	arg1.Keyword = FormatStrForFullTextSearch(arg1.Keyword)
	require.Equal(t, arg1.Keyword, ":*")

	arg2 := strFormatterTest{Keyword: "  Kartika   Sari  "}
	arg2.Keyword = FormatStrForFullTextSearch(arg2.Keyword)
	require.NotEqual(t, arg2.Keyword, ":*")
	require.Equal(t, arg2.Keyword, "Kartika:* & Sari:*")
}
