using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using InsuranceEngine.Beneficiary.Application.Commands;
using InsuranceEngine.Beneficiary.Application.Queries;

namespace InsuranceEngine.Beneficiary.GrpcServices;

public sealed class BeneficiaryGrpcService : BeneficiaryService.BeneficiaryServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<BeneficiaryGrpcService> _logger;

    public BeneficiaryGrpcService(IMediator mediator, ILogger<BeneficiaryGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<CreateIndividualBeneficiaryResponse> CreateIndividualBeneficiary(
        CreateIndividualBeneficiaryRequest request, ServerCallContext context)
    {
        DateTime dob;
        if (!DateTime.TryParse(request.DateOfBirth, out dob))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid DateOfBirth format. Expected ISO 8601 string."));
        }

        var command = new CreateIndividualBeneficiaryCommand(
            request.UserId,
            request.FullName,
            dob,
            request.Gender,
            request.NidNumber,
            request.MobileNumber,
            request.Email,
            request.PartnerId
        );

        var result = await _mediator.Send(command, context.CancellationToken);
        
        if (result.IsFailure)
        {
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Failed to create beneficiary"));
        }

        return new CreateIndividualBeneficiaryResponse
        {
            BeneficiaryId = result.Value,
            Message = "Individual beneficiary created successfully"
        };
    }

    public override async Task<CreateBusinessBeneficiaryResponse> CreateBusinessBeneficiary(
        CreateBusinessBeneficiaryRequest request, ServerCallContext context)
    {
        var command = new CreateBusinessBeneficiaryCommand(
            request.UserId,
            request.BusinessName,
            request.TradeLicenseNumber,
            request.TinNumber,
            request.FocalPersonName,
            request.FocalPersonMobile,
            request.PartnerId
        );

        var result = await _mediator.Send(command, context.CancellationToken);
        
        if (result.IsFailure)
        {
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Failed to create business beneficiary"));
        }

        return new CreateBusinessBeneficiaryResponse
        {
            BeneficiaryId = result.Value,
            Message = "Business beneficiary created successfully"
        };
    }

    public override async Task<GetBeneficiaryResponse> GetBeneficiary(
        GetBeneficiaryRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.BeneficiaryId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Beneficiary ID is required"));
        }

        var query = new GetBeneficiaryQuery(request.BeneficiaryId);
        var result = await _mediator.Send(query, context.CancellationToken);

        if (result == null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, "Beneficiary not found"));
        }

        return result;
    }

    public override async Task<ListBeneficiariesResponse> ListBeneficiaries(
        ListBeneficiariesRequest request, ServerCallContext context)
    {
        var query = new ListBeneficiariesQuery(
            Type: request.Type,
            Status: request.Status,
            Page: request.Page <= 0 ? 1 : request.Page,
            PageSize: request.PageSize <= 0 ? 10 : request.PageSize
        );

        var result = await _mediator.Send(query, context.CancellationToken);

        if (result == null)
        {
            throw new RpcException(new Status(StatusCode.Internal, "Failed to list beneficiaries"));
        }

        return result;
    }

    public override async Task<UpdateBeneficiaryResponse> UpdateBeneficiary(
        UpdateBeneficiaryRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.BeneficiaryId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Beneficiary ID is required"));
        }

        var command = new UpdateBeneficiaryCommand(request);
        var result = await _mediator.Send(command, context.CancellationToken);

        if (result.Error != null)
        {
            throw new RpcException(new Status(
                (StatusCode)result.Error.HttpStatusCode, 
                result.Error.Message));
        }

        return result;
    }
}
