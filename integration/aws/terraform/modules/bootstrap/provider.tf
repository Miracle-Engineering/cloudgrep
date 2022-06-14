terraform {
  required_providers {
    http = {
      source  = "hashicorp/http"
      version = "2.1.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "3.4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}
