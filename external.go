package storage

import (
	"errors"
	"fmt"

	. "github.com/bamgoo/base"
)

func Upload(from Any, opts ...UploadOption) (*File, error) {
	path := ""
	switch vv := from.(type) {
	case string:
		path = vv
	case Map:
		if file, ok := vv["file"].(string); ok {
			path = file
		} else {
			return nil, errors.New("invalid target")
		}
	default:
		path = fmt.Sprintf("%v", vv)
	}
	return module.Upload(path, opts...)
}

func UploadTo(base string, from Any, opts ...UploadOption) (*File, error) {
	path := ""
	switch vv := from.(type) {
	case string:
		path = vv
	case Map:
		if file, ok := vv["file"].(string); ok {
			path = file
		} else {
			return nil, errors.New("invalid target")
		}
	default:
		path = fmt.Sprintf("%v", vv)
	}
	return module.UploadTo(base, path, opts...)
}

func Fetch(code string, opts ...FetchOption) (Stream, error) {
	return module.Fetch(code, opts...)
}

func Download(code string, opts ...DownloadOption) (string, error) {
	return module.Download(code, opts...)
}

func Remove(code string, opts ...RemoveOption) error {
	return module.Remove(code, opts...)
}

func Browse(code string, opts ...BrowseOption) (string, error) {
	return module.Browse(code, opts...)
}

func ThumbnailConfig() string { return module.filecfg.Thumbnail }
func PreviewConfig() string   { return module.filecfg.Preview }
func SaltConfig() string      { return module.filecfg.Salt }
