locals {
  global_vars = yamldecode(file("${find_in_parent_folders("global.yaml")}"))
  env_vars    = yamldecode(file("${find_in_parent_folders("env.yaml")}"))

  default_tags = {
    Name = "roasts"
    Environment = "override",
  }
  override_tags = try(yamldecode(file("${find_in_parent_folders("tags.yaml")}")))
  tags          = merge(local.default_tags, local.override_tags)

  all_inputs = merge(
    local.env_vars, local.global_vars
  )
}

inputs = local.all_inputs

generate "versions" {
  path      = "versions.tf"
  if_exists = "overwrite"
  contents  = <<EOF
terraform {
  required_version = "${local.global_vars.terraform_version}"
  required_providers {
    aws        = {
      source  = "hashicorp/aws"
      version = "${local.global_vars.terraform_provider_aws}"
    }
  }
}
EOF
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite"
  contents  = <<EOF
provider "aws" {
  default_tags {
    tags = ${jsonencode(local.tags)}
  }
}
EOF
}


generate "backend" {
  path = "backend.tf"
  if_exists = "overwrite"
  contents = <<EOF
terraform {
  backend "s3" {
    bucket = "dan-${local.env_vars["env"]}-terraform-state"
    key = "${path_relative_to_include()}/terraform.tfstate"
    region = "eu-west-1"
    encrypt        = true
    dynamodb_table = "my-lock-table"
  }
}
EOF
}
