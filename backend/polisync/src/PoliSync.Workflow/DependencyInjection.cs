using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Options;
using PoliSync.Workflow.Domain;
using PoliSync.Workflow.Infrastructure;
using PoliSync.Workflow.Services;

namespace PoliSync.Workflow;

/// <summary>
/// Extension methods for registering the complete PoliSync.Workflow module.
/// Call from Program.cs: builder.Services.AddWorkflow(builder.Configuration);
/// </summary>
public static class DependencyInjection
{
    public static IServiceCollection AddWorkflow(this IServiceCollection services)
    {
        // ── Options (hot-reloadable via IOptionsMonitor) ──────────────────────
        services.AddOptions<WorkflowTemplateOptions>()
            .BindConfiguration(WorkflowTemplateOptions.SectionName)
            .ValidateOnStart();

        // ── Template Provider (singleton — thread-safe ConcurrentDictionary) ─
        services.AddSingleton<CompositeWorkflowTemplateProvider>();
        services.AddSingleton<IWorkflowTemplateProvider>(sp =>
            sp.GetRequiredService<CompositeWorkflowTemplateProvider>());

        // ── Template Registrar (singleton — hot-reload via IOptionsMonitor) ──
        services.AddSingleton<WorkflowTemplateRegistrar>();

        // ── gRPC client for Go workflow-engine ────────────────────────────────
        services.AddSingleton<WorkflowServiceGrpcClient>();

        // ── Data gateway — communicates with Go workflow-engine via gRPC ──────
        services.AddScoped<IWorkflowDataGateway, GoWorkflowDataGateway>();

        // ── MediatR handlers — registered via AssemblyMarker ─────────────────
        services.AddMediatR(cfg =>
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly));

        // ── Startup seeder — registers canonical templates into the Go engine ─
        services.AddHostedService<WorkflowTemplateSeeder>();

        return services;
    }
}
