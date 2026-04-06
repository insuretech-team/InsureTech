using Microsoft.Extensions.DependencyInjection;
using MediatR;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Quoting;
using InsuranceEngine.Quoting.Infrastructure;

namespace InsuranceEngine.Quoting;

public static class QuotingDependencyInjection
{
    public static IServiceCollection AddQuotingModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(QuotingDbContext).Assembly);
        });

        services.AddDbContext<QuotingDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<IQuotingDataGateway, SqlQuotingDataGateway>();
        
        return services;
    }
}
