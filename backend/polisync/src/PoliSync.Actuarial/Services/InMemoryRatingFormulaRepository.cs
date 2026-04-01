using Google.Protobuf.WellKnownTypes;
using Insuretech.Actuarial.Entity.V1;
using Microsoft.Extensions.Logging;
using System.Collections.Concurrent;

namespace PoliSync.Actuarial.Services;

public class InMemoryRatingFormulaRepository : IRatingFormulaRepository
{
    private readonly ConcurrentDictionary<string, RatingFormula> _formulas = new();
    private readonly ILogger<InMemoryRatingFormulaRepository> _logger;

    public InMemoryRatingFormulaRepository(ILogger<InMemoryRatingFormulaRepository> logger)
    {
        _logger = logger;
        SeedData();
    }

    private void SeedData()
    {
        // Auto Insurance Base Rate Formula
        _formulas.TryAdd("formula-001", new RatingFormula
        {
            FormulaId = "formula-001",
            FormulaCode = "AUTO_BASE_RATE",
            FormulaName = "Auto Insurance Base Rate",
            Description = "Base rate calculation for auto insurance",
            Category = FormulaCategory.BaseRate,
            InsuranceType = "AUTO",
            FormulaExpression = "{SUM_INSURED} * 0.005 * {VEHICLE_AGE_FACTOR} * {LOCATION_FACTOR}",
            VariablesJson = "[{\"name\":\"SUM_INSURED\",\"type\":\"DOUBLE\"},{\"name\":\"VEHICLE_AGE_FACTOR\",\"type\":\"DOUBLE\"},{\"name\":\"LOCATION_FACTOR\",\"type\":\"DOUBLE\"}]",
            SortOrder = 1,
            Version = 1,
            Status = FormulaStatus.Active,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });

        // Property Insurance Base Rate
        _formulas.TryAdd("formula-002", new RatingFormula
        {
            FormulaId = "formula-002",
            FormulaCode = "PROPERTY_BASE_RATE",
            FormulaName = "Property Insurance Base Rate",
            Description = "Base rate calculation for property insurance",
            Category = FormulaCategory.BaseRate,
            InsuranceType = "PROPERTY",
            FormulaExpression = "{PROPERTY_VALUE} * 0.002 * {CONSTRUCTION_TYPE_FACTOR} * {LOCATION_RISK_FACTOR}",
            VariablesJson = "[{\"name\":\"PROPERTY_VALUE\",\"type\":\"DOUBLE\"},{\"name\":\"CONSTRUCTION_TYPE_FACTOR\",\"type\":\"DOUBLE\"},{\"name\":\"LOCATION_RISK_FACTOR\",\"type\":\"DOUBLE\"}]",
            SortOrder = 1,
            Version = 1,
            Status = FormulaStatus.Active,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });

        // Liability Loading Factor
        _formulas.TryAdd("formula-003", new RatingFormula
        {
            FormulaId = "formula-003",
            FormulaCode = "LIABILITY_LOADING",
            FormulaName = "Liability Loading",
            Description = "Additional loading for liability coverage",
            Category = FormulaCategory.Loading,
            InsuranceType = "LIABILITY",
            FormulaExpression = "{BASE_PREMIUM} * 0.15 * {BUSINESS_TYPE_FACTOR}",
            VariablesJson = "[{\"name\":\"BASE_PREMIUM\",\"type\":\"DOUBLE\"},{\"name\":\"BUSINESS_TYPE_FACTOR\",\"type\":\"DOUBLE\"}]",
            SortOrder = 50,
            Version = 1,
            Status = FormulaStatus.Active,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });

        // Safe Driver Discount
        _formulas.TryAdd("formula-004", new RatingFormula
        {
            FormulaId = "formula-004",
            FormulaCode = "SAFE_DRIVER_DISCOUNT",
            FormulaName = "Safe Driver Discount",
            Description = "Discount for drivers with no claims",
            Category = FormulaCategory.Discount,
            InsuranceType = "AUTO",
            FormulaExpression = "{BASE_PREMIUM} * 0.10 * {YEARS_NO_CLAIMS} / 5",
            VariablesJson = "[{\"name\":\"BASE_PREMIUM\",\"type\":\"DOUBLE\"},{\"name\":\"YEARS_NO_CLAIMS\",\"type\":\"DOUBLE\"}]",
            SortOrder = 100,
            Version = 1,
            Status = FormulaStatus.Active,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });

        // Loss Ratio Formula
        _formulas.TryAdd("formula-005", new RatingFormula
        {
            FormulaId = "formula-005",
            FormulaCode = "LOSS_RATIO_CALC",
            FormulaName = "Loss Ratio Calculation",
            Description = "Standard loss ratio calculation",
            Category = FormulaCategory.LossRatio,
            InsuranceType = "ALL",
            FormulaExpression = "({INCURRED_LOSSES} + {LAE}) / {EARNED_PREMIUM}",
            VariablesJson = "[{\"name\":\"INCURRED_LOSSES\",\"type\":\"DOUBLE\"},{\"name\":\"LAE\",\"type\":\"DOUBLE\"},{\"name\":\"EARNED_PREMIUM\",\"type\":\"DOUBLE\"}]",
            SortOrder = 1,
            Version = 1,
            Status = FormulaStatus.Active,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });

        // IBNR Reserve Formula
        _formulas.TryAdd("formula-006", new RatingFormula
        {
            FormulaId = "formula-006",
            FormulaCode = "IBNR_RESERVE",
            FormulaName = "IBNR Reserve Calculation",
            Description = "IBNR reserve using chain ladder method",
            Category = FormulaCategory.Reserve,
            InsuranceType = "ALL",
            FormulaExpression = "{REPORTED_CLAIMS} * {DEVELOPMENT_FACTOR} - {REPORTED_CLAIMS}",
            VariablesJson = "[{\"name\":\"REPORTED_CLAIMS\",\"type\":\"DOUBLE\"},{\"name\":\"DEVELOPMENT_FACTOR\",\"type\":\"DOUBLE\"}]",
            SortOrder = 1,
            Version = 1,
            Status = FormulaStatus.Active,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });

        _logger.LogInformation("Seeded {Count} rating formulas", _formulas.Count);
    }

