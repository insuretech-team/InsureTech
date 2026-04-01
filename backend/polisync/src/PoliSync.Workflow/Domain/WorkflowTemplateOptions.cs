namespace PoliSync.Workflow.Domain;

/// <summary>
/// Configuration options for dynamic workflow templates loaded from appsettings.json.
///
/// Example appsettings.json:
/// <code>
/// "WorkflowTemplates": {
///   "HighValueClaimThresholdPaisa": 10000000,
///   "HighValueRefundThresholdPaisa": 5000000,
///   "AdditionalTemplates": [
///     {
///       "Name": "claim.vip-approval",
///       "EntityType": "CLAIM",
///       "WorkflowType": "APPROVAL",
///       "Description": "VIP customer fast-track claim approval",
///       "Steps": [
///         { "Name": "vip_review", "Type": "APPROVAL", "AssignRole": "vip_manager", "DueHours": 4, "Order": 1 }
///       ]
///     }
///   ]
/// }
/// </code>
/// </summary>
public sealed class WorkflowTemplateOptions
{
    public const string SectionName = "WorkflowTemplates";

    /// <summary>Claims above this threshold use high-value approval flow (paisa, BDT × 100).</summary>
    public long HighValueClaimThresholdPaisa { get; init; } = 10_000_000; // 1,00,000 BDT

    /// <summary>Refunds above this threshold use high-value approval flow.</summary>
    public long HighValueRefundThresholdPaisa { get; init; } = 5_000_000; // 50,000 BDT

    /// <summary>Endorsements marked as "major" types use two-stage approval.</summary>
    public IReadOnlyList<string> MajorEndorsementTypes { get; init; } =
    [
        "SumAssuredChange",
        "PremiumAdjustment",
        "BeneficiaryChange"
    ];

    /// <summary>Additional templates loaded from config (hot-reloadable via IOptionsMonitor).</summary>
    public IReadOnlyList<WorkflowTemplateConfig> AdditionalTemplates { get; init; } = [];
}

/// <summary>
/// JSON-configurable workflow template definition.
/// Loaded from appsettings.json WorkflowTemplates:AdditionalTemplates array.
/// </summary>
public sealed class WorkflowTemplateConfig
{
    public string Name { get; init; } = string.Empty;
    public string EntityType { get; init; } = string.Empty;
    public string WorkflowType { get; init; } = "APPROVAL";
    public string Description { get; init; } = string.Empty;
    public IReadOnlyList<WorkflowStepConfig> Steps { get; init; } = [];
    public WorkflowConditionsConfig Conditions { get; init; } = new();
}

public sealed class WorkflowStepConfig
{
    public string Name { get; init; } = string.Empty;
    public string Type { get; init; } = "APPROVAL";
    public string AssignRole { get; init; } = string.Empty;
    public string AssignTo { get; init; } = string.Empty;
    public int DueHours { get; init; } = 72;
    public int Order { get; init; } = 1;
}

public sealed class WorkflowConditionsConfig
{
    public bool FailFastOnRejection { get; init; } = true;
    public bool RequireAllApprovals { get; init; } = true;
    public int? AutoApproveAfterHours { get; init; }
    public string? EscalateToRole { get; init; }
}
