using FluentValidation;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Persistence;

namespace PoliSync.Products.Application.Commands;

public record AddRiderCommand(
    Guid ProductId,
    string RiderCode,
    string RiderName,
    string Description,
    long PremiumAmountPaisa,
    long SumInsuredPaisa,
    string Category,
    bool IsMandatory,
    string Currency
) : ICommand<Guid>;

public sealed class AddRiderCommandHandler : ICommandHandler<AddRiderCommand, Guid>
{
    private readonly IProductRepository _productRepository;
    private readonly IUnitOfWork _unitOfWork;
    private readonly ICurrentUser _currentUser;
    private readonly ILogger<AddRiderCommandHandler> _logger;

    public AddRiderCommandHandler(
        IProductRepository productRepository,
        IUnitOfWork unitOfWork,
        ICurrentUser currentUser,
        ILogger<AddRiderCommandHandler> logger)
    {
        _productRepository = productRepository;
        _unitOfWork = unitOfWork;
        _currentUser = currentUser;
        _logger = logger;
    }

    public async Task<Result<Guid>> Handle(AddRiderCommand request, CancellationToken ct)
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
            return Result<Guid>.Unauthorized("You do not have permission to add riders to this product");
        }

        // Create domain rider
        var riderResult = Rider.Create(
            request.RiderCode,
            request.RiderName,
            request.Description,
            request.PremiumAmountPaisa,
            request.SumInsuredPaisa,
            request.Category,
            request.IsMandatory,
            request.Currency
        );

        if (!riderResult.IsSuccess)
        {
            _logger.LogWarning("Failed to create rider: {Error}", riderResult.Error?.Message);
            return Result<Guid>.Fail(riderResult.Error!.Code, riderResult.Error.Message);
        }

        var rider = riderResult.Value!;

        // Add rider to product and save
        product.AddRider(rider);
        await _productRepository.UpdateAsync(product, ct);
        await _unitOfWork.CommitAsync(ct);

        _logger.LogInformation("Rider added: {RiderId} to product {ProductId}", rider.Id, request.ProductId);

        return Result<Guid>.Ok(rider.Id);
    }
}

public sealed class AddRiderCommandValidator : AbstractValidator<AddRiderCommand>
{
    public AddRiderCommandValidator()
    {
        RuleFor(x => x.ProductId)
            .NotEmpty().WithMessage("Product ID is required");

        RuleFor(x => x.RiderCode)
            .NotEmpty().WithMessage("Rider code is required")
            .MaximumLength(50).WithMessage("Rider code must not exceed 50 characters")
            .Matches("^[A-Z0-9_-]+$").WithMessage("Rider code must contain only uppercase letters, numbers, underscores, and hyphens");

        RuleFor(x => x.RiderName)
            .NotEmpty().WithMessage("Rider name is required")
            .MaximumLength(200).WithMessage("Rider name must not exceed 200 characters");

        RuleFor(x => x.PremiumAmountPaisa)
            .GreaterThan(0).WithMessage("Premium amount must be greater than 0");

        RuleFor(x => x.SumInsuredPaisa)
            .GreaterThan(0).WithMessage("Sum insured must be greater than 0");

        RuleFor(x => x.Category)
            .NotEmpty().WithMessage("Category is required");

        RuleFor(x => x.Currency)
            .NotEmpty().WithMessage("Currency is required");
    }
}
