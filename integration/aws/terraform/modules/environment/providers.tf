

terraform {
  required_providers {
    tls = {
      source  = "hashicorp/tls"
      version = "3.4.0"
    }
  }
}

provider "tls" {}
