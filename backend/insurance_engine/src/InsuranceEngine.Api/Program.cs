using InsuranceEngine.Application;
using InsuranceEngine.Infrastructure;
using InsuranceEngine.Infrastructure.Persistence;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Clean Architecture Layers
builder.Services.AddApplication();
builder.Services.AddInfrastructure(builder.Configuration);

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

if (app.Environment.IsDevelopment())
{
    app.MapGrpcReflectionService();
}

// Seed Database
await DbInitializer.Initialize(app.Services);

app.MapGrpcService<InsuranceEngine.Api.GrpcServices.InsuranceGrpcService>();

app.Run();
