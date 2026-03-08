# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# Copy solution and project files
COPY PoliSync.sln ./
COPY src/PoliSync.SharedKernel/PoliSync.SharedKernel.csproj src/PoliSync.SharedKernel/
COPY src/PoliSync.Infrastructure/PoliSync.Infrastructure.csproj src/PoliSync.Infrastructure/
COPY src/PoliSync.ApiHost/PoliSync.ApiHost.csproj src/PoliSync.ApiHost/
COPY src/PoliSync.Products/PoliSync.Products.csproj src/PoliSync.Products/
COPY src/PoliSync.Quotes/PoliSync.Quotes.csproj src/PoliSync.Quotes/
COPY src/PoliSync.Orders/PoliSync.Orders.csproj src/PoliSync.Orders/
COPY src/PoliSync.Policy/PoliSync.Policy.csproj src/PoliSync.Policy/
COPY src/PoliSync.Endorsement/PoliSync.Endorsement.csproj src/PoliSync.Endorsement/
COPY src/PoliSync.Renewal/PoliSync.Renewal.csproj src/PoliSync.Renewal/
COPY src/PoliSync.Underwriting/PoliSync.Underwriting.csproj src/PoliSync.Underwriting/
COPY src/PoliSync.Claims/PoliSync.Claims.csproj src/PoliSync.Claims/
COPY src/PoliSync.Commission/PoliSync.Commission.csproj src/PoliSync.Commission/
COPY src/PoliSync.Refund/PoliSync.Refund.csproj src/PoliSync.Refund/

# Restore dependencies
RUN dotnet restore

# Copy source code
COPY src/ src/

# Build
RUN dotnet build -c Release --no-restore

# Publish
RUN dotnet publish src/PoliSync.ApiHost/PoliSync.ApiHost.csproj \
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

# Expose ports (all gRPC + HTTP ports)
EXPOSE 50120 50121 50130 50131 50140 50141 50150 50151 50160 50161 50170 50171 50210 50211

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:50121/health || exit 1

# Run
ENTRYPOINT ["dotnet", "PoliSync.ApiHost.dll"]
