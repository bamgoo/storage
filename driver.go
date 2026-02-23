package storage

import "io"

type (
	Driver interface {
		Connect(*Instance) (Connection, error)
	}

	Health struct {
		Workload int64
	}

	Stream interface {
		io.Reader
		io.Seeker
		io.Closer
		io.ReaderAt
	}

	Connection interface {
		Open() error
		Health() Health
		Close() error

		Upload(string, UploadOption) (*File, error)
		Fetch(*File, FetchOption) (Stream, error)
		Download(*File, DownloadOption) (string, error)
		Remove(*File, RemoveOption) error
		Browse(*File, BrowseOption) (string, error)
	}
)
