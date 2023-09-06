variable "function_name" {
  type        = string
  description = "The name to give the function"
}

variable "function_runtime" {
  type        = string
  description = "Identifier of the function's runtime"
}

variable "function_memory" {
  type        = number
  description = "Amount of memory to allocate to the function"
}

variable "paramstore_prefix" {
  type        = string
  description = "Prefix to allow ssm:GetParameterByPath action to. Alternatively, use the iam_role_name output to attach your own policies"
}

variable "log_retention_in_days" {
  type        = number
  default     = 7
  description = "Number of days the events in the log streams will be retained before being purged. (Default: 7)"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID to attach security groups to"
  default     = ""
}

variable "vpc_subnet_ids" {
  type        = list(string)
  description = "List of Subnet IDs to associate with the lambda"
  default     = []
}

variable "vpc_security_group_ids" {
  type        = list(string)
  description = "List of Security Group IDs to attach to the lambda function"
  default     = []
}

variable "additional_role_policies" {
  type        = map(string)
  default     = {}
  description = "Additional Policies that should be attached to the IAM Execution Role created for this lambda"
}

variable "environment_variables" {
  type        = map(string)
  description = "list of environment variables to inject into the lambda"
}
