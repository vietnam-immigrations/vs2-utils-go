package cv

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssns "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/dchest/uniuri"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/sns"
	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
)

// UploadToS3AndSendSNS uploads PDF to S3 and send SNS. Stage must be in context
func UploadToS3AndSendSNS(ctx context.Context, fileNames []string, fileContents []io.ReadCloser) error {
	log := logger.FromContext(ctx)
	stage, ok := ctx.Value(mycontext.FieldStage).(string)
	if !ok {
		return fmt.Errorf("missing stage in context")
	}

	s3Bucket, err := ssm.GetParameter(ctx, "/result/s3BucketName", false)
	if err != nil {
		log.Errorf("failed to get s3 location: %s", err)
		return err
	}

	result := &db.Result{
		ID:        primitive.NewObjectID(),
		S3DirKey:  uniuri.NewLen(12),
		Files:     make([]db.ResultFile, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// copy files to S3
	for i, fileName := range fileNames {
		_, err := putToS3Bucket(ctx, result, s3Bucket, fileName, fileContents[i])
		if err != nil {
			log.Errorf("failed to put to S3: %s", err)
			return err
		}
		result.Files = append(result.Files, db.ResultFile{
			Name:         fileName,
			Processed:    false,
			ErrorMessage: "",
			OrderNumber:  "",
		})
	}

	// save result to mongo
	colResult, err := db.CollectionResult(ctx)
	if err != nil {
		log.Errorf("failed to get mongodb collection: %s", err)
		return err
	}
	_, err = colResult.InsertOne(ctx, result)
	if err != nil {
		log.Errorf("failed to insert result to mongo: %s", err)
		return err
	}

	// publish SNS
	snsClient, err := sns.NewClient(ctx)
	if err != nil {
		log.Errorf("failed to create SNS client: %s", err)
		return err
	}
	newResultTopic, err := ssm.GetParameter(ctx, "/sns/newResult/arn", false)
	if err != nil {
		log.Errorf("failed to get SNS topic arn: %s", err)
	}
	log.Infof("messages will be published to [%s]", newResultTopic)

	correlationID, ok := ctx.Value(mycontext.FieldCorrelationID).(string)
	if !ok || correlationID == "" {
		correlationID = uuid.New().String()
	}

	for i, fileName := range fileNames {
		_, err := snsClient.Publish(ctx, &awssns.PublishInput{
			Message: aws.String(fmt.Sprintf("New result file [%s]", fileName)),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"filename": {
					DataType:    aws.String("String"),
					StringValue: aws.String(fileName),
				},
				"resultId": {
					DataType:    aws.String("String"),
					StringValue: aws.String(result.ID.Hex()),
				},
				mycontext.FieldCorrelationID: {
					DataType:    aws.String("String"),
					StringValue: aws.String(correlationID),
				},
			},
			MessageDeduplicationId: aws.String(fmt.Sprintf("%s-%d", result.ID.Hex(), i)),
			MessageGroupId:         aws.String(stage),
			Subject:                aws.String("New result file"),
			TopicArn:               aws.String(newResultTopic),
		})
		if err != nil {
			log.Warnf("failed to send SNS result file uploaded for [%s]: %s", fileName, err)
			// TODO: set processing to failed for this file?
		}
	}

	return nil
}
