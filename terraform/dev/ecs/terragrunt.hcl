terraform {
  source = "tfr:///terraform-aws-modules/ecs/aws//modules/service?version=5.11.2"
}

include {
  path = find_in_parent_folders()
}

inputs = {
  cluster_name = "roasts"
  cluster_arn = "arn:aws:ecs:eu-west-1:637423178719:cluster/roasts"
  enable_execute_command = true
  name = "test"
  family = "test"

  subnet_ids = ["subnet-050ed1bda2f079958", "subnet-09dc3fb2f800c00da"]

  container_definitions = {
    ecs-sample = {
      cpu       = 512
      memory    = 1024
      essential = true
      image     = "public.ecr.aws/aws-containers/ecsdemo-frontend:776fd50"
      port_mappings = [
        {
          name          = "ecs-sample"
          containerPort = 80
          protocol      = "tcp"
        }
      ]
      readonly_root_filesystem = false
      enable_cloudwatch_logging = false
      launch_type = "FARGATE"
    }
  }
}
