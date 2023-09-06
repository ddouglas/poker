variable "routes" {
  type = list(string)
}

variable "apigw_id" {
  type = string
}

variable "api_execution_arn" {
  type = string
}

variable "function_name" {
  type = string
}

variable "function_invoke_arn" {
  type = string
}
