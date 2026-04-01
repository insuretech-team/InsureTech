# Start Here — InsureTech Platform

## What This Platform Is

InsureTech is a production-grade insurance commerce platform with a **dual-backend architecture**:

- **InScore** (`backend/inscore/`) — Go microservices: auth, authz, payments, storage, orders data layer, fraud, partner, KYC, notifications, audit, media, docgen
- **PoliSync** (`backend/polisync/`) — C# .NET 8 insurance engine: products, quotations, orders (business logic), underwriting, policy issuance, claims, endorsements, renewals, refunds, commissions

## Key Architectural Rule

`order` (PoliSync, port 50140) and `sync_order` (InScore Go, port 50142) are the **two halves of the same domain**:
- **PoliSync `order`** = business logic, domain rules, CQRS commands
- **InScore `sync_order`** = persistence layer, Kafka events, database

## Primary References (read in order)

1. `documentation/About/ARCHITECTURE_OVERVIEW.md` — full system map
2. `documentation/About/POLISYNC_REFERENCE.md` — PoliSync engine reference
3. `documentation/core_plans/POLISYNC_REFERENCE.md` — PoliSync implementation guide
4. `documentation/core_plans/ACTIVE_WORKSTREAMS.md` — what's built, what needs attention
5. `rules/00-index.md` — normative API rules

## Service Port Map

| Service | Backend | gRPC | HTTP |
|---------|---------|------|------|
| authn | Go/InScore | 50060 | 50061 |
| authz | Go/InScore | 50070 | 50071 |
| audit | Go/InScore | 50080 | 50081 |
| kyc | Go/InScore | 50090 | 50091 |
| partner | Go/InScore | 50100 | 50101 |
| insurance | Go/InScore | 50115 | — |
| product (PoliSync) | C#/PoliSync | 50120 | 50121 |
| quote (PoliSync) | C#/PoliSync | 50130 | 50131 |
| order (PoliSync) | C#/PoliSync | 50140 | 50141 |
| sync_order (Go data layer) | Go/InScore | 50142 | 50143 |
| commission (PoliSync) | C#/PoliSync | 50150 | 50151 |
| policy (PoliSync) | C#/PoliSync | 50160 | 50161 |
| underwriting (PoliSync) | C#/PoliSync | 50170 | 50171 |
| payment | Go/InScore | 50190 | 50191 |
| claim (PoliSync) | C#/PoliSync | 50210 | 50211 |
| fraud | Go/InScore | 50220 | 50221 |
| notification | Go/InScore | 50230 | 50231 |
| docgen | Go/InScore | 50280 | 50281 |
| storage | Go/InScore | 50290 | 50291 |
