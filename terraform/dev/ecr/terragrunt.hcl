terraform {
  source = "tfr:///terraform-aws-modules/ecr/aws?version=2.2.1"
}

include {
  path = find_in_parent_folders()
}

inputs = {
  repository_name = "roasts-api"
  image_tag_mutability = "IMMUTABLE"
  scan_on_push = true
  encryption_type = "AES256"
  repository_lifecycle_policy = jsonencode({
    rules = [
      {
        rulePriority = 1,
        description  = "Keep last 30 images",
        selection = {
          tagStatus     = "tagged",
          tagPrefixList = ["v"],
          countType     = "imageCountMoreThan",
          countNumber   = 30
        },
        action = {
          type = "expire"
        }
      }
    ]
  })
}
