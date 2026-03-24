using MediatR;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.Application.Queries;

public record ListBeneficiariesQuery(string? Type = null, string? Status = null, int Page = 1, int PageSize = 10) : IRequest<ListBeneficiariesResponse>;
