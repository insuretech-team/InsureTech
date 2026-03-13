using System;
using System.Linq;
using System.Threading.Tasks;
using Grpc.Core;
using InsuranceEngine.Underwriting.Application.Interfaces;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Beneficiary.Entity.V1;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Underwriting.GrpcServices;

public sealed class BeneficiaryGrpcService : BeneficiaryService.BeneficiaryServiceBase
{
    private readonly IBeneficiaryRepository _repository;
    private readonly ILogger<BeneficiaryGrpcService> _logger;

    public BeneficiaryGrpcService(IBeneficiaryRepository repository, ILogger<BeneficiaryGrpcService> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public override async Task<GetBeneficiaryResponse> GetBeneficiary(GetBeneficiaryRequest request, ServerCallContext context)
    {
        var beneficiary = await _repository.GetByIdAsync(Guid.Parse(request.BeneficiaryId));

        if (beneficiary == null)
        {
            return new GetBeneficiaryResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "NOT_FOUND", Message = "Beneficiary not found" }
            };
        }

        return new GetBeneficiaryResponse
        {
            Beneficiary = new Beneficiary
            {
                BeneficiaryId = beneficiary.Id.ToString(),
                BeneficiaryCode = beneficiary.Code,
                UserId = beneficiary.UserId.ToString(),
                Type = beneficiary.Type.ToString(),
                Status = beneficiary.Status.ToString()
            }
        };
    }

    public override async Task<ListBeneficiariesResponse> ListBeneficiaries(ListBeneficiariesRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 10 : request.PageSize;

        var items = await _repository.ListAsync(request.Type, request.Status, page, pageSize);
        var total = await _repository.GetTotalCountAsync(request.Type, request.Status);

        var response = new ListBeneficiariesResponse
        {
            TotalCount = total
        };

        response.Beneficiaries.AddRange(items.Select(b => new Beneficiary
        {
            BeneficiaryId = b.Id.ToString(),
            BeneficiaryCode = b.Code,
            UserId = b.UserId.ToString(),
            Type = b.Type.ToString(),
            Status = b.Status.ToString()
        }));

        return response;
    }
}
