using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using PoliSync.Workflow.Domain;

namespace PoliSync.Workflow.Services;

/// <summary>
/// Registers workflow templates from the built-in catalog plus configuration,
/// then refreshes routing rules whenever workflow template options change.
/// </summary>
public sealed class WorkflowTemplateRegistrar
{
    private readonly CompositeWorkflowTemplateProvider _provider;
    private readonly IOptionsMonitor<WorkflowTemplateOptions> _options;
    private readonly ILogger<WorkflowTemplateRegistrar> _logger;

    public WorkflowTemplateRegistrar(
        CompositeWorkflowTemplateProvider provider,
        IOptionsMonitor<WorkflowTemplateOptions> options,
        ILogger<WorkflowTemplateRegistrar> logger)
    {
        _provider = provider;
        _options = options;
        _logger = logger;

        // Hot-reload: re-register config templates when appsettings changes
        _options.OnChange(_ =>
        {
            _logger.LogInformation("WorkflowTemplateOptions changed — reloading config templates...");
            RegisterAll();
        });
    }

    public void RegisterAll()
    {
        RegisterCanonicalTemplates();
        RegisterConfigTemplates();
        RegisterRoutingRules();
        _logger.LogInformation("Workflow templates registered: {Count} total",
            _provider.GetAllTemplates().Count);
    }

    // ── 1. Canonical Templates (always available) ────────────────────────────

    private void RegisterCanonicalTemplates()
    {
        // ── Claims ────────────────────────────────────────────────────────────
        _provider.Register(WorkflowTemplate.TwoStageApproval(
            name: "claim.standard-approval",
            entityType: "CLAIM",
            reviewerRole: "claims_officer",
            approverRole: "claims_manager",
            reviewDueDays: 3,
            approvalDueDays: 2));

        _provider.Register(new WorkflowTemplate
        {
            Name = "claim.high-value-approval",
            EntityType = "CLAIM",
            WorkflowType = "APPROVAL",
            Description = "Three-stage approval for high-value claims (> threshold BDT)",
            Steps =
            [
                new WorkflowStepTemplate { Name = "technical_review",  Type = "REVIEW",    AssignRole = "claims_officer",  DueHours = 48, Order = 1 },
                new WorkflowStepTemplate { Name = "manager_approval",  Type = "APPROVAL",  AssignRole = "claims_manager",  DueHours = 24, Order = 2 },
                new WorkflowStepTemplate { Name = "director_sign_off", Type = "APPROVAL",  AssignRole = "claims_director", DueHours = 24, Order = 3 }
            ],
            Conditions = new WorkflowConditions
            {
                FailFastOnRejection = true,
                RequireAllApprovals = true,
                EscalateToRole = "claims_director"
            }
        });

        // ── Endorsements ──────────────────────────────────────────────────────
        _provider.Register(WorkflowTemplate.SingleApproval(
            name: "endorsement.standard-approval",
            entityType: "ENDORSEMENT",
            assignToRole: "policy_officer",
            dueDays: 3,
            description: "Standard endorsement change approval"));

        _provider.Register(WorkflowTemplate.TwoStageApproval(
            name: "endorsement.major-change-approval",
            entityType: "ENDORSEMENT",
            reviewerRole: "policy_officer",
            approverRole: "underwriting_manager",
            reviewDueDays: 2,
            approvalDueDays: 3));

        // ── Refunds ───────────────────────────────────────────────────────────
        _provider.Register(WorkflowTemplate.SingleApproval(
            name: "refund.standard-approval",
            entityType: "REFUND",
            assignToRole: "finance_officer",
            dueDays: 2,
            description: "Standard refund approval"));

        _provider.Register(WorkflowTemplate.TwoStageApproval(
            name: "refund.high-value-approval",
            entityType: "REFUND",
            reviewerRole: "finance_officer",
            approverRole: "finance_manager",
            reviewDueDays: 2,
            approvalDueDays: 2));

        // ── Underwriting ──────────────────────────────────────────────────────
        _provider.Register(WorkflowTemplate.TwoStageApproval(
            name: "underwriting.manual-review",
            entityType: "UNDERWRITING",
            reviewerRole: "underwriter",
            approverRole: "underwriting_manager",
            reviewDueDays: 5,
            approvalDueDays: 3));

        // ── Policy ────────────────────────────────────────────────────────────
        _provider.Register(WorkflowTemplate.SingleApproval(
            name: "policy.cancellation-approval",
            entityType: "POLICY",
            assignToRole: "policy_manager",
            dueDays: 3,
            description: "Policy cancellation requires manager sign-off"));

        // ── Quotation (B2B) ───────────────────────────────────────────────────
        _provider.Register(WorkflowTemplate.TwoStageApproval(
            name: "quotation.b2b-approval",
            entityType: "QUOTATION",
            reviewerRole: "b2b_officer",
            approverRole: "b2b_manager",
            reviewDueDays: 2,
            approvalDueDays: 2));
    }

