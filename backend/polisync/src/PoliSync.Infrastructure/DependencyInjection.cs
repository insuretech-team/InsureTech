using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using PoliSync.Infrastructure.Auth;
using PoliSync.Infrastructure.Persistence;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.Persistence;

namespace PoliSync.Infrastructure;

/// <summary>
/// Infrastructure DI registration. Called from ApiHost Program.cs.
/// </summary>
public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(this IServiceCollection services, IConfiguration config)
    {
        // ── EF Core: PostgreSQL with snake_case naming ──────────────────
        services.AddDbContext<PoliSyncDbContext>(opts =>
            opts.UseNpgsql(
                config.GetConnectionString("InsuranceDb"),
                npgsql => npgsql.CommandTimeout(30)
            ).UseSnakeCaseNamingConvention()
        );

        // ── Redis cache (fallback to in-memory for local dev) ───────────
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
            services.AddDistributedMemoryCache(); // fallback for local dev
        }

        // ── Persistence ─────────────────────────────────────────────────
        services.AddScoped<IUnitOfWork, UnitOfWork>();

        // ── Identity: CurrentUser is scoped — one instance per gRPC request ─
        services.AddScoped<CurrentUser>();
        services.AddScoped<ICurrentUser>(sp => sp.GetRequiredService<CurrentUser>());

        return services;
    }
}
