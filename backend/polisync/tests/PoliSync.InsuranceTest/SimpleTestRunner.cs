using Grpc.Core;
using Grpc.Net.Client;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Insurance.Services.V1;
using Insuretech.Common.V1;
using ProductEntity = Insuretech.Products.Entity.V1;
using InsurerEntity = Insuretech.Insurer.Entity.V1;
using FraudEntity = Insuretech.Fraud.Entity.V1;

namespace PoliSync.InsuranceTest;

/// <summary>
/// Simple test runner that tests basic CRUD operations for Products table
/// This validates that the Insurance Service is working correctly
/// </summary>
public class SimpleTestRunner
{
    private readonly GrpcChannel _channel;
    private readonly InsuranceService.InsuranceServiceClient _client;
    private readonly string _validUserUuid;
    
    public int TotalTests { get; private set; }
    public int PassedTests { get; private set; }
    public int FailedTests { get; private set; }
    
    private readonly List<string> _createdProductIds = new();
    
    public SimpleTestRunner(GrpcChannel channel, string validUserUuid)
    {
        _channel = channel;
        _client = new InsuranceService.InsuranceServiceClient(channel);
        _validUserUuid = validUserUuid;
    }
    
    private void LogTest(string testName, bool success, string? error = null)
    {
        TotalTests++;
        if (success)
        {
            PassedTests++;
            Console.WriteLine($"  ✓ {testName}");
        }
        else
        {
            FailedTests++;
            Console.WriteLine($"  ✗ {testName}");
            if (error != null)
            {
                Console.WriteLine($"    Error: {error}");
            }
        }
    }
    
