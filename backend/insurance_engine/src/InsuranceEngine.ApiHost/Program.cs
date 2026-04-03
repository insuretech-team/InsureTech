using Microsoft.EntityFrameworkCore;
using Hangfire;
using Hangfire.PostgreSql;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.Infrastructure.Persistence;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Infrastructure.DataGateways;
using InsuranceEngine.Infrastructure.Messaging;
using InsuranceEngine.Beneficiary.GrpcServices;
using InsuranceEngine.Products.GrpcServices;
using InsuranceEngine.Policy.GrpcServices;
using InsuranceEngine.Claims.GrpcServices;
using InsuranceEngine.Renewals.GrpcServices;
// using InsuranceEngine.Endorsements.GrpcServices; // TODO: Re-enable when proto types are available
using InsuranceEngine.FraudDetection.GrpcServices;
using InsuranceEngine.Underwriting.GrpcServices;
using InsuranceEngine.Grpc.Clients;
using InsuranceEngine.FraudDetection;
using InsuranceEngine.Cancellations;
using InsuranceEngine.Renewals;
using InsuranceEngine.Endorsements;
using InsuranceEngine.Products;
using InsuranceEngine.Policy;
using InsuranceEngine.Claims;
using InsuranceEngine.Underwriting;
using InsuranceEngine.Beneficiary;
using InsuranceEngine.Commission;
using InsuranceEngine.Infrastructure.Notifications;
using InsuranceEngine.Infrastructure.Documents;
using InsuranceEngine.Infrastructure.Refunds;
using InsuranceEngine.Infrastructure.Payments;
using InsuranceEngine.Infrastructure.Webhooks;
using InsuranceEngine.Infrastructure;
using InsuranceEngine.Infrastructure.Renewals;
using Microsoft.Extensions.Options;
using Microsoft.Extensions.Logging;

var builder = WebApplication.CreateBuilder(args);

// Add services
builder.Services.AddHttpContextAccessor();
builder.Services.AddSingleton<GrpcClientFactory>();
builder.Services.AddScoped<InsuranceServiceClient>();

builder.Services.AddProductsModule();
builder.Services.AddPolicyModule();
builder.Services.AddClaimsModule();
builder.Services.AddUnderwritingModule();
builder.Services.AddBeneficiaryModule();
builder.Services.AddCommissionModule();
builder.Services.AddFraudDetectionModule();
builder.Services.AddCancellationsModule();
builder.Services.AddRenewalsModule();
builder.Services.AddEndorsementsModule();

builder.Services.AddMediatR(cfg =>
{
    cfg.RegisterServicesFromAssembly(typeof(Program).Assembly);
});

// 1. Hangfire Background Task Setup (Phase 1 / ST-001)
var connectionString = builder.Configuration.GetConnectionString("InsuranceDb");
builder.Services.AddHangfire(config => config
    .SetDataCompatibilityLevel(CompatibilityLevel.Version_180)
    .UseSimpleAssemblyNameTypeSerializer()
    .UseRecommendedSerializerSettings()
    .UsePostgreSqlStorage(options =>
    {
        options.UseNpgsqlConnection(connectionString);
    }));
builder.Services.AddHangfireServer();

builder.Services.AddScoped<PolicyBackgroundJobs>(); // FR-048 Background tasks

builder.Services.AddDistributedMemoryCache(); // FR-028 Product caching (TODO: Replace with Redis)

// M1-3: PDF Generation Service (FR-035)
// Use real DocumentService when Go backend is available, otherwise mock
var useRealDocumentService = builder.Configuration["Features:UseRealDocumentService"]?.ToLower() == "true";
if (useRealDocumentService)
{
    builder.Services.AddScoped<IDocumentService, DocumentService>();
    builder.Services.AddSingleton<IPdfGenerator>(sp =>
    {
        var documentService = sp.GetRequiredService<IDocumentService>();
        var logger = sp.GetRequiredService<ILogger<GoDocumentPdfGenerator>>();
        return new GoDocumentPdfGenerator(documentService, logger);
    });
}
else
{
    builder.Services.AddScoped<IDocumentService, MockDocumentService>();
    builder.Services.AddSingleton<IPdfGenerator, MockPdfGenerator>();
}

// M1-2: Notification Service (FR-136, FR-137)
builder.Services.AddScoped<INotificationService, NotificationService>();

// M1-4: Pro-rata Refund Calculation (FR-053)
builder.Services.AddScoped<IRefundService, RefundService>();

