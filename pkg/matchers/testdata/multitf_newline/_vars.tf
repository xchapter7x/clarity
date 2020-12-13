variable "name" {
  description = "Account name prefix"
  type        = string
}

variable "account" {
  description = "The account type being created (prod, nonprod, tools, sandbox)"
  type        = string
}

variable "parent" {
  description = "The ID of the parent - Root account or OU"
  type        = string
}

variable "ou" {
  description = "The name of the OU - used to create email address"
  type        = string
}

variable "email_prefix" {
  description = "Email to attach to this account"
  type        = string
}

variable "email_domain" {
  description = "Email domain"
  type        = string
}

variable "role_name" {
  description = "Name of the role that will be automatically trusted from the master account"
  type        = string
}

variable "tags" {
  default = []
}

variable "email" {
  description = "Email for account"
  type = string
}