output "realm_issuer_url" {
  value       = "${var.keycloak_url}/realms/${keycloak_realm.plants.realm}"
  description = "OIDC issuer URL for oauth2-proxy."
}

output "client_id" {
  value       = keycloak_openid_client.plant_keeper_proxy.client_id
  description = "Client ID for oauth2-proxy."
}

output "client_secret" {
  value       = keycloak_openid_client.plant_keeper_proxy.client_secret
  description = "Client secret for oauth2-proxy."
  sensitive   = true
}
