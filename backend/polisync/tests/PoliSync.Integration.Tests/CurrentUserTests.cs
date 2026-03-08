using FluentAssertions;
using Grpc.Core;
using PoliSync.Infrastructure.Auth;
using Xunit;

namespace PoliSync.Integration.Tests;

/// <summary>
/// Tests for CurrentUser population from gRPC metadata.
/// Verifies that identity headers injected by the Go gateway are correctly
/// parsed into ICurrentUser by the AuthInterceptor.
/// </summary>
public sealed class CurrentUserTests
{
    private static ServerCallContext BuildContext(Dictionary<string, string> headers)
    {
        var meta = new Metadata();
        foreach (var (k, v) in headers)
            meta.Add(k.ToLowerInvariant(), v);
        return TestServerCallContext.Create(meta);
    }

    [Fact]
    public void Populate_WithAllHeaders_SetsAllProperties()
    {
        var ctx = BuildContext(new()
        {
            ["x-user-id"]    = "11111111-1111-1111-1111-111111111111",
            ["x-tenant-id"]  = "22222222-2222-2222-2222-222222222222",
            ["x-partner-id"] = "33333333-3333-3333-3333-333333333333",
            ["x-token-id"]   = "tok-abc-123",
            ["x-user-type"]  = "agent",
            ["x-portal"]     = "b2b",
            ["x-roles"]      = "product:read,product:write,claim:read",
        });

        var user = new CurrentUser();
        user.Populate(ctx);

        user.UserId.Should().Be(Guid.Parse("11111111-1111-1111-1111-111111111111"));
        user.TenantId.Should().Be(Guid.Parse("22222222-2222-2222-2222-222222222222"));
        user.PartnerId.Should().Be(Guid.Parse("33333333-3333-3333-3333-333333333333"));
        user.TokenId.Should().Be("tok-abc-123");
        user.UserType.Should().Be("agent");
        user.Portal.Should().Be("b2b");
        user.Roles.Should().BeEquivalentTo(["product:read", "product:write", "claim:read"]);
        user.IsAuthenticated.Should().BeTrue();
        user.IsAgent.Should().BeTrue();
        user.IsSystemUser.Should().BeFalse();
        user.IsPartnerUser.Should().BeTrue();
    }

    [Fact]
    public void Populate_WithoutPartnerId_IsPartnerUserFalse()
    {
        var ctx = BuildContext(new()
        {
            ["x-user-id"]   = "11111111-1111-1111-1111-111111111111",
            ["x-tenant-id"] = "22222222-2222-2222-2222-222222222222",
            ["x-user-type"] = "b2c",
        });

        var user = new CurrentUser();
        user.Populate(ctx);

        user.PartnerId.Should().BeNull();
        user.IsPartnerUser.Should().BeFalse();
        user.IsAgent.Should().BeFalse();
    }

    [Fact]
    public void Populate_WithSystemUserType_IsSystemUserTrue()
    {
        var ctx = BuildContext(new()
        {
            ["x-user-id"]   = "11111111-1111-1111-1111-111111111111",
            ["x-tenant-id"] = "22222222-2222-2222-2222-222222222222",
            ["x-user-type"] = "system",
        });

        var user = new CurrentUser();
        user.Populate(ctx);

        user.IsSystemUser.Should().BeTrue();
        user.IsAgent.Should().BeFalse();
    }

    [Fact]
    public void Populate_WithMissingUserId_IsAuthenticatedFalse()
    {
        var ctx = BuildContext(new()
        {
            ["x-tenant-id"] = "22222222-2222-2222-2222-222222222222",
        });

        var user = new CurrentUser();
        user.Populate(ctx);

        user.IsAuthenticated.Should().BeFalse();
        user.UserId.Should().Be(Guid.Empty);
    }

    [Fact]
    public void Populate_WithInvalidGuid_UsesEmptyGuid()
    {
        var ctx = BuildContext(new()
        {
            ["x-user-id"]   = "not-a-valid-guid",
            ["x-tenant-id"] = "also-invalid",
        });

        var user = new CurrentUser();
        user.Populate(ctx);

        user.UserId.Should().Be(Guid.Empty);
        user.TenantId.Should().Be(Guid.Empty);
        user.IsAuthenticated.Should().BeFalse();
    }

    [Fact]
    public void Populate_WithEmptyRoles_ReturnsEmptyList()
    {
        var ctx = BuildContext(new()
        {
            ["x-user-id"]   = "11111111-1111-1111-1111-111111111111",
            ["x-tenant-id"] = "22222222-2222-2222-2222-222222222222",
            ["x-roles"]     = "",
        });

        var user = new CurrentUser();
        user.Populate(ctx);

        user.Roles.Should().BeEmpty();
    }

    [Fact]
    public void Populate_CalledTwice_OverwritesPreviousValues()
    {
        var user = new CurrentUser();

        user.Populate(BuildContext(new()
        {
            ["x-user-id"]   = "11111111-1111-1111-1111-111111111111",
            ["x-tenant-id"] = "22222222-2222-2222-2222-222222222222",
            ["x-user-type"] = "agent",
        }));
        user.UserId.Should().Be(Guid.Parse("11111111-1111-1111-1111-111111111111"));

        user.Populate(BuildContext(new()
        {
            ["x-user-id"]   = "99999999-9999-9999-9999-999999999999",
            ["x-tenant-id"] = "88888888-8888-8888-8888-888888888888",
            ["x-user-type"] = "system",
        }));
        user.UserId.Should().Be(Guid.Parse("99999999-9999-9999-9999-999999999999"));
        user.IsSystemUser.Should().BeTrue();
    }
}

/// <summary>
/// Test helper: creates a minimal ServerCallContext from metadata.
/// </summary>
internal static class TestServerCallContext
{
    public static ServerCallContext Create(Metadata headers)
        => new TestCallContext(headers);

    private sealed class TestCallContext : ServerCallContext
    {
        private readonly Metadata _headers;
        public TestCallContext(Metadata headers) => _headers = headers;

        protected override string MethodCore => "/test/Method";
        protected override string HostCore => "localhost";
        protected override string PeerCore => "127.0.0.1";
        protected override DateTime DeadlineCore => DateTime.MaxValue;
        protected override Metadata RequestHeadersCore => _headers;
        protected override CancellationToken CancellationTokenCore => CancellationToken.None;
        protected override Metadata ResponseTrailersCore => new();
        protected override Status StatusCore { get; set; }
        protected override WriteOptions? WriteOptionsCore { get; set; }
        protected override AuthContext AuthContextCore => new("", []);

        protected override ContextPropagationToken CreatePropagationTokenCore(ContextPropagationOptions? options)
            => throw new NotImplementedException();
        protected override Task WriteResponseHeadersAsyncCore(Metadata responseHeaders)
            => Task.CompletedTask;
    }
}
