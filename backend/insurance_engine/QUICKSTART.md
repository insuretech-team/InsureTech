# Insurance Engine Quick Start

Get the modular Insurance Engine running in 5 minutes.

## Prerequisites

- .NET 8 SDK ([download](https://dotnet.microsoft.com/download/dotnet/8.0))
- Docker Desktop (for PostgreSQL and Kafka)
- EF Core Tools: `dotnet tool install --global dotnet-ef`

## Step 1: Start Shared Infrastructure

The platform uses a shared infrastructure for core services (PostgreSQL, Kafka, Redis, etc.). Start these from the project root:

```bash
cd ../.. # Go to repo root
docker compose up -d postgres kafka redis
```

## Step 2: Configure Database (Migrations)

The system uses the Go-side `dbops` tool to manage schema migrations (PoliSync Pattern).

1.  **Initialize Schema**:
    ```bash
    cd ../.. # Go to repo root
    # Run the official migration service
    docker compose run --rm db-migrate
    ```

2.  **Automatic Seeding**: Once the schema is created, the `ApiHost` will automatically seed initial product data on its first run.

## Step 3: Run the Engine

### Option A: Local Development (IDE/CLI)
```bash
cd src/InsuranceEngine.ApiHost
dotnet run
```

### Option B: Local Automation (PowerShell)
```powershell
./build.ps1 -Configuration Debug
```

### Option C: Docker (Production-like)
```bash
# Build and start the container
docker-compose up --build -d
```

## Step 4: Verify & Test

Once the engine is running, you can verify it using several methods:

### 1. Swagger UI (Recommended)
Navigate to `http://localhost:5135/swagger` in your browser. This provides a full interactive UI to test all modules:
-   **Beneficiary**: Create individuals/businesses, KYC updates.
-   **Policy**: Quote to Issuance lifecycle.
-   **Claims**: Submission and automated ZHTC approval.
-   **Partners/Commission**: Partner management and commission tracking.

### 2. HTTP File
Use the [InsuranceEngine.ApiHost.http](src/InsuranceEngine.ApiHost/InsuranceEngine.ApiHost.http) file in your IDE to run pre-configured test requests.

### 3. Health Checks
```bash
curl http://localhost:5135/health
# Should return: "Healthy"
```

## Module Architecture & Ports

| Module | Purpose | Status |
|--------|---------|--------|
| **ApiHost** | Entry point & gRPC Host | ✅ Active |
| **Beneficiary** | CRM & KYC Management | ✅ Aligned with Docs |
| **Policy** | Lifecycle, Endorsements & Renewals | ✅ Logic Verified |
| **Claims** | Adjudication & ZHTC Auto-Approval | ✅ Matrix Active |
| **Underwriting** | Decision Matrix & Risk Assessment | ✅ Integrated |
| **Partners** | B2B & Agent Management | ✅ M2 Complete |
| **Commission** | 15% Acquisition / 5% Renewal | ✅ Ported |

## Common Commands

```bash
# Build Solution
dotnet build

# Run All Tests
dotnet test

# Hot Reload
dotnet watch run --project src/InsuranceEngine.ApiHost
```

Happy coding! 🚀
