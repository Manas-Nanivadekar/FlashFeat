data "aws_availability_zone" "zone" {
  name = "ap-south-1a"
}

locals {
  az_id       = data.aws_availability_zone.zone.id
  bucket_name = "flashfeat-hot--${local.az_id}--x-s3"
}

resource "aws_s3_directory_bucket" "flashfeat_hot" {
  bucket = local.bucket_name
  location {
    name = local.az_id
  }

  force_destroy   = true
  data_redundancy = "SingleAvailabilityZone" # This is the default, but explicitly set for clarity
}

data "aws_caller_identity" "current" {}
