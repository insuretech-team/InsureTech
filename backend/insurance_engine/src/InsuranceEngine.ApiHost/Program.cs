using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.Infrastructure.Persistence;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Infrastructure.DataGateways;
using InsuranceEngine.Beneficiary.GrpcServices;
using InsuranceEngine.Products.GrpcServices;
using InsuranceEngine.Policy.GrpcServices;
using InsuranceEngine.Claims.GrpcServices;
using InsuranceEngine.Renewals.GrpcServices;
using InsuranceEngine.Endorsements.GrpcServices;
using InsuranceEngine.FraudDetection.GrpcServices;
using InsuranceEngine.Underwriting.GrpcServices;
using InsuranceEngine.Grpc.Clients;
using InsuranceEngine.FraudDetection;
using InsuranceEngine.Cancellations;
using InsuranceEngine.Renewals;
using InsuranceEngine.Endorsements;

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
builder.Services.AddSingleton<IPdfGenerator, MockPdfGenerator>(); // FR-035 PDF generation
builder.Services.AddSingleton<IKafkaPublisher, MockKafkaPublisher>(); // FR-019 Kafka streaming

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
app.MapGrpcService<EndorsementGrpcService>();
app.MapGrpcService<FraudGrpcService>();
app.MapGrpcService<UnderwritingGrpcService>();

// 3. Schedule Recurring Background Jobs (Phase 3 / ST-003)
using (var scope = app.Services.CreateScope())
{
    var recurringJobManager = scope.ServiceProvider.GetRequiredService<IRecurringJobManager>();
    
    // FR-048: Automated policy status check for expiration/lapsing (Daily at 00:00)
    recurringJobManager.AddOrUpdate<PolicyBackgroundJobs>(
        "policy-auto-lapse", 
        job => job.ProcessAutoLapseAsync(), 
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
