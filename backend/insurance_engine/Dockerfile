# Stage 1: Base - ASP.NET runtime for .NET 10
FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS base
WORKDIR /app
EXPOSE 8080
EXPOSE 8081

# Stage 2: Build - SDK for building the code
FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

# Step 1: Copy all csproj files to leverage Docker layer caching for Restoring
# Correct paths from repo root
COPY ["backend/insurance_engine/src/InsuranceEngine.ApiHost/InsuranceEngine.ApiHost.csproj", "backend/insurance_engine/src/InsuranceEngine.ApiHost/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.SharedKernel/InsuranceEngine.SharedKernel.csproj", "backend/insurance_engine/src/InsuranceEngine.SharedKernel/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Beneficiary/InsuranceEngine.Beneficiary.csproj", "backend/insurance_engine/src/InsuranceEngine.Beneficiary/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Products/InsuranceEngine.Products.csproj", "backend/insurance_engine/src/InsuranceEngine.Products/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Policy/InsuranceEngine.Policy.csproj", "backend/insurance_engine/src/InsuranceEngine.Policy/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Underwriting/InsuranceEngine.Underwriting.csproj", "backend/insurance_engine/src/InsuranceEngine.Underwriting/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Claims/InsuranceEngine.Claims.csproj", "backend/insurance_engine/src/InsuranceEngine.Claims/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Commission/InsuranceEngine.Commission.csproj", "backend/insurance_engine/src/InsuranceEngine.Commission/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Cancellations/InsuranceEngine.Cancellations.csproj", "backend/insurance_engine/src/InsuranceEngine.Cancellations/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Renewals/InsuranceEngine.Renewals.csproj", "backend/insurance_engine/src/InsuranceEngine.Renewals/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Endorsements/InsuranceEngine.Endorsements.csproj", "backend/insurance_engine/src/InsuranceEngine.Endorsements/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.FraudDetection/InsuranceEngine.FraudDetection.csproj", "backend/insurance_engine/src/InsuranceEngine.FraudDetection/"]
COPY ["backend/insurance_engine/src/InsuranceEngine.Proto/InsuranceEngine.Proto.csproj", "backend/insurance_engine/src/InsuranceEngine.Proto/"]

# Step 2: Restore dependencies
RUN dotnet restore "backend/insurance_engine/src/InsuranceEngine.ApiHost/InsuranceEngine.ApiHost.csproj"

# Step 3: Copy all source files
COPY ["backend/insurance_engine/src/", "backend/insurance_engine/src/"]

# Step 4: Build the specific project
WORKDIR "/src/backend/insurance_engine/src/InsuranceEngine.ApiHost"
RUN dotnet build "InsuranceEngine.ApiHost.csproj" -c Release -o /app/build

# Stage 3: Publish - Optimize for production
FROM build AS publish
RUN dotnet publish "InsuranceEngine.ApiHost.csproj" -c Release -o /app/publish /p:UseAppHost=false

# Stage 4: Final - Final container with artifacts
FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .

# Environment variables
ENV ASPNETCORE_URLS=http://+:8080
ENV ASPNETCORE_ENVIRONMENT=Production

ENTRYPOINT ["dotnet", "InsuranceEngine.ApiHost.dll"]
