package amplify

const (
	SNSOrderImported      = "/sns/orderImported/arn"
	SNSOrderValidated     = "/sns/orderValidated/arn"
	SNSAttachmentRejected = "/sns/attachmentRejected/arn"
	SNSNewResult          = "/sns/newResult/arn"
	SNSNewOrder           = "/sns/newOrder/arn"
	SNSResultFound        = "/sns/resultFound/arn"

	S3Result     = "/result/s3BucketName"
	S3Attachment = "/attachment/s3BucketName"
	S3Invoice    = "/invoice/s3BucketName"
	S3Temp       = "/temp/s3BucketName"

	SNSShopOrderPaid      = "/shop/sns/order/paid/arn"
	SNSShopPaymentPending = "/shop/sns/order/payment/pending/arn"

	S3ShopFile = "/shop/file/s3BucketName"
)
