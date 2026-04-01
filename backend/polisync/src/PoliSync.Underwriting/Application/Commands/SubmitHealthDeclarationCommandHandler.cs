using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Underwriting.Domain;
using PoliSync.Underwriting.Infrastructure;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;

namespace PoliSync.Underwriting.Application.Commands;

public sealed class SubmitHealthDeclarationCommandHandler
    : IRequestHandler<SubmitHealthDeclarationCommand, Result<SubmitHealthDeclarationResult>>
{
    private readonly IUnderwritingDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly IMediator _mediator;
    private readonly ILogger<SubmitHealthDeclarationCommandHandler> _logger;

    public SubmitHealthDeclarationCommandHandler(
        IUnderwritingDataGateway dataGateway,
        IEventBus eventBus,
        IMediator mediator,
        ILogger<SubmitHealthDeclarationCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result<SubmitHealthDeclarationResult>> Handle(
        SubmitHealthDeclarationCommand request,
        CancellationToken cancellationToken)
    {
        try
        {
            var quote = await _dataGateway.GetQuoteAsync(request.QuoteId, cancellationToken);
            if (quote is null)
                return Result.Fail<SubmitHealthDeclarationResult>("QUOTE_NOT_FOUND", $"Quote {request.QuoteId} not found");

            var aggregateResult = HealthDeclarationAggregate.Create(
                quoteId: request.QuoteId,
                applicantAge: request.ApplicantAge > 0 ? request.ApplicantAge : (int)quote.ApplicantAge,
                heightCm: request.HeightCm,
                weightKg: request.WeightKg,
                hasPreExistingConditions: request.HasPreExistingConditions,
                preExistingConditions: request.PreExistingConditions,
                smoker: request.Smoker,
                alcoholConsumer: request.AlcoholConsumer,
                occupationRiskLevel: request.OccupationRiskLevel,
                isCurrentlyHospitalized: request.IsCurrentlyHospitalized,
                hasFamilyHistory: request.HasFamilyHistory,
                familyHistory: request.FamilyHistory);

            if (aggregateResult.IsFailure)
                return Result.Fail<SubmitHealthDeclarationResult>(
                    aggregateResult.Error!.Code, aggregateResult.Error.Message);

            var aggregate = aggregateResult.Value!;
            var persisted = await _dataGateway.UpsertHealthDeclarationAsync(
                aggregate.Declaration, cancellationToken);

            foreach (var evt in aggregate.DomainEvents)
                await _eventBus.PublishAsync(evt, cancellationToken);

            _logger.LogInformation("Health declaration submitted for quote {QuoteId}. MedExamRequired: {Required}",
                request.QuoteId, persisted.MedicalExamRequired);

            // Trigger underwriting manual review workflow only when medical exam is required
            // (pre-existing conditions, high risk occupation, age > threshold etc.)
            if (persisted.MedicalExamRequired)
            {
                var workflowResult = await _mediator.Send(new TriggerWorkflowCommand(
                    new WorkflowTriggerContext
                    {
                        EntityType  = "UNDERWRITING",
                        EntityId    = persisted.Id,
                        InitiatedBy = "SYSTEM",
                        Portal      = "B2C",
                        Metadata    = new Dictionary<string, string>
                        {
                            ["quote_id"]               = request.QuoteId,
                            ["declaration_id"]         = persisted.Id,
                            ["medical_exam_required"]  = "true",
                            ["has_pre_existing"]       = request.HasPreExistingConditions.ToString().ToLower(),
                            ["occupation_risk"]        = request.OccupationRiskLevel
                        }
                    }), cancellationToken);

                if (workflowResult.IsSuccess && workflowResult.Value!.WasTriggered)
                    _logger.LogInformation(
                        "Underwriting review workflow started: instance={InstanceId} template='{Template}'",
                        workflowResult.Value.WorkflowInstanceId, workflowResult.Value.TemplateName);
            }

            return Result.Ok(new SubmitHealthDeclarationResult(
                persisted.Id,
                persisted.MedicalExamRequired,
                !persisted.MedicalExamRequired));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to submit health declaration for quote {QuoteId}", request.QuoteId);
            return Result.Fail<SubmitHealthDeclarationResult>("SUBMIT_DECLARATION_FAILED", ex.Message);
        }
    }
}
