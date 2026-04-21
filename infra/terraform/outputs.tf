output "realm_issuer_url" {
  value       = "${var.keycloak_url}/realms/${keycloak_realm.plants.realm}"
  description = "OIDC issuer URL for the Go backend and frontend."
}

output "client_id" {
  value       = keycloak_openid_client.plant_keeper_web.client_id
  description = "Public client ID for the SPA login flow."
}
