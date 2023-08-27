resource "aws_dynamodb_table" "sessions" {
  name         = "poker-sessions-${var.region}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ID"

  attribute {
    name = "ID"
    type = "S"
  }
}


resource "aws_dynamodb_table" "users" {
  name         = "poker-users-${var.region}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ID"


  attribute {
    name = "ID"
    type = "S"
  }

  attribute {
    name = "Email"
    type = "S"
  }

  global_secondary_index {
    hash_key        = "Email"
    name            = "email-index"
    projection_type = "ALL"
  }


}


output "session_table_name" {
  value = aws_dynamodb_table.sessions.name
}