    public async Task RunAllTests()
    {
        Console.WriteLine("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
        Console.WriteLine("COMPREHENSIVE CRUD TESTS - Insurance Schema");
        Console.WriteLine("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
        Console.WriteLine();
        
        await TestProductCRUD();
        Console.WriteLine();
        
        await TestProductPlanCRUD();
        Console.WriteLine();
        
        await TestRiderCRUD();
        Console.WriteLine();
        
        await TestInsurerCRUD();
        Console.WriteLine();
        
        await TestFraudRuleCRUD();
        Console.WriteLine();
        
        Console.WriteLine("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
        Console.WriteLine("TEST SUMMARY");
        Console.WriteLine("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
        Console.WriteLine($"Total Tests: {TotalTests}");
        Console.WriteLine($"Passed: {PassedTests}");
        Console.WriteLine($"Failed: {FailedTests}");
        Console.WriteLine($"Success Rate: {(PassedTests * 100.0 / TotalTests):F1}%");
        Console.WriteLine();
    }
    
    private async Task TestProductCRUD()
    {
        Console.WriteLine("Testing Product CRUD Operations...");
        Console.WriteLine();
        
        string? productId = null;
        
        try
        {
            // TEST 1: CREATE
            productId = Guid.NewGuid().ToString();
            var timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() % 1000; // Last 3 digits of timestamp
            var createRequest = new CreateProductRequest
            {
                Product = new ProductEntity.Product
                {
                    ProductId = productId,
                    ProductCode = $"TST-{timestamp:D3}", // Must match pattern: [A-Z]{3}-[0-9]{3}
                    ProductName = "Test Health Insurance Product",
                    Category = ProductEntity.ProductCategory.Health,
                    Description = "Comprehensive health insurance for testing",
                    BasePremium = new Money { Amount = 500000, Currency = "BDT" },
                    MinSumInsured = new Money { Amount = 10000000, Currency = "BDT" },
                    MaxSumInsured = new Money { Amount = 50000000, Currency = "BDT" },
                    MinTenureMonths = 12,
                    MaxTenureMonths = 60,
                    Status = ProductEntity.ProductStatus.Active,
                    CreatedBy = _validUserUuid
                }
            };
            createRequest.Product.Exclusions.Add("Pre-existing conditions");
            createRequest.Product.Exclusions.Add("Cosmetic procedures");
            
            var createResponse = await _client.CreateProductAsync(createRequest);
            bool createSuccess = createResponse?.Product?.ProductId == productId;
            LogTest("CREATE Product", createSuccess);
            
            if (!createSuccess)
            {
                return; // Can't continue without a created product
            }
            
            _createdProductIds.Add(productId);
            
            // TEST 2: READ (Get by ID)
            var getRequest = new GetProductRequest { ProductId = productId };
            var getResponse = await _client.GetProductAsync(getRequest);
            bool getSuccess = getResponse?.Product?.ProductId == productId &&
                            getResponse?.Product?.ProductName == "Test Health Insurance Product";
            LogTest("READ Product (GetByID)", getSuccess);
            
            // TEST 3: UPDATE
            if (getResponse?.Product != null)
            {
                var updateRequest = new UpdateProductRequest
                {
                    Product = getResponse.Product
                };
                updateRequest.Product.ProductName = "Updated Test Health Insurance";
                updateRequest.Product.Description = "Updated description for testing";
                
                var updateResponse = await _client.UpdateProductAsync(updateRequest);
                bool updateSuccess = updateResponse?.Product?.ProductName == "Updated Test Health Insurance" &&
                                   updateResponse?.Product?.Description == "Updated description for testing";
                LogTest("UPDATE Product", updateSuccess);
            }
            else
            {
                LogTest("UPDATE Product", false, "Could not get product for update");
            }
            
            // TEST 4: LIST
            var listRequest = new ListProductsRequest
            {
                TenantId = "",
                Page = 1,
                PageSize = 10
            };
            var listResponse = await _client.ListProductsAsync(listRequest);
            bool listSuccess = listResponse?.Products?.Count > 0;
            LogTest("LIST Products", listSuccess);
            
            if (listSuccess)
            {
                Console.WriteLine($"    Found {listResponse.Products.Count} products (Total: {listResponse.Total})");
            }
            
            // TEST 5: DELETE
            var deleteRequest = new DeleteProductRequest { ProductId = productId };
            await _client.DeleteProductAsync(deleteRequest);
            
            // Verify deletion by trying to get the product
            try
            {
                await _client.GetProductAsync(new GetProductRequest { ProductId = productId });
                LogTest("DELETE Product", false, "Product still exists after deletion");
            }
            catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
            {
                LogTest("DELETE Product", true);
                _createdProductIds.Remove(productId);
            }
        }
        catch (RpcException ex)
        {
            LogTest($"Product CRUD (RPC Error)", false, $"{ex.StatusCode}: {ex.Status.Detail}");
        }
        catch (Exception ex)
        {
            LogTest($"Product CRUD (Exception)", false, ex.Message);
        }
    }
    
    private async Task TestProductPlanCRUD()
    {
        Console.WriteLine("Testing ProductPlan CRUD Operations...");
        
        string? productId = null;
        string? planId = null;
        
        try
        {
            // First create a product to associate the plan with
            productId = Guid.NewGuid().ToString();
            var timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() % 1000;
            var createProductRequest = new CreateProductRequest
            {
                Product = new ProductEntity.Product
                {
                    ProductId = productId,
                    ProductCode = $"PLN-{timestamp:D3}",
                    ProductName = "Test Product for Plans",
                    Category = ProductEntity.ProductCategory.Health,
                    BasePremium = new Money { Amount = 100000, Currency = "BDT" },
                    MinSumInsured = new Money { Amount = 1000000, Currency = "BDT" },
                    MaxSumInsured = new Money { Amount = 5000000, Currency = "BDT" },
                    MinTenureMonths = 12,
                    MaxTenureMonths = 36,
                    Status = ProductEntity.ProductStatus.Active,
                    CreatedBy = _validUserUuid
                }
            };
            await _client.CreateProductAsync(createProductRequest);
            _createdProductIds.Add(productId);
            
            // TEST 1: CREATE ProductPlan
            planId = Guid.NewGuid().ToString();
            var createRequest = new CreateProductPlanRequest
            {
                Plan = new ProductEntity.ProductPlan
                {
                    PlanId = planId,
                    ProductId = productId,
                    PlanName = "Gold Plan",
                    PlanDescription = "Premium coverage plan",
                    PremiumAmount = new Money { Amount = 150000, Currency = "BDT" },
                    MinSumInsured = new Money { Amount = 2000000, Currency = "BDT" },
                    MaxSumInsured = new Money { Amount = 5000000, Currency = "BDT" }
                }
            };
            
            var createResponse = await _client.CreateProductPlanAsync(createRequest);
            bool createSuccess = createResponse?.Plan?.PlanId == planId;
            LogTest("CREATE ProductPlan", createSuccess);
            
            if (!createSuccess) return;
            
            // TEST 2: READ (Get by ID)
            var getRequest = new GetProductPlanRequest { PlanId = planId };
            var getResponse = await _client.GetProductPlanAsync(getRequest);
            bool getSuccess = getResponse?.Plan?.PlanId == planId &&
                            getResponse?.Plan?.PlanName == "Gold Plan";
            LogTest("READ ProductPlan (GetByID)", getSuccess);
            
            // TEST 3: LIST by Product
            var listRequest = new ListProductPlansRequest { ProductId = productId };
            var listResponse = await _client.ListProductPlansAsync(listRequest);
            bool listSuccess = listResponse?.Plans?.Count > 0;
            LogTest("LIST ProductPlans", listSuccess);
            
            if (listSuccess)
            {
                Console.WriteLine($"    Found {listResponse.Plans.Count} plans for product");
            }
            
            // Note: ProductPlan doesn't have Update/Delete methods in the proto
        }
        catch (RpcException ex)
        {
            LogTest($"ProductPlan CRUD (RPC Error)", false, $"{ex.StatusCode}: {ex.Status.Detail}");
        }
        catch (Exception ex)
        {
            LogTest($"ProductPlan CRUD (Exception)", false, ex.Message);
        }
        finally
        {
            // Cleanup: Delete the product (cascade will delete plans)
            if (productId != null)
            {
                try
                {
                    await _client.DeleteProductAsync(new DeleteProductRequest { ProductId = productId });
                    _createdProductIds.Remove(productId);
                }
                catch { /* Ignore cleanup errors */ }
            }
        }
    }
    
    private async Task TestRiderCRUD()
    {
        Console.WriteLine("Testing Rider CRUD Operations...");
        
        string? productId = null;
        string? riderId = null;
        
        try
        {
            // First create a product to associate the rider with
            productId = Guid.NewGuid().ToString();
            var timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() % 1000;
            var createProductRequest = new CreateProductRequest
            {
                Product = new ProductEntity.Product
                {
                    ProductId = productId,
                    ProductCode = $"RDR-{timestamp:D3}",
                    ProductName = "Test Product for Riders",
                    Category = ProductEntity.ProductCategory.Health,
                    BasePremium = new Money { Amount = 100000, Currency = "BDT" },
                    MinSumInsured = new Money { Amount = 1000000, Currency = "BDT" },
                    MaxSumInsured = new Money { Amount = 5000000, Currency = "BDT" },
                    MinTenureMonths = 12,
                    MaxTenureMonths = 36,
                    Status = ProductEntity.ProductStatus.Active,
                    CreatedBy = _validUserUuid
                }
            };
            await _client.CreateProductAsync(createProductRequest);
            _createdProductIds.Add(productId);
            
            // TEST 1: CREATE Rider
            riderId = Guid.NewGuid().ToString();
            var createRequest = new CreateRiderRequest
            {
                Rider = new ProductEntity.Rider
                {
                    RiderId = riderId,
                    ProductId = productId,
                    RiderName = "Critical Illness Rider",
                    Description = "Additional coverage for critical illnesses",
                    PremiumAmount = new Money { Amount = 50000, Currency = "BDT" },
                    CoverageAmount = new Money { Amount = 1000000, Currency = "BDT" },
                    IsMandatory = false
                }
            };
            
            var createResponse = await _client.CreateRiderAsync(createRequest);
            bool createSuccess = createResponse?.Rider?.RiderId == riderId;
            LogTest("CREATE Rider", createSuccess);
            
            if (!createSuccess) return;
            
            // TEST 2: READ (Get by ID)
            var getRequest = new GetRiderRequest { RiderId = riderId };
            var getResponse = await _client.GetRiderAsync(getRequest);
            bool getSuccess = getResponse?.Rider?.RiderId == riderId &&
                            getResponse?.Rider?.RiderName == "Critical Illness Rider";
            LogTest("READ Rider (GetByID)", getSuccess);
            
            // TEST 3: LIST by Product
            var listRequest = new ListRidersRequest { ProductId = productId };
            var listResponse = await _client.ListRidersAsync(listRequest);
            bool listSuccess = listResponse?.Riders?.Count > 0;
            LogTest("LIST Riders", listSuccess);
            
            if (listSuccess)
            {
                Console.WriteLine($"    Found {listResponse.Riders.Count} riders for product");
            }
            
            // Note: Rider doesn't have Update/Delete methods in the proto
        }
        catch (RpcException ex)
        {
            LogTest($"Rider CRUD (RPC Error)", false, $"{ex.StatusCode}: {ex.Status.Detail}");
        }
        catch (Exception ex)
        {
            LogTest($"Rider CRUD (Exception)", false, ex.Message);
        }
        finally
        {
            // Cleanup: Delete the product (cascade will delete riders)
            if (productId != null)
            {
                try
                {
                    await _client.DeleteProductAsync(new DeleteProductRequest { ProductId = productId });
                    _createdProductIds.Remove(productId);
                }
                catch { /* Ignore cleanup errors */ }
            }
        }
    }
    
    private async Task TestInsurerCRUD()
    {
        Console.WriteLine("Testing Insurer CRUD Operations...");
        
        string? insurerId = null;
        
        try
        {
            // TEST 1: CREATE
            insurerId = Guid.NewGuid().ToString();
            var timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() % 1000;
            var createRequest = new CreateInsurerRequest
            {
                Insurer = new InsurerEntity.Insurer
                {
                    Id = insurerId,
                    Name = "Test Insurance Company Ltd",
                    Code = $"TIC{timestamp:D3}",
                    Type = InsurerEntity.InsurerType.NonLife,
                    Status = InsurerEntity.InsurerStatus.Active,
                    TradeLicenseNumber = $"TL-{Guid.NewGuid().ToString().Substring(0, 8).ToUpper()}",
                    TinNumber = $"TIN{timestamp:D9}",
                    ContactInfo = new ContactInfo
                    {
                        Email = "contact@testinsurer.com",
                        MobileNumber = "+8801712345678"
                    },
                    RegisteredAddress = new Address
                    {
                        AddressLine1 = "123 Test Street",
                        City = "Dhaka",
                        Country = "Bangladesh"
                    },
                    HeadOfficeAddress = new Address
                    {
                        AddressLine1 = "123 Test Street",
                        City = "Dhaka",
                        Country = "Bangladesh"
                    },
                    AuditInfo = new AuditInfo
                    {
                        CreatedBy = _validUserUuid
                    }
                }
            };
            
            var createResponse = await _client.CreateInsurerAsync(createRequest);
            bool createSuccess = createResponse?.Insurer?.Id == insurerId;
            LogTest("CREATE Insurer", createSuccess);
            
            if (!createSuccess) return;
            
            // TEST 2: READ (Get by ID)
            var getRequest = new GetInsurerRequest { InsurerId = insurerId };
            var getResponse = await _client.GetInsurerAsync(getRequest);
            bool getSuccess = getResponse?.Insurer?.Id == insurerId &&
                            getResponse?.Insurer?.Name == "Test Insurance Company Ltd";
            LogTest("READ Insurer (GetByID)", getSuccess);
            
            // TEST 3: UPDATE
            if (getResponse?.Insurer != null)
            {
                var updateRequest = new UpdateInsurerRequest
                {
                    Insurer = getResponse.Insurer
                };
                updateRequest.Insurer.Name = "Updated Test Insurance Company";
                updateRequest.Insurer.ContactInfo.Email = "updated@testinsurer.com";
                
                var updateResponse = await _client.UpdateInsurerAsync(updateRequest);
                bool updateSuccess = updateResponse?.Insurer?.Name == "Updated Test Insurance Company" &&
                                   updateResponse?.Insurer?.ContactInfo?.Email == "updated@testinsurer.com";
                LogTest("UPDATE Insurer", updateSuccess);
            }
            else
            {
                LogTest("UPDATE Insurer", false, "Could not get insurer for update");
            }
            
            // TEST 4: LIST
            var listRequest = new ListInsurersRequest
            {
                Page = 1,
                PageSize = 10
            };
            var listResponse = await _client.ListInsurersAsync(listRequest);
            bool listSuccess = listResponse?.Insurers?.Count > 0;
            LogTest("LIST Insurers", listSuccess);
            
            if (listSuccess)
            {
                Console.WriteLine($"    Found {listResponse.Insurers.Count} insurers (Total: {listResponse.Total})");
            }
            
            // TEST 5: DELETE
            var deleteRequest = new DeleteInsurerRequest { InsurerId = insurerId };
            await _client.DeleteInsurerAsync(deleteRequest);
            
            // Verify deletion
            try
            {
                await _client.GetInsurerAsync(new GetInsurerRequest { InsurerId = insurerId });
                LogTest("DELETE Insurer", false, "Insurer still exists after deletion");
            }
            catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
            {
                LogTest("DELETE Insurer", true);
            }
        }
        catch (RpcException ex)
        {
            LogTest($"Insurer CRUD (RPC Error)", false, $"{ex.StatusCode}: {ex.Status.Detail}");
        }
        catch (Exception ex)
        {
            LogTest($"Insurer CRUD (Exception)", false, ex.Message);
        }
    }
    
    private async Task TestFraudRuleCRUD()
    {
        Console.WriteLine("Testing FraudRule CRUD Operations...");
        
        string? ruleId = null;
        
        try
        {
            // TEST 1: CREATE
            ruleId = Guid.NewGuid().ToString();
            var timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds() % 10000;
            var createRequest = new CreateFraudRuleRequest
            {
                Rule = new FraudEntity.FraudRule
                {
                    FraudRuleId = ruleId,
                    Name = $"High Claim Amount Rule {timestamp}",
                    Description = "Flag claims above 500,000 BDT",
                    Category = FraudEntity.RuleCategory.AmountAnomaly,
                    RiskLevel = FraudEntity.RiskLevel.High,
                    ScoreWeight = 50,
                    IsActive = true,
                    Conditions = "{}",
                    AuditInfo = new AuditInfo
                    {
                        CreatedBy = _validUserUuid
                    }
                }
            };
            
            var createResponse = await _client.CreateFraudRuleAsync(createRequest);
            bool createSuccess = createResponse?.Rule?.FraudRuleId == ruleId;
            LogTest("CREATE FraudRule", createSuccess);
            
            if (!createSuccess) return;
            
            // TEST 2: READ (Get by ID)
            var getRequest = new GetFraudRuleRequest { FraudRuleId = ruleId };
            var getResponse = await _client.GetFraudRuleAsync(getRequest);
            bool getSuccess = getResponse?.Rule?.FraudRuleId == ruleId &&
                            getResponse?.Rule?.Name.StartsWith("High Claim Amount Rule") == true;
            LogTest("READ FraudRule (GetByID)", getSuccess);
            
            // TEST 3: UPDATE
            if (getResponse?.Rule != null)
            {
                var updateRequest = new UpdateFraudRuleRequest
                {
                    Rule = getResponse.Rule
                };
                updateRequest.Rule.Name = $"Updated High Claim Rule {timestamp}";
                updateRequest.Rule.RiskLevel = FraudEntity.RiskLevel.Critical;
                
                var updateResponse = await _client.UpdateFraudRuleAsync(updateRequest);
                bool updateSuccess = updateResponse?.Rule?.Name.StartsWith("Updated High Claim Rule") == true &&
                                   updateResponse?.Rule?.RiskLevel == FraudEntity.RiskLevel.Critical;
                LogTest("UPDATE FraudRule", updateSuccess);
            }
            else
            {
                LogTest("UPDATE FraudRule", false, "Could not get fraud rule for update");
            }
            
            // TEST 4: LIST
            var listRequest = new ListFraudRulesRequest
            {
                Page = 1,
                PageSize = 10
            };
            var listResponse = await _client.ListFraudRulesAsync(listRequest);
            bool listSuccess = listResponse?.Rules?.Count > 0;
            LogTest("LIST FraudRules", listSuccess);
            
            if (listSuccess)
            {
                Console.WriteLine($"    Found {listResponse.Rules.Count} fraud rules (Total: {listResponse.Total})");
            }
            
            // TEST 5: LIST ACTIVE
            var listActiveRequest = new ListActiveFraudRulesRequest();
            var listActiveResponse = await _client.ListActiveFraudRulesAsync(listActiveRequest);
            bool listActiveSuccess = listActiveResponse?.Rules?.Count > 0;
            LogTest("LIST Active FraudRules", listActiveSuccess);
            
            if (listActiveSuccess)
            {
                Console.WriteLine($"    Found {listActiveResponse.Rules.Count} active fraud rules");
            }
            
            // TEST 6: DELETE
            var deleteRequest = new DeleteFraudRuleRequest { FraudRuleId = ruleId };
            await _client.DeleteFraudRuleAsync(deleteRequest);
            
            // Verify deletion
            try
            {
                await _client.GetFraudRuleAsync(new GetFraudRuleRequest { FraudRuleId = ruleId });
                LogTest("DELETE FraudRule", false, "FraudRule still exists after deletion");
            }
            catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
            {
                LogTest("DELETE FraudRule", true);
            }
        }
        catch (RpcException ex)
        {
            LogTest($"FraudRule CRUD (RPC Error)", false, $"{ex.StatusCode}: {ex.Status.Detail}");
        }
        catch (Exception ex)
        {
            LogTest($"FraudRule CRUD (Exception)", false, ex.Message);
        }
    }
}
