using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using PoliSync.ApiHost.BackgroundServices;
using PoliSync.ApiHost.Interceptors;
using PoliSync.ApiHost.Services;
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
using PoliSync.Workflow;
using PoliSync.Workflow.Infrastructure;
using Serilog;
using System.Security.Cryptography;
using System.Reflection;

var builder = WebApplication.CreateBuilder(args);

// Configure Serilog
Log.Logger = new LoggerConfiguration()
    .ReadFrom.Configuration(builder.Configuration)
    .CreateLogger();

builder.Host.UseSerilog();

// Add services
builder.Services.AddGrpc(options =>
{
    options.Interceptors.Add<AuthInterceptor>();
    options.Interceptors.Add<JwtAuthInterceptor>();
    options.Interceptors.Add<LoggingInterceptor>();
    options.Interceptors.Add<ValidationInterceptor>();
    options.EnableDetailedErrors = true;
    options.MaxReceiveMessageSize = 16 * 1024 * 1024; // 16MB
    options.MaxSendMessageSize = 16 * 1024 * 1024;
});

builder.Services.AddGrpcReflection();
builder.Services.AddControllers();

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

// GoProductDataGateway — routes ALL product/plan/rider/pricing calls through Go insurance gRPC
// This is the single source of truth pattern: PoliSync → gRPC → Go → DB
builder.Services.AddScoped<PoliSync.Products.Infrastructure.IProductRepository, 
    PoliSync.Products.Infrastructure.GoProductDataGateway>();
builder.Services.AddScoped<PoliSync.Products.Infrastructure.GoProductDataGateway>();

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
builder.Services.AddSingleton<PoliSync.Infrastructure.GrpcClients.GrpcClientFactory>();
builder.Services.AddSingleton<PoliSync.Infrastructure.GrpcClients.DocgenGrpcClient>();
builder.Services.AddSingleton<PoliSync.Infrastructure.Clients.InsuranceServiceClient>();
builder.Services.AddSingleton<PoliSync.Infrastructure.Clients.OrderServiceGrpcClient>();
builder.Services.AddSingleton<PoliSync.Infrastructure.Clients.PaymentServiceGrpcClient>();
builder.Services.AddSingleton<PoliSync.Infrastructure.Clients.CommissionServiceGrpcClient>();
builder.Services.AddScoped<IPolicyDataGateway, GoPolicyDataGateway>();
builder.Services.AddScoped<InsuranceProposalWorkflowService>();
builder.Services.AddScoped<IQuotationDataGateway, GoQuotationDataGateway>();
builder.Services.AddScoped<IClaimDataGateway, GoClaimDataGateway>();
builder.Services.AddScoped<IEndorsementDataGateway, GoEndorsementDataGateway>();
builder.Services.AddScoped<IRenewalDataGateway, GoRenewalDataGateway>();
builder.Services.AddScoped<IOrderDataGateway, GoOrderDataGateway>();
builder.Services.AddScoped<IRefundPaymentGateway, GoRefundPaymentGateway>();
builder.Services.AddScoped<PoliSync.Commission.Infrastructure.ICommissionDataGateway, PoliSync.Commission.Infrastructure.GoCommissionDataGateway>();
builder.Services.AddScoped<IUnderwritingDataGateway, GoUnderwritingDataGateway>();
builder.Services.AddScoped<IWorkflowDataGateway, GoWorkflowDataGateway>();
builder.Services.AddSingleton<IUnderwritingRiskScorer, UnderwritingRiskScorer>();
builder.Services.AddHostedService<UnderwritingQuotationSubmittedConsumer>();
builder.Services.AddHostedService<OrderPaymentConfirmedConsumer>();
builder.Services.AddHostedService<InsuranceProposalDecisionConsumer>();
builder.Services.AddHostedService<QuotationExpiryService>();

// Domain modules
builder.Services.AddProductsModule();

// Workflow engine — IWorkflowDataGateway + WorkflowTemplateSeeder (IHostedService)
builder.Services.AddWorkflow();

// Business Rules Engine (Microsoft RulesEngine)
InvokeModuleRegistration(builder.Services, "PoliSync.RulesEngine", "PoliSync.RulesEngine.DependencyInjection", "AddRulesEngineServices");

// Quoting Service
InvokeModuleRegistration(builder.Services, "PoliSync.Quoting", "PoliSync.Quoting.DependencyInjection", "AddQuotingServices");

// Vehicle Insurance Service
InvokeModuleRegistration(builder.Services, "PoliSync.VehicleInsurance", "PoliSync.VehicleInsurance.DependencyInjection", "AddVehicleInsuranceServices");

