package gsnake

type FileReader interface {
	ReadFile(file string, offset int) (err error)
}
