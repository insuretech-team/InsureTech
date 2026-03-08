# Retrieved Files from insureadmin@146.190.97.242

## Summary

Successfully retrieved the following files from the InsureTech production server:

1. **docker-compose.yml** - Complete Docker Compose configuration (438 lines)
2. **.env.prod** - Production environment variables (370 lines)
3. **.env** - Environment configuration (370 lines - identical to .env.prod)
4. **Directory listing** - /home/insureadmin/insuretech/backend/inscore/configs/
5. **services.yaml** - Configuration file status

---

## File Contents

### 1. docker-compose.yml
Located at: `/home/insureadmin/insuretech/docker-compose.yml`

Full content saved locally in `docker-compose.yml`

**Key Services Configured:**
- **redis**: Redis 7 Alpine - Message caching (port 127.0.0.1:6379)
- **kafka**: Apache Kafka - Event streaming (port 127.0.0.1:9092)
- **db-migrate**: Database migration service (runs once at startup)
- **db-sidecar**: Continuous DB sync and auto-backup companion
- **authn**: Authentication microservice (gRPC 50060)
- **authz**: Authorization microservice (gRPC 50070)
- **gateway**: API Gateway (port 127.0.0.1:8080)
- **b2b**: B2B microservice (gRPC 50112)
- **storage**: Storage microservice (gRPC 50290)
- **b2b_portal**: Next.js B2B Portal (port 3000:8081)

**Commented/Disabled Services:**
- tenant, audit, kyc, partner, beneficiary, workflow, payment, fraud, notification, support, webrtc, media, docgen, iot, analytics, ai

**Volumes:**
- redis_data, kafka_data, go_mod_cache, go_build_cache

---

### 2. .env.prod
Located at: `/home/insureadmin/insuretech/.env.prod`

Full content saved locally in `.env.prod`

**Key Configuration Sections:**

#### Global & Network
- **ENVIRONMENT**: production
- **DOMAIN**: labaidinsuretech.com
- **HOST**: labaidinsuretech.com
- **PORT**: 8081
- **REMOTE_SERVER_PUBLIC_IP**: 146.190.97.242
- **HTTPS_ENABLED**: true
- **PUBLIC_API_BASE_URL**: https://api.labaidinsuretech.com

#### Admin Credentials
- **ADMIN_EMAIL**: faruk.hannan@gmail.com
- **ADMIN_PASSWORD**: 1234567
- **ADMIN_MOBILE**: +8801347-210751

#### Databases
- **Primary (DO)**: db-postgresql-sgp1-03392-do-user-30821835-0.i.db.ondigitalocean.com (port 25061)
- **Secondary (NEON)**: ep-icy-paper-ah86wl8n-pooler.c-3.us-east-1.aws.neon.tech
- **Signal DB (NEON)**: ins-signal database

Database Credentials:
- DO: doadmin / AVNS_Eod4Ng8PzESnHShSQ4I
- NEON: SuperAdmin / npg_eYcAl3msqoN5

#### gRPC Service Addresses
- **AUTHN_GRPC_ADDR**: authn:50060
- **AUTHZ_GRPC_ADDR**: authz:50070
- **B2B_GRPC_ADDR**: b2b:50112
- **STORAGE_GRPC_ADDR**: storage:50290

#### CORS Configuration
- **CORS_ALLOWED_ORIGINS**: https://b2b.labaidinsuretech.com,https://labaidinsuretech.com
- **CORS_ALLOWED_METHODS**: GET,POST,PUT,DELETE,OPTIONS
- **CORS_ALLOW_CREDENTIALS**: true

#### Storage & CDN (DigitalOcean Spaces)
- **SPACES_BUCKET**: lcst
- **SPACES_REGION**: sgp1
- **SPACES_ENDPOINT**: https://sgp1.digitaloceanspaces.com
- **SPACES_CDN_ENDPOINT**: https://lcst.sgp1.cdn.digitaloceanspaces.com
- **Access Key**: DO00DHT8X6EU2DGCWPWV
- **Secret Key**: QNWy2bRQHfwv8apH6b2yUD5wWGFUqIKLkTxfZOb7JNY

