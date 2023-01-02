package resize

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

const (
	megabytes = 1e+6
)

func File(ctx context.Context, fileContent []byte, fileName string) ([]byte, *string, error) {
	log := logger.FromContext(ctx)
	log.Infof("Resize file [%s]", fileName)

	isBigFile := len(fileContent) > megabytes
	if isBigFile {
		log.Infof("File [%s] is [%dMB], bigger than 1MB, will try to resize", fileName, len(fileContent)/megabytes)
	}
	isHEICFile := strings.HasSuffix(strings.ToLower(fileName), ".heic")
	if isHEICFile {
		log.Infof("File [%s] is HEIC, will be converted", fileName)
	}
	isPDFFile := strings.HasSuffix(strings.ToLower(fileName), ".pdf")

	// PDF file
	if isPDFFile {
		newContent, newName, err := pdfToJPG(ctx, fileContent, fileName)
		if err != nil {
			log.Errorf("Failed to convert pdf to image")
			return nil, nil, err
		}
		log.Infof("PDF file converted to image, now will be resized")
		return File(ctx, newContent, *newName)
	}

	// image file
	if isBigFile || isHEICFile {
		log.Infof("[%s] is big file or HEIC, will be converted", fileName)
		if !canBeResize(fileName) {
			log.Errorf("Cannot resize [%s], supported extensions: %s", fileName, strings.Join(extensionsCanBeResize, ", "))
			log.Errorf("Original file will be used")
			return fileContent, lo.ToPtr(fileName), nil
		}

		newContent, newName, err := resizeImage(ctx, fileContent, fileName, isHEICFile)
		if err != nil {
			log.Errorf("Failed to resize image [%s], original file will be used", fileName)
			return fileContent, lo.ToPtr(fileName), nil
		}
		log.Infof("Image [%s] resized successfully", fileName)
		return newContent, newName, nil
	}

	log.Infof("File [%s] is small and is not HEIC, no need to resize", fileName)
	return fileContent, lo.ToPtr(fileName), nil
}
