using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Products.Application.Interfaces;
using Microsoft.Extensions.Logging;
using System.Linq;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Policy.Application.Features.Commands.IssuePolicy;
using InsuranceEngine.Products.Application.Features.Commands.CalculatePremium;
using InsuranceEngine.Policy.Application.Features.Commands.RenewPolicy;
using InsuranceEngine.Policy.Application.Features.Commands.Claims;
using InsuranceEngine.Policy.Application.Features.Commands.CreatePolicy;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using System.Collections.Generic;
using System;
using Google.Protobuf.WellKnownTypes;
using InsuranceEngine.Policy.Application.Features.Queries.GetQuote;
using InsuranceEngine.Policy.Application.Features.Queries.ListQuotes;
using Insuretech.Underwriting.Entity.V1;

namespace InsuranceEngine.ApiHost.GrpcServices;

/// <summary>
/// Standardized gRPC service implementing the master InsuranceService contract.
/// Maps internal domain logic to the project-wide standard proto definitions.
/// </summary>
public class InsuranceGrpcService : InsuranceService.InsuranceServiceBase
{
    private readonly IProductRepository _productRepository;
    private readonly IMediator _mediator;
    private readonly ILogger<InsuranceGrpcService> _logger;

    public InsuranceGrpcService(
        IProductRepository productRepository, 
        IMediator mediator,
        ILogger<InsuranceGrpcService> logger)
    {
        _productRepository = productRepository;
        _mediator = mediator;
        _logger = logger;
    }

    // ===================== PRODUCT CRUD =====================

    public override async Task<GetProductResponse> GetProduct(GetProductRequest request, ServerCallContext context)
    {
        _logger.LogInformation($"Standardized GetProduct called for {request.ProductId}");
        
        if (!Guid.TryParse(request.ProductId, out var productId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid Product ID format"));
        }

        var product = await _productRepository.GetByIdAsync(productId);
        if (product == null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, "Product not found"));
        }

        // Map domain product to proto product
        var protoProduct = new Product
        {
            ProductId = product.Id.ToString(),
            Name = product.Name,
            Code = product.Code,
            Description = product.Description,
            BasePremium = new Money { Amount = product.BasePremium?.Amount ?? 0, Currency = product.BasePremium?.CurrencyCode ?? "BDT" },
            IsActive = product.IsActive
        };

