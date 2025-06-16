terraform {
  backend "s3" {
    bucket         = "flashfeat-tfstate-ap-south-1"
    key            = "dev/terraform.tfstate"
    region         = "ap-south-1"
    dynamodb_table = "tfstate-lock"
    encrypt        = true
  }
}
