package resize

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

const (
	megabytes = 1e+6
)

func File(ctx context.Context, fileContent []byte, fileName string, width int) ([]byte, *string, error) {
	log := logger.FromContext(ctx)
	log.Infof("Resize file [%s]", fileName)

	isHEICFile := strings.HasSuffix(strings.ToLower(fileName), ".heic")
	if isHEICFile {
		log.Infof("File [%s] is HEIC, will be converted", fileName)
	}

	if !canBeResize(fileName) {
		log.Errorf("Cannot resize [%s], supported extensions: %s", fileName, strings.Join(extensionsCanBeResize, ", "))
		log.Errorf("Original file will be used")
		return fileContent, lo.ToPtr(fileName), nil
	}

	newContent, newName, err := resizeImage(ctx, fileContent, fileName, width, isHEICFile)
	if err != nil {
		log.Errorf("Failed to resize image [%s], original file will be used", fileName)
		return fileContent, lo.ToPtr(fileName), nil
	}
	log.Infof("Image [%s] resized successfully", fileName)
	return newContent, newName, nil
}