// M1-6: Payment Verification Workflow (FR-054)
builder.Services.AddScoped<IPaymentService, PaymentService>();

// M1-5: Partner Notification Webhooks (FR-139)
builder.Services.AddHttpClient("WebhookClient");
builder.Services.AddScoped<IWebhookService, WebhookService>();

// M1-7: Grace Period Workflow (FR-047, FR-048)
builder.Services.Configure<GracePeriodSettings>(builder.Configuration.GetSection("GracePeriod"));
builder.Services.AddScoped<IGracePeriodService, GracePeriodService>();
builder.Services.AddScoped<IEventPublisher, GracePeriodEventPublisher>();

// M1-1: Real Kafka Integration (FR-136)
builder.Services.Configure<KafkaSettings>(options =>
{
    options.BootstrapServers = builder.Configuration["Kafka:BootstrapServers"] ?? "localhost:9092";
    options.ClientId = "insurance-engine";
    options.EnableIdempotence = true;
});
builder.Services.AddSingleton<IEventPublisher, KafkaEventPublisher>();

// Register IKafkaPublisher adapter for existing handlers
builder.Services.AddScoped<IKafkaPublisher>(sp =>
{
    var eventPublisher = sp.GetRequiredService<IEventPublisher>();
    var logger = sp.GetRequiredService<ILogger<KafkaPublisherAdapter>>();
    return new KafkaPublisherAdapter(eventPublisher, logger);
});

// gRPC with error interceptor
builder.Services.AddGrpc(options =>
{
    options.Interceptors.Add<GlobalGrpcErrorInterceptor>();
    options.MaxReceiveMessageSize = 16 * 1024 * 1024; // 16MB
    options.MaxSendMessageSize = 16 * 1024 * 1024;
});
builder.Services.AddGrpcReflection();

// Database — Full EF Core (Option A)
// connectionString is already declared above
builder.Services.AddDbContext<InsuranceDbContext>(options =>
    options.UseNpgsql(connectionString)
           .UseSnakeCaseNamingConvention());

// Register bare DbContext for backward compatibility (handlers still reference it)
builder.Services.AddScoped<DbContext>(sp => sp.GetRequiredService<InsuranceDbContext>());

// Repository pattern
builder.Services.AddScoped(typeof(IRepository<>), typeof(EfRepository<>));
builder.Services.AddScoped<ISequenceDataGateway, SqlSequenceDataGateway>();

builder.Services.AddAuthorization();

var app = builder.Build();

// Authorization
app.UseAuthorization();

// 2. Hangfire Dashboard (Monitoring Phase 1 / ST-001)
app.UseHangfireDashboard("/hangfire", new DashboardOptions
{
    DashboardTitle = "LabAid Insurance Engine - Background Jobs",
    Authorization = [] // Note: In production, add a custom IAuthorizationFilter
});

app.MapGrpcReflectionService();

// Map gRPC services
app.MapGrpcService<BeneficiaryGrpcService>();
app.MapGrpcService<ProductGrpcService>();
app.MapGrpcService<PolicyGrpcService>();
app.MapGrpcService<ClaimGrpcService>();
app.MapGrpcService<RenewalGrpcService>();
// app.MapGrpcService<EndorsementGrpcService>(); // TODO: Re-enable when proto types are available
app.MapGrpcService<FraudGrpcService>();
app.MapGrpcService<UnderwritingGrpcService>();

// 3. Schedule Recurring Background Jobs (Phase 3 / ST-003)
using (var scope = app.Services.CreateScope())
{
    var recurringJobManager = scope.ServiceProvider.GetRequiredService<IRecurringJobManager>();
    
    // FR-047 + FR-048: Grace Period workflow (Daily at 00:00)
    // - Move expired policies to GRACE_PERIOD (30 days)
    // - Send daily reminders during grace period
    // - Auto-lapse after grace period expires
    recurringJobManager.AddOrUpdate<PolicyBackgroundJobs>(
        "policy-grace-period-workflow", 
        job => job.ProcessGracePeriodWorkflowAsync(), 
        Cron.Daily(0, 0));

    // FR-045: Renewal reminders for policies expiring in 30 days (Daily at 01:00)
    recurringJobManager.AddOrUpdate<PolicyBackgroundJobs>(
        "policy-renewal-reminders", 
        job => job.ProcessRenewalRemindersAsync(), 
        Cron.Daily(1, 0));
}

app.MapGet("/", () => new
{
    service = "InsuranceEngine",
    version = "1.0.0",
    status = "running"
});

app.Run();
