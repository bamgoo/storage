package storage

import (
	"encoding/base64"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/infrago/infra"
)

type (
	File struct {
		base   string
		prefix string
		key    string
		typee  string
		size   int64

		code   string
		proxy  bool
		remote bool
	}
)

func (f *File) Base() string   { return f.base }
func (f *File) Prefix() string { return f.prefix }
func (f *File) Key() string    { return f.key }
func (f *File) Type() string   { return f.typee }
func (f *File) Size() int64    { return f.size }
func (f *File) Code() string   { return f.code }
func (f *File) Proxy() bool    { return f.proxy }
func (f *File) Remote() bool   { return f.remote }

func (f *File) File() string {
	if f.typee == "" {
		return path.Join(f.prefix, f.key)
	}
	return fmt.Sprintf("%s.%s", path.Join(f.prefix, f.key), f.typee)
}

func (f *File) Name() string {
	if f.typee == "" {
		return path.Base(f.key)
	}
	return fmt.Sprintf("%s.%s", path.Base(f.key), f.typee)
}

func (i *Instance) newFile(prefix, key, typee string, size int64) *File {
	f := &File{
		base:   i.Name,
		prefix: prefix,
		key:    key,
		typee:  typee,
		size:   size,
		proxy:  i.Config.Proxy,
		remote: i.Config.Remote,
	}
	f.code = encodeFile(f)
	return f
}

func encodeFile(file *File) string {
	base := file.base
	if base == infra.DEFAULT {
		base = ""
	}
	raw := fmt.Sprintf("%s\t%s\t%s\t%s\t%d", base, file.prefix, file.key, file.typee, file.size)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func decodeFile(code string) (*File, error) {
	bts, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return nil, errInvalidCode
	}
	args := strings.Split(string(bts), "\t")
	if len(args) != 5 {
		return nil, errInvalidCode
	}

	size, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return nil, errInvalidCode
	}

	info := &File{
		code:   code,
		base:   args[0],
		prefix: args[1],
		key:    args[2],
		typee:  args[3],
		size:   size,
	}
	if info.base == "" {
		info.base = infra.DEFAULT
	}
	if cfg, ok := module.configs[info.base]; ok {
		info.proxy = cfg.Proxy
		info.remote = cfg.Remote
	}
	return info, nil
}

func Decode(code string) (*File, error) {
	return decodeFile(code)
}
