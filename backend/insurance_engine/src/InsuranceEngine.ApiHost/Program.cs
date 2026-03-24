using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Beneficiary.Application.Commands;

var builder = WebApplication.CreateBuilder(args);

// Add services
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Database
var connectionString = builder.Configuration.GetConnectionString("InsuranceDb");
builder.Services.AddDbContext<DbContext>(options =>
    options.UseNpgsql(connectionString));

// MediatR
builder.Services.AddMediatR(cfg =>
{
    cfg.RegisterServicesFromAssembly(typeof(Program).Assembly);
    cfg.RegisterServicesFromAssembly(typeof(CreateIndividualBeneficiaryCommand).Assembly);
});

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseAuthorization();
app.MapControllers();

app.MapGet("/", () => new
{
    service = "InsuranceEngine",
    version = "1.0.0",
    status = "running"
});

app.Run();
