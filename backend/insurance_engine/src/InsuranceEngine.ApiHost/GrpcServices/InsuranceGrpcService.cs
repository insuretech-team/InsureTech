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

        // Map back to proto response
        var protoPolicy = new Policy
        {
            PolicyId = result.Value.PolicyId.ToString(),
            PolicyNumber = result.Value.PolicyNumber,
            Status = PolicyStatus.PolicyStatusActive // Assuming active on creation if logic allows
        };

        return new CreatePolicyResponse { Policy = protoPolicy };
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

        request.Claim.ClaimId = result.Value.ToString();
        request.Claim.Status = ClaimStatus.ClaimStatusSubmitted;

        return new CreateClaimResponse { Claim = request.Claim };
    }
}
