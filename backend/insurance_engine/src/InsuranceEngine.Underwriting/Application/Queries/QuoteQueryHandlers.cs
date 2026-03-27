using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Underwriting.Application.Queries;

public sealed class ListQuotesQueryHandler : IRequestHandler<ListQuotesQuery, Result<ListQuotesResult>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ListQuotesQueryHandler> _logger;

    public ListQuotesQueryHandler(DbContext dbContext, ILogger<ListQuotesQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<ListQuotesResult>> Handle(ListQuotesQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT quote_id, quote_number, beneficiary_id, insurer_product_id, status,
                       sum_assured, term_years, base_premium, total_premium,
                       applicant_age, smoker, valid_until, created_at
                FROM insurance_schema.quotes
                WHERE (@BeneficiaryId IS NULL OR beneficiary_id = @BeneficiaryId)
                  AND (@Status IS NULL OR status = @Status)
                  AND deleted_at IS NULL
                ORDER BY created_at DESC
                LIMIT @PageSize OFFSET @Offset";

            var offset = (request.Page - 1) * request.PageSize;

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var items = await connection.QueryAsync<QuoteDto>(sql, new
            {
                BeneficiaryId = request.BeneficiaryId,
                Status = request.Status,
                PageSize = request.PageSize,
                Offset = offset
            });

            var countSql = @"
                SELECT COUNT(*) FROM insurance_schema.quotes
                WHERE (@BeneficiaryId IS NULL OR beneficiary_id = @BeneficiaryId)
                  AND (@Status IS NULL OR status = @Status)
                  AND deleted_at IS NULL";

            var totalCount = await connection.ExecuteScalarAsync<int>(countSql, new
            {
                BeneficiaryId = request.BeneficiaryId,
                Status = request.Status
            });

            return Result<ListQuotesResult>.Ok(new ListQuotesResult(items.ToList(), totalCount));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list quotes");
            return Result<ListQuotesResult>.Fail("QUOTE_LIST_FAILED", ex.Message);
        }
    }
}

public sealed class GetQuoteQueryHandler : IRequestHandler<GetQuoteQuery, Result<QuoteDto?>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<GetQuoteQueryHandler> _logger;

    public GetQuoteQueryHandler(DbContext dbContext, ILogger<GetQuoteQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<QuoteDto?>> Handle(GetQuoteQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT quote_id, quote_number, beneficiary_id, insurer_product_id, status,
                       sum_assured, term_years, base_premium, total_premium,
                       applicant_age, smoker, valid_until, created_at
                FROM insurance_schema.quotes
                WHERE quote_id = @QuoteId AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var quote = await connection.QueryFirstOrDefaultAsync<QuoteDto>(sql, new
            {
                QuoteId = request.QuoteId
            });

            if (quote == null)
                return Result<QuoteDto?>.NotFound("QUOTE_NOT_FOUND", "Quote not found");

            return Result<QuoteDto?>.Ok(quote);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get quote {QuoteId}", request.QuoteId);
            return Result<QuoteDto?>.Fail("QUOTE_GET_FAILED", ex.Message);
        }
    }
}
