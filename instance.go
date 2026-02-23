package storage

import (
	"fmt"
	"os"
	"path"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
)

type (
	Instance struct {
		conn Connection

		Name    string
		Config  Config
		Setting Map
	}
)

// NewFile creates a storage file metadata object for drivers.
func (i *Instance) NewFile(prefix, key, typee string, size int64) *File {
	return i.newFile(prefix, key, typee, size)
}

func (i *Instance) downloadTarget(file *File) (string, error) {
	name := file.Key()
	if file.Type() != "" {
		name = fmt.Sprintf("%s.%s", file.Key(), file.Type())
	}

	base := file.Base()
	if base == bamgoo.DEFAULT {
		base = ""
	}

	sfile := path.Join(module.filecfg.Download, file.Base(), file.Prefix(), name)
	spath := path.Dir(sfile)
	if err := os.MkdirAll(spath, 0o755); err != nil {
		return "", err
	}
	return sfile, nil
}
