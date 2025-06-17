data "aws_iam_policy_document" "ec2_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "flashfeat_enclave" {
  name               = "flashfeat-ec2-enclave"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume.json

  tags = {
    project     = "flashfeat"
    environment = "dev"
  }
}

resource "aws_iam_role_policy" "flashfeat_enclave" {
  role = aws_iam_role.flashfeat_enclave.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = ["s3:GetObject"],
        Resource = [
          aws_s3_directory_bucket.flashfeat_hot.arn,
          "${aws_s3_directory_bucket.flashfeat_hot.arn}/*"
        ]
      },
      {
        Effect   = "Allow",
        Action   = ["kms:GenerateDataKey*", "kms:Decrypt"],
        Resource = aws_kms_key.flashfeat_sse.arn
      }
    ]
  })
}

resource "aws_iam_instance_profile" "flashfeat_enclave" {
  name = "flashfeat-ec2-enclave"
  role = aws_iam_role.flashfeat_enclave.name
}
