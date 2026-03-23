# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# Copy solution and project files
COPY InsuranceEngine.sln ./
COPY src/InsuranceEngine.SharedKernel/InsuranceEngine.SharedKernel.csproj src/InsuranceEngine.SharedKernel/
COPY src/InsuranceEngine.ApiHost/InsuranceEngine.ApiHost.csproj src/InsuranceEngine.ApiHost/
COPY src/InsuranceEngine.Products/InsuranceEngine.Products.csproj src/InsuranceEngine.Products/
COPY src/InsuranceEngine.Beneficiary/InsuranceEngine.Beneficiary.csproj src/InsuranceEngine.Beneficiary/
COPY src/InsuranceEngine.Policy/InsuranceEngine.Policy.csproj src/InsuranceEngine.Policy/
COPY src/InsuranceEngine.Underwriting/InsuranceEngine.Underwriting.csproj src/InsuranceEngine.Underwriting/
COPY src/InsuranceEngine.Claims/InsuranceEngine.Claims.csproj src/InsuranceEngine.Claims/
COPY src/InsuranceEngine.Fraud/InsuranceEngine.Fraud.csproj src/InsuranceEngine.Fraud/
COPY src/InsuranceEngine.Partners/InsuranceEngine.Partners.csproj src/InsuranceEngine.Partners/
COPY src/InsuranceEngine.Commission/InsuranceEngine.Commission.csproj src/InsuranceEngine.Commission/

# Restore dependencies
RUN dotnet restore

# Copy source code
COPY src/ src/

# Build
RUN dotnet build -c Release --no-restore

# Publish
RUN dotnet publish src/InsuranceEngine.ApiHost/InsuranceEngine.ApiHost.csproj \
    -c Release \
    --no-build \
    -o /app/publish

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app

# Install curl for health checks
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Copy published app
COPY --from=build /app/publish .

# Create logs directory
RUN mkdir -p /app/logs

# Expose ports
EXPOSE 50051 5000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:5000/health || exit 1

# Run
ENTRYPOINT ["dotnet", "InsuranceEngine.ApiHost.dll"]
