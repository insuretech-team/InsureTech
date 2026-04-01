using Grpc.Core;
using Grpc.Net.Client;
using Insuretech.Workflow.Services.V1;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

namespace PoliSync.Workflow.Infrastructure;

/// <summary>
/// Typed gRPC client for the Go workflow-engine service.
/// Owned by PoliSync.Workflow module — not Infrastructure.GrpcClients (which is excluded from compile).
/// Connects to GrpcClients:WorkflowService (default: http://localhost:50180).
/// </summary>
public sealed class WorkflowServiceGrpcClient : IDisposable
{
    private readonly GrpcChannel _channel;
    private readonly WorkflowService.WorkflowServiceClient _client;
    private readonly ILogger<WorkflowServiceGrpcClient> _logger;

    public WorkflowServiceGrpcClient(
        IConfiguration configuration,
        ILogger<WorkflowServiceGrpcClient> logger)
    {
        _logger = logger;
        var url = configuration["GrpcClients:WorkflowService"] ?? "http://localhost:50180";
        _channel = GrpcChannel.ForAddress(url);
        _client = new WorkflowService.WorkflowServiceClient(_channel);
        _logger.LogInformation("WorkflowServiceGrpcClient connected to {Url}", url);
    }

    // ── Definitions ────────────────────────────────────────────────────────────

    public Task<CreateWorkflowDefinitionResponse> CreateWorkflowDefinitionAsync(
        CreateWorkflowDefinitionRequest request,
        CancellationToken ct = default)
        => _client.CreateWorkflowDefinitionAsync(request, cancellationToken: ct).ResponseAsync;

    public Task<GetWorkflowDefinitionResponse> GetWorkflowDefinitionAsync(
        GetWorkflowDefinitionRequest request,
        CancellationToken ct = default)
        => _client.GetWorkflowDefinitionAsync(request, cancellationToken: ct).ResponseAsync;

    // ── Instances ──────────────────────────────────────────────────────────────

    public Task<StartWorkflowResponse> StartWorkflowAsync(
        StartWorkflowRequest request,
        CancellationToken ct = default)
        => _client.StartWorkflowAsync(request, cancellationToken: ct).ResponseAsync;

    public Task<GetWorkflowInstanceResponse> GetWorkflowInstanceAsync(
        GetWorkflowInstanceRequest request,
        CancellationToken ct = default)
        => _client.GetWorkflowInstanceAsync(request, cancellationToken: ct).ResponseAsync;

    public Task<GetWorkflowHistoryResponse> GetWorkflowHistoryAsync(
        GetWorkflowHistoryRequest request,
        CancellationToken ct = default)
        => _client.GetWorkflowHistoryAsync(request, cancellationToken: ct).ResponseAsync;

    // ── Tasks ──────────────────────────────────────────────────────────────────

    public Task<GetMyTasksResponse> GetMyTasksAsync(
        GetMyTasksRequest request,
        Metadata? headers = null,
        CancellationToken ct = default)
        => _client.GetMyTasksAsync(request, headers: headers, cancellationToken: ct).ResponseAsync;

    public Task<CompleteWorkflowTaskResponse> CompleteTaskAsync(
        CompleteWorkflowTaskRequest request,
        Metadata? headers = null,
        CancellationToken ct = default)
        => _client.CompleteTaskAsync(request, headers: headers, cancellationToken: ct).ResponseAsync;

    public void Dispose() => _channel.Dispose();
}
