resource "aws_dynamodb_table" "sessions" {
  name         = "poker-sessions-${var.region}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ID"

  attribute {
    name = "ID"
    type = "S"
  }
}


output "session_table_name" {
  value = aws_dynamodb_table.sessions.name
}
