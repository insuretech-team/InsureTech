using Moq;
using Xunit;
using FluentAssertions;
using Microsoft.Extensions.Logging;
using InsuranceEngine.Policy.Application.Commands;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Infrastructure;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Policy.Tests;

public class ExcelDataIntegrationTests
{
    [Fact]
    public async Task CreatePolicy_FromOverseasMediclaimData()
    {
        var gatewayMock = new Mock<IPolicyDataGateway>();
        var loggerMock = new Mock<ILogger<CreatePolicyCommandHandler>>();
        
        var handler = new CreatePolicyCommandHandler(gatewayMock.Object, loggerMock.Object);
        
        var excelData = ExcelTestData.GetOverseasMediclaimPolicy();
        
        gatewayMock.Setup(g => g.CreatePolicyAsync(It.IsAny<CreatePolicyRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CreatePolicyResponse
            {
                PolicyId = Guid.NewGuid().ToString(),
                PolicyNumber = "POL-2026-0001",
                Message = "Success"
            });

        var command = new CreatePolicyCommand(
            excelData.ProductId,
            excelData.CustomerId,
            null, null, null,
            excelData.PremiumAmount,
            excelData.SumInsured,
            excelData.TenureMonths,
            excelData.StartDate,
            excelData.ProposerDetails,
            null);

        var result = await handler.Handle(command, CancellationToken.None);

        result.PolicyNumber.Should().NotBeEmpty();
        
        gatewayMock.Verify(g => g.CreatePolicyAsync(It.Is<CreatePolicyRequest>(r => 
            r.ProductId == excelData.ProductId &&
            r.CustomerId == excelData.CustomerId
        ), It.IsAny<CancellationToken>()), Times.Once);
    }

    [Fact]
    public async Task CreatePolicy_FromVehicleInsuranceData()
    {
        var gatewayMock = new Mock<IPolicyDataGateway>();
        var loggerMock = new Mock<ILogger<CreatePolicyCommandHandler>>();
        
        var handler = new CreatePolicyCommandHandler(gatewayMock.Object, loggerMock.Object);
        
        var excelData = ExcelTestData.GetVehicleInsurancePolicy();
        
        gatewayMock.Setup(g => g.CreatePolicyAsync(It.IsAny<CreatePolicyRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CreatePolicyResponse
            {
                PolicyId = Guid.NewGuid().ToString(),
                PolicyNumber = "POL-2026-0002",
                Message = "Success"
            });

        var command = new CreatePolicyCommand(
            excelData.ProductId,
            excelData.CustomerId,
            null, null, null,
            excelData.PremiumAmount,
            excelData.SumInsured,
            excelData.TenureMonths,
            excelData.StartDate,
            excelData.ProposerDetails,
            null);

        var result = await handler.Handle(command, CancellationToken.None);

        result.PolicyNumber.Should().NotBeEmpty();
    }

    [Fact]
    public async Task AddNominee_FromExcelBeneficiaryData()
    {
        var gatewayMock = new Mock<IPolicyDataGateway>();
        var loggerMock = new Mock<ILogger<AddNomineeCommandHandler>>();
        
        var handler = new AddNomineeCommandHandler(gatewayMock.Object, loggerMock.Object);
        
        var policyId = "test-policy-001";
        var excelNominee = ExcelTestData.GetBeneficiaryNominee();
        
        var existingPolicy = CreateTestPolicy(policyId);

        gatewayMock.Setup(g => g.GetPolicyAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new GetPolicyResponse { Policy = existingPolicy });

        gatewayMock.Setup(g => g.UpdatePolicyAsync(It.IsAny<string>(), It.IsAny<List<Nominee>?>(), It.IsAny<string?>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new UpdatePolicyResponse { Error = null, Message = "Success" });

        var command = new AddNomineeCommand(
            policyId,
            excelNominee.FullName,
            excelNominee.Relationship,
            excelNominee.SharePercentage,
            excelNominee.DateOfBirth,
            excelNominee.NidNumber,
            excelNominee.PhoneNumber,
            excelNominee.NomineeDobText);

        var result = await handler.Handle(command, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
    }

    [Theory]
    [InlineData(1239, 50000)]
    [InlineData(2499, 50000)]
    [InlineData(8131, 50000)]
    public async Task CreatePolicy_WithVariousPremiumRates(decimal premium, decimal sumInsured)
    {
        var gatewayMock = new Mock<IPolicyDataGateway>();
        var loggerMock = new Mock<ILogger<CreatePolicyCommandHandler>>();
        
        var handler = new CreatePolicyCommandHandler(gatewayMock.Object, loggerMock.Object);
        
        gatewayMock.Setup(g => g.CreatePolicyAsync(It.IsAny<CreatePolicyRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CreatePolicyResponse
            {
                PolicyId = Guid.NewGuid().ToString(),
                PolicyNumber = $"POL-2026-{premium}",
                Message = "Success"
            });

        var command = new CreatePolicyCommand(
            "OMC-PRODUCT-001",
            "CUST-PRAGATI-001",
            null, null, null,
            premium,
            sumInsured,
            1,
            DateTime.UtcNow,
            null,
            null);

        var result = await handler.Handle(command, CancellationToken.None);

        result.PolicyNumber.Should().NotBeEmpty();
    }

    private static Insuretech.Policy.Entity.V1.Policy CreateTestPolicy(string policyId)
    {
        return new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = policyId,
            PolicyNumber = "POL-2026-0001",
            ProductId = "PROD-001",
            CustomerId = "CUST-001",
            Status = Insuretech.Policy.Entity.V1.PolicyStatus.Active,
            PremiumAmount = new Insuretech.Common.V1.Money { Amount = 1239, Currency = "BDT" },
            SumInsured = new Insuretech.Common.V1.Money { Amount = 50000, Currency = "BDT" },
            TenureMonths = 1,
            StartDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            EndDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow.AddMonths(1))
        };
    }
}

public static class ExcelTestData
{
    public static ExcelPolicyData GetOverseasMediclaimPolicy()
    {
        return new ExcelPolicyData
        {
            ProductId = "OMC-PRODUCT-001",
            CustomerId = "CUST-PRAGATI-001",
            PremiumAmount = 1239,
            SumInsured = 50000,
            TenureMonths = 1,
            StartDate = new DateTime(2026, 4, 15),
            ProposerDetails = "Name: Md. Zubayed Ur Rahman, Address: N.B Tower Level-5, 40/7 North Avenue, Gulshan-2, Dhaka-1212, Mobile: 01985700011, Email: Zubayer@ymail.com, Occupation: Service at Medland Bank Plc, Passport: GA-18-6525, Plan: Business & Holiday (14-180 days)"
        };
    }

    public static ExcelPolicyData GetVehicleInsurancePolicy()
    {
        return new ExcelPolicyData
        {
            ProductId = "VEHICLE-PRODUCT-001",
            CustomerId = "CUST-PRAGATI-002",
            PremiumAmount = 5000,
            SumInsured = 500000,
            TenureMonths = 12,
            StartDate = new DateTime(2026, 4, 15),
            ProposerDetails = "Vehicle: Dhaka Metro GA-18-6525, Chassis: NKE165-7216292, Engine: G4NAEM48921, Make: Hyundai, Model: Tucson, Year: 2024"
        };
    }

    public static NomineeData GetBeneficiaryNominee()
    {
        return new NomineeData
        {
            FullName = "Fatema Begum",
            Relationship = "Wife",
            SharePercentage = 50,
            DateOfBirth = new DateTime(1985, 6, 15),
            NidNumber = "198515678900001",
            PhoneNumber = "+88017111234567",
            NomineeDobText = "15-06-1985"
        };
    }
}

public class ExcelPolicyData
{
    public string ProductId { get; set; } = "";
    public string CustomerId { get; set; } = "";
    public decimal PremiumAmount { get; set; }
    public decimal SumInsured { get; set; }
    public int TenureMonths { get; set; }
    public DateTime StartDate { get; set; }
    public string? ProposerDetails { get; set; }
}

public class NomineeData
{
    public string FullName { get; set; } = "";
    public string Relationship { get; set; } = "";
    public int SharePercentage { get; set; }
    public DateTime? DateOfBirth { get; set; }
    public string? NidNumber { get; set; }
    public string? PhoneNumber { get; set; }
    public string? NomineeDobText { get; set; }
}
