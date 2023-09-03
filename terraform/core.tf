terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
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

variable "region" {
  default = "us-east-1"
  type    = string
}
