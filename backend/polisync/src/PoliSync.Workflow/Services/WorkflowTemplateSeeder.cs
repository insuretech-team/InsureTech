using MediatR;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;

namespace PoliSync.Workflow.Services;

/// <summary>
/// Hosted service that runs at startup to:
///   1. Register all canonical + config templates via WorkflowTemplateRegistrar
///   2. Persist each template to the Go workflow-engine idempotently
///
/// The WorkflowTemplateRegistrar handles the in-memory registry and routing rules.
/// This seeder handles persistence to the Go engine so templates survive Go restarts.
/// </summary>
public sealed class WorkflowTemplateSeeder : IHostedService
{
    private readonly IServiceScopeFactory _scopeFactory;
    private readonly WorkflowTemplateRegistrar _registrar;
    private readonly ILogger<WorkflowTemplateSeeder> _logger;

    public WorkflowTemplateSeeder(
        IServiceScopeFactory scopeFactory,
        WorkflowTemplateRegistrar registrar,
        ILogger<WorkflowTemplateSeeder> logger)
    {
        _scopeFactory = scopeFactory;
        _registrar = registrar;
        _logger = logger;
    }

    public async Task StartAsync(CancellationToken cancellationToken)
    {
        _logger.LogInformation("Registering workflow templates in-memory...");

        // Step 1: Register all templates in-memory (fast, synchronous)
        _registrar.RegisterAll();

        _logger.LogInformation("Seeding workflow templates to Go workflow-engine...");

        using var scope = _scopeFactory.CreateScope();
        var mediator = scope.ServiceProvider.GetRequiredService<IMediator>();
        var provider = scope.ServiceProvider.GetRequiredService<IWorkflowTemplateProvider>();

        var templates = provider.GetAllTemplates();
        var succeeded = 0;
        var skipped = 0;

        foreach (var template in templates)
        {
            try
            {
                var result = await mediator.Send(
                    new RegisterWorkflowTemplateCommand(template),
                    cancellationToken);

                if (result.IsSuccess)
                {
                    var action = result.Value!.WasCreated ? "Registered" : "Already exists";
                    _logger.LogDebug("{Action}: '{Name}' → {Id}", action, result.Value.Name, result.Value.DefinitionId);
                    succeeded++;
                }
                else
                {
                    _logger.LogWarning("Failed to seed template '{Name}': {Error}",
                        template.Name, result.Error?.Message);
                    skipped++;
                }
            }
            catch (Exception ex)
            {
                // Don't fail startup if workflow-engine is temporarily unavailable
                _logger.LogWarning(ex, "Could not seed template '{Name}' — workflow-engine may be starting up", template.Name);
                skipped++;
            }
        }

        _logger.LogInformation(
            "Workflow template seeding complete: {Succeeded} persisted, {Skipped} skipped (out of {Total} total)",
            succeeded, skipped, templates.Count);
    }

    public Task StopAsync(CancellationToken cancellationToken) => Task.CompletedTask;
}
