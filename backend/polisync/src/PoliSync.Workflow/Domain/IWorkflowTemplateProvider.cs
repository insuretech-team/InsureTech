namespace PoliSync.Workflow.Domain;

/// <summary>
/// Provides workflow templates for entity types and conditions.
/// Implementations can load templates from: code, JSON config, DB, or remote config.
///
/// The provider is the single source of truth for WHICH template to use for a
/// given entity — callers don't need to know template names.
/// </summary>
public interface IWorkflowTemplateProvider
{
    /// <summary>
    /// Resolves the appropriate workflow template for an entity based on context.
    /// Returns null if no workflow should be triggered for this entity/context.
    /// </summary>
    WorkflowTemplate? Resolve(WorkflowTriggerContext context);

    /// <summary>
    /// Returns all registered templates (for seeding and admin listing).
    /// </summary>
    IReadOnlyList<WorkflowTemplate> GetAllTemplates();

    /// <summary>
    /// Registers a template dynamically at runtime.
    /// Overwrites if a template with the same name already exists.
    /// </summary>
    void Register(WorkflowTemplate template);

    /// <summary>
    /// Removes a template by name (admin / hot-reload use).
    /// </summary>
    bool Remove(string name);
}

/// <summary>
/// Context passed to IWorkflowTemplateProvider.Resolve().
/// Carries all the information needed to select the right template.
/// </summary>
public sealed record WorkflowTriggerContext
{
    /// <summary>Domain entity type: CLAIM, ENDORSEMENT, REFUND, POLICY, UNDERWRITING, QUOTATION</summary>
    public required string EntityType { get; init; }

    /// <summary>UUID of the entity.</summary>
    public required string EntityId { get; init; }

    /// <summary>Who triggered this workflow (user UUID).</summary>
    public required string InitiatedBy { get; init; }

    /// <summary>Monetary amount in paisa (BDT). Used for value-based template routing.</summary>
    public long AmountPaisa { get; init; }

    /// <summary>Sub-type of the entity (e.g. "SumAssuredChange" for endorsements).</summary>
    public string? SubType { get; init; }

    /// <summary>Portal context: B2B, B2C, PARTNER, SYSTEM.</summary>
    public string Portal { get; init; } = "SYSTEM";

    /// <summary>Additional key-value metadata passed through to workflow context.</summary>
    public Dictionary<string, string> Metadata { get; init; } = [];
}
