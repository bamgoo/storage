package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/infrago/infra"
)

func init() {
	module.RegisterDriver(infra.DEFAULT, &defaultDriver{})
}

type (
	defaultDriver struct{}

	defaultConnection struct {
		mutex  sync.RWMutex
		health Health

		instance *Instance
		setting  defaultSetting
	}

	defaultSetting struct {
		Storage string
	}
)

func (d *defaultDriver) Connect(instance *Instance) (Connection, error) {
	setting := defaultSetting{Storage: "store/storage"}
	if vv, ok := instance.Setting["storage"].(string); ok && vv != "" {
		setting.Storage = vv
	}
	return &defaultConnection{instance: instance, setting: setting}, nil
}

func (c *defaultConnection) Open() error { return nil }

func (c *defaultConnection) Health() Health {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.health
}

func (c *defaultConnection) Close() error { return nil }

func (c *defaultConnection) Upload(original string, opt UploadOption) (*File, error) {
	stat, err := os.Stat(original)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errors.New("directory upload not supported")
	}

	ext := infra.Extension(original)
	if opt.Key == "" {
		h, hex, err := hashFile(original)
		if err != nil {
			return nil, err
		}
		opt.Key = h
		if opt.Prefix == "" && len(hex) >= 4 {
			opt.Prefix = hex[0:2] + "/" + hex[2:4]
		}
	}

	file := c.instance.newFile(opt.Prefix, opt.Key, ext, stat.Size())
	_, target, err := c.ensurePath(file)
	if err != nil {
		return nil, err
	}

	src, err := os.Open(original)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	buf := make([]byte, 256*1024)
	if _, err := io.CopyBuffer(dst, src, buf); err != nil {
		return nil, err
	}

	return file, nil
}

func (c *defaultConnection) Fetch(file *File, opt FetchOption) (Stream, error) {
	sfile := c.objectPath(file)
	_, err := os.Stat(sfile)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(sfile)
	if err != nil {
		return nil, err
	}

	if opt.Start > 0 || opt.End > 0 {
		st, err := f.Stat()
		if err != nil {
			_ = f.Close()
			return nil, err
		}
		size := st.Size()
		start := opt.Start
		end := opt.End
		if start < 0 {
			start = 0
		}
		if end <= 0 || end > size {
			end = size
		}
		if end < start {
			end = start
		}
		return &rangeFileReader{file: f, reader: io.NewSectionReader(f, start, end-start)}, nil
	}

	return f, nil
}

func (c *defaultConnection) Download(file *File, opt DownloadOption) (string, error) {
	sfile := c.objectPath(file)
	_, err := os.Stat(sfile)
	if err != nil {
		return "", err
	}
	return sfile, nil
}

func (c *defaultConnection) Remove(file *File, _ RemoveOption) error {
	sfile := c.objectPath(file)
	if err := os.Remove(sfile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (c *defaultConnection) Browse(file *File, _ BrowseOption) (string, error) {
	sfile := c.objectPath(file)
	_, err := os.Stat(sfile)
	if err != nil {
		return "", err
	}
	return sfile, nil
}

func (c *defaultConnection) objectPath(file *File) string {
	name := file.Key()
	if file.Type() != "" {
		name = fmt.Sprintf("%s.%s", file.Key(), file.Type())
	}
	return path.Join(c.setting.Storage, file.Prefix(), name)
}

func (c *defaultConnection) ensurePath(file *File) (string, string, error) {
	sfile := c.objectPath(file)
	spath := path.Dir(sfile)
	if err := os.MkdirAll(spath, 0o755); err != nil {
		return "", "", err
	}
	return spath, sfile, nil
}
