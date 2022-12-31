package sns

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssns "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/google/uuid"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/sns"
	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

func OrderValidated(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)
	log.Infof("Send SNS order validated for order [%s]", id)

	snsARN, err := ssm.GetParameter(ctx, "/sns/orderValidated/arn", false)
	if err != nil {
		return err
	}
	snsClient, err := sns.NewClient(ctx)
	if err != nil {
		return err
	}
	correlationID := ctx.Value(mycontext.FieldCorrelationID)
	if correlationID == nil {
		correlationID = uuid.New().String()
	}
	message := fmt.Sprintf("Order finalized [%s]", id)
	_, err = snsClient.Publish(ctx, &awssns.PublishInput{
		Message: aws.String(message),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(id),
			},
			mycontext.FieldCorrelationID: {
				DataType:    aws.String("String"),
				StringValue: aws.String(correlationID.(string)),
			},
		},
		Subject:  aws.String(message),
		TopicArn: aws.String(snsARN),
	})
	if err != nil {
		log.Errorf("Failed to send SNS order validated [%s]: %s", id, err)
		return err
	}
	return nil
}

func AttachmentRejected(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)
	log.Infof("Send SNS attachment rejected [%s]", id)

	snsARN, err := ssm.GetParameter(ctx, "/sns/attachmentRejected/arn", false)
	if err != nil {
		return err
	}
	snsClient, err := sns.NewClient(ctx)
	if err != nil {
		return err
	}
	correlationID := ctx.Value(mycontext.FieldCorrelationID)
	if correlationID == nil {
		correlationID = uuid.New().String()
	}
	message := fmt.Sprintf("Attachment rejected for order [%s]", id)
	_, err = snsClient.Publish(ctx, &awssns.PublishInput{
		Message: aws.String(message),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(id),
			},
			mycontext.FieldCorrelationID: {
				DataType:    aws.String("String"),
				StringValue: aws.String(correlationID.(string)),
			},
		},
		MessageGroupId: aws.String(id),
		Subject:        aws.String(message),
		TopicArn:       aws.String(snsARN),
	})
	if err != nil {
		log.Errorf("Failed to send SNS attachment rejected [%s]: %s", id, err)
		return err
	}
	return nil
}
