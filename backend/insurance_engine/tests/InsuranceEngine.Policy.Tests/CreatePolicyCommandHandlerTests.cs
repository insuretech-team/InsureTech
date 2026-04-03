using Moq;
using Xunit;
using FluentAssertions;
using Microsoft.Extensions.Logging;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Policy.Application.Commands;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Infrastructure.DataGateways;
using InsuranceEngine.Products.Domain.Entities;
using InsuranceEngine.Infrastructure.Persistence;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Policy.Tests;

public class CreatePolicyCommandHandlerTests : IDisposable
{
    private readonly Mock<IRepository<PolicyEntity>> _policyRepoMock;
    private readonly Mock<IRepository<ProductEntity>> _productRepoMock;
    private readonly Mock<IRepository<PolicyNomineeEntity>> _nomineeRepoMock;
    private readonly Mock<IRepository<PolicyRiderEntity>> _riderRepoMock;
    private readonly Mock<ILogger<CreatePolicyCommandHandler>> _loggerMock;
    private readonly Mock<IPdfGenerator> _pdfGeneratorMock;
    private readonly Mock<IKafkaPublisher> _kafkaPublisherMock;
    private readonly Mock<ISequenceDataGateway> _sequenceGatewayMock;
    private readonly InsuranceDbContext _dbContext;
    private readonly SqliteConnection _connection;
    private readonly CreatePolicyCommandHandler _handler;

    public CreatePolicyCommandHandlerTests()
    {
        _policyRepoMock = new Mock<IRepository<PolicyEntity>>();
        _productRepoMock = new Mock<IRepository<ProductEntity>>();
        _nomineeRepoMock = new Mock<IRepository<PolicyNomineeEntity>>();
        _riderRepoMock = new Mock<IRepository<PolicyRiderEntity>>();
        _loggerMock = new Mock<ILogger<CreatePolicyCommandHandler>>();
        _pdfGeneratorMock = new Mock<IPdfGenerator>();
        _kafkaPublisherMock = new Mock<IKafkaPublisher>();
        _sequenceGatewayMock = new Mock<ISequenceDataGateway>();

        // Use SQLite in-memory for professional relational testing
        _connection = new SqliteConnection("Filename=:memory:");
        _connection.Open();

        var options = new DbContextOptionsBuilder<InsuranceDbContext>()
            .UseSqlite(_connection)
            .Options;

        _dbContext = new InsuranceDbContext(options);
        
        // Essential for SQLite in-memory to create the schema
        _dbContext.Database.EnsureCreated();

        _handler = new CreatePolicyCommandHandler(
            _policyRepoMock.Object,
            _productRepoMock.Object,
            _nomineeRepoMock.Object,
            _riderRepoMock.Object,
            _sequenceGatewayMock.Object,
            _loggerMock.Object,
            _pdfGeneratorMock.Object,
            _kafkaPublisherMock.Object);
    }

    [Fact]
    public async Task Handle_ShouldSuccessfullyCreatePolicy_WithPartnerAndAgentAttribution()
    {
        // Arrange
        var productId = Guid.NewGuid();
        var customerId = Guid.NewGuid();
        var partnerId = Guid.NewGuid();
        var agentId = Guid.NewGuid();
        
        var product = new ProductEntity 
        { 
            ProductId = productId, 
            Status = "ACTIVE",
            ProductName = "Professional Safety Plan" 
        };

        _productRepoMock.Setup(r => r.GetByIdAsync(productId, It.IsAny<CancellationToken>()))
            .ReturnsAsync(product);

        var command = new CreatePolicyCommand(
            productId.ToString(),
            customerId.ToString(),
            partnerId.ToString(),
            agentId.ToString(),
            null,
            1200.50m, // Premium
            100000m,  // Sum Insured
            12,       // Tenure
            DateTime.UtcNow,
            null,
            null);

        // Act
        var response = await _handler.Handle(command, CancellationToken.None);

        // Assert
        response.Error.Should().BeNull();
        response.PolicyNumber.Should().NotBeNullOrEmpty();

        // Verify Repository Persistence (Attribution Check)
        _policyRepoMock.Verify(r => r.AddAsync(It.Is<PolicyEntity>(p => 
            p.ProductId == productId &&
            p.CustomerId == customerId &&
            p.PartnerId == partnerId &&
            p.AgentId == agentId &&
            p.PremiumAmount == 120050), // Correct Paisa Conversion
            It.IsAny<CancellationToken>()), Times.Once);

        // Verify Kafka Event Publication (Unification Check)
        _kafkaPublisherMock.Verify(k => k.PublishAsync(
            "insurance.policy.created", 
            It.Is<PolicyIssuedEvent>(e => 
                e.PolicyId != Guid.Empty &&
                e.PartnerId == partnerId &&
                e.AgentId == agentId)), Times.Once);

        // Verify PDF Generation
        _pdfGeneratorMock.Verify(p => p.GeneratePolicyDocumentAsync(
            It.IsAny<string>(), "N/A", "N/A", 1200.50m), Times.Once);
    }

    [Fact]
    public async Task Handle_ShouldReturnError_WhenProductNotFound()
    {
        // Arrange
        var productId = Guid.NewGuid();
        _productRepoMock.Setup(r => r.GetByIdAsync(productId, It.IsAny<CancellationToken>()))
            .ReturnsAsync((ProductEntity?)null);

        var command = new CreatePolicyCommand(
            productId.ToString(),
            Guid.NewGuid().ToString(),
            null, null, null, 1000, 50000, 12, DateTime.UtcNow, null, null);

        // Act
        var response = await _handler.Handle(command, CancellationToken.None);

        // Assert
        response.Error.Should().NotBeNull();
        response.Error.Code.Should().Be("PRODUCT_NOT_FOUND");
    }

    public void Dispose()
    {
        _dbContext.Dispose();
        _connection.Close();
        _connection.Dispose();
    }
}
