using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using PoliSync.ApiHost.BackgroundServices;
using PoliSync.Infrastructure.Auth;
using PoliSync.Infrastructure.Messaging;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Infrastructure.Pii;
using PoliSync.Claims.Infrastructure;
using PoliSync.Endorsement.Infrastructure;
using PoliSync.Orders.Infrastructure;
using PoliSync.Renewal.Infrastructure;
using PoliSync.Refund.Infrastructure;
using PoliSync.Underwriting.Domain;
using PoliSync.Underwriting.Infrastructure;
using PoliSync.Products;
using PoliSync.Policy.Infrastructure;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.Messaging;
using PoliSync.SharedKernel.Persistence;
using PoliSync.SharedKernel.Pii;
using Serilog;
using System.Security.Cryptography;

var builder = WebApplication.CreateBuilder(args);

// Configure Serilog
Log.Logger = new LoggerConfiguration()
    .ReadFrom.Configuration(builder.Configuration)
    .CreateLogger();

builder.Host.UseSerilog();

// Add services
builder.Services.AddGrpc(options =>
{
    options.EnableDetailedErrors = true;
    options.MaxReceiveMessageSize = 16 * 1024 * 1024; // 16MB
    options.MaxSendMessageSize = 16 * 1024 * 1024;
});

builder.Services.AddGrpcReflection();

// Database
var insuranceConnectionString = builder.Configuration.GetConnectionString("InsuranceDb")!
    .Replace("${DB_PASSWORD}", Environment.GetEnvironmentVariable("DB_PASSWORD") ?? "");

builder.Services.AddDbContext<PoliSyncDbContext>(options =>
    options.UseNpgsql(insuranceConnectionString, npgsqlOptions =>
    {
        npgsqlOptions.MigrationsHistoryTable("__EFMigrationsHistory", "insurance_schema");
        npgsqlOptions.EnableRetryOnFailure(3);
    }));

// Unit of Work
builder.Services.AddScoped<IUnitOfWork, UnitOfWork>();

// Repositories
builder.Services.AddScoped(typeof(PoliSync.SharedKernel.Persistence.IRepository<>), 
    typeof(Repository<>));

// Redis Cache
var redisConnectionString = builder.Configuration.GetConnectionString("Redis")!
    .Replace("${REDIS_PASSWORD}", Environment.GetEnvironmentVariable("REDIS_PASSWORD") ?? "");

builder.Services.AddStackExchangeRedisCache(options =>
{
    options.Configuration = redisConnectionString;
});

// Kafka
builder.Services.Configure<KafkaOptions>(builder.Configuration.GetSection("Kafka"));
builder.Services.AddSingleton<IEventBus, KafkaEventBus>();

// PII Encryption
builder.Services.Configure<PiiEncryptionOptions>(options =>
{
    var keyPath = builder.Configuration["Pii:EncryptionKeyPath"];
    if (!string.IsNullOrEmpty(keyPath) && File.Exists(keyPath))
    {
        options.EncryptionKey = File.ReadAllText(keyPath).Trim();
    }
    else
    {
        options.EncryptionKey = Environment.GetEnvironmentVariable("PII_ENCRYPTION_KEY") ?? "";
    }
});
builder.Services.AddSingleton<IPiiEncryptor, AesGcmPiiEncryptor>();

// Current User
builder.Services.AddHttpContextAccessor();
builder.Services.AddScoped<ICurrentUser, CurrentUser>();
builder.Services.AddSingleton<PoliSync.Infrastructure.Clients.InsuranceServiceClient>();
builder.Services.AddSingleton<PoliSync.Infrastructure.Clients.OrderServiceGrpcClient>();
builder.Services.AddSingleton<PoliSync.Infrastructure.Clients.PaymentServiceGrpcClient>();
builder.Services.AddScoped<IPolicyDataGateway, GoPolicyDataGateway>();
builder.Services.AddScoped<IQuotationDataGateway, GoQuotationDataGateway>();
builder.Services.AddScoped<IClaimDataGateway, GoClaimDataGateway>();
builder.Services.AddScoped<IEndorsementDataGateway, GoEndorsementDataGateway>();
builder.Services.AddScoped<IRenewalDataGateway, GoRenewalDataGateway>();
builder.Services.AddScoped<IOrderDataGateway, GoOrderDataGateway>();
builder.Services.AddScoped<IRefundPaymentGateway, GoRefundPaymentGateway>();
builder.Services.AddScoped<IUnderwritingDataGateway, GoUnderwritingDataGateway>();
builder.Services.AddSingleton<IUnderwritingRiskScorer, UnderwritingRiskScorer>();
builder.Services.AddHostedService<UnderwritingQuotationSubmittedConsumer>();
builder.Services.AddHostedService<OrderPaymentConfirmedConsumer>();

