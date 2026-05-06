variable "keycloak_url" {
  type        = string
  description = "Base URL of Keycloak."
}

variable "keycloak_username" {
  type        = string
  description = "Admin username for Terraform access."
}

variable "keycloak_password" {
  type        = string
  description = "Admin password for Terraform access."
  sensitive   = true
}

variable "realm_name" {
  type        = string
  description = "Realm used by the app."
  default     = "plants"
}

variable "frontend_base_url" {
  type        = string
  description = "Base URL where the web app is served."
  default     = "http://localhost:8080"
}

variable "oauth2_proxy_client_id" {
  type        = string
  description = "Client ID used by oauth2-proxy."
  default     = "plant-keeper-proxy"
}

variable "oauth2_proxy_client_secret" {
  type        = string
  description = "Client secret used by oauth2-proxy."
  sensitive   = true
  default     = "dev-oauth2-proxy-secret"
}
