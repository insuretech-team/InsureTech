using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Newtonsoft.Json;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.ValueObjects;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using InsuranceEngine.SharedKernel.Services;

namespace InsuranceEngine.Policy.Application.Features.Queries;

// --- GetPolicy ---
public record GetPolicyQuery(Guid Id) : IRequest<PolicyDto?>;

public class GetPolicyQueryHandler : IRequestHandler<GetPolicyQuery, PolicyDto?>
{
    private readonly IPolicyRepository _repo;
    private readonly IEncryptionService _encryptionService;

    public GetPolicyQueryHandler(IPolicyRepository repo, IEncryptionService encryptionService)
    {
        _repo = repo;
        _encryptionService = encryptionService;
    }

    public async Task<PolicyDto?> Handle(GetPolicyQuery request, CancellationToken cancellationToken)
    {
        var p = await _repo.GetByIdWithNomineesAsync(request.Id);
        if (p == null) return null;

        ApplicantDto? applicantDto = null;
        if (!string.IsNullOrEmpty(p.ProposerDetailsJson))
        {
            var applicant = JsonConvert.DeserializeObject<Applicant>(p.ProposerDetailsJson);
            if (applicant != null)
            {
                var nid = !string.IsNullOrEmpty(applicant.NidNumber) ? PiiMasker.MaskNid(_encryptionService.Decrypt(applicant.NidNumber)) : null;
                var phone = !string.IsNullOrEmpty(applicant.PhoneNumber) ? PiiMasker.MaskPhone(_encryptionService.Decrypt(applicant.PhoneNumber)) : null;

                applicantDto = new ApplicantDto(
                    applicant.FullName,
                    applicant.DateOfBirth,
                    nid,
                    applicant.Occupation,
                    applicant.AnnualIncome,
                    applicant.Address,
                    phone,
                    null
                );
            }
        }

        return new PolicyDto(
            Id: p.Id, PolicyNumber: p.PolicyNumber, ProductId: p.ProductId,
            CustomerId: p.CustomerId, PartnerId: p.PartnerId, AgentId: p.AgentId,
            Status: p.Status,
            PremiumAmount: new MoneyDto(p.PremiumAmount, p.PremiumCurrency),
            SumInsured: new MoneyDto(p.SumInsuredAmount, p.SumInsuredCurrency),
            VatTax: p.VatTaxAmount > 0 ? new MoneyDto(p.VatTaxAmount) : null,
            ServiceFee: p.ServiceFeeAmount > 0 ? new MoneyDto(p.ServiceFeeAmount) : null,
            TotalPayable: p.TotalPayableAmount > 0 ? new MoneyDto(p.TotalPayableAmount) : null,
            TenureMonths: p.TenureMonths,
            StartDate: p.StartDate, EndDate: p.EndDate, IssuedAt: p.IssuedAt,
            PaymentFrequency: p.PaymentFrequency,
            ProviderName: p.ProviderName,
            ProposerDetails: applicantDto,
            Nominees: p.Nominees.Where(n => !n.IsDeleted).Select(n => new NomineeDto(
                n.Id,
                n.BeneficiaryId,
                n.Relationship,
                n.SharePercentage
            )).ToList(),
            Riders: p.Riders.Select(r => new PolicyRiderDto(
                r.Id, r.RiderName, new MoneyDto(r.PremiumAmount, r.PremiumCurrency),
                new MoneyDto(r.CoverageAmount, r.CoverageCurrency))).ToList(),
            CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt
        );
    }
}

// --- ListPolicies ---
public record ListPoliciesQuery(
    Guid? CustomerId = null, PolicyStatus? Status = null,
    Guid? ProductId = null, int Page = 1, int PageSize = 20
) : IRequest<PaginatedResponse<PolicyListDto>>;

public class ListPoliciesQueryHandler : IRequestHandler<ListPoliciesQuery, PaginatedResponse<PolicyListDto>>
{
    private readonly IPolicyRepository _repo;

    public ListPoliciesQueryHandler(IPolicyRepository repo) => _repo = repo;

    public async Task<PaginatedResponse<PolicyListDto>> Handle(ListPoliciesQuery request, CancellationToken cancellationToken)
    {
        var (items, totalCount) = await _repo.ListAsync(
            request.CustomerId, request.Status, request.ProductId, request.Page, request.PageSize);

        return new PaginatedResponse<PolicyListDto>(
            Items: items.Select(p => new PolicyListDto(
                p.Id, p.PolicyNumber, p.ProductId, p.CustomerId, p.Status,
                new MoneyDto(p.PremiumAmount, p.PremiumCurrency),
                new MoneyDto(p.SumInsuredAmount, p.SumInsuredCurrency),
                p.StartDate, p.EndDate, p.IssuedAt)).ToList(),
            TotalCount: totalCount, Page: request.Page, PageSize: request.PageSize
        );
    }
}

// --- GetGracePeriod ---
public record GetGracePeriodQuery(Guid PolicyId) : IRequest<GracePeriodDto?>;

public class GetGracePeriodQueryHandler : IRequestHandler<GetGracePeriodQuery, GracePeriodDto?>
{
    private readonly IPolicyRepository _repo;
    private const int GracePeriodDays = 30;

    public GetGracePeriodQueryHandler(IPolicyRepository repo) => _repo = repo;

    public async Task<GracePeriodDto?> Handle(GetGracePeriodQuery request, CancellationToken cancellationToken)
    {
        var p = await _repo.GetByIdAsync(request.PolicyId);
        if (p == null) return null;

        var gracePeriodEndDate = p.EndDate.AddDays(GracePeriodDays);
        var isInGracePeriod = DateTime.UtcNow > p.EndDate && DateTime.UtcNow < gracePeriodEndDate;
        var daysRemaining = isInGracePeriod ? (gracePeriodEndDate - DateTime.UtcNow).Days : 0;

        return new GracePeriodDto(
            PolicyId: p.Id, Status: p.Status, EndDate: p.EndDate,
            GracePeriodEndDate: gracePeriodEndDate,
            DaysRemaining: daysRemaining, IsInGracePeriod: isInGracePeriod);
    }
}

// --- GetRenewalSchedule ---
public record GetRenewalScheduleQuery(Guid PolicyId) : IRequest<RenewalScheduleDto?>;

public class GetRenewalScheduleQueryHandler : IRequestHandler<GetRenewalScheduleQuery, RenewalScheduleDto?>
{
    private readonly IPolicyRepository _repo;

    public GetRenewalScheduleQueryHandler(IPolicyRepository repo) => _repo = repo;

    public async Task<RenewalScheduleDto?> Handle(GetRenewalScheduleQuery request, CancellationToken cancellationToken)
    {
        var p = await _repo.GetByIdAsync(request.PolicyId);
        if (p == null) return null;

        var isEligible = p.Status == PolicyStatus.Active ||
                         p.Status == PolicyStatus.GracePeriod ||
                         p.Status == PolicyStatus.Expired;

        return new RenewalScheduleDto(
            PolicyId: p.Id, PolicyNumber: p.PolicyNumber,
            CurrentEndDate: p.EndDate, NextRenewalDate: p.EndDate,
            EstimatedPremium: new MoneyDto(p.PremiumAmount, p.PremiumCurrency),
            IsEligibleForRenewal: isEligible);
    }
}