// Domain modules
builder.Services.AddProductsModule();

// MediatR
builder.Services.AddMediatR(cfg =>
{
    cfg.RegisterServicesFromAssembly(typeof(Program).Assembly);
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Products.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Quotes.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Orders.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Policy.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Claims.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Commission.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Endorsement.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Renewal.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Underwriting.AssemblyMarker>();
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Refund.AssemblyMarker>();
});

// JWT Authentication
var jwtPublicKeyPath = builder.Configuration["Jwt:PublicKeyPath"]!;
RSA? rsa = null;

if (File.Exists(jwtPublicKeyPath))
{
    var publicKeyPem = File.ReadAllText(jwtPublicKeyPath);
    rsa = RSA.Create();
    rsa.ImportFromPem(publicKeyPem);
}

builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true,
            ValidIssuer = builder.Configuration["Jwt:Issuer"],
            ValidAudience = builder.Configuration["Jwt:Audience"],
            IssuerSigningKey = rsa != null ? new RsaSecurityKey(rsa) : null,
            ClockSkew = TimeSpan.FromMinutes(5)
        };
    });

builder.Services.AddAuthorization();

// Health Checks
builder.Services.AddHealthChecks()
    .AddNpgSql(insuranceConnectionString, name: "postgres")
    .AddRedis(redisConnectionString, name: "redis")
    .AddKafka(new Confluent.Kafka.ProducerConfig 
    { 
        BootstrapServers = builder.Configuration["Kafka:BootstrapServers"] 
    }, name: "kafka");

// CORS
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyMethod()
              .AllowAnyHeader()
              .WithExposedHeaders("Grpc-Status", "Grpc-Message", "Grpc-Encoding", "Grpc-Accept-Encoding");
    });
});

var app = builder.Build();

// Configure middleware
app.UseSerilogRequestLogging();

app.UseCors();

app.UseAuthentication();
app.UseAuthorization();

// Map gRPC services
app.MapGrpcService<PoliSync.Products.GrpcServices.ProductGrpcService>();
app.MapGrpcService<PoliSync.Quotes.GrpcServices.QuotesGrpcService>();
app.MapGrpcService<PoliSync.Orders.GrpcServices.OrderGrpcService>();
app.MapGrpcService<PoliSync.Policy.GrpcServices.PolicyGrpcService>();
app.MapGrpcService<PoliSync.Claims.GrpcServices.ClaimGrpcService>();
app.MapGrpcService<PoliSync.Commission.GrpcServices.CommissionGrpcService>();
app.MapGrpcService<PoliSync.Underwriting.GrpcServices.UnderwritingGrpcService>();
app.MapGrpcService<PoliSync.Endorsement.GrpcServices.EndorsementGrpcService>();
app.MapGrpcService<PoliSync.Renewal.GrpcServices.RenewalGrpcService>();
app.MapGrpcService<PoliSync.Refund.GrpcServices.RefundGrpcService>();

app.MapGrpcReflectionService();

// Health checks
app.MapHealthChecks("/health");

// Root endpoint
app.MapGet("/", () => new
{
    service = "PoliSync",
    version = "1.0.0",
    description = "C# .NET 8 Insurance Commerce & Policy Engine",
    status = "running"
});

try
{
    Log.Information("Starting PoliSync ApiHost");
    app.Run();
}
catch (Exception ex)
{
    Log.Fatal(ex, "PoliSync ApiHost terminated unexpectedly");
}
finally
{
    Log.CloseAndFlush();
}
