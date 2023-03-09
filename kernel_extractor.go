package vmlinux

/*
	Original source: https://github.com/Caesurus/extract-vmlinux-v2
*/

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"github.com/itchio/lzma"
	"github.com/pierrec/lz4"
	"github.com/smira/go-xz"
	"io"
	"math"
	"os"
)

type supportedAlgo struct {
	Name        string
	ExtractFunc func(r io.Reader) (reader io.ReadCloser, err error)
	pattern     []byte
}

var algos = []supportedAlgo{
	{
		Name: "GZIP",
		ExtractFunc: func(r io.Reader) (reader io.ReadCloser, err error) {
			return gzip.NewReader(r)
		},
		pattern: []byte("\037\213\010"),
	},
	{
		Name: "BZIP",
		ExtractFunc: func(r io.Reader) (reader io.ReadCloser, err error) {
			return io.NopCloser(bzip2.NewReader(r)), nil
		},
		pattern: []byte("BZh"),
	},
	{
		Name: "LZMA",
		ExtractFunc: func(r io.Reader) (reader io.ReadCloser, err error) {
			return lzma.NewReader(r), nil
		},
		pattern: []byte("\135\000\000\000"),
	},
	{
		Name: "LZ4",
		ExtractFunc: func(r io.Reader) (reader io.ReadCloser, err error) {
			return io.NopCloser(lz4.NewReaderLegacy(r)), nil
		},
		pattern: []byte("\002!L\030"),
	},
	{
		Name: "XZ",
		ExtractFunc: func(r io.Reader) (reader io.ReadCloser, err error) {
			xzIn, err := xz.NewReader(r)
			if err != nil {
				return nil, err
			}
			return io.NopCloser(xzIn), nil
		},
		pattern: []byte("\3757zXZ\000"),
	},
}

func IsKernelImage(r io.ReaderAt) bool {
	var ident [16]uint8
	if _, err := r.ReadAt(ident[0:], 0); err != nil {
		return false
	}

	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		return false
	}

	return true
}

// Extract will attempt to extract the vmlinux file the memory
func Extract(r io.ReaderAt) ([]byte, error) {
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)

	if err := ExtractTo(r, w); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// ExtractTo will attempt to extract the vmlinux file to an io.Writer
func ExtractTo(r io.ReaderAt, w io.Writer) error {
	sectionReader := io.NewSectionReader(r, 0, math.MaxInt64)

	// Search the compressor signature in the first 64K bytes
	headBuffer := make([]byte, 65535)
	if n, err := sectionReader.Read(headBuffer); n < 65535 || err != nil {
		return err
	}

	for _, algo := range algos {
		offset := bytes.Index(headBuffer, algo.pattern)
		if offset == -1 {
			continue
		}

		if _, err := sectionReader.Seek(int64(offset), io.SeekStart); err != nil {
			return err
		}

		buf := bufio.NewReaderSize(sectionReader, os.Getpagesize())
		decompressor, err := algo.ExtractFunc(buf)

		if err == nil {
			magic := make([]byte, 16)
			if _, err := decompressor.Read(magic); err != nil {
				return err
			}

			if !IsKernelImage(bytes.NewReader(magic)) {
				continue
			}

			if _, err := w.Write(magic); err != nil {
				return err
			}

			if _, err := io.Copy(w, decompressor); err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("kernel image is not found")
}
