using PoliSync.SharedKernel.Domain;
using Insuretech.Claims.Entity.V1;

namespace PoliSync.Claims.Domain;

/// <summary>
/// Claim aggregate - implements 4-tier approval matrix and fraud detection
/// </summary>
public class ClaimAggregate
{
    private readonly Claim _claim;
    private readonly List<DomainEvent> _domainEvents = new();
    
    // Approval thresholds in paisa (BDT minor units)
    private const long ZHTC_THRESHOLD = 1_000_000;      // 10,000 BDT
    private const long L1_THRESHOLD = 5_000_000;        // 50,000 BDT
    private const long L2_THRESHOLD = 20_000_000;       // 200,000 BDT
    private const long L3_THRESHOLD = 50_000_000;       // 500,000 BDT
    
    private const double AUTO_APPROVE_FRAUD_SCORE = 0.30;
    private const double FRAUD_FLAG_THRESHOLD = 0.75;
    
    public ClaimAggregate(Claim claim)
    {
        _claim = claim ?? throw new ArgumentNullException(nameof(claim));
    }
    
    public Claim Claim => _claim;
    public string ClaimId => _claim.ClaimId;
    public string ClaimNumber => _claim.ClaimNumber;
    public ClaimStatus Status => _claim.Status;
    public IReadOnlyCollection<DomainEvent> DomainEvents => _domainEvents.AsReadOnly();
    
    public void ClearDomainEvents() => _domainEvents.Clear();
    
    // Factory method - FNOL (First Notice of Loss)
    public static ClaimAggregate FileClaim(
        string policyId,
        string customerId,
        ClaimType claimType,
        long claimedAmountPaisa,
        DateTime incidentDate,
        string incidentDescription,
        string placeOfIncident)
    {
        var claim = new Claim
        {
            ClaimId = Guid.NewGuid().ToString(),
            ClaimNumber = GenerateClaimNumber(),
            PolicyId = policyId,
            CustomerId = customerId,
            Type = claimType,
            Status = ClaimStatus.Submitted,
            ClaimedAmount = new Insuretech.Common.V1.Money 
            { 
                Amount = claimedAmountPaisa,
                Currency = "BDT"
            },
            ApprovedAmount = new Insuretech.Common.V1.Money { Amount = 0, Currency = "BDT" },
            SettledAmount = new Insuretech.Common.V1.Money { Amount = 0, Currency = "BDT" },
            IncidentDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(incidentDate.ToUniversalTime()),
            IncidentDescription = incidentDescription,
            PlaceOfIncident = placeOfIncident,
            SubmittedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            ProcessingType = ClaimProcessingType.Manual
        };
        
        var aggregate = new ClaimAggregate(claim);
        aggregate._domainEvents.Add(new ClaimFiledEvent(claim.ClaimId, claim.PolicyId, claim.CustomerId));
        
        return aggregate;
    }
    
    // Fraud check integration
    public void ApplyFraudCheck(double fraudScore, List<string> fraudFlags)
    {
        var fraudCheck = new FraudCheckResult
        {
            FraudCheckId = Guid.NewGuid().ToString(),
            ClaimId = ClaimId,
            FraudScore = fraudScore,
            CheckedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
        };
        
        fraudCheck.FraudFlags.AddRange(fraudFlags);
        
        _claim.FraudCheck = fraudCheck;
        
        if (fraudScore > FRAUD_FLAG_THRESHOLD)
        {
            _claim.Status = ClaimStatus.FraudCheck;
            _domainEvents.Add(new ClaimFlaggedForFraudEvent(ClaimId, fraudScore));
        }
        else if (fraudScore < AUTO_APPROVE_FRAUD_SCORE && _claim.ClaimedAmount.Amount <= ZHTC_THRESHOLD)
        {
            // Auto-approve for ZHTC
            _claim.ProcessingType = ClaimProcessingType.Automated;
        }
        
        _claim.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
    }
    
    // 4-tier approval matrix
    public int GetRequiredApprovalLevel()
    {
        var amount = _claim.ClaimedAmount.Amount;
        
        if (amount <= ZHTC_THRESHOLD)
            return 0; // ZHTC - can be auto-approved
        else if (amount <= L1_THRESHOLD)
            return 1; // L1 - Claims Officer
        else if (amount <= L2_THRESHOLD)
            return 2; // L2 - Claims Manager
        else if (amount <= L3_THRESHOLD)
            return 3; // L3 - Director
        else
            return 4; // Board approval required
    }
    
