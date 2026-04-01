using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Workflow.Entity.V1;
using Insuretech.Workflow.Services.V1;
using Microsoft.Extensions.Logging;
using PoliSync.Workflow.Application.Queries;
using System.Text.Json;

namespace PoliSync.Workflow.Infrastructure;

/// <summary>
/// Implements IWorkflowDataGateway by calling the Go workflow-engine
/// via gRPC using WorkflowServiceGrpcClient.
/// </summary>
public sealed class GoWorkflowDataGateway : IWorkflowDataGateway
{
    private readonly WorkflowServiceGrpcClient _client;
    private readonly ILogger<GoWorkflowDataGateway> _logger;

    public GoWorkflowDataGateway(
        WorkflowServiceGrpcClient client,
        ILogger<GoWorkflowDataGateway> logger)
    {
        _client = client;
        _logger = logger;
    }

    // ── Definitions ───────────────────────────────────────────────────────────

    public async Task<string?> CreateDefinitionAsync(
        string name,
        string description,
        string workflowType,
        string entityType,
        string stepsJson,
        string conditionsJson,
        CancellationToken cancellationToken = default)
    {
        try
        {
            // Parse stepsJson into Struct for the proto request
            var stepsStruct = ParseJsonToStruct(stepsJson);

            var resp = await _client.CreateWorkflowDefinitionAsync(
                new CreateWorkflowDefinitionRequest
                {
                    Name = name,
                    Description = description,
                    Type = workflowType,
                    EntityType = entityType,
                    Steps = stepsStruct
                },
                cancellationToken);

            if (!string.IsNullOrEmpty(resp?.Error?.Code))
            {
                _logger.LogError("CreateDefinition failed: {Code} {Message}", resp.Error.Code, resp.Error.Message);
                return null;
            }

            return resp?.WorkflowDefinitionId;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.AlreadyExists)
        {
            // Idempotent: template already registered in the Go engine.
            // The Go service does not expose GetByName, so we cache the name as a
            // stable surrogate key so callers can treat this as a successful registration.
            _logger.LogDebug("Workflow definition '{Name}' already exists — treating as registered", name);
            var surrogateId = $"existing:{name}";
            CacheDefinitionId(name, surrogateId);
            return surrogateId;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "gRPC error creating workflow definition '{Name}'", name);
            return null;
        }
    }

    public async Task<string?> ResolveDefinitionIdByNameAsync(
        string name,
        CancellationToken cancellationToken = default)
    {
        // The Go service doesn't expose GetByName directly — we use GetWorkflowDefinition
        // with a special lookup. For now we attempt GetWorkflowDefinition with the name
        // as ID (will 404). Real implementation: cache the name→ID mapping at startup.
        // This method is called by RegisterWorkflowTemplateCommandHandler for idempotency.
        try
        {
            // Try direct lookup by treating name as a potential cache key
            if (_definitionCache.TryGetValue(name, out var cachedId))
                return cachedId;

            // No built-in name lookup in proto — return null to signal "not found"
            // The register handler will create it and cache the result
            return null;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error resolving definition ID for name '{Name}'", name);
            return null;
        }
    }

    /// <summary>Cache definition name → ID to avoid repeated lookups.</summary>
    internal void CacheDefinitionId(string name, string id)
        => _definitionCache[name] = id;

    private readonly Dictionary<string, string> _definitionCache = new();

    // ── Instances ─────────────────────────────────────────────────────────────

    public async Task<string?> StartWorkflowAsync(
        string definitionId,
        string entityType,
        string entityId,
        Struct? context,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var resp = await _client.StartWorkflowAsync(
                new StartWorkflowRequest
                {
                    WorkflowDefinitionId = definitionId,
                    EntityType = entityType,
                    EntityId = entityId,
                    Context = context ?? new Struct()
                },
                cancellationToken);

            if (!string.IsNullOrEmpty(resp?.Error?.Code))
            {
                _logger.LogError("StartWorkflow failed: {Code} {Message}", resp.Error.Code, resp.Error.Message);
                return null;
            }

            return resp?.WorkflowInstanceId;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "gRPC error starting workflow for {EntityType}/{EntityId}", entityType, entityId);
            return null;
        }
    }

    public async Task<GetWorkflowInstanceResult?> GetWorkflowInstanceAsync(
        string instanceId,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var resp = await _client.GetWorkflowInstanceAsync(
                new GetWorkflowInstanceRequest { WorkflowInstanceId = instanceId },
                cancellationToken);

            if (resp?.WorkflowInstance is null)
                return null;

            return new GetWorkflowInstanceResult(
                resp.WorkflowInstance,
                resp.Tasks.ToList());
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "gRPC error getting workflow instance {InstanceId}", instanceId);
            return null;
        }
    }

    public async Task<IReadOnlyList<WorkflowInstance>> GetWorkflowHistoryAsync(
        string entityType,
        string entityId,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var resp = await _client.GetWorkflowHistoryAsync(
                new GetWorkflowHistoryRequest
                {
                    EntityType = entityType,
                    EntityId = entityId
                },
                cancellationToken);

            return resp?.WorkflowInstances?.ToList() ?? [];
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "gRPC error getting workflow history for {EntityType}/{EntityId}", entityType, entityId);
            return [];
        }
    }

    // ── Tasks ─────────────────────────────────────────────────────────────────

    public async Task<GetMyTasksResult> GetMyTasksAsync(
        string userId,
        string? status,
        int page,
        int pageSize,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var headers = new Metadata { { "x-user-id", userId } };
            var resp = await _client.GetMyTasksAsync(
                new GetMyTasksRequest
                {
                    Status = status ?? string.Empty,
                    Page = page,
                    PageSize = pageSize
                },
                headers: headers,
                cancellationToken);

            return new GetMyTasksResult(
                resp?.Tasks?.ToList() ?? [],
                resp?.TotalCount ?? 0);
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "gRPC error getting tasks for user {UserId}", userId);
            return new GetMyTasksResult([], 0);
        }
    }

    public async Task<CompleteTaskResult> CompleteTaskAsync(
        string taskId,
        string decision,
        string comments,
        string completedBy,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var headers = new Metadata { { "x-user-id", completedBy } };
            var resp = await _client.CompleteTaskAsync(
                new CompleteWorkflowTaskRequest
                {
                    TaskId = taskId,
                    Decision = decision,
                    Comments = comments
                },
                headers: headers,
                cancellationToken);

            if (!string.IsNullOrEmpty(resp?.Error?.Code))
                return new CompleteTaskResult(
                    Success: false,
                    ErrorCode: resp.Error.Code.ToString(),
                    ErrorMessage: resp.Error.Message);

            return new CompleteTaskResult(Success: true);
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return new CompleteTaskResult(false, ErrorCode: "NOT_FOUND", ErrorMessage: ex.Status.Detail);
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "gRPC error completing task {TaskId}", taskId);
            return new CompleteTaskResult(false, ErrorCode: "GRPC_ERROR", ErrorMessage: ex.Status.Detail);
        }
    }

    // ── Helpers ───────────────────────────────────────────────────────────────

    private static Struct ParseJsonToStruct(string json)
    {
        if (string.IsNullOrWhiteSpace(json) || json == "{}" || json == "[]")
            return new Struct();

        // Parse JSON first to determine root kind — do NOT let Struct.Parser.ParseJson
        // handle arrays: protobuf C# silently misparses JSON arrays as Structs, producing
        // {"0":{...},"1":{...}} instead of {"steps":[...]}. Detect arrays explicitly first.
        using var document = JsonDocument.Parse(json);

        if (document.RootElement.ValueKind == JsonValueKind.Array)
        {
            // Steps come from C# as a JSON array: [{name,type,assign_role,...}, ...]
            // Wrap as {"steps":[...]} so the Go workflow engine's parseSteps() can
            // unmarshal it via the wrapper path: struct{ Steps []StepDef }.
            var values = new List<Value>();
            foreach (var item in document.RootElement.EnumerateArray())
            {
                values.Add(ConvertJsonElement(item));
            }

            return new Struct
            {
                Fields =
                {
                    ["steps"] = new Value
                    {
                        ListValue = new ListValue { Values = { values } }
                    }
                }
            };
        }

        // JSON object — use the standard protobuf parser.
        return Struct.Parser.ParseJson(json);
    }

    private static Value ConvertJsonElement(JsonElement element)
        => element.ValueKind switch
        {
            JsonValueKind.Object => new Value { StructValue = ConvertObject(element) },
            JsonValueKind.Array => new Value { ListValue = ConvertArray(element) },
            JsonValueKind.String => new Value { StringValue = element.GetString() ?? string.Empty },
            JsonValueKind.Number => element.TryGetInt64(out var integer)
                ? new Value { NumberValue = integer }
                : new Value { NumberValue = element.GetDouble() },
            JsonValueKind.True => new Value { BoolValue = true },
            JsonValueKind.False => new Value { BoolValue = false },
            JsonValueKind.Null or JsonValueKind.Undefined => new Value { NullValue = NullValue.NullValue },
            _ => new Value { NullValue = NullValue.NullValue }
        };

    private static Struct ConvertObject(JsonElement element)
    {
        var result = new Struct();
        foreach (var property in element.EnumerateObject())
        {
            result.Fields[property.Name] = ConvertJsonElement(property.Value);
        }

        return result;
    }

    private static ListValue ConvertArray(JsonElement element)
    {
        var result = new ListValue();
        foreach (var item in element.EnumerateArray())
        {
            result.Values.Add(ConvertJsonElement(item));
        }

        return result;
    }
}
