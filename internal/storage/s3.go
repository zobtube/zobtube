package storage

import (
	"bytes"
	"context"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 implements Storage using S3-compatible object storage.
type S3 struct {
	client *s3.Client
	bucket string
	prefix string // optional key prefix for all objects
}

// NewS3 returns a Storage that uses the given S3 client, bucket, and optional key prefix.
func NewS3(client *s3.Client, bucket, prefix string) *S3 {
	return &S3{client: client, bucket: bucket, prefix: prefix}
}

// contentTypeByPath returns the Content-Type for S3 PutObject based on file extension.
// Covers common image and video types; returns empty string if unknown.
func contentTypeByPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return ""
	}
	// Ensure common video/image types are recognized (some systems' mime.types omit these)
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".m4v":
		return "video/x-m4v"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return mime.TypeByExtension(ext)
	}
}
func (s *S3) key(path string) string {
	path = strings.TrimPrefix(path, "/")
	if s.prefix == "" {
		return path
	}
	if path == "" {
		return strings.TrimSuffix(s.prefix, "/")
	}
	return strings.TrimSuffix(s.prefix, "/") + "/" + path
}

// Open opens the object at path for reading (streams from S3).
func (s *S3) Open(path string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key(path)),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

// s3WriteCloser buffers writes and uploads on Close.
type s3WriteCloser struct {
	s3     *S3
	path   string
	buf    bytes.Buffer
	closed bool
}

func (w *s3WriteCloser) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *s3WriteCloser) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true
	input := &s3.PutObjectInput{
		Bucket: aws.String(w.s3.bucket),
		Key:    aws.String(w.s3.key(w.path)),
		Body:   bytes.NewReader(w.buf.Bytes()),
	}
	if ct := contentTypeByPath(w.path); ct != "" {
		input.ContentType = aws.String(ct)
	}
	_, err := w.s3.client.PutObject(context.Background(), input)
	return err
}

// Create creates a new object at path; writes are buffered and uploaded on Close.
func (s *S3) Create(path string) (io.WriteCloser, error) {
	return &s3WriteCloser{s3: s, path: path}, nil
}

// Delete removes the object at path.
func (s *S3) Delete(path string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key(path)),
	})
	return err
}

// Exists returns true if the object at path exists.
func (s *S3) Exists(path string) (bool, error) {
	_, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key(path)),
	})
	if err != nil {
		// Check for 404
		return false, nil // treat any error as not exists for simplicity
	}
	return true, nil
}

// MkdirAll is a no-op for S3 (keys are flat).
func (s *S3) MkdirAll(path string) error {
	return nil
}

// List returns entries under prefix. Names are relative to the given prefix.
func (s *S3) List(prefix string) ([]Entry, error) {
	fullPrefix := s.key(prefix)
	if fullPrefix != "" && fullPrefix[len(fullPrefix)-1] != '/' {
		fullPrefix += "/"
	}
	var entries []Entry
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(fullPrefix),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			name := key
			if len(fullPrefix) > 0 && len(key) > len(fullPrefix) {
				name = key[len(fullPrefix):]
			}
			modTime := time.Time{}
			if obj.LastModified != nil {
				modTime = *obj.LastModified
			}
			size := int64(0)
			if obj.Size != nil {
				size = *obj.Size
			}
			entries = append(entries, Entry{
				Name:    name,
				Size:    size,
				ModTime: modTime,
				IsDir:   false,
			})
		}
	}
	return entries, nil
}

// PresignGet returns a presigned GET URL for the object at path, valid for expiry.
// Implements PreviewableStorage.
func (s *S3) PresignGet(ctx context.Context, path string, expiry time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client, func(po *s3.PresignOptions) {
		po.Expires = expiry
	})
	result, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key(path)),
	})
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
