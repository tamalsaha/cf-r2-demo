package main

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob/s3blob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	bucketName := "gocloud-demo"

	// Configure to use MinIO Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("CF_ACCESS_KEY_ID"), os.Getenv("CF_SECRET_ACCESS_KEY"), ""),
		Endpoint:         aws.String("https://a46f9a02578d51f3e8e135a14de082a0.r2.cloudflarestorage.com/gocloud-demo"),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess := session.New(s3Config)

	/*
		bucket, err := blob.OpenBucket("s3://mybucket?" +
		    "endpoint=my.minio.local:8080&" +
		    "disableSSL=true&" +
		    "s3ForcePathStyle=true")
	*/

	ctx := context.Background()

	//// blob.OpenBucket creates a *blob.Bucket from a URL.
	//bucket, err := blob.OpenBucket(ctx, "s3://my-bucket?region=us-west-1")
	//if err != nil {
	//	return err
	//}
	//defer bucket.Close()
	//
	//// Forcing AWS SDK V2.
	//bucket, err = blob.OpenBucket(ctx, "s3://my-bucket?region=us-west-1&awssdk=2")
	//if err != nil {
	//	return err
	//}
	//defer bucket.Close()

	// Create a *blob.Bucket.
	bucket, err := s3blob.OpenBucket(ctx, sess, bucketName, nil)
	if err != nil {
		return err
	}
	defer bucket.Close()

	u, err := bucket.SignedURL(ctx, "foo.txt", &blob.SignedURLOptions{
		Expiry: 10 * time.Minute,
		//Method:                   "",
		//ContentType:              "",
		//EnforceAbsentContentType: false,
		//BeforeSign:               nil,
	})
	if err != nil {
		return err
	}
	fmt.Println("Signed URL:", u)

	/*
		// Open the key "foo.txt" for writing with the default options.
		w, err := bucket.NewWriter(ctx, "foo.txt", nil)
		if err != nil {
			return err
		}
		_, writeErr := fmt.Fprintln(w, "Hello, World!")
		// Always check the return value of Close when writing.
		closeErr := w.Close()
		if writeErr != nil {
			return writeErr
		}
		if closeErr != nil {
			return closeErr
		}
	*/

	return nil
}
