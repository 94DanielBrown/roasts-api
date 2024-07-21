terraform {
  source = "tfr:///terraform-aws-modules/ecs/aws?version=5.11.2"
}

include {
  path = find_in_parent_folders()
}

inputs = {
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  fargate_capacity_providers = {
    FARGATE = {
      default_capacity_provider_strategy = {
        weight = 50
        base   = 20
      }
    }
    FARGATE_SPOT = {
      default_capacity_provider_strategy = {
        weight = 50
      }
    }
  }
}
