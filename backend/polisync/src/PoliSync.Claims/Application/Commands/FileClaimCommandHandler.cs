using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Claims.Domain;
using PoliSync.Infrastructure.Persistence;

namespace PoliSync.Claims.Application.Commands;

public sealed class FileClaimCommandHandler : IRequestHandler<FileClaimCommand, Result<string>>
{
    private readonly PoliSyncDbContext _dbContext;
    private readonly IEventBus _eventBus;
    private readonly ILogger<FileClaimCommandHandler> _logger;
    
    public FileClaimCommandHandler(
        PoliSyncDbContext dbContext,
        IEventBus eventBus,
        ILogger<FileClaimCommandHandler> logger)
    {
        _dbContext = dbContext;
        _eventBus = eventBus;
        _logger = logger;
    }
    
    public async Task<Result<string>> Handle(FileClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Create claim aggregate (FNOL - First Notice of Loss)
            var claimAggregate = ClaimAggregate.FileClaim(
                request.PolicyId,
                request.CustomerId,
                request.ClaimType,
                request.ClaimedAmountPaisa,
                request.IncidentDate,
                request.IncidentDescription,
                request.PlaceOfIncident);
            
            // Save to database
            await SaveClaimToDatabase(claimAggregate.Claim, cancellationToken);
            
            // Publish domain events
            foreach (var domainEvent in claimAggregate.DomainEvents)
            {
                await _eventBus.PublishAsync(domainEvent, cancellationToken);
            }
            
            _logger.LogInformation(
                "Claim filed successfully: {ClaimId}, Number: {ClaimNumber}, Policy: {PolicyId}",
                claimAggregate.ClaimId,
                claimAggregate.ClaimNumber,
                request.PolicyId);
            
            return Result.Ok(claimAggregate.ClaimId);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to file claim for policy {PolicyId}", request.PolicyId);
            return Result.Fail<string>("CLAIM_FILING_FAILED", ex.Message);
        }
    }
    
    private async Task SaveClaimToDatabase(Insuretech.Claims.Entity.V1.Claim claim, CancellationToken cancellationToken)
    {
        var sql = @"
            INSERT INTO insurance_schema.claims (
                claim_id, claim_number, policy_id, customer_id, status, type,
                claimed_amount, claimed_currency, approved_amount, approved_currency,
                settled_amount, settled_currency, incident_date, incident_description,
                place_of_incident, submitted_at, processing_type, created_at, updated_at
            ) VALUES (
                @p0, @p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18
            )";
        
        await _dbContext.Database.ExecuteSqlRawAsync(sql,
            new object[]
            {
                claim.ClaimId,
                claim.ClaimNumber,
                claim.PolicyId,
                claim.CustomerId,
                claim.Status.ToString(),
                claim.Type.ToString(),
                claim.ClaimedAmount.Amount,
                claim.ClaimedAmount.Currency,
                claim.ApprovedAmount.Amount,
                claim.ApprovedAmount.Currency,
                claim.SettledAmount.Amount,
                claim.SettledAmount.Currency,
                claim.IncidentDate.ToDateTime(),
                claim.IncidentDescription,
                claim.PlaceOfIncident,
                claim.SubmittedAt.ToDateTime(),
                claim.ProcessingType.ToString(),
                claim.CreatedAt.ToDateTime(),
                claim.UpdatedAt?.ToDateTime() ?? (object)DBNull.Value
            },
            cancellationToken);
    }
}
