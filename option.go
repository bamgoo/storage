package storage

import (
	"time"

	. "github.com/bamgoo/base"
)

type (
	UploadOption struct {
		Key    string
		Prefix string

		Mimetype string
		Metadata Map
		Tags     Map
		Expires  time.Time
	}

	FetchOption struct {
		Start int64
		End   int64
	}

	DownloadOption struct {
		Target string
	}

	RemoveOption struct{}

	BrowseOption struct {
		Headers Map
		Params  Map
		Expires time.Duration
	}
)
