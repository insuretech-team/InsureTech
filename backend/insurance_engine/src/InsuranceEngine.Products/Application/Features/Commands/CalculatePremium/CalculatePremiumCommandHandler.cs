using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Microsoft.Extensions.Caching.Memory;
using InsuranceEngine.Products.Application.DTOs;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain.Services;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Features.Commands.CalculatePremium;

public class CalculatePremiumCommandHandler : IRequestHandler<CalculatePremiumCommand, Result<CalculatePremiumResponse>>
{
    private readonly IProductRepository _productRepository;
    private readonly PricingEngine _pricingEngine;
    private readonly IMemoryCache _cache;

    public CalculatePremiumCommandHandler(IProductRepository productRepository, PricingEngine pricingEngine, IMemoryCache cache)
    {
        _productRepository = productRepository;
        _pricingEngine = pricingEngine;
        _cache = cache;
    }

    public async Task<Result<CalculatePremiumResponse>> Handle(CalculatePremiumCommand request, CancellationToken cancellationToken)
    {
        // Generate Cache Key for fast quote generation (FR-028)
        var applicantDataHash = request.ApplicantData != null 
            ? string.Join(";", request.ApplicantData.OrderBy(k => k.Key).Select(k => $"{k.Key}={k.Value}")).GetHashCode()
            : 0;
            
        var ridersHash = request.RiderIds != null && request.RiderIds.Any()
            ? string.Join(",", request.RiderIds.OrderBy(id => id)).GetHashCode()
            : 0;
            
        var cacheKey = $"Quote_{request.ProductId}_{request.SumInsuredAmount}_{request.TenureMonths}_{ridersHash}_{applicantDataHash}";

        if (_cache.TryGetValue(cacheKey, out CalculatePremiumResponse cachedResponse))
        {
            return Result<CalculatePremiumResponse>.Ok(cachedResponse);
        }

        var product = await _productRepository.GetByIdWithRidersAsync(request.ProductId);
        if (product == null)
            return Result<CalculatePremiumResponse>.Fail(Error.NotFound("Product", request.ProductId.ToString()));

        if (product.Status != Domain.Enums.ProductStatus.Active)
            return Result<CalculatePremiumResponse>.Fail(Error.Validation("Premium can only be calculated for active products."));

        // Get selected riders
        var selectedRiders = new List<Domain.Rider>();
        if (request.RiderIds != null && request.RiderIds.Any())
        {
            selectedRiders = await _productRepository.GetRidersByIdsAsync(request.RiderIds);
        }

        var result = _pricingEngine.Calculate(
            product,
            request.SumInsuredAmount,
            request.TenureMonths,
            selectedRiders,
            request.ApplicantData);

        var response = new CalculatePremiumResponse(
            BasePremium: new MoneyDto(result.BasePremium.Amount, result.BasePremium.CurrencyCode),
            RiderPremium: new MoneyDto(result.RiderPremium.Amount, result.RiderPremium.CurrencyCode),
            Vat: new MoneyDto(result.Vat.Amount, result.Vat.CurrencyCode),
            ServiceFee: new MoneyDto(result.ServiceFee.Amount, result.ServiceFee.CurrencyCode),
            TotalPremium: new MoneyDto(result.TotalPremium.Amount, result.TotalPremium.CurrencyCode),
            Breakdown: result.Breakdown.Select(b => new PremiumBreakdownDto(
                b.Item,
                new MoneyDto(b.Amount, "BDT"),
                b.Description
            )).ToList()
        );

        var cacheOptions = new MemoryCacheEntryOptions()
            .SetSlidingExpiration(TimeSpan.FromMinutes(15))
            .SetAbsoluteExpiration(TimeSpan.FromHours(1));
            
        _cache.Set(cacheKey, response, cacheOptions);

        return Result<CalculatePremiumResponse>.Ok(response);
    }
}
