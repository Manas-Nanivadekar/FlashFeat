provider "aws" {
  region = "ap-south-1"

  default_tags {
    tags = {
      project = "flashfeat"
      environment = "dev"
    }
  }
}