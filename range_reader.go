package storage

import "io"

type rangeFileReader struct {
	file   io.ReadSeekCloser
	reader *io.SectionReader
}

func (r *rangeFileReader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *rangeFileReader) Seek(offset int64, whence int) (int64, error) {
	return r.reader.Seek(offset, whence)
}

func (r *rangeFileReader) ReadAt(p []byte, off int64) (int, error) {
	return r.reader.ReadAt(p, off)
}

func (r *rangeFileReader) Close() error {
	return r.file.Close()
}
