using Insuretech.Document.Services.V1;
using Google.Protobuf.WellKnownTypes;

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
            EntityId = entityId,
            EntityType = "generic",
            TemplateId = templateId,
        };
        foreach (var pair in data)
        {
            req.Data.Fields[pair.Key] = Value.ForString(pair.Value);
        }
        var resp = await Client.GenerateDocumentAsync(req, cancellationToken: ct);
        return string.IsNullOrWhiteSpace(resp.FileUrl) ? resp.DocumentId : resp.FileUrl;
    }

    public async Task<Insuretech.Document.Entity.V1.DocumentGeneration?> GetDocumentAsync(
        string documentId,
        CancellationToken ct = default)
    {
        var resp = await Client.GetDocumentAsync(new GetDocumentRequest
        {
            DocumentId = documentId
        }, cancellationToken: ct);

        return resp.Document;
    }
}