#### Caching & Brokers
- **KAFKA_BROKERS**: kafka:9092
- **REDIS_URL**: redis://redis:6379

#### Email (SMTP)
- **EMAIL_SMTP_HOST**: mx.zoho.com.
- **EMAIL_SMTP_PORT**: 1025
- **EMAIL_FROM**: dev@labaidinsuretech.com

#### SMS Notifications (SSL Wireless)
- **SSLWIRELESS_SID**: labaidinsuretech
- **SSLWIRELESS_API_KEY**: asdfgghytredn
- **SSLWIRELESS_SENDER_ID**: labaidinsuretech

#### KYC & FLVE Integration
- **KYC_SERVICE_ENABLED**: false
- **FLVE_BACKEND**: hybrid
- **FLVE_HF_ENDPOINT**: https://farukhannan-flve.hf.space
- **FLVE_HF_TOKEN**: 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
- **FLVE_MODELS**: yolo-face.onnx, arcface.onnx, liveness.onnx

#### Frontend (B2B Portal)
- **B2B_PORTAL_API_BASE_URL**: https://api.labaidinsuretech.com
- **B2B_PORTAL_INTERNAL_API_BASE_URL**: http://gateway:8080

#### Security
- **JWT_SECRET**: your-super-secret-jwt-key-change-in-production
- **PASSWORD_MIN_LENGTH**: 8
- **PASSWORD_EXPIRY_DAYS**: 90
- **MAX_LOGIN_ATTEMPTS**: 5
- **LOCKOUT_DURATION_MINUTES**: 30

#### Rate Limiting
- **RATE_LIMIT_PER_MINUTE**: 100
- **RATE_LIMIT_LOGIN_PER_MINUTE**: 5
- **RATE_LIMIT_PASSWORD_PER_MINUTE**: 3

---

### 3. .env
Located at: `/home/insureadmin/insuretech/.env`

**Status**: Identical to `.env.prod` (370 lines)

Full content saved locally in `.env`

---

### 4. Directory Listing
**Command**: `ls -la /home/insureadmin/insuretech/backend/inscore/configs/ 2>/dev/null || echo 'No configs dir'`

**Result**: `No configs dir`

The `/home/insureadmin/insuretech/backend/inscore/configs/` directory does not exist on the server.

---

### 5. services.yaml
**Command**: `cat /home/insureadmin/insuretech/backend/inscore/configs/services.yaml 2>/dev/null || echo 'not found'`

**Result**: `not found`

The `services.yaml` file does not exist at the specified location.

---

## Important Security Notes

⚠️ **The following sensitive credentials are exposed in the .env files:**

1. Database passwords for DigitalOcean and Neon databases
2. DigitalOcean Spaces access keys and secrets
3. FLVE Hugging Face API token
4. SSL Wireless API credentials
5. Email and SMS service credentials
6. JWT secret (placeholder, but should be changed)
7. Admin credentials

**Recommendation**: These should be stored in a secure secrets management system (e.g., HashiCorp Vault, AWS Secrets Manager, etc.) and not committed to version control or exposed in plain text.

---

## Architecture Overview

The InsureTech system is a microservices architecture deployed on Docker Compose with:

- **Infrastructure**: Redis (caching), Kafka (event streaming)
- **Core Services**: Authentication (authn), Authorization (authz), Gateway
- **Data Layer**: Multiple PostgreSQL databases (DigitalOcean + Neon)
- **Storage**: DigitalOcean Spaces for CDN
- **Frontend**: Next.js-based B2B Portal
- **Optional Services**: Workflow, Payment, Fraud detection, Notifications, Analytics, AI, etc. (currently disabled)

All services communicate via gRPC internally and are bound to localhost (127.0.0.1) with nginx proxying public traffic.

---

## Files Created Locally

1. `docker-compose.yml` - Full Docker Compose configuration
2. `.env.prod` - Production environment variables
3. `.env` - Environment variables (same as .env.prod)
4. `RETRIEVED_FILES_SUMMARY.md` - This summary document
