module "eventbridge_scheduler" {
  source  = "terraform-aws-modules/eventbridge/aws"
  version = "3.12.0"

  create_bus = false
  // bus_name   = "default" # "default" bus already support schedule_expression in rules

  attach_lambda_policy = true
  lambda_target_arns   = [module.ffnotifier_lambda_function.lambda_function_arn]

  schedule_groups = {
    dev = {
      name_prefix = "tmp-dev-"
    }
    prod = {
      name = "prod"
      tags = {
        Env = "SuperProd"
      }
    }
  }

  schedules = {
    lambda-cron = {
      group_name          = "dev"
      description         = "Trigger for a Lambda"
      schedule_expression = "cron(0/5 * * * ? *)"
      timezone            = "Europe/London"
      arn                 = module.ffnotifier_lambda_function.lambda_function_arn
      input               = jsonencode({ "job" : "cron-by-rate" })
    }
    # prod-lambda-cron = {
    #   group_name          = "prod"
    #   schedule_expression = "rate(10 hours)"
    #   arn                 = module.ffnotifier_lambda_function.lambda_function_arn
    # }
    # kinesis-cron = {
    #   group_name          = "prod"
    #   schedule_expression = "rate(10 hours)"
    #   arn                 = aws_kinesis_stream.this.arn
    #   partition_key       = "foo"
    # }
  }
}