package s3_pkg

import (
	"context"
	"net/http"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/labstack/echo/v4"
)

type S3Service struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
}

func (s *S3Service) HandleInit(ctx echo.Context) error {

	filename := ctx.QueryParam("filename")

	input := s3.CreateMultipartUploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filename),
	}
	out, err := s.Client.CreateMultipartUpload(context.Background(), &input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{"uploadId": *out.UploadId, "key": *out.Key})
}

func (s *S3Service) HandlePresign(ctx echo.Context) error {

	var body struct {
		UploadID   string `json:"uploadId"`
		Key        string `json:"key"`
		PartNumber int32  `json:"partNumber"`
	}

	if err := ctx.Bind(&body); err != nil {
		return err
	}

	presignedReq, err := s.PresignClient.PresignUploadPart(context.Background(), &s3.UploadPartInput{
		Bucket:     aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:        aws.String(body.Key),
		UploadId:   aws.String(body.UploadID),
		PartNumber: aws.Int32(body.PartNumber),
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"url": presignedReq.URL})
}
func (s *S3Service) HandleComplete(c echo.Context) error {
	var body struct {
		UploadID string                `json:"uploadId"`
		Key      string                `json:"key"`
		Parts    []types.CompletedPart `json:"parts"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	// S3 requires parts to be sorted by PartNumber
	sort.Slice(body.Parts, func(i, j int) bool {
		return *body.Parts[i].PartNumber < *body.Parts[j].PartNumber
	})

	res, err := s.Client.CompleteMultipartUpload(context.Background(), &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:      aws.String(body.Key),
		UploadId: aws.String(body.UploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: body.Parts,
		},
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "result": *res})
}
