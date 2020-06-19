package rootmodule

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	discardLogger = log.New(ioutil.Discard, "", 0)
	skipDirNames  = map[string]bool{
		".git":                true,
		".idea":               true,
		".vscode":             true,
		"terraform.tfstate.d": true,
	}
)

type Walker struct {
	logger *log.Logger
}

func NewWalker() *Walker {
	return &Walker{
		logger: discardLogger,
	}
}

func (w *Walker) SetLogger(logger *log.Logger) {
	w.logger = logger
}

type WalkFunc func(rootModulePath string) error

func (w *Walker) WalkInitializedRootModules(path string, wf WalkFunc) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			w.logger.Printf("unable to access %s: %s", path, err.Error())
			return nil
		}

		if isSkippable(info) {
			return filepath.SkipDir
		}

		w.logger.Printf("walking through %s", info.Name())

		if info.Name() == ".terraform" {
			rootDir, err := filepath.Abs(filepath.Dir(path))
			if err != nil {
				return err
			}

			return wf(rootDir)
		}

		return nil
	})
}

func isSkippable(info os.FileInfo) bool {
	if !info.IsDir() {
		// All files are skipped, we only care about dirs
		return true
	}

	_, ok := skipDirNames[info.Name()]
	if ok {
		return true
	}

	return false
}
