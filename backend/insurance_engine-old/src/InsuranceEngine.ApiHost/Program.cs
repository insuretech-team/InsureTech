using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Infrastructure;
using InsuranceEngine.Products.Infrastructure.Persistence;
using InsuranceEngine.Products.Domain.Services;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Infrastructure;
using InsuranceEngine.Policy.Infrastructure.Persistence;
using InsuranceEngine.Policy.Domain.Services;

using InsuranceEngine.SharedKernel.Interfaces;
using InsuranceEngine.SharedKernel.Messaging;
using InsuranceEngine.SharedKernel.Services;
using InsuranceEngine.ApiHost.Persistence;
using InsuranceEngine.SharedKernel.Behaviors;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Serilog;
using InsuranceEngine.Claims;
using InsuranceEngine.Underwriting;
using InsuranceEngine.Claims.GrpcServices;
using InsuranceEngine.Underwriting.GrpcServices;
using InsuranceEngine.Fraud;
using InsuranceEngine.Beneficiary;
using InsuranceEngine.Beneficiary.GrpcServices;
using InsuranceEngine.Partners;
using InsuranceEngine.Commission;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using OpenTelemetry.Metrics;

var builder = WebApplication.CreateBuilder(args);

// Configure Serilog
builder.Host.UseSerilog((context, configuration) =>
    configuration.ReadFrom.Configuration(context.Configuration).WriteTo.Console());

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(options =>
{
    options.CustomSchemaIds(type => type.FullName);
});

// Database Contexts
builder.Services.AddDbContext<ProductsDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")).UseSnakeCaseNamingConvention());

builder.Services.AddDbContext<PolicyDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")).UseSnakeCaseNamingConvention());

// HealthChecks
builder.Services.AddHealthChecks()
    .AddNpgSql(builder.Configuration.GetConnectionString("DefaultConnection") ?? string.Empty);

// Modular Slices Registration
builder.Services.AddClaimsModule(builder.Configuration);
builder.Services.AddUnderwritingModule(builder.Configuration);
builder.Services.AddFraudModule(builder.Configuration);
builder.Services.AddBeneficiaryModule(builder.Configuration);
builder.Services.AddPartnersModule(builder.Configuration);
builder.Services.AddCommissionModule(builder.Configuration);

builder.Services.AddMemoryCache();
builder.Services.AddScoped<IProductRepository, ProductRepository>();
builder.Services.AddScoped<IPolicyRepository, PolicyRepository>();
builder.Services.AddScoped<IEndorsementRepository, EndorsementRepository>();
// IBeneficiaryRepository moved to AddBeneficiaryModule
builder.Services.AddScoped<PolicyNumberGenerator>();
builder.Services.AddScoped<EndorsementNumberGenerator>();
builder.Services.AddSingleton<PricingEngine>();
builder.Services.AddScoped<PolicyDuplicateDetector>();
builder.Services.AddSingleton<IEncryptionService, AesEncryptionService>();


// Messaging
builder.Services.Configure<InsuranceKafkaOptions>(builder.Configuration.GetSection("Kafka"));
builder.Services.AddSingleton<IEventBus, KafkaEventBus>();

// MediatR — register from all modular slices
builder.Services.AddMediatR(cfg =>
{
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Products.Application.DTOs.ProductDto).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Policy.Application.DTOs.PolicyDto).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Claims.ClaimsModule).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Underwriting.UnderwritingModule).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Fraud.DependencyInjection).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Beneficiary.DependencyInjection).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Partners.DependencyInjection).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Commission.DependencyInjection).Assembly);

    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(ValidationBehavior<,>));
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(LoggingBehavior<,>));
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(TransactionBehavior<,>));
});



// gRPC
builder.Services.AddGrpc();
builder.Services.AddGrpcReflection();

// OpenTelemetry
builder.Services.AddOpenTelemetry()
    .WithTracing(tracing => tracing
        .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService("InsuranceEngine.ApiHost"))
        .AddAspNetCoreInstrumentation()
        .AddHttpClientInstrumentation()
        .AddEntityFrameworkCoreInstrumentation()
        .AddConsoleExporter())
    .WithMetrics(metrics => metrics
        .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService("InsuranceEngine.ApiHost"))
        .AddAspNetCoreInstrumentation()
        .AddHttpClientInstrumentation()
        .AddConsoleExporter());

var app = builder.Build();

// Configure the HTTP request pipeline.
app.UseSwagger();
app.UseSwaggerUI();

if (app.Environment.IsDevelopment())
{
}

app.UseAuthorization();
app.MapControllers();
app.MapHealthChecks("/health");

if (app.Environment.IsDevelopment())
{
    app.MapGrpcReflectionService();
}

// Seed Database
await DbInitializer.Initialize(app.Services);

app.MapGrpcService<InsuranceEngine.ApiHost.GrpcServices.InsuranceGrpcService>();
app.MapGrpcService<ClaimsGrpcService>();
app.MapGrpcService<UnderwritingGrpcService>();
app.MapGrpcService<InsuranceEngine.Beneficiary.GrpcServices.BeneficiaryGrpcService>();

app.Run();
