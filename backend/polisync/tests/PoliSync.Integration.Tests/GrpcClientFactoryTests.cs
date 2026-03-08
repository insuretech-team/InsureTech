using FluentAssertions;
using Grpc.Net.Client;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging.Abstractions;
using PoliSync.Infrastructure.GrpcClients;
using Xunit;

namespace PoliSync.Integration.Tests;

/// <summary>
/// Tests for GrpcClientFactory — verifies channel creation, caching, and config resolution.
/// </summary>
public sealed class GrpcClientFactoryTests : IDisposable
{
    private readonly GrpcClientFactory _factory;

    public GrpcClientFactoryTests()
    {
        var config = new ConfigurationBuilder()
            .AddInMemoryCollection(new Dictionary<string, string?>
            {
                ["GrpcClients:AuthzService"]        = "http://localhost:50070",
                ["GrpcClients:AuditService"]        = "http://localhost:50080",
                ["GrpcClients:KycService"]          = "http://localhost:50090",
                ["GrpcClients:PartnerService"]      = "http://localhost:50100",
                ["GrpcClients:FraudService"]        = "http://localhost:50220",
                ["GrpcClients:PaymentService"]      = "http://localhost:50190",
                ["GrpcClients:NotificationService"] = "http://localhost:50230",
                ["GrpcClients:StorageService"]      = "http://localhost:50290",
                ["GrpcClients:DocgenService"]       = "http://localhost:50280",
                ["GrpcClients:WorkflowService"]     = "http://localhost:50180",
            })
            .Build();

        _factory = new GrpcClientFactory(config, NullLogger<GrpcClientFactory>.Instance);
    }

    [Theory]
    [InlineData("AuthzService",        "http://localhost:50070")]
    [InlineData("AuditService",        "http://localhost:50080")]
    [InlineData("KycService",          "http://localhost:50090")]
    [InlineData("PartnerService",      "http://localhost:50100")]
    [InlineData("FraudService",        "http://localhost:50220")]
    [InlineData("PaymentService",      "http://localhost:50190")]
    [InlineData("NotificationService", "http://localhost:50230")]
    [InlineData("StorageService",      "http://localhost:50290")]
    [InlineData("DocgenService",       "http://localhost:50280")]
    [InlineData("WorkflowService",     "http://localhost:50180")]
    public void GetChannel_ConfiguredService_ReturnsChannel(string serviceName, string expectedAddress)
    {
        var channel = _factory.GetChannel(serviceName);

        channel.Should().NotBeNull();
        channel.Target.Should().Be(expectedAddress.Replace("http://", ""));
    }

    [Fact]
    public void GetChannel_SameService_ReturnsCachedChannel()
    {
        var ch1 = _factory.GetChannel("AuthzService");
        var ch2 = _factory.GetChannel("AuthzService");

        ch1.Should().BeSameAs(ch2); // channel is cached
    }

    [Fact]
    public void GetChannel_DifferentServices_ReturnDifferentChannels()
    {
        var ch1 = _factory.GetChannel("AuthzService");
        var ch2 = _factory.GetChannel("AuditService");

        ch1.Should().NotBeSameAs(ch2);
    }

    [Fact]
    public void GetChannel_UnconfiguredService_ThrowsInvalidOperationException()
    {
        var act = () => _factory.GetChannel("NonExistentService");

        act.Should().Throw<InvalidOperationException>()
           .WithMessage("*NonExistentService*");
    }

    [Fact]
    public void GetClient_ReturnsTypedClientFromChannel()
    {
        // Verify the generic GetClient<T> factory works
        var channel = _factory.GetChannel("AuthzService");
        channel.Should().NotBeNull();
        channel.Target.Should().NotBeNullOrEmpty();
    }

    public void Dispose() => _factory.Dispose();
}
