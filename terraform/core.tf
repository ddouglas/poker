terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      Repository = "https://github.com/ddouglas/poker"
    }
  }
}


locals {
  ssm_prefix     = "/poker"
  default_domain = "poker.onetwentyseven.dev"
}

variable "region" {
  default = "us-east-1"
  type    = string
}
