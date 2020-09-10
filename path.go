package python

import (
	"os"
	"path/filepath"
)

const PYTHONPATH = "PYTHONPATH"

func SetPythonPath(path ...string) {
	path_ := filepath.Join(path...)
	os.Setenv(PYTHONPATH, path_)
}

func AppendPythonPath(path ...string) {
	path_ := filepath.SplitList(os.Getenv(PYTHONPATH))
	path_ = append(path_, path...)
	path__ := filepath.Join(path_...)
	os.Setenv(PYTHONPATH, path__)
}

func PrependPythonPath(path ...string) {
	path_ := filepath.SplitList(os.Getenv(PYTHONPATH))
	path_ = append(path, path_...)
	path__ := filepath.Join(path_...)
	os.Setenv(PYTHONPATH, path__)
}
