resource "aws_vpc" "flashfeat_dev" {
  cidr_block           = "10.10.0.0/24"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name        = "flashfeat-dev"
    project     = "flashfeat"
    environment = "dev"
  }
}

resource "aws_subnet" "flashfeat_dev_a" {
  vpc_id                  = aws_vpc.flashfeat_dev.id
  availability_zone       = "ap-south-1a"
  cidr_block              = "10.10.0.0/26"
  map_public_ip_on_launch = true
  tags = {
    Name        = "flashfeat-dev-a"
    project     = "flashfeat"
    environment = "dev"
  }
}

resource "aws_internet_gateway" "flashfeat_dev_igw" {
  vpc_id = aws_vpc.flashfeat_dev.id
  tags   = merge(aws_vpc.flashfeat_dev.tags, { Name = "flashfeat-dev-igw" })
}

resource "aws_route_table" "flashfeat_dev" {
  vpc_id = aws_vpc.flashfeat_dev.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.flashfeat_dev_igw.id
  }
  tags = merge(aws_vpc.flashfeat_dev.tags, { Name = "flashfeat-dev-rt" })

}
