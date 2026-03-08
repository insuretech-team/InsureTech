using Insuretech.Workflow.Services.V1;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Workflow service gRPC client.
/// </summary>
public sealed class WorkflowGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public WorkflowGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private WorkflowService.WorkflowServiceClient Client =>
        _factory.GetClient("WorkflowService", ch => new WorkflowService.WorkflowServiceClient(ch));

    public async Task<string> StartAsync(
        string workflowType, string entityId, string entityType,
        Dictionary<string, string> metadata,
        CancellationToken ct = default)
    {
        var req = new StartWorkflowRequest
        {
            WorkflowType = workflowType,
            EntityId     = entityId,
            EntityType   = entityType,
        };
        req.Metadata.Add(metadata);
        var resp = await Client.StartWorkflowAsync(req, cancellationToken: ct);
        return resp.WorkflowId;
    }
}
