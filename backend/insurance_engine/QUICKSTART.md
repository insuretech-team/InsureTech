# Insurance Engine - Quickstart Guide

This guide will help you set up and run the **Insurance Engine** (.NET 10) on your local machine.

---

## 🇧🇩 বাংলায় নির্দেশাবলী (Setup in Bengali)

### ১. প্রি-রিকুইজিট (Prerequisites)
- **.NET 10 SDK**: [এখান থেকে ডাউনলোড করুন](https://dotnet.microsoft.com/download/dotnet/10.0)
- **Go 1.25**: মাইগ্রেশন চালানোর জন্য প্রয়োজন। [ডাউনলোড করুন](https://go.dev/dl/)
- **PostgreSQL**: [এখান থেকে ডাউনলোড করুন](https://www.postgresql.org/download/windows/) অথবা ডকার ব্যবহার করুন।
- **Docker Desktop**: (ঐচ্ছিক, কন্টেইনারে চালানোর জন্য)

### ২. ডাটাবেজ সেটআপ (Database & Migrations)
ইঞ্জিন রান করার আগে ডাটাবেজ টেবিলগুলো তৈরি করা প্রয়োজন:
১. রুট ফোল্ডারে টার্মিনাল ওপেন করুন।
২. এনভায়রনমেন্ট ফাইল কপি করুন: `copy .env.example .env` (এবং আপনার ডাটাবেজ ডিটেইলস দিন)।
৩. প্রিকুইজিট চেক করুন: `./scripts/bootstrap.ps1`
৪. মাইগ্রেশন রান করুন:
   ```powershell
   ./run_migration.ps1 -Target primary
   ```

### ৩. লোকাল রান (Local Run - Dotnet)
১. টার্মিনালে প্রজেক্টের রুটে যান: `cd backend/insurance_engine/src/InsuranceEngine.ApiHost`
২. অ্যাপটি রান করুন: `dotnet run`
৩. অ্যাপটি ডিফল্টভাবে `https://localhost:5001` এ চলবে।

### ৩. ডকার দিয়ে রান (Docker Run)
১. ইন্স্যুরেন্স ইঞ্জিনের ফোল্ডারে যান: `cd backend/insurance_engine`
২. নিচের কমান্ডটি দিন (এটি ডাটাবেজ এবং ইঞ্জিন একসাথে চালু করবে):
   ```bash
   docker compose up -d
   ```
৩. ইঞ্জিনটি `http://localhost:5001` এ অ্যাভেলেবল হবে।

### ৪. সার্ভিস টেস্টিং (Testing)
১. **Postman** ওপেন করুন।
২. সার্ভিস হিসেবে `https://localhost:5001` দিন।
৩. পেমেন্ট বা প্রোডাক্ট ক্রিয়েট টেস্ট করার জন্য আমাদের [Postman Testing Guide](file:///C:/Users/Rizve%20Rahman%20Reza/.gemini/antigravity/brain/2038c740-cc80-4024-9cd3-7885a6d13fc2/grpc_postman_testing_guide.md) ফলো করুন।

---

## 🇺🇸 English Instructions (Setup in English)

### 1. Prerequisites
- **.NET 10 SDK**: [Download here](https://dotnet.microsoft.com/download/dotnet/10.0)
- **Go 1.25**: Required for database migrations. [Download here](https://go.dev/dl/)
- **PostgreSQL**: [Download here](https://www.postgresql.org/download/) or use Docker.
- **Docker Desktop**: (Optional, for containerized run)

### 2. Database Setup & Migrations
Before running the engine, ensure the database schema is up to date:
1. Open terminal in the **repository root**.
2. Setup environment: `cp .env.example .env` (Update DB credentials if needed).
3. Bootstrap dependencies: `./scripts/bootstrap.ps1`
4. Run migrations:
   ```powershell
   ./run_migration.ps1 -Target primary
   ```

### 3. Running Locally (Dotnet CLI)
1. Navigate to the API Host directory: `cd backend/insurance_engine/src/InsuranceEngine.ApiHost`
2. Run the application: `dotnet run`
3. The service will be available at `https://localhost:5001`.

### 3. Running with Docker
1. Navigate to the insurance engine folder: `cd backend/insurance_engine`
2. Spin up the engine and its dependencies (Postgres/Redis):
   ```bash
   docker compose up -d
   ```
3. The engine will be mapped to `http://localhost:5001`.

### 4. Verification
1. Ensure the build is successful: `dotnet build backend/insurance_engine/InsuranceEngine.sln`
2. Use gRPC reflection in Postman to discover services automatically.

---

> [!TIP]
> **Configuration**: All database connection strings are managed via `appsettings.json` or Environment Variables (`ConnectionStrings__InsuranceDb`).
