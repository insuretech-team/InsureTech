using Google.Protobuf.WellKnownTypes;
using InsuranceEngine.Grpc.Clients;
using InsuranceEngine.SharedKernel.Infrastructure;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Infrastructure.Documents;

public interface IDocumentService
{
    Task<string> GeneratePolicyDocumentAsync(
        string policyId,
        string policyNumber,
        string customerName,
        string productName,
        decimal premium,
        string startDate,
        string endDate,
        bool includeQrCode = true,
        CancellationToken ct = default);

    Task<string> GenerateClaimDocumentAsync(
        string claimId,
        string policyNumber,
        string claimType,
        decimal claimAmount,
        CancellationToken ct = default);

    Task<byte[]> DownloadDocumentAsync(string documentId, CancellationToken ct = default);

    Task<List<DocumentInfo>> ListDocumentsForEntityAsync(
        string entityType,
        string entityId,
        CancellationToken ct = default);

    Task DeleteDocumentAsync(string documentId, CancellationToken ct = default);
}

public class DocumentInfo
{
    public string DocumentId { get; set; } = string.Empty;
    public string EntityType { get; set; } = string.Empty;
    public string EntityId { get; set; } = string.Empty;
    public string DocumentType { get; set; } = string.Empty;
    public string FileUrl { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; }
}

public class DocumentService : IDocumentService
{
    private readonly InsuranceServiceClient _client;
    private readonly IQrCodeService _qrCodeService;
    private readonly ILogger<DocumentService> _logger;

    private const string PolicyDocumentTemplateId = "policy-document-v1";
    private const string ClaimDocumentTemplateId = "claim-receipt-v1";

    public DocumentService(
        InsuranceServiceClient client,
        IQrCodeService qrCodeService,
        ILogger<DocumentService> logger)
    {
        _client = client;
        _qrCodeService = qrCodeService;
        _logger = logger;
    }

    public async Task<string> GeneratePolicyDocumentAsync(
        string policyId,
        string policyNumber,
        string customerName,
        string productName,
        decimal premium,
        string startDate,
        string endDate,
        bool includeQrCode = true,
        CancellationToken ct = default)
    {
        try
        {
            var qrUrl = $"https://insuretech.labaid.com/verify/{policyNumber}";
            var qrBase64 = _qrCodeService.GenerateQrCodeBase64(qrUrl);

            var request = new Insuretech.Document.Services.V1.GenerateDocumentRequest
            {
                TemplateId = PolicyDocumentTemplateId,
                EntityType = "policy",
                EntityId = policyId,
                IncludeQrCode = includeQrCode,
                OutputFormat = "pdf",
                Data = new Struct
                {
                    Fields =
                    {
                        ["policy_number"] = Value.ForString(policyNumber),
                        ["customer_name"] = Value.ForString(customerName),
                        ["product_name"] = Value.ForString(productName),
                        ["premium_amount"] = Value.ForNumber((double)premium),
                        ["premium_text"] = Value.ForString($"{premium:N2} BDT"),
                        ["start_date"] = Value.ForString(startDate),
                        ["end_date"] = Value.ForString(endDate),
                        ["issued_date"] = Value.ForString(DateTime.UtcNow.ToString("dd MMM yyyy")),
                        ["qr_code_data"] = Value.ForString(qrBase64),
                        ["qr_verification_url"] = Value.ForString(qrUrl)
                    }
                }
            };

            var response = await _client.Documents.GenerateDocumentAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to generate policy document: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Document generation failed: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Policy document generated: {DocumentId} for policy {PolicyNumber}",
                response.DocumentId, policyNumber);

            return response.DocumentId;
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error generating policy document for {PolicyNumber}", policyNumber);
            throw;
        }
    }

    public async Task<string> GenerateClaimDocumentAsync(
        string claimId,
        string policyNumber,
        string claimType,
        decimal claimAmount,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Document.Services.V1.GenerateDocumentRequest
            {
                TemplateId = ClaimDocumentTemplateId,
                EntityType = "claim",
                EntityId = claimId,
                IncludeQrCode = false,
                OutputFormat = "pdf",
                Data = new Struct
                {
                    Fields =
                    {
                        ["claim_id"] = Value.ForString(claimId),
                        ["policy_number"] = Value.ForString(policyNumber),
                        ["claim_type"] = Value.ForString(claimType),
                        ["claim_amount"] = Value.ForNumber((double)claimAmount),
                        ["claim_amount_text"] = Value.ForString($"{claimAmount:N2} BDT"),
                        ["claim_date"] = Value.ForString(DateTime.UtcNow.ToString("dd MMM yyyy")),
                        ["received_date"] = Value.ForString(DateTime.UtcNow.ToString("dd MMM yyyy HH:mm"))
                    }
                }
            };

            var response = await _client.Documents.GenerateDocumentAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to generate claim document: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Document generation failed: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Claim document generated: {DocumentId} for claim {ClaimId}",
                response.DocumentId, claimId);

            return response.DocumentId;
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error generating claim document for {ClaimId}", claimId);
            throw;
        }
    }

    public async Task<byte[]> DownloadDocumentAsync(string documentId, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Document.Services.V1.DownloadDocumentRequest
            {
                DocumentId = documentId
            };

            var response = await _client.Documents.DownloadDocumentAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to download document: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Document download failed: {response.Error.Message}");
            }

            return response.Content.ToByteArray();
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error downloading document {DocumentId}", documentId);
            throw;
        }
    }

    public async Task<List<DocumentInfo>> ListDocumentsForEntityAsync(
        string entityType,
        string entityId,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Document.Services.V1.ListDocumentsRequest
            {
                EntityType = entityType,
                EntityId = entityId,
                Page = 1,
                PageSize = 100
            };

            var response = await _client.Documents.ListDocumentsAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to list documents: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Document listing failed: {response.Error.Message}");
            }

            return response.Documents.Select(d => new DocumentInfo
            {
                DocumentId = d.Id,
                EntityType = d.EntityType,
                EntityId = d.EntityId,
                DocumentType = d.DocumentTemplateId,
                FileUrl = d.FileUrl,
                Status = d.Status.ToString(),
                CreatedAt = d.GeneratedAt?.ToDateTime() ?? DateTime.UtcNow
            }).ToList();
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error listing documents for {EntityType}/{EntityId}", entityType, entityId);
            throw;
        }
    }

    public async Task DeleteDocumentAsync(string documentId, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Document.Services.V1.DeleteDocumentRequest
            {
                DocumentId = documentId
            };

            var response = await _client.Documents.DeleteDocumentAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to delete document: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Document deletion failed: {response.Error.Message}");
            }

            _logger.LogInformation("Document deleted: {DocumentId}", documentId);
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error deleting document {DocumentId}", documentId);
            throw;
        }
    }
}

