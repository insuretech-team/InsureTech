using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Infrastructure.Persistence;
using InsuranceEngine.Infrastructure.Repositories;
using InsuranceEngine.Infrastructure.Messaging;

namespace InsuranceEngine.Infrastructure;

public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(this IServiceCollection services, IConfiguration configuration)
    {
        services.AddDbContext<InsuranceDbContext>(options =>
            options.UseNpgsql(configuration.GetConnectionString("DefaultConnection"))
                   .UseSnakeCaseNamingConvention());

        services.AddScoped<IProductRepository, ProductRepository>();
        services.AddScoped<IInsurerRepository, InsurerRepository>();
        
        var kafkaSection = configuration.GetSection("Kafka");
        services.Configure<InsuranceKafkaOptions>(options => 
        {
            options.BootstrapServers = kafkaSection["BootstrapServers"] ?? string.Empty;
            options.GroupId = kafkaSection["GroupId"] ?? "insurance-engine-group";
        });
        services.AddSingleton<IEventBus, KafkaEventBus>();

        return services;
    }
}
