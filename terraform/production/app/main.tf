provider "aws" {
  access_key = var.access_key_id
  secret_key = var.secret_access_id
  region     = var.region
}

data "aws_alb_target_group" "selected" {
  name = "${var.product_area}-tg"
}

data "aws_subnet" "selected" {
  filter {
    name   = "tag:Name"
    values = ["${var.product_area}-public-subnet-zone-a"]
  }
}

data "aws_security_group" "selected" {
  filter {
    name   = "tag:Name"
    values = [var.product_area]
  }
}

resource "aws_alb_target_group_attachment" "tg_attachment" {
  target_group_arn = data.aws_alb_target_group.selected.arn
  target_id        = aws_instance.ec2.id
  port             = 80
}

resource "aws_instance" "ec2" {
  associate_public_ip_address = true
  subnet_id                   = data.aws_subnet.selected.id
  security_groups             = [data.aws_security_group.selected.id]
  ami                         = "ami-0f7719e8b7ba25c61"
  instance_type               = "t2.medium"
  key_name                    = "morchana-static-qr-code"
  iam_instance_profile        = aws_iam_instance_profile.ec2_calling_services_for_deployment.name
  user_data                   = file("bootstrap.sh")

  root_block_device {
    volume_type = "gp2"
    volume_size = 30
  }

  volume_tags = {
    Name         = var.product_area
    environment  = var.environment
    product-area = var.product_area
  }

  tags = {
    Name                  = var.product_area
    InfrastructureVersion = var.infrastructure_version
    environment           = var.environment
    product-area          = var.product_area
    team                  = var.team
  }

  lifecycle {
    create_before_destroy = true
  }
}
