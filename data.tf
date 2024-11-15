data "aws_organizations_organization" "this" {}

data "aws_caller_identity" "current" {}

data "archive_file" "zip" {
  type        = "zip"
  source_dir  = "${path.module}/ffnotifier"
  output_path = "${path.module}/ffnotifier.zip"
}