using System.Collections.Concurrent;
using Microsoft.Extensions.Logging;

namespace PoliSync.Workflow.Domain;

/// <summary>
/// Default implementation of IWorkflowTemplateProvider.
///
/// Supports THREE sources of templates, applied in priority order:
///   1. Runtime-registered templates (highest priority — admin/hot-reload)
///   2. JSON config templates (appsettings.json WorkflowTemplates section)
///   3. Code-defined templates (lowest priority — always available)
///
/// Template resolution uses routing rules:
///   - EntityType must match exactly
///   - AmountPaisa thresholds select standard vs high-value templates
///   - SubType can route to specialised templates (e.g. major endorsement)
///   - Portal can route to B2B-specific templates
///
/// Thread-safe via ConcurrentDictionary — safe for hot-reload without restart.
/// </summary>
public sealed class CompositeWorkflowTemplateProvider : IWorkflowTemplateProvider
{
    // name → template registry (thread-safe for hot-reload)
    private readonly ConcurrentDictionary<string, WorkflowTemplate> _registry = new();

    // routing rules: (entityType, predicate) → template name
    private readonly List<RoutingRule> _rules = [];

    private readonly object _rulesLock = new();
    private readonly ILogger<CompositeWorkflowTemplateProvider> _logger;

    public CompositeWorkflowTemplateProvider(ILogger<CompositeWorkflowTemplateProvider> logger)
    {
        _logger = logger;
    }

    // ── Registration ──────────────────────────────────────────────────────────

    public void Register(WorkflowTemplate template)
    {
        _registry[template.Name] = template;
        _logger.LogDebug("Registered workflow template '{Name}' for {EntityType}", template.Name, template.EntityType);
    }

    public bool Remove(string name)
    {
        var removed = _registry.TryRemove(name, out _);
        if (removed) _logger.LogInformation("Removed workflow template '{Name}'", name);
        return removed;
    }

    /// <summary>
    /// Adds a routing rule that maps a context predicate to a template name.
    /// Rules are evaluated in the order they are added — first match wins.
    /// </summary>
    public void AddRule(string entityType, Func<WorkflowTriggerContext, bool> predicate, string templateName)
    {
        lock (_rulesLock)
        {
            _rules.Add(new RoutingRule(entityType.ToUpperInvariant(), predicate, templateName));
        }
    }

    // ── Resolution ────────────────────────────────────────────────────────────

    public WorkflowTemplate? Resolve(WorkflowTriggerContext context)
    {
        var entityType = context.EntityType.ToUpperInvariant();

        List<RoutingRule> snapshot;
        lock (_rulesLock) { snapshot = [.._rules]; }

        // Evaluate rules in order — first match wins
        foreach (var rule in snapshot)
        {
            if (rule.EntityType != entityType) continue;
            if (!rule.Predicate(context)) continue;
            if (_registry.TryGetValue(rule.TemplateName, out var template))
            {
                _logger.LogDebug("Resolved template '{Name}' for {EntityType}/{EntityId}",
                    rule.TemplateName, context.EntityType, context.EntityId);
                return template;
            }
            _logger.LogWarning("Rule matched template '{Name}' but it is not registered", rule.TemplateName);
        }

        // Fallback: find any active template for this entity type
        var fallback = _registry.Values
            .Where(t => t.EntityType.Equals(entityType, StringComparison.OrdinalIgnoreCase))
            .MinBy(t => t.Name); // deterministic fallback

        if (fallback is not null)
        {
            _logger.LogDebug("Using fallback template '{Name}' for {EntityType}", fallback.Name, entityType);
        }
        else
        {
            _logger.LogDebug("No workflow template found for entity type '{EntityType}' — skipping workflow", entityType);
        }

        return fallback;
    }

    public IReadOnlyList<WorkflowTemplate> GetAllTemplates()
        => [.._registry.Values.OrderBy(t => t.EntityType).ThenBy(t => t.Name)];

    // ── Inner types ───────────────────────────────────────────────────────────

    private sealed record RoutingRule(
        string EntityType,
        Func<WorkflowTriggerContext, bool> Predicate,
        string TemplateName);
}