    // ── 2. Config Templates (from appsettings — hot-reloadable) ──────────────

    private void RegisterConfigTemplates()
    {
        var opts = _options.CurrentValue;
        foreach (var cfg in opts.AdditionalTemplates)
        {
            if (string.IsNullOrWhiteSpace(cfg.Name) || string.IsNullOrWhiteSpace(cfg.EntityType))
            {
                _logger.LogWarning("Skipping config template with missing Name or EntityType");
                continue;
            }

            _provider.Register(new WorkflowTemplate
            {
                Name = cfg.Name,
                EntityType = cfg.EntityType,
                WorkflowType = cfg.WorkflowType,
                Description = cfg.Description,
                Steps = cfg.Steps.Select(s => new WorkflowStepTemplate
                {
                    Name = s.Name,
                    Type = s.Type,
                    AssignRole = s.AssignRole,
                    AssignTo = s.AssignTo,
                    DueHours = s.DueHours,
                    Order = s.Order
                }).ToList(),
                Conditions = new WorkflowConditions
                {
                    FailFastOnRejection = cfg.Conditions.FailFastOnRejection,
                    RequireAllApprovals = cfg.Conditions.RequireAllApprovals,
                    AutoApproveAfterHours = cfg.Conditions.AutoApproveAfterHours,
                    EscalateToRole = cfg.Conditions.EscalateToRole
                }
            });

            _logger.LogInformation("Loaded config template '{Name}' for {EntityType}", cfg.Name, cfg.EntityType);
        }
    }

    // ── 3. Routing Rules (amount/type/portal based routing) ───────────────────

    private void RegisterRoutingRules()
    {
        // Clear existing rules by re-building via the provider's AddRule
        // (rules are additive — we register once at startup and on hot-reload)
        var opts = _options.CurrentValue;

        // CLAIM: high-value claims get 3-stage approval
        _provider.AddRule(
            entityType: "CLAIM",
            predicate: ctx => ctx.AmountPaisa >= opts.HighValueClaimThresholdPaisa,
            templateName: "claim.high-value-approval");

        // CLAIM: standard approval for all others
        _provider.AddRule(
            entityType: "CLAIM",
            predicate: _ => true,
            templateName: "claim.standard-approval");

        // ENDORSEMENT: major change types → two-stage
        _provider.AddRule(
            entityType: "ENDORSEMENT",
            predicate: ctx => ctx.SubType is not null &&
                              opts.MajorEndorsementTypes.Contains(ctx.SubType, StringComparer.OrdinalIgnoreCase),
            templateName: "endorsement.major-change-approval");

        // ENDORSEMENT: standard for all others
        _provider.AddRule(
            entityType: "ENDORSEMENT",
            predicate: _ => true,
            templateName: "endorsement.standard-approval");

        // REFUND: high-value refunds → two-stage
        _provider.AddRule(
            entityType: "REFUND",
            predicate: ctx => ctx.AmountPaisa >= opts.HighValueRefundThresholdPaisa,
            templateName: "refund.high-value-approval");

        // REFUND: standard for all others
        _provider.AddRule(
            entityType: "REFUND",
            predicate: _ => true,
            templateName: "refund.standard-approval");

        // UNDERWRITING: always manual review
        _provider.AddRule(
            entityType: "UNDERWRITING",
            predicate: _ => true,
            templateName: "underwriting.manual-review");

        // POLICY: cancellation
        _provider.AddRule(
            entityType: "POLICY",
            predicate: ctx => ctx.SubType == "CANCELLATION",
            templateName: "policy.cancellation-approval");

        // QUOTATION: B2B only
        _provider.AddRule(
            entityType: "QUOTATION",
            predicate: ctx => ctx.Portal == "B2B",
            templateName: "quotation.b2b-approval");
    }
}
