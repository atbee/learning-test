variable "access_key_id" {
  type        = string
  description = "AWS Access Key ID"
}

variable "secret_access_id" {
  type        = string
  description = "AWS Secret"
}

variable "region" {
  type        = string
  description = "AWS Region"
  default     = "ap-southeast-1"
}

variable "product_area" {
  type        = string
  description = "Product Area"
  default     = "morchana-static-qr-code"
}

variable "environment" {
  type        = string
  description = "Product Environment"
  default     = "production"
}

variable "team" {
  type        = string
  description = "Team"
  default     = "ODDS"
}

variable "ec2_deployment_iam_policy_arn" {
  type        = list
  description = "IAM Policy to be Attached to EC2 Deployment Role"
}

variable "infrastructure_version" {
  type        = string
  description = "Infrastructure Version"
}
