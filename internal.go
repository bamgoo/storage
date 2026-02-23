package storage

import (
	"errors"

	"github.com/bamgoo/bamgoo"
)

func (m *Module) instance(code string) (*Instance, *File, error) {
	file, err := decodeFile(code)
	if err != nil {
		return nil, nil, errInvalidCode
	}
	base := file.Base()
	if base == "" {
		base = bamgoo.DEFAULT
	}
	inst, ok := m.instances[base]
	if !ok {
		return nil, nil, errInvalidConnection
	}
	return inst, file, nil
}

func (m *Module) UploadTo(base string, original string, opts ...UploadOption) (*File, error) {
	if base == "" {
		base = bamgoo.DEFAULT
	}
	inst, ok := m.instances[base]
	if !ok {
		return nil, errInvalidConnection
	}

	opt := UploadOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if inst.Config.Prefix != "" && opt.Prefix == "" && opt.Key == "" {
		opt.Prefix = inst.Config.Prefix
	}

	if opt.Key == "" {
		hash, hex, err := hashFile(original)
		if err != nil {
			return nil, err
		}
		opt.Key = hash
		if len(hex) >= 4 {
			if opt.Prefix == "" {
				opt.Prefix = hex[0:2] + "/" + hex[2:4]
			} else {
				opt.Prefix = opt.Prefix + "/" + hex[0:2] + "/" + hex[2:4]
			}
		}
	}

	return inst.conn.Upload(original, opt)
}

func (m *Module) Upload(original string, opts ...UploadOption) (*File, error) {
	if m.hashring == nil {
		return nil, errInvalidConnection
	}

	opt := UploadOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	hash, hex, err := hashFile(original)
	if err != nil {
		return nil, err
	}
	if opt.Key == "" {
		opt.Key = hash
	}
	if opt.Prefix == "" && len(hex) >= 4 {
		opt.Prefix = hex[0:2] + "/" + hex[2:4]
	}

	base := m.hashring.Locate(hash)
	if base == "" {
		return nil, errors.New("no available storage instance")
	}
	return m.UploadTo(base, original, opt)
}

func (m *Module) Fetch(code string, opts ...FetchOption) (Stream, error) {
	inst, file, err := m.instance(code)
	if err != nil {
		return nil, err
	}
	opt := FetchOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return inst.conn.Fetch(file, opt)
}

func (m *Module) Download(code string, opts ...DownloadOption) (string, error) {
	inst, file, err := m.instance(code)
	if err != nil {
		return "", err
	}
	opt := DownloadOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.Target == "" {
		target, err := inst.downloadTarget(file)
		if err != nil {
			return "", err
		}
		opt.Target = target
	}
	return inst.conn.Download(file, opt)
}

func (m *Module) Remove(code string, opts ...RemoveOption) error {
	inst, file, err := m.instance(code)
	if err != nil {
		return err
	}
	opt := RemoveOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return inst.conn.Remove(file, opt)
}

func (m *Module) Browse(code string, opts ...BrowseOption) (string, error) {
	inst, file, err := m.instance(code)
	if err != nil {
		return "", err
	}
	opt := BrowseOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return inst.conn.Browse(file, opt)
}
