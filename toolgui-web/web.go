package toolguiweb

import (
	"embed"
	"io/fs"
	"path"
)

//go:embed build/index.html
var IndexBody string

//go:embed build/static/*
var staticDir embed.FS

//go:embed build/*
var rootAssets embed.FS

func GetStaticDir() fs.FS {
	fsys, err := fs.Sub(staticDir, "build")
	if err != nil {
		panic(err)
	}
	return fsys
}

func GetRootAssets() map[string][]byte {
	entries, err := rootAssets.ReadDir("build")
	if err != nil {
		panic(err)
	}

	files := map[string][]byte{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		bs, err := rootAssets.ReadFile(path.Join("build", entry.Name()))
		if err != nil {
			panic(err)
		}

		files[entry.Name()] = bs
	}
	return files
}
