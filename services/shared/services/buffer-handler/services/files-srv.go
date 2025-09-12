package services

import "os"

type FilesSrv struct{}

func NewFilesSrv() *FilesSrv {
	return &FilesSrv{}
}

func (srv *FilesSrv) CreateAndWriteFile(bytes []byte) error {
	pwd, _ := os.Getwd()
	path := pwd + "/services/shared/services/buffer-handler/assets/clone-image.jpg"
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)

	return err
}
