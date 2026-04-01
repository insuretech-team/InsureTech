using System.Text.Json;
using Confluent.Kafka;
using PoliSync.ApiHost.Services;

namespace PoliSync.ApiHost.BackgroundServices;

public sealed class InsuranceProposalDecisionConsumer : BackgroundService
{
    private readonly ILogger<InsuranceProposalDecisionConsumer> _logger;
    private readonly IServiceScopeFactory _scopeFactory;
    private readonly IConsumer<Ignore, string> _consumer;
    private readonly string _approvedTopic;
    private readonly string _rejectedTopic;

    public InsuranceProposalDecisionConsumer(
        IConfiguration configuration,
        ILogger<InsuranceProposalDecisionConsumer> logger,
        IServiceScopeFactory scopeFactory)
    {
        _logger = logger;
        _scopeFactory = scopeFactory;

        _approvedTopic = configuration["Kafka:Topics:InsurerProposalApprovedInbound"] ?? "insuretech.proposal.approved.v1";
        _rejectedTopic = configuration["Kafka:Topics:InsurerProposalRejectedInbound"] ?? "insuretech.proposal.rejected.v1";
        var bootstrapServers = configuration["Kafka:BootstrapServers"] ?? "localhost:9092";
        var groupId = configuration["Kafka:Consumer:InsuranceProposalDecision:GroupId"] ?? "polisync-insurance-proposal-decision";

        var consumerConfig = new ConsumerConfig
        {
            BootstrapServers = bootstrapServers,
            GroupId = groupId,
            AutoOffsetReset = AutoOffsetReset.Earliest,
            EnableAutoCommit = false
        };

        _consumer = new ConsumerBuilder<Ignore, string>(consumerConfig).Build();
    }

    protected override Task ExecuteAsync(CancellationToken stoppingToken)
    {
        _consumer.Subscribe([_approvedTopic, _rejectedTopic]);
        _logger.LogInformation(
            "Subscribed to insurer proposal decision topics {ApprovedTopic} and {RejectedTopic}",
            _approvedTopic,
            _rejectedTopic);

        return Task.Run(async () =>
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                ConsumeResult<Ignore, string>? result = null;
                try
                {
                    result = _consumer.Consume(stoppingToken);
                }
                catch (OperationCanceledException)
                {
                    break;
                }
                catch (ConsumeException ex)
                {
                    _logger.LogError(ex, "Kafka consume error on insurer proposal decision topics");
                    continue;
                }

                if (result is null)
                {
                    continue;
                }

                var processed = await TryProcessMessageAsync(result.Topic, result.Message.Value, stoppingToken);
                if (processed)
                {
                    _consumer.Commit(result);
                }
            }
        }, stoppingToken);
    }

    public override void Dispose()
    {
        _consumer.Close();
        _consumer.Dispose();
        base.Dispose();
    }

    private async Task<bool> TryProcessMessageAsync(
        string topic,
        string payload,
        CancellationToken cancellationToken)
    {
        if (string.IsNullOrWhiteSpace(payload))
        {
            return true;
        }

        ProposalDecisionPayload? evt;
        try
        {
            using var doc = JsonDocument.Parse(payload);
            evt = ParsePayload(doc.RootElement, payload);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to parse insurer proposal decision payload from topic {Topic}", topic);
            return false;
        }

        if (evt is null || string.IsNullOrWhiteSpace(evt.ProposalId))
        {
            _logger.LogWarning("Proposal decision event on topic {Topic} is missing proposal_id/proposalId. Payload ignored.", topic);
            return true;
        }

        try
        {
            using var scope = _scopeFactory.CreateScope();
            var workflowService = scope.ServiceProvider.GetRequiredService<InsuranceProposalWorkflowService>();
            var reviewedByUserId = FirstNonEmpty(evt.ReviewedByUserId, evt.ActorUserId, "insurer-event")!;

            if (string.Equals(topic, _approvedTopic, StringComparison.Ordinal))
            {
                await workflowService.ApproveProposalAsync(
                    evt.ProposalId,
                    reviewedByUserId,
                    evt.InsurerResponsePayload,
                    evt.DecisionReason,
                    cancellationToken);

                _logger.LogInformation("Processed insurer approval event for proposal {ProposalId}", evt.ProposalId);
                return true;
            }

            if (string.Equals(topic, _rejectedTopic, StringComparison.Ordinal))
            {
                await workflowService.RejectProposalAsync(
                    evt.ProposalId,
                    reviewedByUserId,
                    evt.InsurerResponsePayload,
                    evt.DecisionReason,
                    cancellationToken);

                _logger.LogInformation("Processed insurer rejection event for proposal {ProposalId}", evt.ProposalId);
                return true;
            }

            _logger.LogWarning("Received proposal decision message on unexpected topic {Topic}", topic);
            return true;
        }
        catch (KeyNotFoundException ex)
        {
            _logger.LogWarning(ex, "Proposal decision ignored because the target proposal was not found");
            return true;
        }
        catch (InvalidOperationException ex)
        {
            _logger.LogWarning(ex, "Proposal decision ignored because the proposal was already closed or invalid for transition");
            return true;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to process insurer proposal decision for proposal {ProposalId}", evt.ProposalId);
            return false;
        }
    }

    private static ProposalDecisionPayload ParsePayload(JsonElement root, string rawPayload)
        => new(
            ProposalId: FirstString(root, "proposal_id", "proposalId"),
            ReviewedByUserId: FirstString(root, "reviewed_by_user_id", "reviewedByUserId"),
            ActorUserId: FirstString(root, "actor_user_id", "actorUserId", "user_id", "userId"),
            DecisionReason: FirstString(root, "decision_reason", "decisionReason", "rejection_reason", "rejectionReason", "reason"),
            InsurerResponsePayload: FirstString(root, "insurer_response_payload", "insurerResponsePayload") ?? rawPayload);

    private static string? FirstString(JsonElement root, params string[] names)
    {
        foreach (var name in names)
        {
            if (root.TryGetProperty(name, out var element) && element.ValueKind == JsonValueKind.String)
            {
                return element.GetString();
            }
        }

        return null;
    }

    private static string? FirstNonEmpty(params string?[] values)
        => values.FirstOrDefault(value => !string.IsNullOrWhiteSpace(value));

    private sealed record ProposalDecisionPayload(
        string? ProposalId,
        string? ReviewedByUserId,
        string? ActorUserId,
        string? DecisionReason,
        string? InsurerResponsePayload);
}
