package service

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInlineTemplateLocalAssets_InlinesLogoFile(t *testing.T) {
	root := t.TempDir()
	logoDir := filepath.Join(root, "logos")
	require.NoError(t, os.MkdirAll(logoDir, 0o755))

	pngData, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/w8AAn8BT0J7ewAAAABJRU5ErkJggg==")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(logoDir, "tiny.png"), pngData, 0o644))

	svc := &DocumentService{templateDirPath: root}
	rendered := svc.inlineTemplateLocalAssets(`<img src="../logos/tiny.png" alt="Labaid InsureTech">`)

	require.NotContains(t, rendered, "../logos/tiny.png")
	require.Contains(t, rendered, `src="data:image/png;base64,`)
}

func TestInlineTemplateLocalAssets_LeavesRemoteURLsUntouched(t *testing.T) {
	svc := &DocumentService{templateDirPath: t.TempDir()}
	in := `<img src="https://example.com/logo.png" alt="remote">`

	out := svc.inlineTemplateLocalAssets(in)

	require.Equal(t, in, out)
	require.True(t, strings.Contains(out, "https://example.com/logo.png"))
}
