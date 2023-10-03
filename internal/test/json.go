package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/kien-manabie/rio-mocker/internal/types"
	"github.com/stretchr/testify/require"
)

func ParseJSONFileToMap(t *testing.T, filePath string) types.Map {
	rawJSON, err := os.ReadFile(filePath)
	require.NoError(t, err)

	data, err := types.CreateMapFromReader(bytes.NewReader(rawJSON))
	require.NoError(t, err)
	return data
}
