using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using PoliSync.Infrastructure.Auth;
using PoliSync.Infrastructure.Cache;
using PoliSync.Infrastructure.GrpcClients;
using PoliSync.Infrastructure.Messaging;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Infrastructure.Pii;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.Messaging;
using PoliSync.SharedKernel.Persistence;

namespace PoliSync.Infrastructure;

/// <summary>
/// Infrastructure DI registration. Called from ApiHost Program.cs.
/// </summary>
public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(this IServiceCollection services, IConfiguration config)
    {
        // ── EF Core: insurance_schema ──────────────────────────────────────
        services.AddDbContext<PoliSyncDbContext>(opts =>
            opts.UseNpgsql(
                config.GetConnectionString("InsuranceDb"),
                npgsql => npgsql
                    .CommandTimeout(30)
            ).UseSnakeCaseNamingConvention()
        );

        // ── EF Core: commission_schema ─────────────────────────────────────
        services.AddDbContext<CommissionDbContext>(opts =>
            opts.UseNpgsql(
                config.GetConnectionString("CommissionDb"),
                npgsql => npgsql
                    .CommandTimeout(30)
            ).UseSnakeCaseNamingConvention()
        );

        // ── Redis cache ────────────────────────────────────────────────────
        var redisConn = config.GetConnectionString("Redis");
        if (!string.IsNullOrEmpty(redisConn))
        {
            services.AddStackExchangeRedisCache(opts =>
            {
                opts.Configuration = redisConn;
                opts.InstanceName = "polisync:";
            });
        }
        else
        {
            services.AddDistributedMemoryCache(); // fallback for local dev without Redis
        }

        // ── Messaging ──────────────────────────────────────────────────────
        services.AddSingleton<IDomainEventDispatcher, DomainEventDispatcher>();
        services.AddSingleton<IEventBus, KafkaEventBus>();

        // ── Persistence ────────────────────────────────────────────────────
        services.AddScoped<IUnitOfWork, UnitOfWork>();

        // ── Identity: CurrentUser is scoped — one instance per gRPC request ─
        // AuthInterceptor calls currentUser.Populate(context) on the SAME
        // scoped instance that command/query handlers receive via DI.
        services.AddScoped<CurrentUser>();
        services.AddScoped<ICurrentUser>(sp => sp.GetRequiredService<CurrentUser>());

        // ── PII encryption ─────────────────────────────────────────────────
        services.AddSingleton<IPiiEncryptor, AesGcmPiiEncryptor>();

        // ── Caches ─────────────────────────────────────────────────────────
        services.AddSingleton<RedisProductCache>();

        // ── gRPC client factory (singleton — channels are reused) ──────────
        services.AddSingleton<GrpcClientFactory>();

        // ── Typed gRPC clients for upstream Go services ────────────────────
        services.AddSingleton<AuthzGrpcClient>();
        services.AddSingleton<AuditGrpcClient>();
        services.AddSingleton<KycGrpcClient>();
        services.AddSingleton<PartnerGrpcClient>();
        services.AddSingleton<FraudGrpcClient>();
        services.AddSingleton<NotificationGrpcClient>();
        services.AddSingleton<PaymentGrpcClient>();
        services.AddSingleton<StorageGrpcClient>();
        services.AddSingleton<DocgenGrpcClient>();
        services.AddSingleton<WorkflowGrpcClient>();

        return services;
    }
}
