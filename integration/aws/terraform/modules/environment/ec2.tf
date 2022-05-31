locals {
  ec2_instance_count = 2
}

resource "aws_launch_template" "amz_arm" {
  count = local.ec2_instance_count

  name_prefix   = "testing-"
  image_id      = "ami-0e449176cecc3e577"
  instance_type = "t4g.nano"

  tag_specifications {
    resource_type = "instance"

    tags = {
      // provider default_tags doesn't support aws_autoscaling_group
      IntegrationTest = "true"
      test            = "ec2-instance-${count.index}"
    }
  }

  tag_specifications {
    resource_type = "volume"

    tags = {
      IntegrationTest = "true"
      test            = "ec2-volume-${count.index}"
    }
  }
}


resource "aws_autoscaling_group" "test" {
  count = local.ec2_instance_count

  name_prefix        = "testing-${count.index}-"
  availability_zones = ["us-east-1a"]
  desired_capacity   = 1
  max_size           = 1
  min_size           = 1

  launch_template {
    id      = aws_launch_template.amz_arm[count.index].id
    version = aws_launch_template.amz_arm[count.index].latest_version
  }

  instance_refresh {
    strategy = "Rolling"
    preferences {
      instance_warmup        = 1
      min_healthy_percentage = 0
    }
    triggers = ["tag"]
  }
}
