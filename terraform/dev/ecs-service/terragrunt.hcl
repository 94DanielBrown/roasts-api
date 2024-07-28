terraform {
  source = "tfr:///terraform-aws-modules/ecs/aws//modules/service?version=5.11.2"
}

include {
  path = find_in_parent_folders()
}

inputs = {
  enable_execute_command = true
  name = "roasts-api"
  family = "roasts"
}
