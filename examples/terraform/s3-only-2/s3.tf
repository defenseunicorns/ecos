
# Create the S3 bucket with versioning and server-side  encryption
resource "aws_s3_bucket" "ecos_test_bucket_2" {
  bucket = "ecos-test-2"
}

# Block all public access
resource "aws_s3_bucket_public_access_block" "ecos_test_bucket_block_public" {
  bucket = aws_s3_bucket.ecos_test_bucket_2.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
