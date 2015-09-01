package gsnake

type FileReader interface {
	ReadFile(file string, pos int) (err error)
}
