package api

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/evilsocket/shieldwall/frontend"
	"net/http"
	"strings"
)

type mockFS struct {
	fs http.FileSystem
}

func (b *mockFS) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *mockFS) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func MockFS() *mockFS {
	return &mockFS{
		fs: &assetfs.AssetFS{
			Asset:    frontend.Asset,
			AssetDir: frontend.AssetDir,
		},
	}
}
