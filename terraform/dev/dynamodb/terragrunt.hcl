terraform {
  source = "tfr:///terraform-aws-modules/dynamodb-table/aws?version=4.0.0"
}

include {
  path = find_in_parent_folders()
}

// TODO - Apply creates new table but doesn't destroy old one
inputs = {
  name       = "roasts"
  hash_key   = "PK"
  range_key = "SK"
  attributes = [
    {
      name = "PK"
      type = "S"
    },
    {
      name = "SK"
      type = "S"
    }
  ]
}