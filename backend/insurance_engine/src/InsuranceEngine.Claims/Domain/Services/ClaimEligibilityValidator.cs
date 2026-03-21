using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Claims.Domain.Services;

/// <summary>
/// Domain service for claim eligibility validation (FR-042).
/// Validates: policy active, coverage period, claim type, duplicate submission.
/// </summary>
public class ClaimEligibilityValidator
{
    private readonly IClaimsRepository _claimsRepository;

    public ClaimEligibilityValidator(IClaimsRepository claimsRepository)
    {
        _claimsRepository = claimsRepository;
    }

    public async Task<Result> ValidateAsync(
        PolicyDto policy,
        ClaimType claimType,
        DateTime incidentDate,
        CancellationToken cancellationToken = default)
    {
        var errors = new List<string>();

        // 1. Policy Active Check
        if (policy.Status != InsuranceEngine.Policy.Domain.Enums.PolicyStatus.Active &&
            policy.Status != InsuranceEngine.Policy.Domain.Enums.PolicyStatus.GracePeriod)
        {
            errors.Add($"Policy is not active (current status: {policy.Status}). Claims can only be submitted for Active or GracePeriod policies.");
        }

        // 2. Coverage Period Check
        if (incidentDate.Date < policy.StartDate.Date)
        {
            errors.Add($"Incident date ({incidentDate:yyyy-MM-dd}) is before the policy start date ({policy.StartDate:yyyy-MM-dd}).");
        }

        if (incidentDate.Date > policy.EndDate.Date)
        {
            errors.Add($"Incident date ({incidentDate:yyyy-MM-dd}) is after the policy end date ({policy.EndDate:yyyy-MM-dd}).");
        }

        // 3. Claim Type Coverage Check (M1: simplified — all types allowed for now)
        // In M2, this will check against product configuration for allowed claim types.
        if (claimType == ClaimType.Unspecified)
        {
            errors.Add("Claim type must be specified.");
        }

        // 4. Duplicate Submission Check
        var isDuplicate = await _claimsRepository.ExistsAsync(
            policy.Id, claimType, incidentDate, cancellationToken);

        if (isDuplicate)
        {
            errors.Add($"A claim of type '{claimType}' for incident date {incidentDate:yyyy-MM-dd} already exists for this policy.");
        }

        if (errors.Count > 0)
        {
            return Result.Fail(Error.Validation(string.Join(" | ", errors)));
        }

        return Result.Ok();
    }
}
