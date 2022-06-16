locals {
  ec2_instance_count     = 2
  elastic_ip_count       = 1
  ami_count              = 1
  ec2_keypair_count      = 1
  ec2_eni_count          = 1
  ec2_sg_count           = 1
  ec2_ebs_snapshot_count = 1
}

resource "aws_launch_template" "amz_arm" {
  count = local.ec2_instance_count

  name_prefix   = "testing-"
  image_id      = "ami-0e449176cecc3e577"
  instance_type = "t4g.nano"

  network_interfaces {
    associate_public_ip_address = false
    security_groups             = [module.vpc.default_sg_id]
  }

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

  tags = {
    "test" : "ec2-launch-template-${count.index}"
  }
}

resource "aws_autoscaling_group" "test" {
  count = local.ec2_instance_count

  name_prefix = "testing-${count.index}-"
  vpc_zone_identifier = [
    module.vpc.private_subnet_az_map["us-east-1a"]
  ]
  desired_capacity = 1
  max_size         = 1
  min_size         = 1

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

resource "aws_eip" "test" {
  count = local.elastic_ip_count

  vpc = true
  tags = {
    "test" : "ec2-ip-${count.index}"
  }
}

resource "aws_ami_copy" "test" {
  count = local.ami_count

  name              = "test-ami-copy-${count.index}"
  source_ami_id     = "ami-0cff7528ff583bf9a" // Amazon Linux 2 AMI (HVM) - Kernel 5.10, SSD Volume Type, 64-bit x86
  source_ami_region = "us-east-1"

  tags = {
    "test" : "ec2-ami-${count.index}"
  }
}

resource "aws_ebs_snapshot_copy" "test" {
  count = local.ec2_ebs_snapshot_count

  source_snapshot_id = "snap-08f1069dfde2007ba" // EBS snapshot for AMI ami-0cff7528ff583bf9a (above)
  source_region      = "us-east-1"

  tags = {
    "test" : "ec2-ebs-snapshot-${count.index}"
  }
}

resource "tls_private_key" "ec2_keypair" {
  count = local.ec2_keypair_count

  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "aws_key_pair" "test" {
  count = local.ec2_keypair_count

  key_name_prefix = "test-${count.index}-"
  public_key      = tls_private_key.ec2_keypair[count.index].public_key_openssh

  tags = {
    "test" : "ec2-keypair-${count.index}"
  }
}

resource "aws_network_interface" "test" {
  count = local.ec2_eni_count

  subnet_id = module.vpc.private_subnet_az_map["us-east-1a"]

  tags = {
    "test" : "ec2-eni-${count.index}"
  }
}

resource "aws_security_group" "test" {
  count = local.ec2_sg_count

  name_prefix = "test-${count.index}-"
  vpc_id      = module.vpc.id

  tags = {
    "test" : "ec2-sg-${count.index}"
  }
}
