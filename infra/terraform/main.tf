provider "keycloak" {
  client_id = "admin-cli"
  username  = var.keycloak_username
  password  = var.keycloak_password
  realm     = "master"
  url       = var.keycloak_url
}

locals {
  redirect_uris = [
    "${var.frontend_base_url}/oauth2/callback",
  ]
}

resource "keycloak_realm" "plants" {
  realm        = var.realm_name
  enabled      = true
  display_name = "Plants"
}

resource "keycloak_openid_client" "plant_keeper_proxy" {
  realm_id                     = keycloak_realm.plants.id
  client_id                    = var.oauth2_proxy_client_id
  name                         = "Plant Keeper OAuth2 Proxy"
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  client_secret                = var.oauth2_proxy_client_secret
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  service_accounts_enabled     = false
  valid_redirect_uris          = local.redirect_uris
  web_origins                  = [var.frontend_base_url]
  base_url                     = var.frontend_base_url
}

resource "keycloak_role" "app_user" {
  realm_id = keycloak_realm.plants.id
  name     = "app-user"
}
