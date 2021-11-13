package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyPath             = errors.New("empty path")
	ErrEmptySize             = errors.New("empty size")
)

// буффер для копирования файла.
const bufferMaxSize = 1024 * 1024 * 10

func Copy(fromPath, toPath string, offset, limit int64) error {
	// проверка на пустые пути к файлам
	if fromPath == "" || toPath == "" {
		return ErrEmptyPath
	}

	// проверка файла fromPath
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return ErrUnsupportedFile
	}
	size := fileInfo.Size()
	if size == 0 {
		return ErrEmptySize
	}

	// если offset больше, чем размер файла - невалидная ситуация;
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	copySize := size - offset
	if limit != 0 && copySize > limit {
		copySize = limit
	}

	// определение буфера для прогресса копирования в процентах (%)
	bufferSize := copySize / 4
	switch {
	case bufferSize == 0:
		bufferSize = 1
	case bufferSize > bufferMaxSize:
		bufferSize = bufferMaxSize
	}
	buf := make([]byte, bufferSize)

	// файл from
	fi, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fi.Close()

	_, err = fi.Seek(offset, 0)
	if err != nil {
		return err
	}

	// файл to
	fo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fo.Close()

	// процесс копирования
	copyLeft := copySize
	var copyCount int64

	for copyCount != copySize {
		if copyLeft < bufferSize {
			buf = buf[:copyLeft]
		}

		n, err := io.ReadFull(fi, buf)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
			return err
		}

		_, err = fo.Write(buf[:n])
		if err != nil {
			return err
		}
		copyCount += int64(n)

		// прогресс
		copyLeft -= bufferSize
		fmt.Println("progress", 100*copyCount/copySize, "%")
	}

	return nil
}
