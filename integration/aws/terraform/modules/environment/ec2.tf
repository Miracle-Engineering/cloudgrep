locals {
  ec2_instance_groups = 2
}

resource "aws_launch_template" "amz_arm" {
  count = local.ec2_instance_groups

  name_prefix   = "testing-"
  image_id      = "ami-0e449176cecc3e577"
  instance_type = "t4g.nano"

  tag_specifications {
    resource_type = "instance"

    tags = {
      test = "ec2-instance-${count.index}"
    }
  }
}


resource "aws_autoscaling_group" "test" {
  count = local.ec2_instance_groups

  name_prefix        = "testing-${count.index}-"
  availability_zones = ["us-east-1a"]
  desired_capacity   = 1
  max_size           = 1
  min_size           = 1

  launch_template {
    id      = aws_launch_template.amz_arm[count.index].id
    version = "$Latest"
  }

  tag {
    // provider default_tags doesn't support aws_autoscaling_group
    key                 = "IntegrationTest"
    value               = "true"
    propagate_at_launch = true
  }

  instance_refresh {
    strategy = "Rolling"
    preferences {
      min_healthy_percentage = 0
    }
    triggers = ["tag"]
  }
}
