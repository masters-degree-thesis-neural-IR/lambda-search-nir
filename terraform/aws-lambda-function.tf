resource "archive_file" "lambda" {
  source_dir  = local.lambdas_path
  output_path = "files/${local.lambda_local_name}.zip"
  type        = "zip"
}

resource "aws_lambda_function" "lambda" {

  filename      = archive_file.lambda.output_path
  function_name = local.lambda_local_name
  role          = aws_iam_role.lambda-role.arn
  handler       = "app"

  description = "Exemplo de uma lambda com terraform e GO"

  source_code_hash = archive_file.lambda.output_base64sha256

  runtime = var.lambda_runtime

#  environment {
#    variables = {
#      REGION = var.aws_region
#      SURVEY_URL = var.survey_url
#      CONVERSATION_URL = var.conversation_url
#      SURVEY_RESPONSE_COLLECTOR_TABLE = local.survey-response-collector-table
#      CALLBACK_REQUEST_TABLE = local.callback-request-table
#      CONVERSATION_TABLE = local.conversation-table
#    }
#  }

  tags = local.common_tags
}
