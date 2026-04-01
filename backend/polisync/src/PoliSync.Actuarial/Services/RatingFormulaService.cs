using Google.Protobuf.WellKnownTypes;
using Insuretech.Actuarial.Entity.V1;
using Microsoft.Extensions.Logging;
using System.Text.Json;

namespace PoliSync.Actuarial.Services;

public class RatingFormulaService : IRatingFormulaService
{
    private readonly IRatingFormulaRepository _repository;
    private readonly IFormulaEvaluator _evaluator;
    private readonly ILogger<RatingFormulaService> _logger;

    public RatingFormulaService(
        IRatingFormulaRepository repository,
        IFormulaEvaluator evaluator,
        ILogger<RatingFormulaService> logger)
    {
        _repository = repository;
        _evaluator = evaluator;
        _logger = logger;
    }

    public async Task<RatingFormula?> GetFormulaAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        return await _repository.GetByIdAsync(formulaId, cancellationToken);
    }

    public async Task<RatingFormula?> GetFormulaByCodeAsync(string formulaCode, CancellationToken cancellationToken = default)
    {
        return await _repository.GetByCodeAsync(formulaCode, cancellationToken);
    }

    public async Task<IEnumerable<RatingFormula>> GetFormulasAsync(
        string? insuranceType = null,
        FormulaCategory? category = null,
        FormulaStatus? status = null,
        CancellationToken cancellationToken = default)
    {
        return await _repository.GetByFiltersAsync(insuranceType, category, status, cancellationToken);
    }

    public async Task<RatingFormula> CreateFormulaAsync(RatingFormula formula, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Creating rating formula: {FormulaCode}", formula.FormulaCode);
        
        // Validate formula expression
        var isValid = _evaluator.ValidateExpression(formula.FormulaExpression, out var errors);
        if (!isValid)
        {
            _logger.LogWarning("Formula validation failed: {Errors}", string.Join(", ", errors));
        }
        
        formula.FormulaId = Guid.NewGuid().ToString();
        formula.Status = FormulaStatus.Draft;
        formula.Version = 1;
        formula.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        formula.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        return await _repository.CreateAsync(formula, cancellationToken);
    }

    public async Task<RatingFormula?> UpdateFormulaAsync(string formulaId, RatingFormula formula, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Updating rating formula: {FormulaId}", formulaId);
        
        formula.FormulaId = formulaId;
        formula.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        return await _repository.UpdateAsync(formula, cancellationToken);
    }

    public async Task<bool> DeleteFormulaAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Deleting rating formula: {FormulaId}", formulaId);
        return await _repository.DeleteAsync(formulaId, cancellationToken);
    }

    public async Task<RatingFormula?> ActivateFormulaAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Activating rating formula: {FormulaId}", formulaId);
        
        var formula = await _repository.GetByIdAsync(formulaId, cancellationToken);
        if (formula == null)
        {
            return null;
        }
        
        formula.Status = FormulaStatus.Active;
        formula.ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow);
        formula.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        return await _repository.UpdateAsync(formula, cancellationToken);
    }
}
