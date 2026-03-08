using Grpc.Net.Client;
using Insuretech.Insurance.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;
using Google.Protobuf.WellKnownTypes;
using Npgsql;

Console.WriteLine("=== Insurance Service CRUD Test ===\n");

// Build connection string from environment variables
var pgHost = Environment.GetEnvironmentVariable("PGHOST") ?? "localhost";
var pgPort = Environment.GetEnvironmentVariable("PGPORT") ?? "5432";
var pgDatabase = Environment.GetEnvironmentVariable("PGDATABASE") ?? "insuretech_db";
var pgUser = Environment.GetEnvironmentVariable("PGUSER") ?? "insuretech_user";
var pgPassword = Environment.GetEnvironmentVariable("PGPASSWORD") ?? "insuretech_pass_2024";
var pgSslMode = Environment.GetEnvironmentVariable("PGSSLMODE") ?? "disable";

var connectionString = $"Host={pgHost};Port={pgPort};Database={pgDatabase};Username={pgUser};Password={pgPassword};SSL Mode={pgSslMode}";
string validUserUuid;

using (var conn = new NpgsqlConnection(connectionString))
{
    await conn.OpenAsync();
    using var cmd = new NpgsqlCommand("SELECT user_id FROM authn_schema.users LIMIT 1", conn);
    var result = await cmd.ExecuteScalarAsync();
    
    if (result == null)
    {
        Console.WriteLine("ERROR: No users found in database. Please create a user first.");
        return;
    }
    
    validUserUuid = result.ToString()!;
    Console.WriteLine($"Using user UUID: {validUserUuid}\n");
}

// Create gRPC channel
var channel = GrpcChannel.ForAddress("http://localhost:50115");
var client = new InsuranceService.InsuranceServiceClient(channel);

