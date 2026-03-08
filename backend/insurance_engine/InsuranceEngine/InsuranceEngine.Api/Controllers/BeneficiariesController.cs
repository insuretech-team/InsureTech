using InsuranceEngine.Api.DTOs;
using InsuranceEngine.Api.RequestModels;
using InsuranceEngine.Api.ResponseModels;
using InsuranceEngine.Application.Interfaces;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.ModelBinding;

namespace InsuranceEngine.Api.Controllers;

[ApiController]
[Route("v1/beneficiaries")]
public class BeneficiariesController(IBeneficiaryService beneficiaryService) : ControllerBase
{
    private readonly IBeneficiaryService _beneficiaryService = beneficiaryService;

    [HttpGet]
    public async Task<ActionResult<ListBeneficiariesResponseV1>> GetBeneficiaries(
        [FromQuery(Name = "page")] int page = 1,
        [FromQuery(Name = "page_size")] int pageSize = 20,
        CancellationToken cancellationToken = default)
    {
        if (page < 1 || pageSize < 1 || pageSize > 100)
        {
            return BadRequest(new ListBeneficiariesResponseV1
            {
                Error = Error.Create(
                    "INVALID_PAGINATION",
                    "Page must be >= 1 and page_size must be between 1 and 100.",
                    StatusCodes.Status400BadRequest)
            });
        }

        ListBeneficiariesResponseV1 response =
            await _beneficiaryService.GetBeneficiariesAsync(page, pageSize, cancellationToken);

        return Ok(response);
    }

    [HttpGet("{beneficiaryId}")]
    public async Task<ActionResult<GetBeneficiaryResponseV1>> GetBeneficiaryById(
        string beneficiaryId,
        CancellationToken cancellationToken = default)
    {
        (GetBeneficiaryResponseV1 response, bool notFound) =
            await _beneficiaryService.GetBeneficiaryByIdAsync(beneficiaryId, cancellationToken);

        if (response.Error is not null)
        {
            return response.Error.HttpStatusCode switch
            {
                400 => BadRequest(response),
                404 => NotFound(response),
                _ => BadRequest(response)
            };
        }

        if (notFound)
        {
            return NotFound(response);
        }

        return Ok(response);
    }

    [HttpPost("individual")]
    public async Task<ActionResult<CreateBeneficiaryResponseV1>> CreateIndividualBeneficiary(
        [FromBody] CreateIndividualBeneficiaryRequestV1 request,
        CancellationToken cancellationToken = default)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest(new CreateBeneficiaryResponseV1
            {
                Error = BuildValidationError(ModelState)
            });
        }

        CreateBeneficiaryResponseV1 response =
            await _beneficiaryService.CreateIndividualBeneficiaryAsync(request, cancellationToken);

        return CreatedAtAction(nameof(GetBeneficiaryById), new { beneficiaryId = response.BeneficiaryId }, response);
    }

    [HttpPost("business")]
    public async Task<ActionResult<CreateBeneficiaryResponseV1>> CreateBusinessBeneficiary(
        [FromBody] CreateBusinessBeneficiaryRequestV1 request,
        CancellationToken cancellationToken = default)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest(new CreateBeneficiaryResponseV1
            {
                Error = BuildValidationError(ModelState)
            });
        }

        CreateBeneficiaryResponseV1 response =
            await _beneficiaryService.CreateBusinessBeneficiaryAsync(request, cancellationToken);

        return CreatedAtAction(nameof(GetBeneficiaryById), new { beneficiaryId = response.BeneficiaryId }, response);
    }

    [HttpPatch("{beneficiaryId}")]
    public async Task<ActionResult<UpdateBeneficiaryResponseV1>> UpdateBeneficiary(
        string beneficiaryId,
        [FromBody] UpdateBeneficiaryRequestV1 request,
        CancellationToken cancellationToken = default)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest(new UpdateBeneficiaryResponseV1
            {
                Error = BuildValidationError(ModelState)
            });
        }

        if (!string.Equals(beneficiaryId, request.BeneficiaryId, StringComparison.OrdinalIgnoreCase))
        {
            return BadRequest(new UpdateBeneficiaryResponseV1
            {
                Error = Error.Create(
                    "BENEFICIARY_ID_MISMATCH",
                    "beneficiary_id must match the path parameter.",
                    StatusCodes.Status400BadRequest)
            });
        }

        (UpdateBeneficiaryResponseV1 response, bool notFound) =
            await _beneficiaryService.UpdateBeneficiaryAsync(beneficiaryId, request, cancellationToken);

        if (response.Error is not null)
        {
            return response.Error.HttpStatusCode switch
            {
                400 => BadRequest(response),
                404 => NotFound(response),
                _ => BadRequest(response)
            };
        }

        if (notFound)
        {
            return NotFound(response);
        }

        return Ok(response);
    }

    private static Error BuildValidationError(ModelStateDictionary modelState)
    {
        List<FieldViolation> fieldViolations = modelState
            .Where(entry => entry.Value?.Errors.Count > 0)
            .SelectMany(entry =>
            {
                string fieldName = entry.Key;
                return entry.Value!.Errors.Select(error => new FieldViolation
                {
                    Field = fieldName,
                    Code = "INVALID_FIELD",
                    Description = string.IsNullOrWhiteSpace(error.ErrorMessage)
                        ? "Invalid input parameters."
                        : error.ErrorMessage,
                    RejectedValue = entry.Value.AttemptedValue ?? string.Empty
                });
            })
            .ToList();

        return Error.Create(
            "INVALID_REQUEST",
            "Invalid input parameters.",
            StatusCodes.Status400BadRequest,
            fieldViolations: fieldViolations.Count > 0 ? fieldViolations : null);
    }
}
