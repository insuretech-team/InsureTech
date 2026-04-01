using MediatR;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Workflow.Domain;

namespace PoliSync.Workflow.Application.Commands;

/// <summary>
/// Registers (creates or updates) a workflow template in the Go workflow engine.
/// This allows dynamic template registration at startup or via admin API without
/// requiring code deployment or database migrations.
/// </summary>
public sealed record RegisterWorkflowTemplateCommand(
    WorkflowTemplate Template
) : IRequest<Result<RegisterWorkflowTemplateResult>>;

public sealed record RegisterWorkflowTemplateResult(
    string DefinitionId,
    string Name,
    bool WasCreated  // false = already existed
);
