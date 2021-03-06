package export

import (
	"io"
	"os"
	"path/filepath"
)

type copyInfo struct {
	source, target string
}

func copyDir(src, dst string, excludes []string) error {
	// Sweep for sources and targets
	var cpInfos []copyInfo
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, exclude := range excludes {
			if path == exclude {
				return nil
			}
		}
		cpInfo := copyInfo{
			source: path,
			target: filepath.Join(dst, filepath.Base(path)),
		}
		cpInfos = append(cpInfos, cpInfo)
		return nil
	})
	if err != nil {
		return err
	}

	// Copy sources to targets
	for _, cpInfo := range cpInfos {
		if dir := filepath.Dir(cpInfo.target); dir != "." {
			if err := os.MkdirAll(dir, MODE_DIR); err != nil {
				return err
			}
		}
		source, err := os.Open(cpInfo.source)
		if err != nil {
			return err
		}
		target, err := os.Create(cpInfo.target)
		if err != nil {
			return err
		}
		if _, err := io.Copy(target, source); err != nil {
			return err
		}
		source.Close()
		target.Close()
	}
	return nil
}
