using Microsoft.Extensions.Logging;
using PoliSync.Products.Application.Mappers;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public record GetProductPlanQuery(
    Guid PlanId,
    Guid ProductId
) : IQuery<Insuretech.Products.Entity.V1.ProductPlan>;

public sealed class GetProductPlanQueryHandler : IQueryHandler<GetProductPlanQuery, Insuretech.Products.Entity.V1.ProductPlan>
{
    private readonly IProductRepository _productRepository;
    private readonly ILogger<GetProductPlanQueryHandler> _logger;

    public GetProductPlanQueryHandler(
        IProductRepository productRepository,
        ILogger<GetProductPlanQueryHandler> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public async Task<Result<Insuretech.Products.Entity.V1.ProductPlan>> Handle(GetProductPlanQuery request, CancellationToken ct)
    {
        // Load product
        var product = await _productRepository.GetByIdAsync(request.ProductId, ct);
        if (product is null)
        {
            _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
            return Result<Insuretech.Products.Entity.V1.ProductPlan>.NotFound($"Product '{request.ProductId}' not found");
        }

        // Find plan
        var plan = product.Plans.FirstOrDefault(p => p.Id == request.PlanId);
        if (plan is null)
        {
            _logger.LogWarning("Plan not found: {PlanId} in product {ProductId}", request.PlanId, request.ProductId);
            return Result<Insuretech.Products.Entity.V1.ProductPlan>.NotFound($"Plan '{request.PlanId}' not found");
        }

        // Map to proto
        var protoPlan = plan.ToProto();

        _logger.LogInformation("Product plan retrieved: {PlanId}", request.PlanId);

        return Result<Insuretech.Products.Entity.V1.ProductPlan>.Ok(protoPlan);
    }
}
