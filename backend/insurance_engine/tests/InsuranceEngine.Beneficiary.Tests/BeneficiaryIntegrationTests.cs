using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using InsuranceEngine.Beneficiary;
using InsuranceEngine.Beneficiary.Application.Commands;
using InsuranceEngine.Beneficiary.Application.Queries;
using InsuranceEngine.Beneficiary.Infrastructure;
using Insuretech.Beneficiary.Services.V1;
using Xunit;

namespace InsuranceEngine.Beneficiary.Tests;

public class BeneficiaryIntegrationTests
{
    private readonly IServiceProvider _serviceProvider;
    private readonly BeneficiaryDbContext _context;

    public BeneficiaryIntegrationTests()
    {
        var services = new ServiceCollection();
        services.AddLogging(builder => builder.AddConsole());

        var connectionString = "Host=localhost;Port=5432;Database=insuretech_primary;Username=postgres;Password=12345678;SearchPath=insurance_schema,public";

        services.AddDbContext<BeneficiaryDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<IBeneficiaryDataGateway, SqlBeneficiaryDataGateway>();

        _serviceProvider = services.BuildServiceProvider();
        _context = _serviceProvider.GetRequiredService<BeneficiaryDbContext>();
    }

    [Fact]
    public async Task CreateAndRetrieveIndividualBeneficiary_ShouldWork()
    {
        using var scope = _serviceProvider.CreateScope();
        var gateway = scope.ServiceProvider.GetRequiredService<IBeneficiaryDataGateway>();

        var createRequest = new CreateIndividualBeneficiaryRequest
        {
            UserId = Guid.NewGuid().ToString(),
            FullName = "Test Individual",
            DateOfBirth = "1990-01-15",
            Gender = "MALE",
            NidNumber = "1234567890",
            MobileNumber = "+8801712345678",
            Email = "test@example.com",
            PartnerId = "PARTNER001"
        };

        var createResponse = await gateway.CreateIndividualBeneficiaryAsync(createRequest);

        Assert.NotNull(createResponse);
        Assert.Empty(createResponse.Error?.Code ?? "");
        Assert.NotEmpty(createResponse.BeneficiaryId);
        Assert.NotEmpty(createResponse.BeneficiaryCode);
        Assert.Contains("IND-", createResponse.BeneficiaryCode);

        var getRequest = new GetBeneficiaryRequest
        {
            BeneficiaryId = createResponse.BeneficiaryId
        };

        var getResponse = await gateway.GetBeneficiaryAsync(getRequest);

        Assert.NotNull(getResponse);
        Assert.Null(getResponse.Error);
        Assert.NotNull(getResponse.Beneficiary);
        Assert.Equal(createResponse.BeneficiaryId, getResponse.Beneficiary.BeneficiaryId);
    }

    [Fact]
    public async Task CreateAndRetrieveBusinessBeneficiary_ShouldWork()
    {
        using var scope = _serviceProvider.CreateScope();
        var gateway = scope.ServiceProvider.GetRequiredService<IBeneficiaryDataGateway>();

        var createRequest = new CreateBusinessBeneficiaryRequest
        {
            UserId = Guid.NewGuid().ToString(),
            BusinessName = "Test Business Ltd",
            TradeLicenseNumber = "TL123456",
            TinNumber = "TIN789012",
            FocalPersonName = "John Doe",
            FocalPersonMobile = "+8801712345679",
            PartnerId = "PARTNER001"
        };

        var createResponse = await gateway.CreateBusinessBeneficiaryAsync(createRequest);

        Assert.NotNull(createResponse);
        Assert.Empty(createResponse.Error?.Code ?? "");
        Assert.NotEmpty(createResponse.BeneficiaryId);
        Assert.NotEmpty(createResponse.BeneficiaryCode);
        Assert.Contains("BUS-", createResponse.BeneficiaryCode);

        var getRequest = new GetBeneficiaryRequest
        {
            BeneficiaryId = createResponse.BeneficiaryId
        };

        var getResponse = await gateway.GetBeneficiaryAsync(getRequest);

        Assert.NotNull(getResponse);
        Assert.Null(getResponse.Error);
        Assert.NotNull(getResponse.Beneficiary);
        Assert.Equal(createResponse.BeneficiaryId, getResponse.Beneficiary.BeneficiaryId);
    }

    [Fact]
    public async Task GetNonExistentBeneficiary_ShouldReturnError()
    {
        using var scope = _serviceProvider.CreateScope();
        var gateway = scope.ServiceProvider.GetRequiredService<IBeneficiaryDataGateway>();

        var request = new GetBeneficiaryRequest
        {
            BeneficiaryId = Guid.NewGuid().ToString()
        };

        var response = await gateway.GetBeneficiaryAsync(request);

        Assert.NotNull(response);
        Assert.NotNull(response.Error);
        Assert.Equal("NOT_FOUND", response.Error.Code);
    }
}
