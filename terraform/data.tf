data "aws_iam_policy_document" "allow_dynamodb_basic" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:DeleteItem",
      "dynamodb:Query",
    ]
    resources = [
      aws_dynamodb_table.sessions.arn,
      "${aws_dynamodb_table.sessions.arn}/*",
      aws_dynamodb_table.timers.arn,
      "${aws_dynamodb_table.timers.arn}/*",
      aws_dynamodb_table.users.arn,
      "${aws_dynamodb_table.users.arn}/*",
    ]
  }
}
