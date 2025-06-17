resource "aws_kms_key" "flashfeat_sse" {
  description             = "Flashfeat envelope-encryption key"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid       = "AllowAccountUse",
        Effect    = "Allow",
        Principal = { AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root" },
        Action    = "kms:*",
        Resource  = "*"
      },
      {
        Sid       = "DenyCrossAccount",
        Effect    = "Deny",
        Principal = "*",
        Action    = "kms:*",
        Resource  = "*",
        Condition = {
          StringNotEquals = {
            "aws:PrincipalAccount" = data.aws_caller_identity.current.account_id
          }
        }
      }
    ]
  })

  tags = {
    project     = "flashfeat"
    environment = "dev"
  }
}

resource "aws_kms_alias" "flashfeat_sse" {
  name          = "alias/flashfeat-sse"
  target_key_id = aws_kms_key.flashfeat_sse.id
}
