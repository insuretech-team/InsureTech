using System.Text.Json;

namespace PoliSync.Workflow.Domain;

/// <summary>
/// Represents a configurable dynamic workflow template that drives
/// the orchestration logic without requiring code changes.
///
/// Templates are stored as workflow definitions in the Go workflow-engine
/// (workflow_schema.workflow_definitions). This class is the C# domain model
/// used by PoliSync to reason about and execute workflow logic.
/// </summary>
public sealed class WorkflowTemplate
{
    public string DefinitionId { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string EntityType { get; init; } = string.Empty;  // CLAIM, ENDORSEMENT, REFUND, QUOTATION, POLICY
    public string WorkflowType { get; init; } = "APPROVAL";  // APPROVAL, REVIEW, ESCALATION, NOTIFICATION
    public string Description { get; init; } = string.Empty;
    public IReadOnlyList<WorkflowStepTemplate> Steps { get; init; } = [];
    public WorkflowConditions Conditions { get; init; } = new();

    /// <summary>Serialise steps to JSON for sending to the Go service.</summary>
    public string SerializeSteps()
        => JsonSerializer.Serialize(Steps, JsonOptions);

    /// <summary>Serialise conditions to JSON.</summary>
    public string SerializeConditions()
        => JsonSerializer.Serialize(Conditions, JsonOptions);

    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        WriteIndented = false
    };

    /// <summary>
    /// Factory: build a standard single-approver approval template.
    /// </summary>
    public static WorkflowTemplate SingleApproval(
        string name,
        string entityType,
        string assignToRole,
        int dueDays = 3,
        string description = "")
        => new()
        {
            Name = name,
            EntityType = entityType,
            WorkflowType = "APPROVAL",
            Description = description,
            Steps =
            [
                new WorkflowStepTemplate
                {
                    Name = "initial_review",
                    Type = "APPROVAL",
                    AssignRole = assignToRole,
                    DueHours = dueDays * 24,
                    Order = 1
                }
            ]
        };

    /// <summary>
    /// Factory: build a two-stage approval template (review then approve).
    /// </summary>
    public static WorkflowTemplate TwoStageApproval(
        string name,
        string entityType,
        string reviewerRole,
        string approverRole,
        int reviewDueDays = 2,
        int approvalDueDays = 3)
        => new()
        {
            Name = name,
            EntityType = entityType,
            WorkflowType = "APPROVAL",
            Steps =
            [
                new WorkflowStepTemplate
                {
                    Name = "technical_review",
                    Type = "REVIEW",
                    AssignRole = reviewerRole,
                    DueHours = reviewDueDays * 24,
                    Order = 1
                },
                new WorkflowStepTemplate
                {
                    Name = "final_approval",
                    Type = "APPROVAL",
                    AssignRole = approverRole,
                    DueHours = approvalDueDays * 24,
                    Order = 2
                }
            ]
        };
}

/// <summary>
/// A single step within a workflow template.
/// Serialized as JSONB into the Go workflow_definitions.steps column.
/// </summary>
public sealed class WorkflowStepTemplate
{
    public string Name { get; init; } = string.Empty;
    public string Type { get; init; } = "APPROVAL"; // APPROVAL, REVIEW, NOTIFICATION, ACTION
    public string AssignRole { get; init; } = string.Empty;
    public string AssignTo { get; init; } = string.Empty; // explicit user UUID (overrides role)
    public int DueHours { get; init; } = 72;
    public int Order { get; init; } = 1;
}

/// <summary>
/// Conditional routing rules for the workflow template.
/// Evaluated by the Go engine to determine step execution order.
/// </summary>
public sealed class WorkflowConditions
{
    /// <summary>If true, any rejection immediately fails the entire workflow.</summary>
    public bool FailFastOnRejection { get; init; } = true;

    /// <summary>If true, all approvers must approve (AND gate). Otherwise first approval wins (OR gate).</summary>
    public bool RequireAllApprovals { get; init; } = true;

    /// <summary>Auto-approve if no action taken within escalation hours.</summary>
    public int? AutoApproveAfterHours { get; init; }

    /// <summary>Escalate to role if no action within escalation hours.</summary>
    public string? EscalateToRole { get; init; }

    /// <summary>Additional custom rules as key-value pairs.</summary>
    public Dictionary<string, string> CustomRules { get; init; } = [];
}
