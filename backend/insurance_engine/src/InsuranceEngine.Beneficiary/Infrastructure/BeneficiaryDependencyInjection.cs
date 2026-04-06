using Microsoft.Extensions.DependencyInjection;
using MediatR;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Beneficiary.Application;
using InsuranceEngine.Beneficiary.Infrastructure;

namespace InsuranceEngine.Beneficiary;

public static class BeneficiaryDependencyInjection
{
    public static IServiceCollection AddBeneficiaryModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<BeneficiaryDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<IBeneficiaryDataGateway, SqlBeneficiaryDataGateway>();
        
        return services;
    }
}
