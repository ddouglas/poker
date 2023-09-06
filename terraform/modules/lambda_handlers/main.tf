resource "aws_lambda_permission" "allow_apigateway" {
  for_each      = toset(var.routes)
  statement_id  = sha1(each.value)
  function_name = var.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.api_execution_arn}/*/${replace(each.value, " ", "")}"
  action        = "lambda:InvokeFunction"
}

resource "aws_apigatewayv2_route" "lambda" {
  for_each = toset(var.routes)

  api_id    = var.apigw_id
  route_key = each.value

  target = "integrations/${aws_apigatewayv2_integration.lambda.id}"

  depends_on = [
    aws_apigatewayv2_integration.lambda,
  ]
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id      = var.apigw_id
  description = "Integration for ${var.function_name}"

  integration_type   = "AWS_PROXY"
  integration_uri    = var.function_invoke_arn
  integration_method = "POST"

  passthrough_behavior = "WHEN_NO_MATCH"

}
