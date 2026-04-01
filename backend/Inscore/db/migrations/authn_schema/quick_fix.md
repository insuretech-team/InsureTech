# Quick Fix Notes — Live ALTER (Proto-First Migration Warning)

> ⚠️ **IMPORTANT**: This project uses proto-first schema generation.
> Normal migration files that `ALTER`, `CREATE`, or `ADD COLUMN` will **conflict with or be ignored by**
> the auto-generated schema from proto definitions. Always apply structural changes as a **live ALTER**
> via `dbsql` and document them here.

---

## 2026-03-16 — `authn_schema.otps.user_id` — DROP NOT NULL

**Applied live via dbsql:**
```sql
ALTER TABLE authn_schema.otps ALTER COLUMN user_id DROP NOT NULL;
```

**Why:**
- OTP is sent *before* the user is known (pre-login flow: SendOTP → VerifyOTP → Login).
- The proto entity `OTP.user_id` was set to `""` (empty string) in `otp_service.go` with a comment
  "Will be set on verification if needed".
- PostgreSQL's `uuid` column type rejects `""` as `invalid input syntax for type uuid (SQLSTATE 22P02)`.
- The column was `NOT NULL` with no default, causing every `SendOTP` call to fail at DB insert
  after the SMS was already sent — SMS cost wasted, user gets 500.

**Fix:**
1. `ALTER TABLE authn_schema.otps ALTER COLUMN user_id DROP NOT NULL;` — applied live ✅
2. `otp_service.go` — `UserId` field omitted (left as proto zero value `""`) — GORM now passes NULL ✅
   Actually the proto string zero-value `""` still fails — set `UserId` to a sentinel or omit via
   GORM tag. See code fix below.

**Proto-first note:**
- When the proto schema is next regenerated and migrations re-run, ensure the `otps` table DDL
  reflects `user_id UUID NULL` not `user_id UUID NOT NULL`.
- The proto field `OTP.user_id` should be marked `optional string user_id = 2;` in the `.proto` file
  so the generated Go struct uses `*string` (pointer) allowing nil.
- Until proto is updated, the repository `Create()` call must explicitly set `user_id = NULL` when
  the value is empty — done by setting `UserId` to a special omit tag or by patching the GORM insert.

**Affected files:**
- `backend/inscore/microservices/authn/internal/service/otp_service.go` — `SendOTP()` + `SendEmailOTP()`
- `backend/inscore/microservices/authn/internal/repository/otp_repository.go` — `Create()` uses GORM

---

## 2026-03-16 — `authn_schema.otps.sender_id` — VARCHAR(11) → VARCHAR(64)

**Applied live via dbsql:**
```sql
ALTER TABLE authn_schema.otps ALTER COLUMN sender_id TYPE VARCHAR(64);
```

**Why:**
- `sender_id` was `VARCHAR(11)` but `SSLWIRELESS_SENDER_ID=LIFEPLUSBDBRAND` is 15 chars.
- Every OTP send failed at DB insert with `value too long for type character varying(11)`.
- Standard alphanumeric sender IDs can be up to 11 chars per GSM spec, but provider-specific
  IDs (e.g. `LIFEPLUSBDBRAND`) exceed this. 64 chars gives plenty of headroom.

**Proto-first note:**
- `otp.proto` `sender_id` sql_type updated from `VARCHAR(11)` → `VARCHAR(64)`.

**Affected files:**
- `proto/insuretech/authn/entity/v1/otp.proto`

---

## 2026-03-16 — Email SMTP blocked on DigitalOcean

**Issue:**
- DigitalOcean blocks outbound SMTP on ports 25, 587, 465 at the network level (spam prevention).
- Gmail (`smtp.gmail.com:587`) is unreachable from the droplet.
- `DEADLINE_EXCEEDED` on all email OTP / password reset requests.

**Fix applied:**
- Switched SMTP relay to **smtp2go** (`mail.smtp2go.com:2525`) — port 2525 is unblocked.
- Updated `.env` and `.env.prod` `EMAIL_SMTP_HOST`, `EMAIL_SMTP_PORT`.

**Action required:**
1. Sign up at https://app.smtp2go.com (free tier: 1000 emails/month)
2. Add sender domain `labaidinsuretech.com` and verify DNS (SPF/DKIM)
3. Create SMTP user — get username + password
4. Update `.env` and `.env.prod`:
   ```
   EMAIL_USERNAME=<smtp2go-username>
   EMAIL_PASSWORD=<smtp2go-password>
   EMAIL_INFO_USERNAME=<smtp2go-username>
   EMAIL_INFO_PASSWORD=<smtp2go-password>
   ```
5. Redeploy authn: `bash scripts/quickerdeploy.sh --services=authn`

**Alternative:** Contact DigitalOcean support to enable outbound SMTP for droplet #538539366.

---

## Template for future live fixes

```
## YYYY-MM-DD — <schema>.<table>.<column> — <operation>

**Applied live via dbsql:**
\`\`\`sql
-- SQL here
\`\`\`

**Why:** ...
**Proto-first note:** ...
**Affected files:** ...
```
