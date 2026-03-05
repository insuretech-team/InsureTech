using PoliSync.SharedKernel.CQRS;
using PoliSync.Underwriting.Domain;

namespace PoliSync.Underwriting.Application.Commands;

// ── Request Quote ───────────────────────────────────────────────────

public record RequestQuoteCommand(
    Guid BeneficiaryId,
    Guid InsurerProductId,
    long SumAssured,
    int TermYears,
    string PremiumPaymentMode,
    int ApplicantAgeDays,
    bool IsSmoker,
    long BasePremiumAmount,   // typically passed from frontend or calculated via policy service
    long RiderPremiumAmount,
    long TaxAmount
) : ICommand<Guid>;

public class RequestQuoteHandler : ICommandHandler<RequestQuoteCommand, Guid>
{
    private readonly IQuoteRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public RequestQuoteHandler(IQuoteRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result<Guid>> Handle(RequestQuoteCommand cmd, CancellationToken ct)
    {
        var quote = Quote.Create(
            cmd.BeneficiaryId,
            cmd.InsurerProductId,
            cmd.SumAssured,
            cmd.TermYears,
            cmd.PremiumPaymentMode,
            cmd.BasePremiumAmount,
            cmd.RiderPremiumAmount,
            cmd.TaxAmount,
            cmd.ApplicantAgeDays,
            cmd.IsSmoker
        );

        await _repo.AddAsync(quote, ct);
        await _uow.SaveChangesAsync(ct);

        return Result<Guid>.Ok(quote.QuoteId);
    }
}

// ── Submit Health Declaration ───────────────────────────────────────

public record SubmitHealthDeclarationCommand(
    Guid QuoteId,
    float HeightCm,
    float WeightKg,
    bool IsSmoker,
    bool ConsumesAlcohol,
    bool HasPreExistingConditions,
    string? ConditionDetails,
    bool HasFamilyHistoryOfCriticalIllness,
    string OccupationRiskLevel
) : ICommand<Guid>;

public class SubmitHealthDeclarationHandler : ICommandHandler<SubmitHealthDeclarationCommand, Guid>
{
    private readonly IQuoteRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public SubmitHealthDeclarationHandler(IQuoteRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result<Guid>> Handle(SubmitHealthDeclarationCommand cmd, CancellationToken ct)
    {
        var quote = await _repo.GetByIdAsync(cmd.QuoteId, ct);
        if (quote is null) return Result<Guid>.NotFound("Quote not found");

        var existingHmd = await _repo.GetHealthDeclarationByQuoteIdAsync(cmd.QuoteId, ct);
        if (existingHmd is not null) return Result<Guid>.Conflict("Health declaration already submitted for this quote");

        var hd = HealthDeclaration.Create(
            cmd.QuoteId, cmd.HeightCm, cmd.WeightKg, cmd.IsSmoker, cmd.ConsumesAlcohol,
            cmd.HasPreExistingConditions, cmd.ConditionDetails, cmd.HasFamilyHistoryOfCriticalIllness,
            cmd.OccupationRiskLevel
        );

        await _repo.AddHealthDeclarationAsync(hd, ct);
        await _uow.SaveChangesAsync(ct);

        return Result<Guid>.Ok(hd.DeclarationId);
    }
}

// ── Approve/Reject Underwriting ──────────────────────────────────────

public record ApproveUnderwritingCommand(
    Guid QuoteId,
    float RiskScore,
    RiskLevel RiskLevel,
    string? Reason,
    string? Conditions,
    string? RiskFactors,
    long? AdjustedPremiumAmount,
    Guid? UnderwriterId
) : ICommand;

public class ApproveUnderwritingHandler : ICommandHandler<ApproveUnderwritingCommand>
{
    private readonly IQuoteRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public ApproveUnderwritingHandler(IQuoteRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(ApproveUnderwritingCommand cmd, CancellationToken ct)
    {
        var quote = await _repo.GetByIdAsync(cmd.QuoteId, ct);
        if (quote is null) return Result.NotFound("Quote not found");

        quote.Approve(); // changes quote status
        
        var decision = UnderwritingDecision.Create(
            cmd.QuoteId, DecisionType.Approved, DecisionMethod.Manual, cmd.RiskScore, cmd.RiskLevel,
            cmd.Reason, cmd.Conditions, cmd.RiskFactors, cmd.AdjustedPremiumAmount, cmd.UnderwriterId
        );

        _repo.Update(quote);
        await _repo.AddDecisionAsync(decision, ct);
        await _uow.SaveChangesAsync(ct);

        return Result.Ok();
    }
}

public record RejectUnderwritingCommand(
    Guid QuoteId,
    float RiskScore,
    RiskLevel RiskLevel,
    string Reason,
    string? RiskFactors,
    Guid? UnderwriterId
) : ICommand;

public class RejectUnderwritingHandler : ICommandHandler<RejectUnderwritingCommand>
{
    private readonly IQuoteRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public RejectUnderwritingHandler(IQuoteRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(RejectUnderwritingCommand cmd, CancellationToken ct)
    {
        var quote = await _repo.GetByIdAsync(cmd.QuoteId, ct);
        if (quote is null) return Result.NotFound("Quote not found");

        quote.Reject(); // changes quote status
        
        var decision = UnderwritingDecision.Create(
            cmd.QuoteId, DecisionType.Rejected, DecisionMethod.Manual, cmd.RiskScore, cmd.RiskLevel,
            cmd.Reason, null, cmd.RiskFactors, null, cmd.UnderwriterId
        );

        _repo.Update(quote);
        await _repo.AddDecisionAsync(decision, ct);
        await _uow.SaveChangesAsync(ct);

        return Result.Ok();
    }
}

// ── Convert to Policy ───────────────────────────────────────────────

public record ConvertQuoteToPolicyCommand(Guid QuoteId) : ICommand;

public class ConvertQuoteToPolicyHandler : ICommandHandler<ConvertQuoteToPolicyCommand>
{
    // The actual policy creation logic will be in the Policies module listening to an event,
    // or coordinated visually. Here we just update the quote status.
    private readonly IQuoteRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public ConvertQuoteToPolicyHandler(IQuoteRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(ConvertQuoteToPolicyCommand cmd, CancellationToken ct)
    {
        var quote = await _repo.GetByIdAsync(cmd.QuoteId, ct);
        if (quote is null) return Result.NotFound("Quote not found");

        quote.ConvertToPolicy();
        _repo.Update(quote);
        await _uow.SaveChangesAsync(ct);

        // A QuoteConvertedEvent could trigger Policy module creation via messaging later
        return Result.Ok();
    }
}
