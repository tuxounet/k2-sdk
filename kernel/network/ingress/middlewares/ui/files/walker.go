package files

import (
	"embed"
	"path"
)

func WalkFolder(dir string, fs *embed.FS) ([]string, error) {
	files, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range files {
		if file.IsDir() {
			subFiles, err := WalkFolder(path.Join(dir, file.Name()), fs)
			if err != nil {
				return nil, err
			}
			result = append(result, subFiles...)
		} else {
			result = append(result, path.Join(dir, file.Name()))
		}
	}
	return result, nil
}
