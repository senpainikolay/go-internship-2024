// Denumirile de package nu le scrim cu _/- sau alte simboluri, ex: imagehandler
package image_handler

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

// Dupa parerea mea subiectiva, poate nu am dreptate, dar ce ai tu aici e cam un usecase
// si ar trebui sa fie in sertvicii
// sau nu ?

func CreateImageFile(file *multipart.File, imgFileName string) error {

	f, err := os.OpenFile("./uploads/"+imgFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("could not create image")
	}
	defer f.Close()

	_, err = io.Copy(f, *file)
	if err != nil {
		return errors.New("could not copy image from request to the server")
	}

	return nil
}

func ReadImageFile(imgFileName string) (*[]byte, error) {

	file, err := os.OpenFile("./uploads/"+imgFileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.New("could not open the image")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	binaryData := make([]byte, fileInfo.Size())
	_, err = file.Read(binaryData)
	if err != nil {
		return nil, err
	}

	return &binaryData, nil
}
