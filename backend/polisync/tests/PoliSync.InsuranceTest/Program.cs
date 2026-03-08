using Grpc.Net.Client;
using Microsoft.Extensions.Configuration;
using PoliSync.InsuranceTest;
using Npgsql;

Console.WriteLine("==============================================");
Console.WriteLine("Insurance Service CRUD Test");
Console.WriteLine("==============================================");
Console.WriteLine();

// Load configuration
var configuration = new ConfigurationBuilder()
    .SetBasePath(Directory.GetCurrentDirectory())
    .AddJsonFile("appsettings.json", optional: false)
    .Build();

var serviceUrl = configuration["InsuranceService:Url"] ?? "http://localhost:50115";

// Build connection string from environment variables (like .env)
var pgHost = Environment.GetEnvironmentVariable("PGHOST") ?? "localhost";
var pgPort = Environment.GetEnvironmentVariable("PGPORT") ?? "5432";
var pgDatabase = Environment.GetEnvironmentVariable("PGDATABASE") ?? "insuretech_db";
var pgUser = Environment.GetEnvironmentVariable("PGUSER") ?? "insuretech_user";
var pgPassword = Environment.GetEnvironmentVariable("PGPASSWORD") ?? "insuretech_pass_2024";
var pgSslMode = Environment.GetEnvironmentVariable("PGSSLMODE") ?? "disable";

var connectionString = $"Host={pgHost};Port={pgPort};Database={pgDatabase};Username={pgUser};Password={pgPassword};SSL Mode={pgSslMode}";

Console.WriteLine($"Connecting to Insurance Service at: {serviceUrl}");
Console.WriteLine($"Database: {pgHost}:{pgPort}/{pgDatabase}");
Console.WriteLine();

// Get valid user UUID from database
string validUserUuid;
try
{
    using var conn = new NpgsqlConnection(connectionString);
    await conn.OpenAsync();
    using var cmd = new NpgsqlCommand("SELECT user_id FROM authn_schema.users LIMIT 1", conn);
    var result = await cmd.ExecuteScalarAsync();
    
    if (result == null)
    {
        Console.WriteLine("❌ ERROR: No users found in database.");
        Console.WriteLine("Please create a user first or check database connection.");
        Environment.Exit(1);
        return;
    }
    
    validUserUuid = result.ToString()!;
    Console.WriteLine($"✓ Using user UUID: {validUserUuid}");
    Console.WriteLine();
}
catch (Exception ex)
{
    Console.WriteLine($"❌ Database connection failed: {ex.Message}");
    Console.WriteLine("Please ensure PostgreSQL is running and accessible.");
    Environment.Exit(1);
    return;
}

// Create gRPC channel
using var channel = GrpcChannel.ForAddress(serviceUrl);

// Run simple tests (comprehensive test needs proto field verification)
var testRunner = new SimpleTestRunner(channel, validUserUuid);

try
{
    await testRunner.RunAllTests();
    
    if (testRunner.FailedTests > 0)
    {
        Console.WriteLine("⚠️  Some tests failed. Please check the errors above.");
        Environment.Exit(1);
    }
    else
    {
        Console.WriteLine("✅ All tests passed successfully!");
    }
}
catch (Exception ex)
{
    Console.WriteLine();
    Console.WriteLine($"❌ FATAL ERROR: {ex.Message}");
    Console.WriteLine($"Stack Trace: {ex.StackTrace}");
    Environment.Exit(1);
}
