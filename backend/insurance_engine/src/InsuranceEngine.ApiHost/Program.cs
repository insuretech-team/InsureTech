using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using InsuranceEngine.Beneficiary.Application.Commands;
using InsuranceEngine.Products.Application.Commands;
using InsuranceEngine.Policy.Application.Commands;
using InsuranceEngine.Underwriting.Application.Commands;
using InsuranceEngine.Claims.Application.Commands; // ClaimCommands
using InsuranceEngine.Commission.Application.Commands; // CalculateCommissionCommand
using InsuranceEngine.Beneficiary.GrpcServices;

var builder = WebApplication.CreateBuilder(args);

// Add services
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddDistributedMemoryCache(); // FR-028 Product caching
builder.Services.AddSingleton<InsuranceEngine.SharedKernel.Infrastructure.IPdfGenerator, InsuranceEngine.SharedKernel.Infrastructure.MockPdfGenerator>(); // FR-035 PDF generation
builder.Services.AddSingleton<InsuranceEngine.SharedKernel.Infrastructure.IKafkaPublisher, InsuranceEngine.SharedKernel.Infrastructure.MockKafkaPublisher>(); // FR-019 Kafka streaming
builder.Services.AddGrpc().AddJsonTranscoding();
builder.Services.AddGrpcReflection();
builder.Services.AddGrpcSwagger();
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new Microsoft.OpenApi.Models.OpenApiInfo { Title = "Insurance Engine API", Version = "v1" });
    c.IgnoreObsoleteProperties();
    c.IgnoreObsoleteActions();
    c.CustomSchemaIds(type => type.FullName);
    
    // Fix Resolver error for well-known types
    c.MapType<Google.Protobuf.WellKnownTypes.Value>(() => new OpenApiSchema { Type = "object", AdditionalPropertiesAllowed = true });
    c.MapType<Google.Protobuf.WellKnownTypes.Struct>(() => new OpenApiSchema { Type = "object", AdditionalPropertiesAllowed = true });
    c.MapType<Google.Protobuf.WellKnownTypes.ListValue>(() => new OpenApiSchema { Type = "array", Items = new OpenApiSchema { Type = "object", AdditionalPropertiesAllowed = true } });
});

// Database
var connectionString = builder.Configuration.GetConnectionString("InsuranceDb");
builder.Services.AddDbContext<DbContext>(options =>
    options.UseNpgsql(connectionString));

// MediatR
builder.Services.AddMediatR(cfg =>
{
    cfg.RegisterServicesFromAssembly(typeof(Program).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(CreateIndividualBeneficiaryCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(CreateProductCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(CreatePolicyCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(CreateQuoteCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(SubmitClaimCommand).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(CalculateCommissionCommand).Assembly);
});

var app = builder.Build();

// Configure the HTTP request pipeline.
app.UseSwagger();
app.UseSwaggerUI();

app.UseAuthorization();
app.MapControllers();
app.MapGrpcReflectionService();

// Map gRPC services
app.MapGrpcService<BeneficiaryGrpcService>();
app.MapGrpcService<InsuranceEngine.Products.GrpcServices.ProductGrpcService>();
app.MapGrpcService<InsuranceEngine.Policy.GrpcServices.PolicyGrpcService>();
app.MapGrpcService<InsuranceEngine.Claims.GrpcServices.ClaimGrpcService>();

app.MapGet("/", () => new
{
    service = "InsuranceEngine",
    version = "1.0.0",
    status = "running"
});

app.Run();
