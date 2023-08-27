terraform {
  backend "s3" {
    region         = "us-east-1"
    bucket         = "ots-terraform-state-us-east-1"
    key            = "poker.tfstate"
    dynamodb_table = "ots-terraform-state-us-east-1"
  }
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
      Repository = "https://github.com/ddouglas/rv-poker"
    }
  }
}

variable "region" {
  default = "us-east-1"
  type    = string
}
