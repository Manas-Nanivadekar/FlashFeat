resource "aws_security_group" "flashfeat_enclave_sg" {
  name        = "flashfeat-enclave-sg"
  vpc_id      = aws_vpc.flashfeat_dev.id
  description = "Allow HTTP/s from given ip(office); vsock (CID 3) local only"

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [var.office_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    project     = "flashfeat"
    environment = "dev"
  }
}

variable "office_cidr" {
  description = "Public IP/CIDR of office VPN"
  type        = string
}
