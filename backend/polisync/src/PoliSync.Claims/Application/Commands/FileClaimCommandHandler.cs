using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Claims.Domain;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;

namespace PoliSync.Claims.Application.Commands;

public class FileClaimCommandHandler : IRequestHandler<FileClaimCommand, Result<string>>
{
    private readonly PoliSyncDbContext _dbContext;
    private readonly IEventBus _eventBus;
    private readonly IMediator _mediator;
    private readonly ILogger<FileClaimCommandHandler> _logger;

    public FileClaimCommandHandler(
        PoliSyncDbContext dbContext,
        IEventBus eventBus,
        IMediator mediator,
        ILogger<FileClaimCommandHandler> logger)
    {
        _dbContext = dbContext;
        _eventBus = eventBus;
        _mediator = mediator;
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

            // Trigger approval workflow — template resolved dynamically by amount
            var workflowResult = await _mediator.Send(new TriggerWorkflowCommand(
                new WorkflowTriggerContext
                {
                    EntityType  = "CLAIM",
                    EntityId    = claimAggregate.ClaimId,
                    InitiatedBy = request.CustomerId,
                    AmountPaisa = request.ClaimedAmountPaisa,
                    Portal      = "B2C",
                    Metadata    = new Dictionary<string, string>
                    {
                        ["policy_id"]    = request.PolicyId,
                        ["claim_number"] = claimAggregate.ClaimNumber,
                        ["claim_type"]   = request.ClaimType.ToString().Replace("ClaimType", "").Trim()
                    }
                }), cancellationToken);

            if (workflowResult.IsSuccess && workflowResult.Value!.WasTriggered)
                _logger.LogInformation(
                    "Claim approval workflow started: instance={InstanceId} template='{Template}'",
                    workflowResult.Value.WorkflowInstanceId, workflowResult.Value.TemplateName);

            return Result.Ok(claimAggregate.ClaimId);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to file claim for policy {PolicyId}", request.PolicyId);
            return Result.Fail<string>("CLAIM_FILING_FAILED", ex.Message);
        }
    }
    
    protected virtual async Task SaveClaimToDatabase(Insuretech.Claims.Entity.V1.Claim claim, CancellationToken cancellationToken)
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
                Guid.Parse(claim.ClaimId),                     // uuid: claim_id
                claim.ClaimNumber,                              // varchar: claim_number
                Guid.Parse(claim.PolicyId),                    // uuid: policy_id
                Guid.Parse(claim.CustomerId),                  // uuid: customer_id
                // Claim status/type must be uppercase to match DB check constraint
                // proto enum .ToString() gives "Submitted", "Health" etc — convert to DB format
                claim.Status.ToString().ToUpperInvariant(),    // varchar: status (e.g. "SUBMITTED")
                claim.Type.ToString().ToUpperInvariant(),      // varchar: type (e.g. "HEALTH")
                claim.ClaimedAmount?.Amount ?? 0L,             // bigint: claimed_amount
                claim.ClaimedAmount?.Currency ?? "BDT",        // varchar: claimed_currency
                claim.ApprovedAmount?.Amount ?? 0L,            // bigint: approved_amount
                claim.ApprovedAmount?.Currency ?? "BDT",       // varchar: approved_currency
                claim.SettledAmount?.Amount ?? 0L,             // bigint: settled_amount
                claim.SettledAmount?.Currency ?? "BDT",        // varchar: settled_currency
                claim.IncidentDate?.ToDateTime(),               // timestamp: incident_date
                claim.IncidentDescription,                      // text: incident_description
                claim.PlaceOfIncident,                          // text: place_of_incident
                claim.SubmittedAt?.ToDateTime() ?? DateTime.UtcNow, // timestamp: submitted_at
                claim.ProcessingType.ToString(),                // varchar: processing_type
                claim.CreatedAt?.ToDateTime() ?? DateTime.UtcNow,  // timestamp: created_at
                claim.UpdatedAt?.ToDateTime() ?? DateTime.UtcNow  // timestamp: updated_at (defaults to now if not set)
            },
            cancellationToken);
    }
}
