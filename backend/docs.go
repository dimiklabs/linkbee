// Package main is the entry point for the Shortlink API server.
//
//	@title					Shortlink API
//	@version				1.0
//	@description			A SaaS URL shortener with analytics, QR codes, geo-routing, A/B split testing, webhooks, and API key access.
//	@contact.name			Shortlink Support
//	@license.name			MIT
//	@host					localhost:8080
//	@BasePath				/
//	@schemes				http https
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT access token. Format: Bearer <token>
//
//	@securityDefinitions.apikey	APIKeyAuth
//	@in							header
//	@name						X-API-Key
//	@description				API key for programmatic access. Format: sl_<hex>
package main
