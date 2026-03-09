using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Infrastructure;
using InsuranceEngine.Products.Infrastructure.Persistence;
using InsuranceEngine.SharedKernel.Interfaces;
using InsuranceEngine.SharedKernel.Messaging;
using InsuranceEngine.ApiHost.Persistence;
using InsuranceEngine.SharedKernel.Behaviors;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Serilog;

var builder = WebApplication.CreateBuilder(args);

// Configure Serilog
builder.Host.UseSerilog((context, configuration) => 
    configuration.ReadFrom.Configuration(context.Configuration).WriteTo.Console());

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddScoped<ITenantService, DefaultTenantService>();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Database Contexts
builder.Services.AddDbContext<ProductsDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")).UseSnakeCaseNamingConvention());

builder.Services.AddDbContext<InsurersDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")).UseSnakeCaseNamingConvention());

// HealthChecks
builder.Services.AddHealthChecks()
    .AddNpgSql(builder.Configuration.GetConnectionString("DefaultConnection") ?? string.Empty);

// Repositories
builder.Services.AddScoped<IProductRepository, ProductRepository>();

// Messaging
builder.Services.Configure<InsuranceKafkaOptions>(builder.Configuration.GetSection("Kafka"));
builder.Services.AddSingleton<IEventBus, KafkaEventBus>();

// MediatR
builder.Services.AddMediatR(cfg => {
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Products.Application.DTOs.ProductDto).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Products.Application.Features.Queries.ListInsurers.ListInsurersQuery).Assembly);
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(ValidationBehavior<,>));
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(LoggingBehavior<,>));
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(TransactionBehavior<,>));
});

// gRPC
builder.Services.AddGrpc();
builder.Services.AddGrpcReflection();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
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

app.MapGrpcService<InsuranceEngine.Products.GrpcServices.InsuranceGrpcService>();

app.Run();



