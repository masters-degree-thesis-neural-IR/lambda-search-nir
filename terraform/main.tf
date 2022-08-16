terraform {
  required_version = "0.14.8"

#  backend "s3" {}

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.32.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.1.0"
    }
    template = {
      source  = "hashicorp/template"
      version = "2.2.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.2.0"
    }

  }
}

provider "aws" {
  region = var.aws_region
}

resource "random_pet" "pet" {}
data "aws_caller_identity" "current" {}
