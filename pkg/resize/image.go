package resize

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/s3"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
	vs2ssm "github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/aws/ssm"
)

var extensionsCanBeResize = []string{
	".jpeg", ".jpg", ".png", ".webp", ".gif", ".tiff", ".heic",
}

type imageConvertRequest struct {
	Source string  `json:"source"`
	From   *string `json:"from,omitempty"`
	To     string  `json:"to"`
	Width  int     `json:"width"`
}

type imageConvertResponse struct {
	Base64 string `json:"base64"`
}

func canBeResize(filename string) bool {
	for _, ext := range extensionsCanBeResize {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}

func resizeImage(ctx context.Context, fileContent []byte, fileName string, isHEIC bool) ([]byte, *string, error) {
	log := logger.FromContext(ctx)
	log.Infof("Resize image [%s]", fileName)

	tempKey := fmt.Sprintf("%s/%s_%s", time.Now().Format("2006/01/02"), random.String(5, lo.AlphanumericCharset), fileName)
	err := s3.WriteFileBucketSSM(ctx, vs2ssm.S3BucketTemp, tempKey, fileContent)
	if err != nil {
		log.Errorf("Failed to create temp file [%s]", tempKey)
		return nil, nil, err
	}
	tempFileURL, err := s3.PublicURLSSMBucket(ctx, vs2ssm.S3BucketTemp, tempKey)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	toolsURL, err := ssm.GetParameter(ctx, "/external/tools/url", false)
	if err != nil {
		return nil, nil, err
	}
	toolsAPIKey, err := ssm.GetParameter(ctx, "/external/tools/apikey", true)
	if err != nil {
		return nil, nil, err
	}
	resize := imageConvertRequest{
		Source: *tempFileURL,
		To:     "jpeg",
		Width:  1000,
	}
	if isHEIC {
		resize.From = lo.ToPtr("heic")
	}
	body, err := json.Marshal(resize)
	if err != nil {
		log.Errorf("failed to marshal resize request body: %s", err)
		return nil, nil, err
	}
	resizeURL := fmt.Sprintf("%s/%s", toolsURL, "image/v1/convert")
	log.Infof("Resize image [%s]", resizeURL)
	req, err := http.NewRequest(http.MethodPost, resizeURL, bytes.NewReader(body))
	if err != nil {
		log.Errorf("failed to create resize request: %s", err)
		return nil, nil, err
	}
	req.Header.Set("X-Api-Key", toolsAPIKey)
	req.Header.Set("Content-Type", "application/json")
	correlationID, ok := ctx.Value(mycontext.FieldCorrelationID).(string)
	if !ok {
		correlationID = uuid.New().String()
	}
	req.Header.Set("X-Correlation-ID", correlationID)
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("http request error: %s", err)
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("failed to close request body: %s", err)
		}
	}(res.Body)
	if res.StatusCode != 200 {
		log.Errorf("status code [%d]", res.StatusCode)
		return nil, nil, fmt.Errorf("status code [%d]", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("failed to read response body: %s", err)
		return nil, nil, err
	}
	resizeRes := new(imageConvertResponse)
	err = json.Unmarshal(resBody, resizeRes)
	if err != nil {
		log.Errorf("failed to parse response JSON: %s", err)
		return nil, nil, err
	}
	convertedFile, err := base64.StdEncoding.DecodeString(resizeRes.Base64)
	if err != nil {
		log.Errorf("failed to base64 decode converted image: %s", err)
		return nil, nil, err
	}

	newFileName := fmt.Sprintf("%s.jpg", fileName)
	return convertedFile, lo.ToPtr(newFileName), nil
}
