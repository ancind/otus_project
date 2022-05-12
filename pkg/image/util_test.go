package image

import (
	"bufio"
	"fmt"
	"os"
)

func loadImage(imgName string) []byte {
	fileToBeUploaded := "../../tests/static/" + imgName
	file, err := os.Open(fileToBeUploaded)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	return bytes
}
