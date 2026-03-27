using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed class AddNomineeCommandHandler : IRequestHandler<AddNomineeCommand, Result<string>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<AddNomineeCommandHandler> _logger;

    public AddNomineeCommandHandler(DbContext dbContext, ILogger<AddNomineeCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(AddNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var nomineeId = Guid.NewGuid();
            var now = DateTime.UtcNow;

            var sql = @"
                INSERT INTO insurance_schema.policy_nominees (
                    nominee_id, policy_id, full_name, relationship, share_percentage,
                    date_of_birth, nid_number, phone_number, nominee_dob_text, created_at
                ) VALUES (
                    @NomineeId, @PolicyId, @FullName, @Relationship, @SharePercentage,
                    @DateOfBirth, @NidNumber, @PhoneNumber, @NomineeDobText, @CreatedAt
                )";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            await connection.ExecuteAsync(sql, new
            {
                NomineeId = nomineeId,
                PolicyId = request.PolicyId,
                FullName = request.FullName,
                Relationship = request.Relationship,
                SharePercentage = request.SharePercentage,
                DateOfBirth = request.DateOfBirth,
                NidNumber = request.NidNumber,
                PhoneNumber = request.PhoneNumber,
                NomineeDobText = request.NomineeDobText,
                CreatedAt = now
            });

            _logger.LogInformation("Nominee added: {NomineeId} to Policy: {PolicyId}", nomineeId, request.PolicyId);
            return Result<string>.Ok(nomineeId.ToString());
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to add nominee to policy {PolicyId}", request.PolicyId);
            return Result<string>.Fail("NOMINEE_ADD_FAILED", ex.Message);
        }
    }
}

public sealed class UpdateNomineeCommandHandler : IRequestHandler<UpdateNomineeCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<UpdateNomineeCommandHandler> _logger;

    public UpdateNomineeCommandHandler(DbContext dbContext, ILogger<UpdateNomineeCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(UpdateNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.policy_nominees
                SET full_name = COALESCE(@FullName, full_name),
                    relationship = COALESCE(@Relationship, relationship),
                    share_percentage = COALESCE(@SharePercentage, share_percentage),
                    date_of_birth = COALESCE(@DateOfBirth, date_of_birth),
                    nid_number = COALESCE(@NidNumber, nid_number),
                    phone_number = COALESCE(@PhoneNumber, phone_number),
                    updated_at = @UpdatedAt
                WHERE nominee_id = @NomineeId::uuid AND policy_id = @PolicyId::uuid";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var rows = await connection.ExecuteAsync(sql, new
            {
                NomineeId = request.NomineeId,
                PolicyId = request.PolicyId,
                FullName = request.FullName,
                Relationship = request.Relationship,
                SharePercentage = request.SharePercentage,
                DateOfBirth = request.DateOfBirth,
                NidNumber = request.NidNumber,
                PhoneNumber = request.PhoneNumber,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.Fail("NOMINEE_NOT_FOUND", "Nominee not found");

            _logger.LogInformation("Nominee updated: {NomineeId}", request.NomineeId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update nominee {NomineeId}", request.NomineeId);
            return Result<bool>.Fail("NOMINEE_UPDATE_FAILED", ex.Message);
        }
    }
}

public sealed class DeleteNomineeCommandHandler : IRequestHandler<DeleteNomineeCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<DeleteNomineeCommandHandler> _logger;

    public DeleteNomineeCommandHandler(DbContext dbContext, ILogger<DeleteNomineeCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DeleteNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                DELETE FROM insurance_schema.policy_nominees
                WHERE nominee_id = @NomineeId::uuid AND policy_id = @PolicyId::uuid";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var rows = await connection.ExecuteAsync(sql, new
            {
                NomineeId = request.NomineeId,
                PolicyId = request.PolicyId
            });

            if (rows == 0)
                return Result<bool>.Fail("NOMINEE_NOT_FOUND", "Nominee not found");

            _logger.LogInformation("Nominee deleted: {NomineeId}", request.NomineeId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to delete nominee {NomineeId}", request.NomineeId);
            return Result<bool>.Fail("NOMINEE_DELETE_FAILED", ex.Message);
        }
    }
}
