using FluentValidation;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Persistence;

namespace PoliSync.Products.Application.Commands;

public record CreateProductPlanCommand(
    Guid ProductId,
    string PlanCode,
    string PlanName,
    string Description,
    long BasePremiumPaisa,
    long SumInsuredPaisa,
    List<string> Features,
    string Currency
) : ICommand<Guid>;

public sealed class CreateProductPlanCommandHandler : ICommandHandler<CreateProductPlanCommand, Guid>
{
    private readonly IProductRepository _productRepository;
    private readonly IUnitOfWork _unitOfWork;
    private readonly ICurrentUser _currentUser;
    private readonly ILogger<CreateProductPlanCommandHandler> _logger;

    public CreateProductPlanCommandHandler(
        IProductRepository productRepository,
        IUnitOfWork unitOfWork,
        ICurrentUser currentUser,
        ILogger<CreateProductPlanCommandHandler> logger)
    {
        _productRepository = productRepository;
        _unitOfWork = unitOfWork;
        _currentUser = currentUser;
        _logger = logger;
    }

    public async Task<Result<Guid>> Handle(CreateProductPlanCommand request, CancellationToken ct)
    {
        // Load product
        var product = await _productRepository.GetByIdAsync(request.ProductId, ct);
        if (product is null)
        {
            _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
            return Result<Guid>.NotFound($"Product '{request.ProductId}' not found");
        }

        // Check tenant authorization
        if (product.TenantId != _currentUser.TenantId)
        {
            _logger.LogWarning("Unauthorized access to product {ProductId} by tenant {TenantId}", 
                request.ProductId, _currentUser.TenantId);
            return Result<Guid>.Unauthorized("You do not have permission to add plans to this product");
        }

        // Create domain plan
        var planResult = ProductPlan.Create(
            request.PlanCode,
            request.PlanName,
            request.Description,
            request.BasePremiumPaisa,
            request.SumInsuredPaisa,
            request.Features,
            request.Currency
        );

        if (!planResult.IsSuccess)
        {
            _logger.LogWarning("Failed to create product plan: {Error}", planResult.Error?.Message);
            return Result<Guid>.Fail(planResult.Error!.Code, planResult.Error.Message);
        }

        var plan = planResult.Value!;

        // Add plan to product and save
        product.AddPlan(plan);
        await _productRepository.UpdateAsync(product, ct);
        await _unitOfWork.CommitAsync(ct);

        _logger.LogInformation("Product plan created: {PlanId} for product {ProductId}", plan.Id, request.ProductId);

        return Result<Guid>.Ok(plan.Id);
    }
}

public sealed class CreateProductPlanCommandValidator : AbstractValidator<CreateProductPlanCommand>
{
    public CreateProductPlanCommandValidator()
    {
        RuleFor(x => x.ProductId)
            .NotEmpty().WithMessage("Product ID is required");

        RuleFor(x => x.PlanCode)
            .NotEmpty().WithMessage("Plan code is required")
            .MaximumLength(50).WithMessage("Plan code must not exceed 50 characters")
            .Matches("^[A-Z0-9_-]+$").WithMessage("Plan code must contain only uppercase letters, numbers, underscores, and hyphens");

        RuleFor(x => x.PlanName)
            .NotEmpty().WithMessage("Plan name is required")
            .MaximumLength(200).WithMessage("Plan name must not exceed 200 characters");

        RuleFor(x => x.BasePremiumPaisa)
            .GreaterThan(0).WithMessage("Base premium must be greater than 0");

        RuleFor(x => x.SumInsuredPaisa)
            .GreaterThan(0).WithMessage("Sum insured must be greater than 0");

        RuleFor(x => x.Currency)
            .NotEmpty().WithMessage("Currency is required");
    }
}