        return new GetProductResponse { Product = protoProduct };
    }

    // ===================== POLICY CRUD =====================

    public override async Task<GetPolicyResponse> GetPolicy(GetPolicyRequest request, ServerCallContext context)
    {
        _logger.LogInformation($"Standardized GetPolicy called for {request.PolicyId}");

        if (!Guid.TryParse(request.PolicyId, out var policyId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid Policy ID format"));

        var policyDto = await _mediator.Send(new GetPolicyQuery(policyId));
        if (policyDto == null)
            throw new RpcException(new Status(StatusCode.NotFound, "Policy not found"));

        return new GetPolicyResponse { Policy = MapToProtoPolicy(policyDto) };
    }

    public override async Task<ListPoliciesResponse> ListPolicies(ListPoliciesRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Standardized ListPolicies called");

        Guid? customerId = string.IsNullOrEmpty(request.CustomerId) ? null : Guid.Parse(request.CustomerId);
        Guid? productId = string.IsNullOrEmpty(request.ProductId) ? null : Guid.Parse(request.ProductId);
        
        var query = new ListPoliciesQuery(
            customerId,
            null, // Status mapping could be added if needed
            productId,
            request.Page > 0 ? (int)request.Page : 1,
            request.PageSize > 0 ? (int)request.PageSize : 20
        );

        var result = await _mediator.Send(query);

        var response = new ListPoliciesResponse
        {
            TotalCount = (uint)result.TotalCount,
            Page = (uint)result.Page,
            PageSize = (uint)result.PageSize
        };

        foreach (var p in result.Items)
        {
            response.Policies.Add(new Policy
            {
                PolicyId = p.Id.ToString(),
                PolicyNumber = p.PolicyNumber,
                ProductId = p.ProductId.ToString(),
                CustomerId = p.CustomerId.ToString(),
                Status = (PolicyStatus)p.Status,
                PremiumAmount = new Money { Amount = p.PremiumAmount.Amount, Currency = p.PremiumAmount.CurrencyCode },
                SumInsured = new Money { Amount = p.SumInsured.Amount, Currency = p.SumInsured.CurrencyCode },
                StartDate = p.StartDate.ToTimestamp(),
                EndDate = p.EndDate.ToTimestamp(),
                IssuedAt = p.IssuedAt?.ToTimestamp()
            });
        }

        return response;
    }

    private static Policy MapToProtoPolicy(PolicyDto dto)
    {
        var policy = new Policy
        {
            PolicyId = dto.Id.ToString(),
            PolicyNumber = dto.PolicyNumber,
            ProductId = dto.ProductId.ToString(),
            CustomerId = dto.CustomerId.ToString(),
            PartnerId = dto.PartnerId?.ToString() ?? "",
            AgentId = dto.AgentId?.ToString() ?? "",
            Status = (PolicyStatus)dto.Status,
            PremiumAmount = new Money { Amount = dto.PremiumAmount.Amount, Currency = dto.PremiumAmount.CurrencyCode },
            SumInsured = new Money { Amount = dto.SumInsured.Amount, Currency = dto.SumInsured.CurrencyCode },
            TenureMonths = dto.TenureMonths,
            StartDate = dto.StartDate.ToTimestamp(),
            EndDate = dto.EndDate.ToTimestamp(),
            IssuedAt = dto.IssuedAt?.ToTimestamp(),
            CreatedAt = dto.CreatedAt.ToTimestamp(),
            UpdatedAt = dto.UpdatedAt.ToTimestamp(),
            PaymentFrequency = dto.PaymentFrequency ?? "",
            ProviderName = dto.ProviderName ?? ""
        };

        if (dto.VatTax != null) policy.VatTax = new Money { Amount = dto.VatTax.Amount, Currency = dto.VatTax.CurrencyCode };
        if (dto.ServiceFee != null) policy.ServiceFee = new Money { Amount = dto.ServiceFee.Amount, Currency = dto.ServiceFee.CurrencyCode };
        if (dto.TotalPayable != null) policy.TotalPayable = new Money { Amount = dto.TotalPayable.Amount, Currency = dto.TotalPayable.CurrencyCode };

        if (dto.ProposerDetails != null)
        {
            policy.ProposerDetails = new Applicant
            {
                FullName = dto.ProposerDetails.FullName,
                DateOfBirth = dto.ProposerDetails.DateOfBirth?.ToTimestamp(),
                NidNumber = dto.ProposerDetails.NidNumber ?? "",
                Occupation = dto.ProposerDetails.Occupation ?? "",
                Address = dto.ProposerDetails.Address ?? ""
            };
        }

        if (dto.Nominees != null)
        {
            foreach (var n in dto.Nominees)
            {
                policy.Nominees.Add(new Nominee
                {
                    NomineeId = n.Id?.ToString() ?? "",
                    BeneficiaryId = n.BeneficiaryId.ToString(),
                    Relationship = n.Relationship,
                    SharePercentage = n.SharePercentage
                });
            }
        }

        if (dto.Riders != null)
        {
            foreach (var r in dto.Riders)
            {
                policy.Riders.Add(new Rider
                {
                    RiderId = r.Id.ToString(),
                    RiderName = r.RiderName,
                    PremiumAmount = new Money { Amount = r.PremiumAmount.Amount, Currency = r.PremiumAmount.CurrencyCode },
                    CoverageAmount = new Money { Amount = r.CoverageAmount.Amount, Currency = r.CoverageAmount.CurrencyCode }
                });
            }
        }

        return policy;
    }

    public override async Task<CreatePolicyResponse> CreatePolicy(CreatePolicyRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Standardized CreatePolicy called via gRPC");

        if (request.Policy == null)
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy data is missing"));

        // Map proto policy request to internal command
        var command = new CreatePolicyCommand(
            Guid.Parse(request.Policy.ProductId),
            Guid.Parse(request.Policy.CustomerId),
            (int)request.Policy.TenureMonths,
            request.Policy.SumInsured?.Amount ?? 0,
            new List<Guid>() // Riders mapping omitted for brevity in first pass
        );

        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Failed to create policy"));

        // Fetch the created policy to return full entity
        var policyDto = await _mediator.Send(new GetPolicyQuery(result.Value.PolicyId));
        
        return new CreatePolicyResponse { Policy = MapToProtoPolicy(policyDto) };
    }

    public override async Task<UpdatePolicyResponse> UpdatePolicy(UpdatePolicyRequest request, ServerCallContext context)
    {
        _logger.LogInformation($"Standardized UpdatePolicy called for {request.Policy?.PolicyId}");
        
        // This RPC handles business transitions like "Issue" or "Renew" based on status changes in the request
        if (request.Policy?.Status == PolicyStatus.PolicyStatusActive)
        {
             // Mapping to existing IssuePolicy logic
             var result = await _mediator.Send(new IssuePolicyCommand(Guid.Parse(request.Policy.PolicyId)));
             
             if (!result.IsSuccess)
                throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Failed to issue policy"));

             return new UpdatePolicyResponse { Policy = request.Policy };
        }

        throw new RpcException(new Status(StatusCode.Unimplemented, "Generic UpdatePolicy transitions not yet fully implemented"));
    }

    // ===================== CLAIM CRUD =====================
    
    public override async Task<GetClaimResponse> GetClaim(GetClaimRequest request, ServerCallContext context)
    {
        _logger.LogInformation($"Standardized GetClaim called for {request.ClaimId}");

        if (!Guid.TryParse(request.ClaimId, out var claimId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid Claim ID format"));

        var result = await _mediator.Send(new GetClaimByIdQuery(claimId));
        if (!result.IsSuccess)
            throw new RpcException(new Status(StatusCode.NotFound, result.Error?.Message ?? "Claim not found"));

        return new GetClaimResponse { Claim = MapToProtoClaim(result.Value) };
    }

    public override async Task<ListClaimsResponse> ListClaims(ListClaimsRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Standardized ListClaims called");

        if (!Guid.TryParse(request.CustomerId, out var customerId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid Customer ID format"));

        var result = await _mediator.Send(new ListClaimsByCustomerQuery(
            customerId, 
            request.Page > 0 ? (int)request.Page : 1, 
            request.PageSize > 0 ? (int)request.PageSize : 10
        ));

        if (!result.IsSuccess)
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Failed to list claims"));

        var response = new ListClaimsResponse();
        foreach (var c in result.Value)
        {
            response.Claims.Add(MapToProtoClaim(c));
        }

        return response;
    }

    private static Claim MapToProtoClaim(ClaimResponseDto dto)
    {
        var claim = new Claim
        {
            ClaimId = dto.Id.ToString(),
            ClaimNumber = dto.ClaimNumber,
            PolicyId = dto.PolicyId.ToString(),
            CustomerId = dto.CustomerId.ToString(),
            Status = Enum.TryParse<ClaimStatus>($"ClaimStatus{dto.Status}", out var status) ? status : ClaimStatus.ClaimStatusUnspecified,
            Type = Enum.TryParse<ClaimType>($"ClaimType{dto.ClaimType}", out var type) ? type : ClaimType.ClaimTypeUnspecified,
            ClaimedAmount = new Money { Amount = dto.ClaimedAmount, Currency = dto.Currency },
            IncidentDate = dto.IncidentDate.ToTimestamp(),
            IncidentDescription = dto.IncidentDescription,
            PlaceOfIncident = dto.PlaceOfIncident ?? "",
            SubmittedAt = dto.SubmittedAt.ToTimestamp()
        };

        foreach (var a in dto.Approvals)
        {
            claim.Approvals.Add(new ClaimApproval
            {
                ApprovalId = a.Id.ToString(),
                Decision = Enum.TryParse<ClaimApprovalDecision>($"ClaimApprovalDecision{a.Decision}", out var decision) ? decision : ClaimApprovalDecision.ClaimApprovalDecisionUnspecified,
                ApprovalLevel = (uint)a.Level,
                Notes = a.Notes ?? "",
                DecidedAt = a.DecidedAt?.ToTimestamp()
            });
        }

        foreach (var d in dto.Documents)
        {
            claim.Documents.Add(new ClaimDocument
            {
                DocumentId = d.Id.ToString(),
                DocumentType = d.DocumentType,
                FileUrl = d.FileUrl,
                IsVerified = d.IsVerified,
                UploadedAt = d.UploadedAt.ToTimestamp()
            });
        }

        return claim;
    }

    public override async Task<CreateClaimResponse> CreateClaim(CreateClaimRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Standardized CreateClaim called via gRPC");

        if (request.Claim == null)
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim data is missing"));

        // Map proto claim to SubmitClaimCommand
        var command = new SubmitClaimCommand(
            Guid.Parse(request.Claim.PolicyId),
            Guid.Parse(request.Claim.CustomerId),
            (InsuranceEngine.Policy.Domain.Enums.ClaimType)request.Claim.Type, 
            request.Claim.ClaimedAmount?.Amount ?? 0,
            request.Claim.IncidentDate?.ToDateTime() ?? DateTime.UtcNow,
            request.Claim.IncidentDescription,
            request.Claim.PlaceOfIncident
        );

        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Claim submission failed"));

        // Fetch the created claim to return full entity
        var claimResult = await _mediator.Send(new GetClaimByIdQuery(result.Value));
        
        return new CreateClaimResponse { Claim = MapToProtoClaim(claimResult.Value) };
    }

    // ===================== QUOTE CRUD =====================

    public override async Task<GetQuoteResponse> GetQuote(GetQuoteRequest request, ServerCallContext context)
    {
        _logger.LogInformation($"Standardized GetQuote called for {request.QuoteId}");

        if (!Guid.TryParse(request.QuoteId, out var quoteId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid Quote ID format"));

        var result = await _mediator.Send(new GetQuoteQuery(quoteId));
        if (!result.IsSuccess)
            throw new RpcException(new Status(StatusCode.NotFound, result.Error?.Message ?? "Quote not found"));

        return new GetQuoteResponse { Quote = MapToProtoQuote(result.Value) };
    }

    public override async Task<ListQuotesResponse> ListQuotes(ListQuotesRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Standardized ListQuotes called");

        Guid? beneficiaryId = string.IsNullOrEmpty(request.BeneficiaryId) ? null : Guid.Parse(request.BeneficiaryId);
        
        var query = new ListQuotesQuery(
            beneficiaryId,
            null, // Status mapping could be added
            request.Page > 0 ? (int)request.Page : 1,
            request.PageSize > 0 ? (int)request.PageSize : 20
        );

        var result = await _mediator.Send(query);

        var response = new ListQuotesResponse
        {
            TotalCount = (uint)result.TotalCount,
            Page = (uint)result.Page,
            PageSize = (uint)result.PageSize
        };

        foreach (var q in result.Items)
        {
            response.Quotes.Add(MapToProtoQuote(q));
        }

        return response;
    }

    private static Quote MapToProtoQuote(QuoteDto dto)
    {
        return new Quote
        {
            Id = dto.Id.ToString(),
            QuoteNumber = dto.QuoteNumber,
            BeneficiaryId = dto.BeneficiaryId.ToString(),
            InsurerProductId = dto.InsurerProductId.ToString(),
            Status = (Insuretech.Underwriting.Entity.V1.QuoteStatus)dto.Status,
            SumAssured = new Money { Amount = dto.SumAssured.Amount, Currency = dto.SumAssured.CurrencyCode },
            TermYears = dto.TermYears,
            PremiumPaymentMode = dto.PremiumPaymentMode,
            BasePremium = new Money { Amount = dto.BasePremium.Amount, Currency = dto.BasePremium.CurrencyCode },
            RiderPremium = new Money { Amount = dto.RiderPremium.Amount, Currency = dto.RiderPremium.CurrencyCode },
            TotalPremium = new Money { Amount = dto.TotalPremium.Amount, Currency = dto.TotalPremium.CurrencyCode },
            ApplicantAge = dto.ApplicantAge,
            ApplicantOccupation = dto.ApplicantOccupation ?? "",
            Smoker = dto.IsSmoker,
            ValidUntil = dto.ValidUntil.ToTimestamp()
        };
    }
}
