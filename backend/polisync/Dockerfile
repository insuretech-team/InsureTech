# Build stage — context is project root (E:\Projects\InsureTech)
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# Copy generated proto C# files (required by PoliSync.Proto.csproj)
COPY gen/csharp/ gen/csharp/

# Copy all PoliSync project files for restore
COPY backend/polisync/src/PoliSync.SharedKernel/PoliSync.SharedKernel.csproj       backend/polisync/src/PoliSync.SharedKernel/
COPY backend/polisync/src/PoliSync.Infrastructure/PoliSync.Infrastructure.csproj    backend/polisync/src/PoliSync.Infrastructure/
COPY backend/polisync/src/PoliSync.Proto/PoliSync.Proto.csproj                      backend/polisync/src/PoliSync.Proto/
COPY backend/polisync/src/PoliSync.ApiHost/PoliSync.ApiHost.csproj                  backend/polisync/src/PoliSync.ApiHost/
COPY backend/polisync/src/PoliSync.Products/PoliSync.Products.csproj                backend/polisync/src/PoliSync.Products/
COPY backend/polisync/src/PoliSync.Quotes/PoliSync.Quotes.csproj                    backend/polisync/src/PoliSync.Quotes/
COPY backend/polisync/src/PoliSync.Orders/PoliSync.Orders.csproj                    backend/polisync/src/PoliSync.Orders/
COPY backend/polisync/src/PoliSync.Policy/PoliSync.Policy.csproj                    backend/polisync/src/PoliSync.Policy/
COPY backend/polisync/src/PoliSync.Endorsement/PoliSync.Endorsement.csproj          backend/polisync/src/PoliSync.Endorsement/
COPY backend/polisync/src/PoliSync.Renewal/PoliSync.Renewal.csproj                  backend/polisync/src/PoliSync.Renewal/
COPY backend/polisync/src/PoliSync.Underwriting/PoliSync.Underwriting.csproj        backend/polisync/src/PoliSync.Underwriting/
COPY backend/polisync/src/PoliSync.Claims/PoliSync.Claims.csproj                    backend/polisync/src/PoliSync.Claims/
COPY backend/polisync/src/PoliSync.Commission/PoliSync.Commission.csproj            backend/polisync/src/PoliSync.Commission/
COPY backend/polisync/src/PoliSync.Refund/PoliSync.Refund.csproj                    backend/polisync/src/PoliSync.Refund/
COPY backend/polisync/src/PoliSync.Workflow/PoliSync.Workflow.csproj                backend/polisync/src/PoliSync.Workflow/
# These 6 were missing — all referenced by PoliSync.ApiHost.csproj
COPY backend/polisync/src/PoliSync.RulesEngine/PoliSync.RulesEngine.csproj          backend/polisync/src/PoliSync.RulesEngine/
COPY backend/polisync/src/PoliSync.Quoting/PoliSync.Quoting.csproj                  backend/polisync/src/PoliSync.Quoting/
COPY backend/polisync/src/PoliSync.VehicleInsurance/PoliSync.VehicleInsurance.csproj backend/polisync/src/PoliSync.VehicleInsurance/
COPY backend/polisync/src/PoliSync.LifeInsurance/PoliSync.LifeInsurance.csproj       backend/polisync/src/PoliSync.LifeInsurance/
COPY backend/polisync/src/PoliSync.CRM/PoliSync.CRM.csproj                          backend/polisync/src/PoliSync.CRM/
COPY backend/polisync/src/PoliSync.Actuarial/PoliSync.Actuarial.csproj              backend/polisync/src/PoliSync.Actuarial/

# Restore from ApiHost project (resolves all transitive deps, skips missing test stubs)
RUN dotnet restore backend/polisync/src/PoliSync.ApiHost/PoliSync.ApiHost.csproj

# Copy all source code
COPY backend/polisync/src/ backend/polisync/src/

# Build and publish
RUN dotnet publish backend/polisync/src/PoliSync.ApiHost/PoliSync.ApiHost.csproj \
    -c Release \
    -o /app/publish

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app

# Copy published app
COPY --from=build /app/publish .

# Create logs directory
RUN mkdir -p /app/logs

# Expose gRPC + HTTP ports for all PoliSync services
EXPOSE 50120 50121 50130 50131 50140 50141 50150 50151 50160 50161 50170 50171 50210 50211

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD bash -lc "exec 3<>/dev/tcp/127.0.0.1/50121 && printf 'GET /health HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n' >&3 && grep -q '200' <&3" || exit 1

# Run
ENTRYPOINT ["dotnet", "PoliSync.ApiHost.dll"]
