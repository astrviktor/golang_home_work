package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrNotDir              = errors.New("path is not dir")
	ErrUnsupportedFileName = errors.New("file not supported (have = in name)")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// проверка папка ли это
	ok, err := isDir(dir)
	if !ok {
		return nil, err
	}

	// получение списка файлов в папке
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	EnvMap := make(Environment, len(files))

	for _, file := range files {
		EnvName, EnvValue, err := getEnvValueFromFile(dir, file.Name())
		if err != nil {
			return nil, err
		}
		EnvMap[EnvName] = EnvValue
	}

	return EnvMap, nil
}

func isDir(dir string) (bool, error) {
	// проверка что это не папка
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return false, err
	}
	if !fileInfo.IsDir() {
		return false, ErrNotDir
	}

	return true, nil
}

func getEnvValueFromFile(dir, file string) (string, EnvValue, error) {
	// имя файла не должно содержать `=`
	if strings.Contains(file, "=") {
		return "", EnvValue{}, ErrUnsupportedFileName
	}

	// проверка на пустой файл
	fileInfo, err := os.Stat(filepath.Join(dir, file))
	if err != nil {
		return "", EnvValue{}, err
	}
	if fileInfo.Size() == 0 {
		return file, EnvValue{"", true}, nil
	}

	// читаем весь файл
	f, err := os.Open(filepath.Join(dir, file))
	if err != nil {
		return "", EnvValue{}, err
	}
	defer f.Close()

	if err != nil {
		return "", EnvValue{}, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", EnvValue{}, err
	}

	// читаем первую строку
	bytesReader := bytes.NewReader(b)
	bufReader := bufio.NewReader(bytesReader)
	line, _, err := bufReader.ReadLine()
	if err != nil || len(line) == 0 {
		return "", EnvValue{}, err
	}

	// терминальные нули (`0x00`) заменяются на перевод строки (`\n`)
	line = bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))

	// пробелы и табуляция в конце строки удаляются
	str := strings.TrimRight(string(line), " "+"\t")

	return file, EnvValue{str, false}, nil
}
