using MediatR;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.Application.Queries;

public record GetBeneficiaryQuery(string BeneficiaryId) : IRequest<GetBeneficiaryResponse>;

public record ListBeneficiariesQuery(
    int Page,
    int PageSize,
    string? Type,
    string? Status
) : IRequest<ListBeneficiariesResponse>;
