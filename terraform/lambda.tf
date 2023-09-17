locals {
  default_envs = {
    SSM_PREFIX = local.ssm_prefix

    MODE               = "lambda"
    APP_URL            = "https://poker.onetwentyseven.dev"
    AUTH0_CLIENT_ID    = "TBxKbUrFIIROYisVOTlHq1lb2BojSBFd"
    AUTH0_DOMAIN       = "onetwentyseven.us.auth0.com"
    AUTH0_CALLBACK_URL = "https://poker.onetwentyseven.dev/login"
    ENVIRONMENT        = "production"
  }



}

module "lambda" {
  source = "./modules/lambda"

  environment_variables = local.default_envs
  function_memory       = 128
  function_name         = "poker-handler"
  function_runtime      = "provided.al2"
  log_retention_in_days = 3
  paramstore_prefix     = local.ssm_prefix
  additional_role_policies = {
    allow_dynamodb_basic   = data.aws_iam_policy_document.allow_dynamodb_basic.json
    allow_s3_full          = data.aws_iam_policy_document.allow_s3_full.json
    allow_polly_synthesize = data.aws_iam_policy_document.allow_polly_synthesize.json
  }

}

module "routes" {
  depends_on = [module.lambda]
  source     = "./modules/lambda_handlers"

  apigw_id            = aws_apigatewayv2_api.poker.id
  api_execution_arn   = aws_apigatewayv2_api.poker.execution_arn
  function_name       = module.lambda.function_name
  function_invoke_arn = module.lambda.function_invoke_arn

  routes = [
    "GET /",
    "GET /login",
    "GET /logout",
    "GET /static/{proxy+}",

    "GET /dashboard",
    "GET /dashboard/timers",
    "GET /dashboard/timers/new",
    "POST /dashboard/timers/new",

    "GET /dashboard/timers/{timerID}",
    "DELETE /dashboard/timers/{timerID}",

    "GET /play/{timerID}",

    "GET /play/{timerID}/levels/reset",
    "GET /play/{timerID}/levels/next",
    "GET /play/{timerID}/levels/previous",

    "GET /dashboard/timers/{timerID}/levels/new",
    "POST /dashboard/timers/{timerID}/levels/new",

    "GET /dashboard/timers/{timerID}/levels/{levelID}",
    "POST /dashboard/timers/{timerID}/levels/{levelID}",
    "DELETE /dashboard/timers/{timerID}/levels/{levelID}",

  ]


}
