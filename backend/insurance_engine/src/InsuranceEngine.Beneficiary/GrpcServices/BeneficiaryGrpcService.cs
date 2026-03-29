using Grpc.Core;
using MediatR;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.GrpcServices;

public class BeneficiaryGrpcService : Insuretech.Beneficiary.Services.V1.BeneficiaryService.BeneficiaryServiceBase
{
    private readonly IMediator _mediator;

    public BeneficiaryGrpcService(IMediator mediator) => _mediator = mediator;

    public override async Task<CreateIndividualBeneficiaryResponse> CreateIndividualBeneficiary(CreateIndividualBeneficiaryRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new Application.Commands.CreateIndividualBeneficiaryCommand(
            request.UserId,
            request.FullName,
            DateTime.Parse(request.DateOfBirth),
            request.Gender,
            request.NidNumber,
            request.MobileNumber,
            request.Email,
            request.PartnerId
        ));
    }

    public override async Task<CreateBusinessBeneficiaryResponse> CreateBusinessBeneficiary(CreateBusinessBeneficiaryRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new Application.Commands.CreateBusinessBeneficiaryCommand(
            request.UserId,
            request.BusinessName,
            request.TradeLicenseNumber,
            request.TinNumber,
            request.FocalPersonName,
            request.FocalPersonMobile,
            request.PartnerId
        ));
    }

    public override async Task<GetBeneficiaryResponse> GetBeneficiary(GetBeneficiaryRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new Application.Queries.GetBeneficiaryQuery(request.BeneficiaryId));
    }

    public override async Task<UpdateBeneficiaryResponse> UpdateBeneficiary(UpdateBeneficiaryRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new Application.Commands.UpdateBeneficiaryCommand(
            request.BeneficiaryId,
            request.MobileNumber,
            request.Email,
            request.Address
        ));
    }

    public override async Task<CompleteKYCResponse> CompleteKYC(CompleteKYCRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new Application.Commands.CompleteKYCCommand(
            request.BeneficiaryId,
            request.NidFrontUrl,
            request.NidBackUrl,
            request.SelfieUrl,
            request.PorichoyVerificationId
        ));
    }

    public override async Task<UpdateRiskScoreResponse> UpdateRiskScore(UpdateRiskScoreRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new Application.Commands.UpdateRiskScoreCommand(
            request.BeneficiaryId,
            request.RiskScore,
            request.Reason
        ));
    }

    public override async Task<ListBeneficiariesResponse> ListBeneficiaries(ListBeneficiariesRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new Application.Queries.ListBeneficiariesQuery(
            request.Page,
            request.PageSize,
            request.Type,
            request.Status
        ));
    }
}
