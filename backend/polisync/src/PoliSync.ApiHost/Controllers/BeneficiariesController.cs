using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Beneficiaries.Application.Commands;
using PoliSync.Beneficiaries.Application.Queries;
using PoliSync.Beneficiaries.Domain;

namespace PoliSync.ApiHost.Controllers;

[ApiController]
[Route("v1/beneficiaries")]
public class BeneficiariesController : ControllerBase
{
    private readonly IMediator _mediator;

    public BeneficiariesController(IMediator mediator) => _mediator = mediator;

    [HttpPost("individual")]
    public async Task<IActionResult> CreateIndividual(CreateIndividualBeneficiaryRequest request)
    {
        var cmd = new CreateIndividualBeneficiaryCommand(
            request.UserId,
            request.FullName,
            request.DateOfBirth,
            request.Gender,
            request.NidNumber,
            request.MobileNumber,
            request.Email,
            request.PartnerId
        );

        var result = await _mediator.Send(cmd);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error);
    }

    [HttpPost("business")]
    public async Task<IActionResult> CreateBusiness(CreateBusinessBeneficiaryRequest request)
    {
        var cmd = new CreateBusinessBeneficiaryCommand(
            request.UserId,
            request.BusinessName,
            request.TradeLicenseNumber,
            request.TinNumber,
            request.FocalPersonName,
            request.FocalPersonMobile,
            request.PartnerId
        );

        var result = await _mediator.Send(cmd);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error);
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetBeneficiaryQuery(id));
        if (!result.IsSuccess) return BadRequest(result.Error);
        if (result.Value is null) return NotFound();

        return Ok(MapBeneficiaryDto(result.Value));
    }

    [HttpGet]
    public async Task<IActionResult> List([FromQuery] BeneficiaryType? type, [FromQuery] BeneficiaryStatus? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListBeneficiariesQuery(type, status, page, pageSize));
        return Ok(new
        {
            items = result.Value.Items.Select(MapBeneficiaryDto),
            totalCount = result.Value.TotalCount,
            page = result.Value.Page,
            pageSize = result.Value.PageSize
        });
    }

    [HttpPost("{id}/kyc")]
    public async Task<IActionResult> CompleteKyc(Guid id, [FromBody] KycStatus status)
    {
        var result = await _mediator.Send(new CompleteKycCommand(id, status));
        return result.IsSuccess ? Ok() : BadRequest(result.Error);
    }

    private static object MapBeneficiaryDto(Beneficiary b) => new
    {
        beneficiary_id = b.BeneficiaryId,
        user_id = b.UserId,
        type = b.Type.ToString(),
        code = b.Code,
        status = b.Status.ToString(),
        kyc_status = b.KycStatus.ToString(),
        kyc_completed_at = b.KycCompletedAt,
        risk_score = b.RiskScore,
        referral_code = b.ReferralCode,
        individual_details = b.IndividualDetails == null ? null : new
        {
            full_name = b.IndividualDetails.FullName,
            date_of_birth = b.IndividualDetails.DateOfBirth,
            gender = b.IndividualDetails.Gender.ToString(),
            nid_number = b.IndividualDetails.NidNumber
        },
        business_details = b.BusinessDetails == null ? null : new
        {
            business_name = b.BusinessDetails.BusinessName,
            focal_person_name = b.BusinessDetails.FocalPersonName
        }
    };
}

public record CreateIndividualBeneficiaryRequest(
    Guid UserId,
    string FullName,
    DateTime DateOfBirth,
    Gender Gender,
    string NidNumber,
    string MobileNumber,
    string? Email,
    Guid? PartnerId
);

public record CreateBusinessBeneficiaryRequest(
    Guid UserId,
    string BusinessName,
    string TradeLicenseNumber,
    string TinNumber,
    string FocalPersonName,
    string FocalPersonMobile,
    Guid? PartnerId
);