// Life Insurance Service
InvokeModuleRegistration(builder.Services, "PoliSync.LifeInsurance", "PoliSync.LifeInsurance.DependencyInjection", "AddLifeInsuranceServices");

// CRM Service
InvokeModuleRegistration(builder.Services, "PoliSync.CRM", "PoliSync.CRM.DependencyInjection", "AddCrmServices");

// Actuarial Service
InvokeModuleRegistration(builder.Services, "PoliSync.Actuarial", "PoliSync.Actuarial.DependencyInjection", "AddActuarialServices");

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
    cfg.RegisterServicesFromAssemblyContaining<PoliSync.Workflow.AssemblyMarker>();
    cfg.RegisterServicesFromAssembly(LoadModuleAssembly("PoliSync.RulesEngine"));
    cfg.RegisterServicesFromAssembly(LoadModuleAssembly("PoliSync.Quoting"));
    cfg.RegisterServicesFromAssembly(LoadModuleAssembly("PoliSync.VehicleInsurance"));
    cfg.RegisterServicesFromAssembly(LoadModuleAssembly("PoliSync.LifeInsurance"));
    cfg.RegisterServicesFromAssembly(LoadModuleAssembly("PoliSync.CRM"));
    cfg.RegisterServicesFromAssembly(LoadModuleAssembly("PoliSync.Actuarial"));
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
builder.Services.AddScoped<AuthInterceptor>();
builder.Services.AddScoped<JwtAuthInterceptor>();
builder.Services.AddScoped<LoggingInterceptor>();
builder.Services.AddScoped<ValidationInterceptor>();

// Health Checks — DB health check is conditional (PoliSync delegates CRUD to Go insurance service)
var healthChecks = builder.Services.AddHealthChecks();
if (!string.IsNullOrEmpty(insuranceConnectionString))
{
    healthChecks.AddNpgSql(insuranceConnectionString, name: "postgres");
}
healthChecks
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
app.MapGrpcService<PoliSync.Workflow.GrpcServices.WorkflowGrpcService>();
MapGrpcServiceByName(app, "PoliSync.RulesEngine", "PoliSync.RulesEngine.GrpcServices.BusinessWorkflowGrpcService");
MapGrpcServiceByName(app, "PoliSync.Quoting", "PoliSync.Quoting.GrpcServices.QuotingGrpcService");
MapGrpcServiceByName(app, "PoliSync.VehicleInsurance", "PoliSync.VehicleInsurance.GrpcServices.VehicleGrpcService");
MapGrpcServiceByName(app, "PoliSync.LifeInsurance", "PoliSync.LifeInsurance.GrpcServices.LifeInsuranceGrpcService");
MapGrpcServiceByName(app, "PoliSync.CRM", "PoliSync.CRM.GrpcServices.CrmGrpcService");
MapGrpcServiceByName(app, "PoliSync.Actuarial", "PoliSync.Actuarial.GrpcServices.ActuarialGrpcService");

app.MapGrpcReflectionService();

// Health checks
app.MapHealthChecks("/health");
app.MapControllers();

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

static Assembly LoadModuleAssembly(string assemblyName)
{
    return AppDomain.CurrentDomain.GetAssemblies()
        .FirstOrDefault(a => string.Equals(a.GetName().Name, assemblyName, StringComparison.Ordinal))
        ?? Assembly.Load(assemblyName);
}

static void InvokeModuleRegistration(IServiceCollection services, string assemblyName, string typeName, string methodName)
{
    var assembly = LoadModuleAssembly(assemblyName);
    var type = assembly.GetType(typeName, throwOnError: true)!;
    var method = type.GetMethod(methodName, BindingFlags.Public | BindingFlags.Static, [typeof(IServiceCollection)]);
    if (method is null)
    {
        throw new InvalidOperationException($"Could not find {typeName}.{methodName}(IServiceCollection).");
    }

    method.Invoke(null, [services]);
}

static void MapGrpcServiceByName(IEndpointRouteBuilder app, string assemblyName, string serviceTypeName)
{
    var assembly = LoadModuleAssembly(assemblyName);
    var serviceType = assembly.GetType(serviceTypeName, throwOnError: true)!;

    var mapGrpcService = typeof(Microsoft.AspNetCore.Builder.GrpcEndpointRouteBuilderExtensions)
        .GetMethods(BindingFlags.Public | BindingFlags.Static)
        .Single(m => m.Name == "MapGrpcService" && m.IsGenericMethodDefinition && m.GetParameters().Length == 1);

    mapGrpcService.MakeGenericMethod(serviceType).Invoke(null, [app]);
}
