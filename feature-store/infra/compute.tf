data "aws_ssm_parameter" "al2023_ami" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-arm64"
}

resource "aws_launch_template" "flashfeat_enclave_lt" {
  name_prefix   = "flashfeat-enclave"
  image_id      = data.aws_ssm_parameter.al2023_ami.value
  instance_type = "c7g.large"

  enclave_options {
    enabled = true
  }

  iam_instance_profile {
    name = aws_iam_instance_profile.flashfeat_enclave.name
  }

  network_interfaces {
    associate_public_ip_address = true
    subnet_id                   = aws_subnet.flashfeat_dev_a.id
    security_groups             = [aws_security_group.flashfeat_enclave_sg.id]
  }

  user_data = base64encode(templatefile("${path.module}/user-data.sh.tmpl", {
    s3_bucket = aws_s3_directory_bucket.flashfeat_hot.bucket
    region    = "ap-south-1"
  }))

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name        = "flashfeat-enclave"
      project     = "flashfeat"
      environment = "dev"
    }
  }
}
