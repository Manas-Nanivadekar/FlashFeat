# feature-store/infra/asg.tf
resource "aws_autoscaling_group" "flashfeat_enclave_asg" {
  name                = "flashfeat-enclave-asg"
  desired_capacity    = 1
  max_size            = 1
  min_size            = 1
  vpc_zone_identifier = [aws_subnet.flashfeat_dev_a.id]
  health_check_type   = "EC2"
  launch_template {
    id      = aws_launch_template.flashfeat_enclave_lt.id
    version = "$Latest"
  }

  tag {
    key                 = "project"
    value               = "flashfeat"
    propagate_at_launch = true
  }

  tag {
    key                 = "environment"
    value               = "dev"
    propagate_at_launch = true
  }

}
