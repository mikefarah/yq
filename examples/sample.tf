# main.tf

# Define required providers and minimum Terraform version
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 1.2"
}

# Configure the AWS provider
provider "aws" {
  region = var.aws_region
}

# Define an S3 bucket resource
resource "aws_s3_bucket" "example_bucket" {
  bucket = var.bucket_name

  tags = {
    Environment = "Development"
    Project     = "TerraformExample"
  }
}
