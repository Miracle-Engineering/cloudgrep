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
    bucket         = "316817240772-terraform-state"
    key            = "bootstrap/tfstate"
    dynamodb_table = "terraform-locks"
  }
}

module "bootstrap" {
  source = "../modules/bootstrap"
}

output "github_actions_iam_role" {
  value = module.bootstrap.github_actions_iam_role
}