public class MockDocumentService : IDocumentService
{
    private readonly IQrCodeService _qrCodeService;
    private readonly ILogger<MockDocumentService> _logger;

    public MockDocumentService(IQrCodeService qrCodeService, ILogger<MockDocumentService> logger)
    {
        _qrCodeService = qrCodeService;
        _logger = logger;
    }

    public Task<string> GeneratePolicyDocumentAsync(
        string policyId,
        string policyNumber,
        string customerName,
        string productName,
        decimal premium,
        string startDate,
        string endDate,
        bool includeQrCode = true,
        CancellationToken ct = default)
    {
        var documentId = Guid.NewGuid().ToString();
        var qrUrl = $"https://insuretech.labaid.com/verify/{policyNumber}";
        var qrBase64 = _qrCodeService.GenerateQrCodeBase64(qrUrl);

        _logger.LogInformation(
            "[MOCK] Policy document generated: {DocumentId} for policy {PolicyNumber} (QR: {QrUrl})",
            documentId, policyNumber, qrUrl);

        var docContent = $"--- LABAID INSURETECH POLICY ---\n" +
                         $"Policy Number: {policyNumber}\n" +
                         $"Customer: {customerName}\n" +
                         $"Product: {productName}\n" +
                         $"Premium: {premium:N2} BDT\n" +
                         $"Start Date: {startDate}\n" +
                         $"End Date: {endDate}\n" +
                         $"QR Verification: {qrUrl}\n" +
                         $"QR Base64: [TRUNCATED]";

        _logger.LogDebug("[MOCK] Document content:\n{Content}", docContent);

        return Task.FromResult(documentId);
    }

    public Task<string> GenerateClaimDocumentAsync(
        string claimId,
        string policyNumber,
        string claimType,
        decimal claimAmount,
        CancellationToken ct = default)
    {
        var documentId = Guid.NewGuid().ToString();

        _logger.LogInformation(
            "[MOCK] Claim document generated: {DocumentId} for claim {ClaimId}",
            documentId, claimId);

        return Task.FromResult(documentId);
    }

    public Task<byte[]> DownloadDocumentAsync(string documentId, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Document downloaded: {DocumentId}", documentId);
        var content = $"[MOCK PDF CONTENT FOR DOCUMENT {documentId}]";
        return Task.FromResult(System.Text.Encoding.UTF8.GetBytes(content));
    }

    public Task<List<DocumentInfo>> ListDocumentsForEntityAsync(
        string entityType,
        string entityId,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Listed documents for {EntityType}/{EntityId}",
            entityType, entityId);

        return Task.FromResult(new List<DocumentInfo>
        {
            new()
            {
                DocumentId = Guid.NewGuid().ToString(),
                EntityType = entityType,
                EntityId = entityId,
                DocumentType = "policy_document",
                Status = "completed",
                CreatedAt = DateTime.UtcNow.AddDays(-7)
            }
        });
    }

    public Task DeleteDocumentAsync(string documentId, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Document deleted: {DocumentId}", documentId);
        return Task.CompletedTask;
    }
}

public class GoDocumentPdfGenerator : IPdfGenerator
{
    private readonly IDocumentService _documentService;
    private readonly ILogger<GoDocumentPdfGenerator> _logger;

    public GoDocumentPdfGenerator(IDocumentService documentService, ILogger<GoDocumentPdfGenerator> logger)
    {
        _documentService = documentService;
        _logger = logger;
    }

    public async Task<byte[]> GeneratePolicyDocumentAsync(
        string policyNumber,
        string customerName,
        string productName,
        decimal premium)
    {
        try
        {
            var documentId = await _documentService.GeneratePolicyDocumentAsync(
                policyId: policyNumber,
                policyNumber: policyNumber,
                customerName: customerName,
                productName: productName,
                premium: premium,
                startDate: DateTime.UtcNow.ToString("dd MMM yyyy"),
                endDate: DateTime.UtcNow.AddYears(1).ToString("dd MMM yyyy"),
                includeQrCode: true);

            _logger.LogInformation(
                "Generated policy document: {DocumentId} for {PolicyNumber}",
                documentId, policyNumber);

            var bytes = await _documentService.DownloadDocumentAsync(documentId);
            return bytes;
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Failed to generate real PDF, falling back to mock");
            throw;
        }
    }
}
