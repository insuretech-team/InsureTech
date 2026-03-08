using Insuretech.Storage.Service.V1;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Storage service gRPC client.
/// </summary>
public sealed class StorageGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public StorageGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private StorageService.StorageServiceClient Client =>
        _factory.GetClient("StorageService", ch => new StorageService.StorageServiceClient(ch));

    public async Task<GetFileResponse> GetFileAsync(string fileId, CancellationToken ct = default)
        => await Client.GetFileAsync(new GetFileRequest { FileId = fileId }, cancellationToken: ct);
}
