# Linkbee

A high-performance, production-ready URL shortener platform with QR code generation, analytics, and subscription-based access tiers.

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Subscription Tiers](#subscription-tiers)
- [Functional Requirements](#functional-requirements)
  - [URL Shortening](#url-shortening)
  - [QR Code Generation](#qr-code-generation)
  - [Analytics & Reporting](#analytics--reporting)
  - [User Management & Authentication](#user-management--authentication)
  - [Subscription & Billing](#subscription--billing)
  - [Admin Panel](#admin-panel)
  - [Landing Page](#landing-page)
  - [Frontend Application](#frontend-application)
  - [API](#api)
- [Non-Functional Requirements](#non-functional-requirements)
- [Architecture Overview](#architecture-overview)
- [Data Models](#data-models)
- [API Design](#api-design)
- [Security Requirements](#security-requirements)
- [Infrastructure & DevOps](#infrastructure--devops)
- [Development Roadmap](#development-roadmap)

---

## Overview

**Linkbee** is a SaaS URL shortener platform that allows individuals, teams, and enterprises to create, manage, and track short URLs and QR codes. It provides a branded experience, deep analytics, and a subscription model to monetize access to advanced features.

---

## Tech Stack

| Layer            | Technology                                                |
|------------------|-----------------------------------------------------------|
| Backend          | Go (Golang)                                               |
| Frontend         | Vue.js 3 (Composition API)                                |
| Database         | PostgreSQL 18                                             |
| Cache / Queue    | Valkey (Redis-compatible)                                 |
| Auth             | JWT + OAuth2 (Google, GitHub)                             |
| Payments         | Lemon Squeezy                                             |
| Storage          | S3-compatible (QR code images)                            |
| Containerization | Docker + Docker Compose                                   |
| Reverse Proxy    | Nginx                                                     |
| CI/CD            | GitHub Actions                                            |
| Analytics        | Low cost time series database may something like postgres |
| Logging          | Structured JSON logs (Zap)                                |

---


Backend directory and code sturcture should be like this /Users/shafikshaon/workplace/development/projects/kothagpt/auth
and frontend code structure like this /Users/shafikshaon/workplace/development/projects/kothagpt/client

You can copy user registration and sigin with google, github, facebook from /Users/shafikshaon/workplace/development/projects/kothagpt/auth

## Features

### Core
- Create short URLs with auto-generated (length 5) or custom slugs
- Custom branded domains (BYOD — Bring Your Own Domain)
- QR code generation (PNG, SVG) for any short URL
- Bulk URL creation via CSV/EXCEL import
- Link expiration (date/time or click count limit)
- Password-protected links
- Link preview before redirect
- UTM parameter builder

### Analytics
- Click tracking (total clicks, unique clicks)
- Geographic breakdown (country, city)
- Device & browser breakdown
- Referrer tracking
- Time-series click graphs
- Export analytics as CSV/PDF

### Management
- Dashboard to manage all links
- Link tagging and folder/campaign organization
- Team workspaces with role-based access control (RBAC)
- Link health monitoring (detect broken destination URLs)

### Monetization
- Subscription plans with tiered feature access
- Per-seat billing for team plans
- Usage-based overage billing
- Stripe integration for recurring payments and invoicing

### Developer
- REST API with API key authentication
- Webhooks for link events (click, creation, expiry)
- OpenAPI / Swagger documentation

---

## Subscription Tiers

### Free
- 30 short links per month
- 1,000 clicks tracked per month
- Basic click analytics (total clicks only)
- QR code generation (PNG only)
- No custom domain
- Linkbee-branded links only
- 30-day link history

### Starter — $9/month
- 500 short links per month
- 50,000 clicks tracked per month
- Full analytics (geo, device, referrer)
- QR code generation (PNG + SVG)
- 1 custom domain
- Link expiration
- Password-protected links
- 1-year link history
- API access (1,000 requests/day)

### Pro — $29/month
- 5,000 short links per month
- 500,000 clicks tracked per month
- Full analytics + CSV export
- QR code customization (logo, colors, shapes)
- 3 custom domains
- UTM builder
- Bulk URL import
- Link previews
- Webhooks
- API access (10,000 requests/day)
- 1-year link history

### Business — $99/month
- Unlimited short links
- 5,000,000 clicks tracked per month
- Full analytics + PDF/CSV export
- Advanced QR code customization
- 10 custom domains
- Team workspace (up to 10 seats)
- RBAC (Admin, Editor, Viewer roles)
- Link health monitoring
- Priority support
- API access (100,000 requests/day)
- Unlimited link history

### Enterprise — Custom pricing
- Everything in Business
- Unlimited seats
- Unlimited custom domains
- SSO (SAML 2.0 / OIDC)
- Dedicated infrastructure option
- SLA guarantee (99.99% uptime)
- Custom data retention policies
- Audit logs
- Dedicated account manager
- Custom contract & invoicing

---

## Functional Requirements

### URL Shortening

- **FR-URL-01**: The system shall generate a unique short code (5 characters, alphanumeric) for each submitted URL.
- **FR-URL-02**: Users shall be able to provide a custom slug (alias) subject to availability and plan limits.
- **FR-URL-03**: The system shall validate destination URLs for format correctness and optionally for reachability.
- **FR-URL-04**: The system shall redirect visitors from short URL to destination URL with HTTP 301 (permanent) or 302 (temporary) redirects, configurable per link.
- **FR-URL-05**: Links shall support an optional expiration date/time after which they return HTTP 410 Gone.
- **FR-URL-06**: Links shall support an optional maximum click count after which they expire.
- **FR-URL-07**: Password-protected links shall prompt the visitor for a password before redirecting.
- **FR-URL-08**: The system shall support bulk URL creation via CSV file upload.
- **FR-URL-09**: Users shall be able to edit the destination URL of an existing short link without changing the slug.
- **FR-URL-10**: Users shall be able to delete links, which disables the redirect and frees the slug after a grace period.
- **FR-URL-11**: The system shall support custom branded domains — users can configure their own domain to serve short links.
- **FR-URL-12**: The system shall append UTM parameters to destination URLs via a built-in UTM builder.

### QR Code Generation

- **FR-QR-01**: The system shall generate a QR code for every short URL on demand.
- **FR-QR-02**: QR codes shall be available in PNG and SVG formats.
- **FR-QR-03**: Pro and above plans shall support QR code customization: foreground color, background color, embedded logo, and corner styles.
- **FR-QR-04**: QR codes shall be downloadable directly from the dashboard.
- **FR-QR-05**: QR codes shall remain valid as long as the short URL exists; updates to the destination URL shall not require a new QR code.
- **FR-QR-06**: The system shall store generated QR code images in object storage (S3-compatible) and serve them via CDN.
- **FR-QR-07**: QR code scans shall be tracked as click events with source tagged as `qr`.

### Analytics & Reporting

- **FR-AN-01**: Every redirect shall record a click event capturing: timestamp, IP address (hashed for privacy), country, city, device type, OS, browser, and referrer.
- **FR-AN-02**: The dashboard shall display total clicks and unique clicks per link.
- **FR-AN-03**: The system shall provide time-series click data aggregated by hour, day, week, and month.
- **FR-AN-04**: Geographic analytics shall display click distribution by country and city on a map and table.
- **FR-AN-05**: Device analytics shall break down clicks by device type (desktop, mobile, tablet), OS, and browser.
- **FR-AN-06**: Referrer analytics shall list the top referrer domains driving traffic.
- **FR-AN-07**: Users shall be able to filter analytics by date range.
- **FR-AN-08**: Pro and above users shall be able to export link analytics as CSV.
- **FR-AN-09**: Business and above users shall be able to export analytics as PDF reports.
- **FR-AN-10**: The system shall aggregate analytics asynchronously to avoid adding latency to redirects.

### User Management & Authentication

- **FR-AUTH-01**: Users shall be able to register with email and password.
- **FR-AUTH-02**: The system shall support OAuth2 login via Google and GitHub.
- **FR-AUTH-03**: Email verification shall be required before account activation.
- **FR-AUTH-04**: The system shall support password reset via email.
- **FR-AUTH-05**: Sessions shall be managed with short-lived JWT access tokens and rotating refresh tokens stored in Valkey.
- **FR-AUTH-06**: The system shall support multi-factor authentication (TOTP/Authenticator app).
- **FR-AUTH-07**: Enterprise plans shall support SSO via SAML 2.0 or OIDC.
- **FR-AUTH-08**: Team workspaces shall support RBAC with roles: Owner, Admin, Editor, Viewer.
- **FR-AUTH-09**: Owners shall be able to invite members by email with role assignment.
- **FR-AUTH-10**: Audit logs shall record all sensitive actions (login, link deletion, member removal) for Business and Enterprise plans.

### Subscription & Billing

- **FR-BILL-01**: The system shall integrate with Stripe for subscription management.
- **FR-BILL-02**: Users shall be able to subscribe, upgrade, downgrade, and cancel plans from the billing dashboard.
- **FR-BILL-03**: Downgrades shall take effect at the end of the current billing period; upgrades shall be prorated immediately.
- **FR-BILL-04**: The system shall enforce plan limits (link count, click tracking, API quota) and notify users as they approach limits (at 80% and 100%).
- **FR-BILL-05**: Upon exceeding limits, the system shall block new link creation (or charge overage per plan configuration) and notify the user to upgrade.
- **FR-BILL-06**: The system shall send billing notifications: payment successful, payment failed, subscription renewing, trial ending.
- **FR-BILL-07**: Invoices shall be available for download from the billing dashboard.
- **FR-BILL-08**: The system shall support a 14-day free trial for Starter and Pro plans without requiring a credit card.
- **FR-BILL-09**: Team plans shall support per-seat billing with the ability to add and remove seats.
- **FR-BILL-10**: Stripe webhooks shall be used to handle payment lifecycle events reliably.

### Admin Panel

- **FR-ADMIN-01**: A superadmin dashboard shall provide oversight of all users, links, and subscriptions.
- **FR-ADMIN-02**: Admins shall be able to search, view, suspend, and delete user accounts.
- **FR-ADMIN-03**: Admins shall be able to manually override subscription plans.
- **FR-ADMIN-04**: Admins shall be able to view platform-wide analytics: total links, total clicks, active users, revenue.
- **FR-ADMIN-05**: Admins shall be able to flag or disable malicious/spam links.
- **FR-ADMIN-06**: Admins shall be able to manage a blocklist of destination URL patterns (domains, paths).
- **FR-ADMIN-07**: Admins shall be able to configure system-wide settings: default link TTL, slug length, rate limits.

### Landing Page

- **FR-LAND-01**: The landing page shall clearly communicate the product value proposition above the fold.
- **FR-LAND-02**: The landing page shall include a live URL shortener demo widget (no login required, free tier).
- **FR-LAND-03**: The landing page shall display pricing plans with a clear comparison table.
- **FR-LAND-04**: The landing page shall include a features section, testimonials, and FAQ.
- **FR-LAND-05**: The landing page shall be fully responsive (mobile, tablet, desktop).
- **FR-LAND-06**: The landing page shall have optimized SEO meta tags, Open Graph tags, and a sitemap.
- **FR-LAND-07**: The landing page shall include a blog/changelog section for product updates.
- **FR-LAND-08**: Page load time shall be under 2 seconds on a 4G connection (Core Web Vitals: LCP < 2.5s, CLS < 0.1, FID < 100ms).

### Frontend Application

- **FR-FE-01**: The frontend shall be built with Vue.js 3 using the Composition API and TypeScript.
- **FR-FE-02**: The application shall use Vue Router for client-side routing with code splitting per route.
- **FR-FE-03**: State management shall use Pinia.
- **FR-FE-04**: The UI component library shall be based on a design system (e.g., shadcn-vue or PrimeVue).
- **FR-FE-05**: The dashboard shall provide a paginated, searchable, and filterable table of all user links.
- **FR-FE-06**: A link creation modal shall allow users to configure slug, expiry, password, UTM, and QR options in a single form.
- **FR-FE-07**: Charts and graphs (analytics) shall use Chart.js or ECharts.
- **FR-FE-08**: The application shall support light and dark mode.
- **FR-FE-09**: All forms shall provide inline validation with meaningful error messages.
- **FR-FE-10**: The application shall be a Progressive Web App (PWA) with offline capability for the dashboard.

### API

- **FR-API-01**: The system shall expose a RESTful API following OpenAPI 3.1 specification.
- **FR-API-02**: API authentication shall use Bearer tokens (API keys) generated per user.
- **FR-API-03**: The API shall support full CRUD operations for links, tags, and domains.
- **FR-API-04**: API rate limits shall be enforced per plan tier and tracked via Valkey.
- **FR-API-05**: The API shall support webhook registration for events: `link.created`, `link.clicked`, `link.expired`, `link.deleted`.
- **FR-API-06**: API responses shall follow a consistent envelope format with `data`, `error`, and `meta` fields.
- **FR-API-07**: The system shall provide interactive API documentation (Swagger UI / Redoc).
- **FR-API-08**: API versioning shall follow URI versioning (`/api/v1/`).

---

## Non-Functional Requirements

### Performance
- **NFR-PERF-01**: Redirect latency (p99) shall be under 50ms for cached links.
- **NFR-PERF-02**: Short URL creation API shall respond within 200ms under normal load.
- **NFR-PERF-03**: The system shall handle at least 10,000 redirects per second horizontally scaled.
- **NFR-PERF-04**: Valkey shall be the primary lookup for hot links; PostgreSQL is the source of truth.
- **NFR-PERF-05**: Analytics writes shall be asynchronous (event queue) to ensure redirect latency is unaffected.

### Availability & Reliability
- **NFR-REL-01**: The platform shall target 99.9% uptime (SLA) for standard plans; 99.99% for Enterprise.
- **NFR-REL-02**: The system shall be stateless at the application tier to enable horizontal scaling.
- **NFR-REL-03**: Valkey and PostgreSQL shall be deployed with replication and automatic failover.
- **NFR-REL-04**: All background jobs shall be idempotent and support retry with exponential backoff.

### Scalability
- **NFR-SCALE-01**: The application shall be horizontally scalable via container orchestration (Kubernetes-ready).
- **NFR-SCALE-02**: Database reads shall be offloaded to read replicas for analytics queries.
- **NFR-SCALE-03**: The system shall support multi-region deployments with geo-routing.

### Security
- **NFR-SEC-01**: All data in transit shall be encrypted with TLS 1.2+.
- **NFR-SEC-02**: All passwords shall be hashed with bcrypt (cost factor ≥ 12).
- **NFR-SEC-03**: IP addresses in analytics shall be hashed (SHA-256 + salt) before storage to comply with GDPR.
- **NFR-SEC-04**: The system shall implement OWASP Top 10 mitigations.
- **NFR-SEC-05**: Rate limiting shall be applied to all public endpoints (redirect, API, auth).
- **NFR-SEC-06**: The system shall scan destination URLs against a known malware/phishing blocklist (Google Safe Browsing API).
- **NFR-SEC-07**: CSRF protection shall be enforced on all state-changing API endpoints.
- **NFR-SEC-08**: Content Security Policy (CSP) headers shall be configured on all web responses.

### Compliance
- **NFR-COMP-01**: The system shall comply with GDPR — users can export and delete all their data.
- **NFR-COMP-02**: A cookie consent banner shall be present on the landing page.
- **NFR-COMP-03**: The platform shall maintain a privacy policy and terms of service.
- **NFR-COMP-04**: PCI DSS compliance is delegated to Stripe; no card data is stored by the platform.

### Maintainability
- **NFR-MAINT-01**: Backend code coverage shall be at least 80%.
- **NFR-MAINT-02**: All services shall expose `/health` and `/readiness` endpoints.
- **NFR-MAINT-03**: All deployments shall be zero-downtime (rolling updates).
- **NFR-MAINT-04**: Infrastructure shall be defined as code (Docker Compose for local, Helm charts for production).

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                          Client Layer                           │
│    Landing Page (Vue.js SSG)   │   Dashboard SPA (Vue.js 3)    │
└─────────────────┬───────────────────────────┬───────────────────┘
                  │                           │
                  ▼                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Nginx / API Gateway                       │
│          TLS Termination · Rate Limiting · CORS Headers          │
└────────┬──────────────────────────────────────┬─────────────────┘
         │                                      │
         ▼                                      ▼
┌─────────────────┐                  ┌──────────────────────┐
│  Redirect API   │                  │    Application API   │
│  (Go — ultra    │                  │  (Go — REST /api/v1) │
│   low latency)  │                  │                      │
└────────┬────────┘                  └──────────┬───────────┘
         │                                      │
         ▼                                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                           Valkey                                │
│   Link Cache · Session Store · Rate Limit Counters · Job Queue  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         PostgreSQL                              │
│    Users · Links · Analytics · Subscriptions · Audit Logs      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Background Workers                         │
│   Analytics Processor · Link Health Checker · Email Sender     │
│   QR Code Generator · Subscription Sync · Webhook Dispatcher   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      External Services                          │
│   Lemon Squezzy · Google Safe Browsing · S3 · SMTP · OAuth Providers  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Data Models

### Users
| Field            | Type        | Notes                          |
|------------------|-------------|--------------------------------|
| id               | UUID        | Primary key                    |
| email            | VARCHAR     | Unique, indexed                |
| password_hash    | VARCHAR     | bcrypt                         |
| name             | VARCHAR     |                                |
| avatar_url       | VARCHAR     |                                |
| email_verified   | BOOLEAN     |                                |
| mfa_enabled      | BOOLEAN     |                                |
| plan_id          | FK          | References subscriptions       |
| created_at       | TIMESTAMPTZ |                                |
| updated_at       | TIMESTAMPTZ |                                |

### Links
| Field            | Type        | Notes                          |
|------------------|-------------|--------------------------------|
| id               | UUID        | Primary key                    |
| user_id          | FK          | References users               |
| workspace_id     | FK          | References workspaces (nullable)|
| slug             | VARCHAR     | Unique, indexed                |
| destination_url  | TEXT        |                                |
| domain_id        | FK          | Custom domain (nullable)       |
| password_hash    | VARCHAR     | Nullable                       |
| expires_at       | TIMESTAMPTZ | Nullable                       |
| max_clicks       | INTEGER     | Nullable                       |
| click_count      | INTEGER     | Cached counter                 |
| redirect_type    | SMALLINT    | 301 or 302                     |
| is_active        | BOOLEAN     |                                |
| tags             | TEXT[]      |                                |
| utm_source       | VARCHAR     | Nullable                       |
| utm_medium       | VARCHAR     | Nullable                       |
| utm_campaign     | VARCHAR     | Nullable                       |
| created_at       | TIMESTAMPTZ |                                |
| updated_at       | TIMESTAMPTZ |                                |

### Click Events
| Field            | Type        | Notes                          |
|------------------|-------------|--------------------------------|
| id               | UUID        | Primary key                    |
| link_id          | FK          | References links               |
| clicked_at       | TIMESTAMPTZ | Indexed                        |
| ip_hash          | VARCHAR     | SHA-256 hashed                 |
| country          | CHAR(2)     | ISO 3166-1 alpha-2             |
| city             | VARCHAR     |                                |
| device_type      | VARCHAR     | desktop / mobile / tablet      |
| os               | VARCHAR     |                                |
| browser          | VARCHAR     |                                |
| referrer         | TEXT        |                                |
| source           | VARCHAR     | web / qr / api                 |

### Subscriptions
| Field            | Type        | Notes                          |
|------------------|-------------|--------------------------------|
| id               | UUID        | Primary key                    |
| user_id          | FK          | References users               |
| plan             | VARCHAR     | free / starter / pro / business / enterprise |
| status           | VARCHAR     | active / trialing / past_due / canceled |
| stripe_customer_id | VARCHAR   |                                |
| stripe_sub_id    | VARCHAR     |                                |
| current_period_start | TIMESTAMPTZ |                           |
| current_period_end   | TIMESTAMPTZ |                           |
| cancel_at_period_end | BOOLEAN   |                                |

---

## API Design

### Base URL
```
https://api.linkbee.io/v1
```

### Authentication
```
Authorization: Bearer <api_key>
```

### Key Endpoints

```
POST   /auth/register
POST   /auth/login
POST   /auth/refresh
POST   /auth/logout
POST   /auth/password-reset

GET    /links
POST   /links
GET    /links/:id
PUT    /links/:id
DELETE /links/:id
GET    /links/:id/analytics
GET    /links/:id/qr

GET    /workspaces
POST   /workspaces
GET    /workspaces/:id/members
POST   /workspaces/:id/members/invite

GET    /domains
POST   /domains
DELETE /domains/:id

GET    /billing/plans
GET    /billing/subscription
POST   /billing/checkout
POST   /billing/portal

GET    /user/profile
PUT    /user/profile
GET    /user/api-keys
POST   /user/api-keys
DELETE /user/api-keys/:id

GET    /admin/users
GET    /admin/stats
POST   /admin/links/:id/disable
```

---

## Security Requirements

- All API keys stored as hashed values; prefix shown to user for identification (e.g., `sl_live_xxxx...`)
- Secrets managed via environment variables; never committed to source control
- `.env.example` provided with placeholder values
- All user-supplied URL destinations validated and scanned against Safe Browsing API before storage
- Input sanitization on all endpoints to prevent XSS and injection attacks
- Database queries use parameterized statements exclusively (no raw string interpolation)
- Stripe webhook signatures verified on every incoming webhook request

---

## Infrastructure & DevOps

### Local Development
```
docker-compose up
```
Starts: PostgreSQL, Valkey, Go API, Vue.js dev server, Nginx

### Environment Variables
| Variable                  | Description                        |
|---------------------------|------------------------------------|
| DATABASE_URL              | PostgreSQL connection string       |
| VALKEY_URL                | Valkey connection string           |
| JWT_SECRET                | JWT signing secret                 |
| STRIPE_SECRET_KEY         | Stripe secret key                  |
| STRIPE_WEBHOOK_SECRET     | Stripe webhook signing secret      |
| GOOGLE_CLIENT_ID          | OAuth2 Google client ID            |
| GITHUB_CLIENT_ID          | OAuth2 GitHub client ID            |
| SAFE_BROWSING_API_KEY     | Google Safe Browsing API key       |
| S3_BUCKET                 | Object storage bucket name         |
| S3_ENDPOINT               | S3-compatible endpoint URL         |
| BASE_URL                  | Public base URL for short links    |
| SMTP_HOST                 | Email server host                  |

---

## Development Roadmap

### Phase 1 — MVP
- [ ] Go project scaffold (modules, routing, middleware)
- [ ] PostgreSQL schema and migrations
- [ ] Valkey integration and link cache
- [ ] URL shortening (create, redirect, delete)
- [ ] QR code generation (PNG)
- [ ] User auth (email/password, JWT)
- [ ] Basic click tracking
- [ ] Vue.js dashboard (link CRUD, basic analytics)
- [ ] Landing page with demo widget
- [ ] Docker Compose local dev setup

### Phase 2 — Monetization
- [ ] Lemon squezzy subscription integration
- [ ] Plan enforcement (limits and quotas)
- [ ] Billing dashboard (upgrade, cancel, invoices)
- [ ] Free trial (14 days, no credit card)
- [ ] Email notifications (billing, link expiry)

### Phase 3 — Growth Features
- [ ] Custom domains (BYOD)
- [ ] OAuth2 login (Google, GitHub)
- [ ] Advanced analytics (geo, device, referrer, time-series)
- [ ] CSV export of analytics
- [ ] Bulk URL import
- [ ] Link health monitoring
- [ ] UTM builder

### Phase 4 — Teams & Enterprise
- [ ] Team workspaces + RBAC
- [ ] Per-seat billing
- [ ] SSO (SAML 2.0 / OIDC)
- [ ] Audit logs
- [ ] QR code customization (logo, colors)
- [ ] Webhooks
- [ ] REST API with API key auth
- [ ] OpenAPI documentation

### Phase 5 — Scale & Reliability
- [ ] Kubernetes Helm charts
- [ ] Multi-region support
- [ ] Read replica for analytics queries
- [ ] Prometheus + Grafana monitoring
- [ ] Load testing and performance tuning
- [ ] GDPR data export / deletion tooling
