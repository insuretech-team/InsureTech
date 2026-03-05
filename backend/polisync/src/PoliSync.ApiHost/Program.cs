using Serilog;
using PoliSync.Infrastructure;
using PoliSync.Products;
using PoliSync.Beneficiaries;

// ── Serilog bootstrap ───────────────────────────────────────────────
Log.Logger = new LoggerConfiguration()
    .WriteTo.Console()
    .WriteTo.File("logs/polisync-.log", rollingInterval: RollingInterval.Day)
    .CreateBootstrapLogger();

try
{
    Log.Information("Starting PoliSync API Host...");

    var builder = WebApplication.CreateBuilder(args);

    // ── Serilog ─────────────────────────────────────────────────────
    builder.Host.UseSerilog((ctx, lc) => lc
        .ReadFrom.Configuration(ctx.Configuration)
        .WriteTo.Console()
        .WriteTo.File("logs/polisync-.log", rollingInterval: RollingInterval.Day)
    );

    // ── Infrastructure (EF Core, Redis, CurrentUser) ────────────────
    builder.Services.AddInfrastructure(builder.Configuration);

    // ── Module registrations ────────────────────────────────────────
    builder.Services.AddProductsModule();
    builder.Services.AddBeneficiaryModule();

    // ── MediatR — scan all module assemblies ────────────────────────
    builder.Services.AddMediatR(cfg =>
    {
        cfg.RegisterServicesFromAssembly(typeof(PoliSync.Infrastructure.DependencyInjection).Assembly);
        cfg.RegisterServicesFromAssembly(typeof(PoliSync.Products.ProductsModule).Assembly);
        cfg.RegisterServicesFromAssembly(typeof(PoliSync.Beneficiaries.BeneficiaryModule).Assembly);
    });

    // ── REST API Controllers ────────────────────────────────────────
    builder.Services.AddControllers()
        .AddJsonOptions(opts =>
        {
            opts.JsonSerializerOptions.PropertyNamingPolicy = System.Text.Json.JsonNamingPolicy.SnakeCaseLower;
        });

    // ── Swagger ─────────────────────────────────────────────────────
    builder.Services.AddEndpointsApiExplorer();
    builder.Services.AddSwaggerGen(c =>
    {
        c.SwaggerDoc("v1", new Microsoft.OpenApi.Models.OpenApiInfo 
        { 
            Title = "PoliSync API", 
            Version = "v1",
            Description = "Digital Insurance Platform API" 
        });
        c.EnableAnnotations();
    });

    // ── gRPC ────────────────────────────────────────────────────────
    builder.Services.AddGrpc(opts =>
    {
        opts.EnableDetailedErrors = builder.Environment.IsDevelopment();
        opts.MaxReceiveMessageSize = 10 * 1024 * 1024; // 10 MB
    });
    builder.Services.AddGrpcReflection();

    // ── Health checks ───────────────────────────────────────────────
    builder.Services.AddHealthChecks()
        .AddNpgSql(
            builder.Configuration.GetConnectionString("InsuranceDb")!,
            name: "postgresql",
            tags: ["db", "ready"]
        );

    var app = builder.Build();

    // ── Middleware pipeline ──────────────────────────────────────────
    app.UseSerilogRequestLogging();

    // ── Swagger UI (Home Page) ──────────────────────────────────────
    if (app.Environment.IsDevelopment() || true)
    {
        app.UseSwagger();
        app.UseSwaggerUI(c =>
        {
            c.SwaggerEndpoint("/swagger/v1/swagger.json", "PoliSync API V1");
            c.RoutePrefix = string.Empty; // Serve Swagger at the root
        });
    }

    // ── Health check endpoint ───────────────────────────────────────
    app.MapHealthChecks("/health");

    // ── REST API endpoints ──────────────────────────────────────────
    app.MapControllers();

    // ── gRPC reflection (dev only) ──────────────────────────────────
    if (app.Environment.IsDevelopment())
    {
        app.MapGrpcReflectionService();
    }

    Log.Information("PoliSync API Host started successfully");
    app.Run();
}
catch (Exception ex)
{
    Log.Fatal(ex, "PoliSync API Host terminated unexpectedly");
}
finally
{
    Log.CloseAndFlush();
}