try
{
    // ============================================
    // PRODUCT CRUD TESTS
    // ============================================
    Console.WriteLine("=== PRODUCT CRUD TESTS ===\n");
    
    var productId = Guid.NewGuid().ToString();
    
    // CREATE Product
    Console.WriteLine("1. Creating Product...");
    var createProductRequest = new CreateProductRequest
    {
        Product = new Product
        {
            ProductId = productId,
            ProductCode = "TEST-HEALTH-001",
            ProductName = "Test Health Insurance",
            Category = ProductCategory.Health,
            Description = "Comprehensive health insurance for testing",
            BasePremium = new Money { Amount = 500000, Currency = "BDT" }, // 5000 BDT
            MinSumInsured = new Money { Amount = 10000000, Currency = "BDT" }, // 100,000 BDT
            MaxSumInsured = new Money { Amount = 50000000, Currency = "BDT" }, // 500,000 BDT
            MinTenureMonths = 12,
            MaxTenureMonths = 60,
            Status = ProductStatus.Active,
            CreatedBy = validUserUuid,
            ProductAttributes = "{\"features\": [\"hospitalization\", \"surgery\"]}"
        }
    };
    createProductRequest.Product.Exclusions.Add("Pre-existing conditions");
    createProductRequest.Product.Exclusions.Add("Cosmetic surgery");
    
    var createdProduct = await client.CreateProductAsync(createProductRequest);
    Console.WriteLine($"✓ Product created: {createdProduct.Product.ProductId}");
    Console.WriteLine($"  Name: {createdProduct.Product.ProductName}");
    Console.WriteLine($"  Premium: {createdProduct.Product.BasePremium.Amount / 100.0} {createdProduct.Product.BasePremium.Currency}");
    Console.WriteLine($"  Exclusions: {string.Join(", ", createdProduct.Product.Exclusions)}");
    
    // READ Product
    Console.WriteLine("\n2. Reading Product...");
    var getProductRequest = new GetProductRequest { ProductId = productId };
    var retrievedProduct = await client.GetProductAsync(getProductRequest);
    Console.WriteLine($"✓ Product retrieved: {retrievedProduct.Product.ProductName}");
    Console.WriteLine($"  Status: {retrievedProduct.Product.Status}");
    Console.WriteLine($"  Category: {retrievedProduct.Product.Category}");
    
    // UPDATE Product
    Console.WriteLine("\n3. Updating Product...");
    retrievedProduct.Product.ProductName = "Updated Test Health Insurance";
    retrievedProduct.Product.BasePremium.Amount = 600000; // 6000 BDT
    retrievedProduct.Product.Description = "Updated comprehensive health insurance";
    var updateProductRequest = new UpdateProductRequest { Product = retrievedProduct.Product };
    var updatedProduct = await client.UpdateProductAsync(updateProductRequest);
    Console.WriteLine($"✓ Product updated: {updatedProduct.Product.ProductName}");
    Console.WriteLine($"  New Premium: {updatedProduct.Product.BasePremium.Amount / 100.0} {updatedProduct.Product.BasePremium.Currency}");
    
    // LIST Products
    Console.WriteLine("\n4. Listing Products...");
    var listProductsRequest = new ListProductsRequest { Page = 1, PageSize = 10 };
    var listProductsResponse = await client.ListProductsAsync(listProductsRequest);
    Console.WriteLine($"✓ Found {listProductsResponse.Total} products");
    foreach (var p in listProductsResponse.Products)
    {
        Console.WriteLine($"  - {p.ProductName} ({p.ProductCode})");
    }
    
    // ============================================
    // POLICY CRUD TESTS
    // ============================================
    Console.WriteLine("\n=== POLICY CRUD TESTS ===\n");
    
    var policyId = Guid.NewGuid().ToString();
    var policyNumber = $"LBT-{DateTime.Now.Year}-TEST-{new Random().Next(100000, 999999)}";
    
    // CREATE Policy
    Console.WriteLine("1. Creating Policy...");
    var createPolicyRequest = new CreatePolicyRequest
    {
        Policy = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = policyId,
            PolicyNumber = policyNumber,
            ProductId = productId,
            CustomerId = validUserUuid,
            Status = PolicyStatus.Active,
            PremiumAmount = new Money { Amount = 600000, Currency = "BDT" },
            SumInsured = new Money { Amount = 20000000, Currency = "BDT" },
            TenureMonths = 12,
            StartDate = Timestamp.FromDateTime(DateTime.UtcNow),
            EndDate = Timestamp.FromDateTime(DateTime.UtcNow.AddMonths(12)),
            IssuedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            PaymentFrequency = "YEARLY",
            VatTax = new Money { Amount = 90000, Currency = "BDT" },
            ServiceFee = new Money { Amount = 10000, Currency = "BDT" },
            TotalPayable = new Money { Amount = 700000, Currency = "BDT" },
            ProviderName = "LabAid InsureTech",
            UnderwritingData = "{\"risk_score\": 75, \"approved_by\": \"system\"}"
        }
    };
    
    var createdPolicy = await client.CreatePolicyAsync(createPolicyRequest);
    Console.WriteLine($"✓ Policy created: {createdPolicy.Policy.PolicyNumber}");
    Console.WriteLine($"  Premium: {createdPolicy.Policy.PremiumAmount.Amount / 100.0} {createdPolicy.Policy.PremiumAmount.Currency}");
    Console.WriteLine($"  Sum Insured: {createdPolicy.Policy.SumInsured.Amount / 100.0} {createdPolicy.Policy.SumInsured.Currency}");
    Console.WriteLine($"  Status: {createdPolicy.Policy.Status}");
    
    // READ Policy
    Console.WriteLine("\n2. Reading Policy...");
    var getPolicyRequest = new GetPolicyRequest { PolicyId = policyId };
    var retrievedPolicy = await client.GetPolicyAsync(getPolicyRequest);
    Console.WriteLine($"✓ Policy retrieved: {retrievedPolicy.Policy.PolicyNumber}");
    Console.WriteLine($"  Customer: {retrievedPolicy.Policy.CustomerId}");
    Console.WriteLine($"  Tenure: {retrievedPolicy.Policy.TenureMonths} months");
    
    // UPDATE Policy
    Console.WriteLine("\n3. Updating Policy...");
    retrievedPolicy.Policy.Status = PolicyStatus.Active;
    retrievedPolicy.Policy.PremiumAmount.Amount = 650000;
    retrievedPolicy.Policy.TotalPayable.Amount = 750000;
    var updatePolicyRequest = new UpdatePolicyRequest { Policy = retrievedPolicy.Policy };
    var updatedPolicy = await client.UpdatePolicyAsync(updatePolicyRequest);
    Console.WriteLine($"✓ Policy updated: {updatedPolicy.Policy.PolicyNumber}");
    Console.WriteLine($"  New Premium: {updatedPolicy.Policy.PremiumAmount.Amount / 100.0} {updatedPolicy.Policy.PremiumAmount.Currency}");
    Console.WriteLine($"  Status: {updatedPolicy.Policy.Status}");
    
    // LIST Policies
    Console.WriteLine("\n4. Listing Policies...");
    var listPoliciesRequest = new ListPoliciesRequest { Page = 1, PageSize = 10 };
    var listPoliciesResponse = await client.ListPoliciesAsync(listPoliciesRequest);
    Console.WriteLine($"✓ Found {listPoliciesResponse.Total} policies");
    foreach (var pol in listPoliciesResponse.Policies)
    {
        Console.WriteLine($"  - {pol.PolicyNumber} (Status: {pol.Status})");
    }
    
    // ============================================
    // CLAIM CRUD TESTS
    // ============================================
    Console.WriteLine("\n=== CLAIM CRUD TESTS ===\n");
    
    var claimId = Guid.NewGuid().ToString();
    var claimNumber = $"CLM-{DateTime.Now.Year}-TEST-{new Random().Next(100000, 999999)}";
    
    // CREATE Claim
    Console.WriteLine("1. Creating Claim...");
    var createClaimRequest = new CreateClaimRequest
    {
        Claim = new Claim
        {
            ClaimId = claimId,
            ClaimNumber = claimNumber,
            PolicyId = policyId,
            CustomerId = validUserUuid,
            Status = ClaimStatus.Submitted,
            Type = ClaimType.HealthHospitalization,
            ClaimedAmount = new Money { Amount = 5000000, Currency = "BDT" },
            IncidentDate = Timestamp.FromDateTime(DateTime.UtcNow.AddDays(-7)),
            IncidentDescription = "Emergency hospitalization for testing",
            SubmittedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            PlaceOfIncident = "Dhaka Medical College Hospital",
            ProcessingType = ClaimProcessingType.Manual,
            InAppMessages = "[{\"message\": \"Claim submitted successfully\", \"timestamp\": \"" + DateTime.UtcNow.ToString("o") + "\"}]"
        }
    };
    
    var createdClaim = await client.CreateClaimAsync(createClaimRequest);
    Console.WriteLine($"✓ Claim created: {createdClaim.Claim.ClaimNumber}");
    Console.WriteLine($"  Claimed Amount: {createdClaim.Claim.ClaimedAmount.Amount / 100.0} {createdClaim.Claim.ClaimedAmount.Currency}");
    Console.WriteLine($"  Status: {createdClaim.Claim.Status}");
    Console.WriteLine($"  Type: {createdClaim.Claim.Type}");
    
    // READ Claim
    Console.WriteLine("\n2. Reading Claim...");
    var getClaimRequest = new GetClaimRequest { ClaimId = claimId };
    var retrievedClaim = await client.GetClaimAsync(getClaimRequest);
    Console.WriteLine($"✓ Claim retrieved: {retrievedClaim.Claim.ClaimNumber}");
    Console.WriteLine($"  Policy: {retrievedClaim.Claim.PolicyId}");
    Console.WriteLine($"  Incident: {retrievedClaim.Claim.IncidentDescription}");
    
    // UPDATE Claim
    Console.WriteLine("\n3. Updating Claim...");
    retrievedClaim.Claim.Status = ClaimStatus.UnderReview;
    retrievedClaim.Claim.ApprovedAmount = new Money { Amount = 4500000, Currency = "BDT" };
    retrievedClaim.Claim.ProcessorNotes = "Claim under review by claims officer";
    var updateClaimRequest = new UpdateClaimRequest { Claim = retrievedClaim.Claim };
    var updatedClaim = await client.UpdateClaimAsync(updateClaimRequest);
    Console.WriteLine($"✓ Claim updated: {updatedClaim.Claim.ClaimNumber}");
    Console.WriteLine($"  New Status: {updatedClaim.Claim.Status}");
    Console.WriteLine($"  Approved Amount: {updatedClaim.Claim.ApprovedAmount.Amount / 100.0} {updatedClaim.Claim.ApprovedAmount.Currency}");
    
    // LIST Claims
    Console.WriteLine("\n4. Listing Claims...");
    var listClaimsRequest = new ListClaimsRequest { Page = 1, PageSize = 10 };
    var listClaimsResponse = await client.ListClaimsAsync(listClaimsRequest);
    Console.WriteLine($"✓ Found {listClaimsResponse.Total} claims");
    foreach (var clm in listClaimsResponse.Claims)
    {
        Console.WriteLine($"  - {clm.ClaimNumber} (Status: {clm.Status}, Amount: {clm.ClaimedAmount.Amount / 100.0})");
    }
    
    // ============================================
    // PRODUCT PLAN CRUD TESTS (Table 4)
    // ============================================
    Console.WriteLine("\n=== PRODUCT PLAN CRUD TESTS ===\n");
    
    var planId = Guid.NewGuid().ToString();
    
    // CREATE ProductPlan
    Console.WriteLine("1. Creating ProductPlan...");
    var createPlanRequest = new CreateProductPlanRequest
    {
        Plan = new Insuretech.Products.Entity.V1.ProductPlan
        {
            PlanId = planId,
            ProductId = productId,
            PlanName = "Gold Plan",
            PlanDescription = "Premium coverage with maximum benefits",
            PremiumAmount = new Money { Amount = 800000, Currency = "BDT" },
            MinSumInsured = new Money { Amount = 20000000, Currency = "BDT" },
            MaxSumInsured = new Money { Amount = 100000000, Currency = "BDT" }
        }
    };
    
    var createdPlan = await client.CreateProductPlanAsync(createPlanRequest);
    Console.WriteLine($"✓ ProductPlan created: {createdPlan.Plan.PlanName}");
    Console.WriteLine($"  Premium: {createdPlan.Plan.PremiumAmount.Amount / 100.0} {createdPlan.Plan.PremiumAmount.Currency}");
    
    // READ ProductPlan
    Console.WriteLine("\n2. Reading ProductPlan...");
    var getPlanRequest = new GetProductPlanRequest { PlanId = planId };
    var retrievedPlan = await client.GetProductPlanAsync(getPlanRequest);
    Console.WriteLine($"✓ ProductPlan retrieved: {retrievedPlan.Plan.PlanName}");
    
    // UPDATE ProductPlan
    Console.WriteLine("\n3. Updating ProductPlan...");
    retrievedPlan.Plan.PlanName = "Platinum Plan";
    retrievedPlan.Plan.PremiumAmount.Amount = 900000;
    var updatePlanRequest = new UpdateProductPlanRequest { Plan = retrievedPlan.Plan };
    var updatedPlan = await client.UpdateProductPlanAsync(updatePlanRequest);
    Console.WriteLine($"✓ ProductPlan updated: {updatedPlan.Plan.PlanName}");
    Console.WriteLine($"  New Premium: {updatedPlan.Plan.PremiumAmount.Amount / 100.0} {updatedPlan.Plan.PremiumAmount.Currency}");
    
    // LIST ProductPlans
    Console.WriteLine("\n4. Listing ProductPlans...");
    var listPlansRequest = new ListProductPlansRequest { Page = 1, PageSize = 10 };
    var listPlansResponse = await client.ListProductPlansAsync(listPlansRequest);
    Console.WriteLine($"✓ Found {listPlansResponse.Total} product plans");
    
    // ============================================
    // DELETE ADDITIONAL TESTS
    // ============================================
    Console.WriteLine("\n=== DELETING ADDITIONAL TEST DATA ===\n");
    
    // DELETE ProductPlan
    Console.WriteLine("1. Deleting ProductPlan...");
    await client.DeleteProductPlanAsync(new DeleteProductPlanRequest { PlanId = planId });
    Console.WriteLine($"✓ ProductPlan deleted");
    
    // ============================================
    // DELETE TESTS
    // ============================================
    Console.WriteLine("\n=== DELETE CORE TEST DATA ===\n");
    
    // DELETE Claim
    Console.WriteLine("1. Deleting Claim...");
    var deleteClaimRequest = new DeleteClaimRequest { ClaimId = claimId };
    await client.DeleteClaimAsync(deleteClaimRequest);
    Console.WriteLine($"✓ Claim deleted: {claimNumber}");
    
    // DELETE Policy
    Console.WriteLine("\n2. Deleting Policy...");
    var deletePolicyRequest = new DeletePolicyRequest { PolicyId = policyId };
    await client.DeletePolicyAsync(deletePolicyRequest);
    Console.WriteLine($"✓ Policy deleted: {policyNumber}");
    
    // DELETE Product
    Console.WriteLine("\n3. Deleting Product...");
    var deleteProductRequest = new DeleteProductRequest { ProductId = productId };
    await client.DeleteProductAsync(deleteProductRequest);
    Console.WriteLine($"✓ Product deleted: {productId}");
    
    Console.WriteLine("\n=== ALL TESTS PASSED ===");
}
catch (Exception ex)
{
    Console.WriteLine($"\n❌ ERROR: {ex.Message}");
    Console.WriteLine($"Stack Trace: {ex.StackTrace}");
    if (ex.InnerException != null)
    {
        Console.WriteLine($"Inner Exception: {ex.InnerException.Message}");
    }
}
finally
{
    await channel.ShutdownAsync();
}
