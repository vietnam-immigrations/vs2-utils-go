package resize

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/random"
	vs2ssm "github.com/vietnam-immigrations/vs2-utils-go/pkg/aws/ssm"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/s3"
	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

type pdfConvertRequest struct {
	Source string `json:"source"`
}

type pdfConvertResponse struct {
	Pages []pdfConvertResponsePage `json:"pages"`
}

type pdfConvertResponsePage struct {
	Base64 string `json:"base64"`
}

func pdfToJPG(ctx context.Context, fileContent []byte, fileName string) ([]byte, *string, error) {
	log := logger.FromContext(ctx)
	log.Infof("Convert PDF file [%s] to JPG", fileName)

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

	client := http.Client{Timeout: 5 * time.Minute}
	convert := pdfConvertRequest{
		Source: *tempFileURL,
	}
	body, err := json.Marshal(convert)
	if err != nil {
		log.Errorf("failed to marshal request body: %s", err)
		return nil, nil, err
	}

	toolsURL, err := ssm.GetParameter(ctx, "/external/tools/url", false)
	if err != nil {
		return nil, nil, err
	}
	toolsAPIKey, err := ssm.GetParameter(ctx, "/external/tools/apikey", true)
	if err != nil {
		return nil, nil, err
	}

	convertURL := fmt.Sprintf("%s/%s", toolsURL, "pdf/v1/convert")
	log.Infof("convert pdf to image: %s", convertURL)
	req, err := http.NewRequest(http.MethodPost, convertURL, bytes.NewReader(body))
	if err != nil {
		log.Errorf("failed to create request: %s", err)
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
		log.Errorf("http request failed: %s", err)
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("failed to close request body: %s", err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		log.Errorf("http status code [%s]", res.Status)
		return nil, nil, fmt.Errorf("http status code [%s]", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("failed to read response body: %s", err)
		return nil, nil, err
	}
	convertRes := new(pdfConvertResponse)
	err = json.Unmarshal(resBody, convertRes)
	if err != nil {
		log.Errorf("failed to read response JSON: %s", err)
		return nil, nil, err
	}

	images := make([][]byte, 0)
	for _, page := range convertRes.Pages {
		converted, err := base64.StdEncoding.DecodeString(page.Base64)
		if err != nil {
			return nil, nil, err
		}
		images = append(images, converted)
	}

	newFileName := fmt.Sprintf("%s.jpg", fileName)
	log.Infof("PDF file [%s] converted to [%s]", fileName, newFileName)
	return images[0], lo.ToPtr(newFileName), nil
}
