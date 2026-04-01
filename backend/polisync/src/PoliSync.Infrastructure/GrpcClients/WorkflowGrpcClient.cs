using Google.Protobuf.WellKnownTypes;
using Insuretech.Workflow.Services.V1;
using Grpc.Core;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go workflow-engine gRPC service.
/// Exposes all WorkflowService operations used by PoliSync.
/// </summary>
public sealed class WorkflowGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public WorkflowGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private WorkflowService.WorkflowServiceClient Client =>
        _factory.GetClient("WorkflowService", ch => new WorkflowService.WorkflowServiceClient(ch));

    // ── Definitions ────────────────────────────────────────────────────────────

    public Task<CreateWorkflowDefinitionResponse> CreateWorkflowDefinitionAsync(
        CreateWorkflowDefinitionRequest request,
        CancellationToken ct = default)
        => Client.CreateWorkflowDefinitionAsync(request, cancellationToken: ct).ResponseAsync;

    public Task<GetWorkflowDefinitionResponse> GetWorkflowDefinitionAsync(
        GetWorkflowDefinitionRequest request,
        CancellationToken ct = default)
        => Client.GetWorkflowDefinitionAsync(request, cancellationToken: ct).ResponseAsync;

    // ── Instances ──────────────────────────────────────────────────────────────

    public Task<StartWorkflowResponse> StartWorkflowAsync(
        StartWorkflowRequest request,
        CancellationToken ct = default)
        => Client.StartWorkflowAsync(request, cancellationToken: ct).ResponseAsync;

    public Task<GetWorkflowInstanceResponse> GetWorkflowInstanceAsync(
        GetWorkflowInstanceRequest request,
        CancellationToken ct = default)
        => Client.GetWorkflowInstanceAsync(request, cancellationToken: ct).ResponseAsync;

    public Task<GetWorkflowHistoryResponse> GetWorkflowHistoryAsync(
        GetWorkflowHistoryRequest request,
        CancellationToken ct = default)
        => Client.GetWorkflowHistoryAsync(request, cancellationToken: ct).ResponseAsync;

    // ── Tasks ──────────────────────────────────────────────────────────────────

    public Task<GetMyTasksResponse> GetMyTasksAsync(
        GetMyTasksRequest request,
        Metadata? headers = null,
        CancellationToken ct = default)
        => Client.GetMyTasksAsync(request, headers: headers, cancellationToken: ct).ResponseAsync;

    public Task<CompleteWorkflowTaskResponse> CompleteTaskAsync(
        CompleteWorkflowTaskRequest request,
        Metadata? headers = null,
        CancellationToken ct = default)
        => Client.CompleteTaskAsync(request, headers: headers, cancellationToken: ct).ResponseAsync;
}
