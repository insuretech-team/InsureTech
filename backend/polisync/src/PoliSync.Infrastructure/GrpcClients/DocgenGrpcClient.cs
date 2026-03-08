using Insuretech.Document.Services.V1;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Document generation service gRPC client.
/// </summary>
public sealed class DocgenGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public DocgenGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private DocumentService.DocumentServiceClient Client =>
        _factory.GetClient("DocgenService", ch => new DocumentService.DocumentServiceClient(ch));

    public async Task<string> GenerateAsync(
        string entityId, string templateId,
        Dictionary<string, string> data,
        CancellationToken ct = default)
    {
        var req = new GenerateDocumentRequest
        {
            EntityId   = entityId,
            TemplateId = templateId,
        };
        req.Data.Add(data);
        var resp = await Client.GenerateDocumentAsync(req, cancellationToken: ct);
        return resp.DocumentUrl;
    }
}
