using FluentValidation;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.SharedKernel.Persistence;

namespace PoliSync.Products.Application.Commands;

public record RuleConditionDto(
    string Field,
    string Operator,
    string Value
);

public record RuleActionDto(
    string ActionType,
    decimal Value
);

public record PricingRuleDto(
    string RuleId,
    string RuleName,
    string RuleType,
    int Priority,
    bool ApplyAll,
    List<RuleConditionDto> Conditions,
    RuleActionDto Action
);

public record CreatePricingConfigCommand(
    Guid ProductId,
    List<PricingRuleDto> Rules,
    DateTimeOffset EffectiveFrom,
    DateTimeOffset? EffectiveTo
) : ICommand<Guid>;

public sealed class CreatePricingConfigCommandHandler : ICommandHandler<CreatePricingConfigCommand, Guid>
{
    private readonly IProductRepository _productRepository;
    private readonly IUnitOfWork _unitOfWork;
    private readonly ICurrentUser _currentUser;
    private readonly IEventBus _eventBus;
    private readonly ILogger<CreatePricingConfigCommandHandler> _logger;

    public CreatePricingConfigCommandHandler(
        IProductRepository productRepository,
        IUnitOfWork unitOfWork,
        ICurrentUser currentUser,
        IEventBus eventBus,
        ILogger<CreatePricingConfigCommandHandler> logger)
    {
        _productRepository = productRepository;
        _unitOfWork = unitOfWork;
        _currentUser = currentUser;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<Guid>> Handle(CreatePricingConfigCommand request, CancellationToken ct)
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
            return Result<Guid>.Unauthorized("You do not have permission to configure pricing for this product");
        }

        // Create domain pricing rules from DTOs
        var rules = new List<PricingRule>();
        foreach (var ruleDto in request.Rules)
        {
            var conditions = ruleDto.Conditions
                .Select(c => new RuleCondition(c.Field, c.Operator, c.Value))
                .ToList();

            var action = new RuleAction(ruleDto.Action.ActionType, ruleDto.Action.Value);

            var rule = new PricingRule(
                Guid.Parse(ruleDto.RuleId),
                ruleDto.RuleName,
                ruleDto.RuleType,
                ruleDto.Priority,
                ruleDto.ApplyAll,
                conditions,
                action
            );

            rules.Add(rule);
        }

        // Create domain pricing config
        var configResult = PricingConfig.Create(
            request.ProductId,
            rules,
            request.EffectiveFrom,
            request.EffectiveTo
        );

        if (!configResult.IsSuccess)
        {
            _logger.LogWarning("Failed to create pricing config: {Error}", configResult.Error?.Message);
            return Result<Guid>.Fail(configResult.Error!.Code, configResult.Error.Message);
        }

        var config = configResult.Value!;

        // Add config to product and save
        product.SetPricingConfig(config);
        await _productRepository.UpdateAsync(product, ct);
        await _unitOfWork.CommitAsync(ct);

        _logger.LogInformation("Pricing config created: {ConfigId} for product {ProductId}", config.Id, request.ProductId);

        // Publish Kafka event
        var kafkaEvent = new
        {
            product_id = product.Id.ToString(),
            tenant_id = product.TenantId.ToString(),
            partner_id = product.PartnerId?.ToString(),
            config_id = config.Id.ToString(),
            rule_count = rules.Count,
            effective_from = request.EffectiveFrom.ToUnixTimeSeconds(),
            effective_to = request.EffectiveTo?.ToUnixTimeSeconds(),
            updated_at = DateTimeOffset.UtcNow.ToUnixTimeSeconds()
        };

        try
        {
            await _eventBus.PublishAsync("insuretech.product.pricing_updated.v1", product.Id.ToString(), kafkaEvent, ct);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to publish pricing updated event for {ProductId}", product.Id);
            // Don't fail the request if event publish fails
        }

        return Result<Guid>.Ok(config.Id);
    }
}

public sealed class CreatePricingConfigCommandValidator : AbstractValidator<CreatePricingConfigCommand>
{
    public CreatePricingConfigCommandValidator()
    {
        RuleFor(x => x.ProductId)
            .NotEmpty().WithMessage("Product ID is required");

        RuleFor(x => x.Rules)
            .NotEmpty().WithMessage("At least one pricing rule is required");

        RuleFor(x => x.EffectiveFrom)
            .NotEmpty().WithMessage("Effective from date is required");

        RuleFor(x => x.EffectiveTo)
            .GreaterThan(x => x.EffectiveFrom)
            .When(x => x.EffectiveTo.HasValue)
            .WithMessage("Effective to date must be after effective from date");
    }
}
