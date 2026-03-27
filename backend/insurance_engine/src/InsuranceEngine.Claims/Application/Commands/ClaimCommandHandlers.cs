using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using Dapper;
using System.Data;
using InsuranceEngine.Claims.Domain;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Domain.Events;

namespace InsuranceEngine.Claims.Application.Commands;

public sealed class SubmitClaimCommandHandler : IRequestHandler<SubmitClaimCommand, SubmitClaimResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<SubmitClaimCommandHandler> _logger;
    private readonly IKafkaPublisher _kafkaPublisher;

    public SubmitClaimCommandHandler(DbContext dbContext, ILogger<SubmitClaimCommandHandler> logger, IKafkaPublisher kafkaPublisher)
    {
        _dbContext = dbContext;
        _logger = logger;
        _kafkaPublisher = kafkaPublisher;
    }

    public async Task<SubmitClaimResponse> Handle(SubmitClaimCommand request, CancellationToken cancellationToken)
    {
        var connection = _dbContext.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        try
        {
            // 1. Domain Logic: Create Aggregate (DDD)
            var claim = ClaimAggregate.Submit(
                policyId: Guid.Parse(request.PolicyId),
                type: request.ClaimType,
                amount: request.ClaimAmount,
                description: request.Description ?? "",
                documentContent: request.DocumentContent
            );

            // 1.1 Support for Tiered Approval Matrix (FR-086)
            claim.RequestApproval();

            // 2. Persist Aggregate State using Dapper
            var insertClaimSql = @"
                INSERT INTO insurance_schema.claims (
                    claim_id, claim_number, policy_id, claim_type, claim_amount,
                    description, status, document_hash, created_at
                ) VALUES (
                    @ClaimId, @ClaimNumber, @PolicyId, @ClaimType, @ClaimAmount,
                    @Description, @Status, @DocumentHash, @CreatedAt
                )";

            await connection.ExecuteAsync(insertClaimSql, new
            {
                ClaimId = claim.Id,
                ClaimNumber = claim.ClaimNumber,
                PolicyId = claim.PolicyId,
                ClaimType = claim.ClaimType,
                ClaimAmount = claim.ClaimAmount.Amount,
                Description = claim.Description,
                Status = claim.Status,
                DocumentHash = claim.DocumentHash,
                CreatedAt = claim.CreatedAt
            }, transaction);

            await transaction.CommitAsync(cancellationToken);

            // FR-019: Kafka Event Streaming
            var claimEvent = new ClaimSubmittedEvent(claim.Id, claim.ClaimNumber, claim.PolicyId, claim.ClaimAmount.Amount);
            await _kafkaPublisher.PublishAsync("insurance.claims.submitted", claimEvent);

            _logger.LogInformation("Claim submitted and event published: {ClaimNumber} (Status: {Status})", claim.ClaimNumber, claim.Status);

            return new SubmitClaimResponse
            {
                ClaimId = claim.Id.ToString(),
                Message = "Claim submitted successfully"
            };
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            _logger.LogError(ex, "Failed to submit claim");
            throw;
        }
    }
}

public sealed class ApproveClaimCommandHandler : IRequestHandler<ApproveClaimCommand, ApproveClaimResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ApproveClaimCommandHandler> _logger;

    public ApproveClaimCommandHandler(DbContext dbContext, ILogger<ApproveClaimCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<ApproveClaimResponse> Handle(ApproveClaimCommand request, CancellationToken cancellationToken)
    {
        var connection = _dbContext.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        // In a real DDD flow, we would load the Aggregate here. 
        // For now, we perform the update directly but align with the Status logic.
        try
        {
            var sql = @"
                UPDATE insurance_schema.claims
                SET status = 'APPROVED', 
                    approved_amount = @ApprovedAmount, 
                    updated_at = @UpdatedAt
                WHERE claim_id = @ClaimId AND deleted_at IS NULL";

            var rows = await connection.ExecuteAsync(sql, new 
            { 
                ClaimId = Guid.Parse(request.ClaimId), 
                ApprovedAmount = (long)(request.ApprovedAmount * 100), 
                UpdatedAt = DateTime.UtcNow 
            });

            if (rows == 0) throw new Exception("Claim not found or already processed");

            return new ApproveClaimResponse
            {
                Message = "Claim approved successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve claim {ClaimId}", request.ClaimId);
            throw;
        }
    }
}
