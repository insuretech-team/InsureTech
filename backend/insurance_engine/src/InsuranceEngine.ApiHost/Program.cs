using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.Beneficiary.GrpcServices;
using InsuranceEngine.Products.GrpcServices;
using InsuranceEngine.Policy.GrpcServices;
using InsuranceEngine.Claims.GrpcServices;
using InsuranceEngine.Renewals.GrpcServices;
using InsuranceEngine.Endorsements.GrpcServices;
using InsuranceEngine.FraudDetection.GrpcServices;
using InsuranceEngine.Underwriting.GrpcServices;

var builder = WebApplication.CreateBuilder(args);

// Add services
builder.Services.AddMediatR(cfg =>
{
    cfg.RegisterServicesFromAssembly(typeof(Program).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Beneficiary.Application.Commands.CreateIndividualBeneficiaryCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Products.Application.Commands.CreateProductCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Policy.Application.Commands.CreatePolicyCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Underwriting.Application.Commands.RequestQuoteCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Claims.Application.Commands.SubmitClaimCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Commission.Application.Commands.CalculateCommissionCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Cancellations.Application.Commands.CancelPolicyCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Renewals.Application.Commands.RenewPolicyCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.Endorsements.Application.Commands.UpdatePolicyCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(InsuranceEngine.FraudDetection.Application.Commands.CheckFraudCommand).Assembly);
});

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
var connectionString = builder.Configuration.GetConnectionString("InsuranceDb");
builder.Services.AddDbContext<InsuranceDbContext>(options =>
    options.UseNpgsql(connectionString)
           .UseSnakeCaseNamingConvention());

// Register bare DbContext for backward compatibility (handlers still reference it)
builder.Services.AddScoped<DbContext>(sp => sp.GetRequiredService<InsuranceDbContext>());

// Repository pattern
builder.Services.AddScoped(typeof(IRepository<>), typeof(EfRepository<>));

builder.Services.AddAuthorization();

var app = builder.Build();

// Authorization
app.UseAuthorization();
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

app.MapGet("/", () => new
{
    service = "InsuranceEngine",
    version = "1.0.0",
    status = "running"
});

app.Run();
