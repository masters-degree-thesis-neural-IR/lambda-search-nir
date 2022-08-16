locals {

  lambdas_path      = "${path.module}/../bin"
  lambda_local_name = "${var.lambda_name}-${var.local_environment}"

  common_tags = {
    Environment = var.local_environment
    ManagedBy   = "Terraform"
    Service     = "Interagir"
  }

}
