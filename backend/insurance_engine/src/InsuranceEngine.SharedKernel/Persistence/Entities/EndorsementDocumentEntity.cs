namespace InsuranceEngine.SharedKernel.Persistence.Entities;

public class EndorsementDocumentEntity
{
    public string DocumentId { get; set; } = string.Empty;
    public string EndorsementId { get; set; } = string.Empty;
    public string DocumentType { get; set; } = string.Empty;
    public string DocumentNumber { get; set; } = string.Empty;
    public string FileUrl { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public DateTime GeneratedAt { get; set; }
    public DateTime CreatedAt { get; set; }
    
    public EndorsementEntity? Endorsement { get; set; }
}
