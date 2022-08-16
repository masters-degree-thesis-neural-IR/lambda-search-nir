variable "aws_region" {
  type        = string
  description = "Região da AWS"
  default = "us-east-1"
}

variable "local_environment" {
  type        = string
  description = "type of profile"
  default = "development"
}

variable "lambda_name" {
  type        = string
  description = "name of lambda"
  default = "indexgolambda"
}

variable "lambda_runtime" {
  type        = string
  description = "lambda runtime"
  default     = "go1.x"
}