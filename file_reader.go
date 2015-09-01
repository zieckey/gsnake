package gsnake

type FileReader interface {
    LoadFile(file string, pos int) (err error)
}

type TextFileReader interface {
    FileReader
    ReadLine() ([]byte, error)
}

