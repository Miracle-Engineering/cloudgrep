terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.15.1"
    }
  }

  required_version = ">= 1.0.0"

  backend "s3" {
    region         = "us-east-1"
    bucket         = "438881294876-terraform-state"
    key            = "environment/tfstate"
    dynamodb_table = "terraform-locks"
  }
}

provider "aws" {
  region              = "us-east-1"
  allowed_account_ids = ["438881294876"]

  default_tags {
    tags = {
      IntegrationTest = "true"
    }
  }
}

module "environment" {
  source = "../modules/environment"
}
