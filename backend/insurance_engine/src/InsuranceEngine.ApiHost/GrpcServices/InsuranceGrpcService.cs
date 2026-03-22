using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Policy.Application.Features.Queries;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Claims.Application.DTOs;
using InsuranceEngine.Claims.Application.Features.Queries.Claims;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Application.Features.Queries;
using InsuranceEngine.Products.Application.DTOs;
using Insuretech.Underwriting.Entity.V1;
using MediatR;
using System;
using System.Collections.Generic;
using Google.Protobuf.WellKnownTypes;
using InsuranceEngine.Policy.Application.Features.Commands.CreatePolicy;
using InsuranceEngine.Policy.Application.Features.Commands.IssuePolicy;
using InsuranceEngine.Claims.Application.Features.Commands.Claims;
using InsuranceEngine.Underwriting.Application.Features.Queries.GetQuote;
using InsuranceEngine.Underwriting.Application.Features.Queries.ListQuotes;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using Microsoft.Extensions.Logging;
using System.Linq;
using QuoteProto = Insuretech.Underwriting.Entity.V1;
using ProductProto = Insuretech.Products.Entity.V1;

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
        var protoProduct = new global::Insuretech.Products.Entity.V1.Product
        {
            ProductId = product.Id.ToString(),
            ProductName = product.ProductName,
            ProductCode = product.ProductCode,
            Description = product.Description ?? "",
            BasePremium = new Money { Amount = product.BasePremiumAmount, Currency = product.BasePremiumCurrency },
            Status = product.Status == InsuranceEngine.Products.Domain.Enums.ProductStatus.Active 
                ? ProductProto.ProductStatus.Active 
                : ProductProto.ProductStatus.Inactive
        };

        return new GetProductResponse { Product = protoProduct };
    }

    private static ProductProto.Product MapToProtoProduct(ProductDto dto)
    {
        var product = new ProductProto.Product
        {
            ProductId = dto.Id.ToString(),
            ProductCode = dto.ProductCode,
            ProductName = dto.ProductName,
            Description = dto.Description ?? "",
            Status = dto.Status == InsuranceEngine.Products.Domain.Enums.ProductStatus.Active 
                ? ProductProto.ProductStatus.Active 
                : ProductProto.ProductStatus.Inactive
        };

        if (dto.BasePremium != null)
            product.BasePremium = new Money { Amount = dto.BasePremium.Amount, Currency = dto.BasePremium.CurrencyCode };

        return product;
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
        
        var query = new ListPoliciesQuery(
            customerId,
            null, // Status filter not in proto request
            null, // ProductId not in proto request
            request.Page > 0 ? (int)request.Page : 1,
            request.PageSize > 0 ? (int)request.PageSize : 20
        );

        var result = await _mediator.Send(query);

        var response = new ListPoliciesResponse
        {
            Total = result.TotalCount
        };

        foreach (var p in result.Items)
        {
            response.Policies.Add(new global::Insuretech.Policy.Entity.V1.Policy
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

    private static global::Insuretech.Policy.Entity.V1.Policy MapToProtoPolicy(PolicyDto dto)
    {
        var policy = new global::Insuretech.Policy.Entity.V1.Policy
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
                    NomineeId = n.Id?.ToString() ?? Guid.NewGuid().ToString(),
                    PolicyId = dto.Id.ToString(),
                    FullName = n.FullName,
                    Relationship = n.Relationship,
                    SharePercentage = n.SharePercentage
                });
            }
        }

        if (dto.Riders != null)
        {
            foreach (var r in dto.Riders)
            {
                policy.Riders.Add(new Insuretech.Policy.Entity.V1.Rider
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

    public override async Task<global::Insuretech.Insurance.Services.V1.CreatePolicyResponse> CreatePolicy(CreatePolicyRequest request, ServerCallContext context)
    {
        _logger.LogInformation("Standardized CreatePolicy called via gRPC");

        if (request.Policy == null)
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy data is missing"));

        var command = new CreatePolicyCommand(
            Guid.Parse(request.Policy.ProductId),
            Guid.Parse(request.Policy.CustomerId),
            string.IsNullOrEmpty(request.Policy.PartnerId) ? null : Guid.Parse(request.Policy.PartnerId),
            string.IsNullOrEmpty(request.Policy.AgentId) ? null : Guid.Parse(request.Policy.AgentId),
            new ApplicantDto(
                request.Policy.ProposerDetails?.FullName ?? "",
                request.Policy.ProposerDetails?.DateOfBirth?.ToDateTime(),
                request.Policy.ProposerDetails?.NidNumber,
                request.Policy.ProposerDetails?.Occupation,
                0, // AnnualIncome not in proto
                request.Policy.ProposerDetails?.Address,
                null, // PhoneNumber
                null // HealthDeclaration
            ),
            new List<NomineeDto>(), // Simplified Nominees mapping
            new List<PolicyRiderDto>(), // Simplified Riders mapping
            request.Policy.PremiumAmount?.Amount ?? 0,
            request.Policy.SumInsured?.Amount ?? 0,
            request.Policy.TenureMonths,
            request.Policy.StartDate?.ToDateTime() ?? DateTime.UtcNow
        );

        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Failed to create policy"));

        // Fetch the created policy to return full entity
        var policyDto = await _mediator.Send(new GetPolicyQuery(result.Value.PolicyId));
        
        return new global::Insuretech.Insurance.Services.V1.CreatePolicyResponse { Policy = MapToProtoPolicy(policyDto) };
    }

    public override async Task<UpdatePolicyResponse> UpdatePolicy(UpdatePolicyRequest request, ServerCallContext context)
    {
        _logger.LogInformation($"Standardized UpdatePolicy called for {request.Policy?.PolicyId}");
        
        // This RPC handles business transitions like "Issue" or "Renew" based on status changes in the request
        if (request.Policy?.Status == PolicyStatus.Active)
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
        foreach (var c in result.Value.Items)
        {
            response.Claims.Add(MapToProtoClaim(c));
        }

        return response;
    }

    private static global::Insuretech.Claims.Entity.V1.Claim MapToProtoClaim(InsuranceEngine.Claims.Application.DTOs.ClaimResponseDto dto)
    {
        ClaimStatus status = ClaimStatus.Unspecified;
        if (!string.IsNullOrEmpty(dto.Status.ToString()))
        {
            if (System.Enum.TryParse<ClaimStatus>(dto.Status.ToString(), true, out var parsedStatus))
                status = parsedStatus;
        }

        ClaimType type = ClaimType.Unspecified;
        if (!string.IsNullOrEmpty(dto.Type.ToString()))
        {
            if (System.Enum.TryParse<ClaimType>(dto.Type.ToString(), true, out var parsedType))
                type = parsedType;
        }

        var claim = new global::Insuretech.Claims.Entity.V1.Claim
        {
            ClaimId = dto.Id.ToString(),
            ClaimNumber = dto.ClaimNumber,
            PolicyId = dto.PolicyId.ToString(),
            CustomerId = dto.CustomerId.ToString(),
            Status = status,
            Type = type,
            ClaimedAmount = new Money { Amount = dto.ClaimedAmount.Amount, Currency = dto.ClaimedAmount.CurrencyCode },
            IncidentDate = dto.IncidentDate.ToTimestamp(),
            IncidentDescription = dto.IncidentDescription ?? "",
            PlaceOfIncident = dto.PlaceOfIncident ?? "",
            SubmittedAt = dto.SubmittedAt.ToTimestamp(),
            ApprovedAt = dto.ApprovedAt?.ToTimestamp(),
            SettledAt = dto.SettledAt?.ToTimestamp()
        };

        foreach (var a in dto.Approvals)
        {
            claim.Approvals.Add(new ClaimApproval
            {
                ApprovalId = a.Id.ToString(),
                Decision = System.Enum.TryParse<ApprovalDecision>(a.Decision.ToString(), true, out var decision) ? decision : ApprovalDecision.Unspecified,
                ApprovalLevel = a.ApprovalLevel,
                Notes = a.Notes ?? "",
                ApprovedAt = a.ApprovedAt.ToTimestamp()
            });
        }

        foreach (var d in dto.Documents)
        {
            claim.Documents.Add(new ClaimDocument
            {
                DocumentId = d.Id.ToString(),
                DocumentType = d.DocumentType,
                FileUrl = d.FileUrl,
                Verified = d.Verified,
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
            PolicyId: Guid.Parse(request.Claim.PolicyId),
            CustomerId: Guid.Parse(request.Claim.CustomerId),
            Type: (InsuranceEngine.Claims.Domain.Enums.ClaimType)request.Claim.Type, 
            ClaimedAmount: request.Claim.ClaimedAmount?.Amount ?? 0,
            IncidentDate: request.Claim.IncidentDate?.ToDateTime() ?? DateTime.UtcNow,
            IncidentDescription: request.Claim.IncidentDescription,
            PlaceOfIncident: request.Claim.PlaceOfIncident,
            BankDetailsForPayout: null,
            Documents: new List<ClaimDocumentDto>()
        );

        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Claim submission failed"));

        // Fetch the created claim to return full entity
        var claimResult = await _mediator.Send(new GetClaimByIdQuery(result.Value));
        if (!claimResult.IsSuccess)
            throw new RpcException(new Status(StatusCode.Internal, "Claim created but failed to retrieve"));

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
            Total = (int)result.TotalCount
        };

        foreach (var q in result.Items)
        {
            response.Quotes.Add(MapToProtoQuote(q));
        }

        return response;
    }

    private static global::Insuretech.Underwriting.Entity.V1.Quote MapToProtoQuote(InsuranceEngine.Underwriting.Application.DTOs.QuoteDto dto)
    {
        return new global::Insuretech.Underwriting.Entity.V1.Quote
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
