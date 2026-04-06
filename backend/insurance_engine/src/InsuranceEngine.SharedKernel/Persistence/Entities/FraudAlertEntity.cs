namespace InsuranceEngine.SharedKernel.Persistence.Entities;

public class FraudAlertEntity
{
    public string AlertId { get; set; } = string.Empty;
    public string AlertNumber { get; set; } = string.Empty;
    public string EntityType { get; set; } = string.Empty;
    public string EntityId { get; set; } = string.Empty;
    public string AlertType { get; set; } = string.Empty;
    public string Severity { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public decimal FraudScore { get; set; }
    public string RecommendedAction { get; set; } = string.Empty;
    public string ResolvedBy { get; set; } = string.Empty;
    public string ResolutionNotes { get; set; } = string.Empty;
    public DateTime? ResolvedAt { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