    public void AddApproval(
        string approverId,
        string approverRole,
        int approvalLevel,
        ApprovalDecision decision,
        long approvedAmountPaisa,
        string notes)
    {
        var approval = new ClaimApproval
        {
            ApprovalId = Guid.NewGuid().ToString(),
            ClaimId = ClaimId,
            ApproverId = approverId,
            ApproverRole = approverRole,
            ApprovalLevel = approvalLevel,
            Decision = decision,
            ApprovedAmount = new Insuretech.Common.V1.Money 
            { 
                Amount = approvedAmountPaisa,
                Currency = "BDT"
            },
            Notes = notes,
            ApprovedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
        };
        
        _claim.Approvals.Add(approval);
        
        if (decision == ApprovalDecision.Approved)
        {
            _claim.ApprovedAmount = new Insuretech.Common.V1.Money 
            { 
                Amount = approvedAmountPaisa,
                Currency = "BDT"
            };
            
            var requiredLevel = GetRequiredApprovalLevel();
            if (approvalLevel >= requiredLevel)
            {
                _claim.Status = ClaimStatus.Approved;
                _claim.ApprovedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
                _domainEvents.Add(new ClaimApprovedEvent(ClaimId, approvedAmountPaisa));
            }
            else
            {
                _claim.Status = ClaimStatus.UnderReview;
                _domainEvents.Add(new ClaimEscalatedEvent(ClaimId, approvalLevel + 1));
            }
        }
        else if (decision == ApprovalDecision.Rejected)
        {
            _claim.Status = ClaimStatus.Rejected;
            _claim.RejectionReason = notes;
            _domainEvents.Add(new ClaimRejectedEvent(ClaimId, notes));
        }
        else if (decision == ApprovalDecision.Escalated)
        {
            _claim.Status = ClaimStatus.UnderReview;
            _domainEvents.Add(new ClaimEscalatedEvent(ClaimId, approvalLevel + 1));
        }
        
        _claim.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
    }
    
    public void Settle(long settledAmountPaisa, string paymentMethod, string paymentReference)
    {
        if (Status != ClaimStatus.Approved)
            throw new InvalidOperationException($"Cannot settle claim in status {Status}");
        
        _claim.SettledAmount = new Insuretech.Common.V1.Money 
        { 
            Amount = settledAmountPaisa,
            Currency = "BDT"
        };
        _claim.Status = ClaimStatus.Settled;
        _claim.SettledAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        _claim.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        
        _domainEvents.Add(new ClaimSettledEvent(ClaimId, settledAmountPaisa, paymentMethod));
    }
    
    public void AddDocument(string documentType, string fileUrl, string fileHash)
    {
        var document = new ClaimDocument
        {
            DocumentId = Guid.NewGuid().ToString(),
            ClaimId = ClaimId,
            DocumentType = documentType,
            FileUrl = fileUrl,
            FileHash = fileHash,
            Verified = false,
            UploadedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
        };
        
        _claim.Documents.Add(document);
        _claim.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
    }
    
    private static string GenerateClaimNumber()
    {
        var year = DateTime.UtcNow.Year;
        var random = new Random().Next(1000, 9999);
        var sequence = new Random().Next(100000, 999999);
        return $"CLM-{year}-{random:D4}-{sequence:D6}";
    }
}

// Domain Events
public sealed record ClaimFiledEvent(string ClaimId, string PolicyId, string CustomerId) : DomainEvent;
public sealed record ClaimFlaggedForFraudEvent(string ClaimId, double FraudScore) : DomainEvent;
public sealed record ClaimApprovedEvent(string ClaimId, long ApprovedAmount) : DomainEvent;
public sealed record ClaimRejectedEvent(string ClaimId, string Reason) : DomainEvent;
public sealed record ClaimEscalatedEvent(string ClaimId, int ToLevel) : DomainEvent;
public sealed record ClaimSettledEvent(string ClaimId, long SettledAmount, string PaymentMethod) : DomainEvent;
