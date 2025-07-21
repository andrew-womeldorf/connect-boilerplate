variable "name" {
  description = "name for resources in this module"
  type = string
  default = "users"
}

variable "tags" {
  description = "Optional map of tags to apply to resources"
  type = map(string)
  default = {}
}
