resource "aws_iam_role" "ec2_calling_services_for_deployment" {
  name               = "EC2CallingServicesForDeployment"
  description        = "Allow EC2 instances to call AWS services for deployment."
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "role_policy_attachment" {
  count      = length(var.ec2_deployment_iam_policy_arn)
  role       = aws_iam_role.ec2_calling_services_for_deployment.name
  policy_arn = element(var.ec2_deployment_iam_policy_arn, count.index)
}

resource "aws_iam_instance_profile" "ec2_calling_services_for_deployment" {
  name = "ec2_calling_services_for_deployment_instance_profile"
  role = aws_iam_role.ec2_calling_services_for_deployment.name
}
