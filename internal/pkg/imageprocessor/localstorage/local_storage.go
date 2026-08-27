package localstorage

type LocalStorage struct {
	subDir string 
}

func New (subDir string) LocalStorage {
	return LocalStorage{
		subDir: subDir,
	}
}