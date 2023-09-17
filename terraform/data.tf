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

data "aws_iam_policy_document" "allow_s3_full" {
  statement {
    effect = "Allow"
    actions = [
      "s3:*",
    ]
    resources = [
      aws_s3_bucket.poker_audio_cache.arn,
      "${aws_s3_bucket.poker_audio_cache.arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "allow_polly_synthesize" {
  statement {
    effect = "Allow"
    actions = [
      "polly:SynthesizeSpeech",
    ]
    resources = [
      "*",
    ]
  }
}