    public Task<RatingFormula?> GetByIdAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        _formulas.TryGetValue(formulaId, out var formula);
        return Task.FromResult(formula);
    }

    public Task<RatingFormula?> GetByCodeAsync(string formulaCode, CancellationToken cancellationToken = default)
    {
        var formula = _formulas.Values.FirstOrDefault(f => f.FormulaCode == formulaCode);
        return Task.FromResult(formula);
    }

    public Task<IEnumerable<RatingFormula>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_formulas.Values.AsEnumerable());
    }

    public Task<IEnumerable<RatingFormula>> GetByFiltersAsync(
        string? insuranceType, 
        FormulaCategory? category, 
        FormulaStatus? status, 
        CancellationToken cancellationToken = default)
    {
        var query = _formulas.Values.AsEnumerable();
        
        if (!string.IsNullOrEmpty(insuranceType))
        {
            query = query.Where(f => f.InsuranceType == insuranceType || f.InsuranceType == "ALL");
        }
        
        if (category.HasValue)
        {
            query = query.Where(f => f.Category == category.Value);
        }
        
        if (status.HasValue)
        {
            query = query.Where(f => f.Status == status.Value);
        }
        
        return Task.FromResult(query);
    }

    public Task<RatingFormula> CreateAsync(RatingFormula formula, CancellationToken cancellationToken = default)
    {
        _formulas.TryAdd(formula.FormulaId, formula);
        return Task.FromResult(formula);
    }

    public Task<RatingFormula?> UpdateAsync(RatingFormula formula, CancellationToken cancellationToken = default)
    {
        if (_formulas.ContainsKey(formula.FormulaId))
        {
            _formulas[formula.FormulaId] = formula;
            return Task.FromResult<RatingFormula?>(formula);
        }
        return Task.FromResult<RatingFormula?>(null);
    }

    public Task<bool> DeleteAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_formulas.TryRemove(formulaId, out _));
    }
}
