using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Actuarial.Entity.V1;
using Insuretech.Actuarial.Services.V1;
using Microsoft.Extensions.Logging;
using PoliSync.Actuarial.Services;

namespace PoliSync.Actuarial.GrpcServices;

public class ActuarialGrpcService : ActuarialService.ActuarialServiceBase
{
    private readonly IActuarialService _actuarialService;
    private readonly IRatingFormulaService _formulaService;
    private readonly IReserveCalculationService _reserveService;
    private readonly ILossRatioService _lossRatioService;
    private readonly IFormulaEvaluator _formulaEvaluator;
    private readonly ILogger<ActuarialGrpcService> _logger;

    public ActuarialGrpcService(
        IActuarialService actuarialService,
        IRatingFormulaService formulaService,
        IReserveCalculationService reserveService,
        ILossRatioService lossRatioService,
        IFormulaEvaluator formulaEvaluator,
        ILogger<ActuarialGrpcService> logger)
    {
        _actuarialService = actuarialService;
        _formulaService = formulaService;
        _reserveService = reserveService;
        _lossRatioService = lossRatioService;
        _formulaEvaluator = formulaEvaluator;
        _logger = logger;
    }

    // Premium calculations
    public override async Task<CalculatePremiumResponse> CalculatePremium(
        CalculatePremiumRequest request,
        ServerCallContext context)
    {
        try
        {
            var result = await _actuarialService.CalculatePremiumAsync(request.Input, context.CancellationToken);
            
            return new CalculatePremiumResponse
            {
                CalculationId = Guid.NewGuid().ToString(),
                CalculationReference = request.CalculationReference ?? $"ACT-{DateTime.UtcNow:yyyy}-{Guid.NewGuid().ToString()[..6].ToUpper()}",
                Result = result,
                Success = true,
                CalculatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error calculating premium");
            return new CalculatePremiumResponse
            {
                Success = false,
                Errors = { ex.Message }
            };
        }
    }

    public override async Task<CalculatePremiumResponse> CalculatePurePremium(
        CalculatePurePremiumRequest request,
        ServerCallContext context)
    {
        try
        {
            var result = await _actuarialService.CalculatePurePremiumAsync(
                request.ExpectedClaims,
                request.ClaimSeverity,
                request.ExposureUnits,
                request.RiskAdjustmentFactor);
            
            return new CalculatePremiumResponse
            {
                CalculationId = Guid.NewGuid().ToString(),
                CalculationReference = request.CalculationReference ?? $"ACT-{DateTime.UtcNow:yyyy}-{Guid.NewGuid().ToString()[..6].ToUpper()}",
                Result = new PremiumCalculationResult
                {
                    GrossPremium = result,
                    NetPremium = result,
                    BasePremium = result,
                    Currency = "USD"
                },
                Success = true,
                CalculatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error calculating pure premium");
            return new CalculatePremiumResponse
            {
                Success = false,
                Errors = { ex.Message }
            };
        }
    }

    // Rating formula evaluation
    public override async Task<EvaluateRatingFormulaResponse> EvaluateRatingFormula(
        EvaluateRatingFormulaRequest request,
        ServerCallContext context)
    {
        try
        {
            var stopwatch = System.Diagnostics.Stopwatch.StartNew();
            
            // Get formula
            var formula = !string.IsNullOrEmpty(request.FormulaId) 
                ? await _formulaService.GetFormulaAsync(request.FormulaId, context.CancellationToken)
                : await _formulaService.GetFormulaByCodeAsync(request.FormulaCode, context.CancellationToken);
            
            if (formula == null)
            {
                return new EvaluateRatingFormulaResponse
                {
                    Success = false,
                    Errors = { "Formula not found" }
                };
            }
            
            // Convert Struct to Dictionary
            var variables = request.Variables.Fields.ToDictionary(
                f => f.Key,
                f => f.Value.KindCase == Value.KindOneofCase.NumberValue ? f.Value.NumberValue : 0.0);
            
            // Evaluate formula
            var result = _formulaEvaluator.Evaluate(formula.FormulaExpression, variables);
            
            stopwatch.Stop();
            
            return new EvaluateRatingFormulaResponse
            {
                CalculationId = Guid.NewGuid().ToString(),
                Success = true,
                Result = result,
                OutputVariables = request.Variables,
                ExecutionTimeMs = stopwatch.ElapsedMilliseconds,
                CalculatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error evaluating rating formula");
            return new EvaluateRatingFormulaResponse
            {
                Success = false,
                Errors = { ex.Message }
            };
        }
    }

    public override Task<ValidateFormulaExpressionResponse> ValidateFormulaExpression(
        ValidateFormulaExpressionRequest request,
        ServerCallContext context)
    {
        try
        {
            var isValid = _formulaEvaluator.ValidateExpression(request.FormulaExpression, out var errors);
            var variables = _formulaEvaluator.ExtractVariables(request.FormulaExpression);
            
            return Task.FromResult(new ValidateFormulaExpressionResponse
            {
                IsValid = isValid,
                Errors = { errors },
                ParsedVariables = { variables }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error validating formula expression");
            return Task.FromResult(new ValidateFormulaExpressionResponse
            {
                IsValid = false,
                Errors = { ex.Message }
            });
        }
    }

    // Reserve calculations
    public override async Task<CalculateReservesResponse> CalculateReserves(
        CalculateReservesRequest request,
        ServerCallContext context)
    {
        try
        {
            var result = await _reserveService.CalculateReservesAsync(request.Input, context.CancellationToken);
            
            return new CalculateReservesResponse
            {
                ReserveId = Guid.NewGuid().ToString(),
                CalculationReference = request.CalculationReference ?? $"RES-{DateTime.UtcNow:yyyy}-{Guid.NewGuid().ToString()[..6].ToUpper()}",
                Result = result,
                Success = true,
                CalculatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error calculating reserves");
            return new CalculateReservesResponse
            {
                Success = false,
                Errors = { ex.Message }
            };
        }
    }

    public override async Task<GetReserveCalculationResponse> GetReserveCalculation(
        GetReserveCalculationRequest request,
        ServerCallContext context)
    {
        try
        {
            ReserveCalculation? reserve = null;
            
            if (!string.IsNullOrEmpty(request.ReserveId))
            {
                reserve = await _reserveService.GetReserveAsync(request.ReserveId, context.CancellationToken);
            }
            else if (!string.IsNullOrEmpty(request.ClaimId))
            {
                reserve = await _reserveService.GetReserveByClaimAsync(request.ClaimId, context.CancellationToken);
            }
            
            return new GetReserveCalculationResponse
            {
                Reserve = reserve,
                Found = reserve != null
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting reserve calculation");
            return new GetReserveCalculationResponse
            {
                Found = false
            };
        }
    }

    // Loss ratio analysis
    public override async Task<CalculateLossRatioResponse> CalculateLossRatio(
        CalculateLossRatioRequest request,
        ServerCallContext context)
    {
        try
        {
            var result = await _lossRatioService.CalculateLossRatioAsync(request.Input, context.CancellationToken);
            
            return new CalculateLossRatioResponse
            {
                LossRatioId = Guid.NewGuid().ToString(),
                CalculationReference = request.CalculationReference ?? $"LR-{DateTime.UtcNow:yyyy}-{Guid.NewGuid().ToString()[..6].ToUpper()}",
                Result = result,
                Success = true,
                CalculatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error calculating loss ratio");
            return new CalculateLossRatioResponse
            {
                Success = false,
                Errors = { ex.Message }
            };
        }
    }

    public override Task<AnalyzeLossTrendsResponse> AnalyzeLossTrends(
        AnalyzeLossTrendsRequest request,
        ServerCallContext context)
    {
        // Simplified implementation - in production, this would analyze actual historical data
        var segments = new List<LossTrendSegment>
        {
            new LossTrendSegment
            {
                SegmentKey = "overall",
                SegmentValue = "all_products",
                CurrentLossRatio = 0.68,
                PreviousLossRatio = 0.72,
                ChangePercentage = -5.6,
                TrendDirection = "IMPROVING",
                Details = new Struct()
            }
        };
        
        return Task.FromResult(new AnalyzeLossTrendsResponse
        {
            Success = true,
            Segments = { segments },
            AnalysisDate = Timestamp.FromDateTime(DateTime.UtcNow)
        });
    }

    // Rating formula management
    public override async Task<CreateRatingFormulaResponse> CreateRatingFormula(
        CreateRatingFormulaRequest request,
        ServerCallContext context)
    {
        try
        {
            var formula = new RatingFormula
            {
                FormulaId = Guid.NewGuid().ToString(),
                FormulaCode = request.FormulaCode,
                FormulaName = request.FormulaName,
                Description = request.Description,
                Category = request.Category,
                InsuranceType = request.InsuranceType,
                FormulaExpression = request.FormulaExpression,
                VariablesJson = System.Text.Json.JsonSerializer.Serialize(request.Variables),
                SortOrder = request.SortOrder,
                Status = FormulaStatus.Draft,
                ValidFrom = request.ValidFrom ?? Timestamp.FromDateTime(DateTime.UtcNow),
                ValidUntil = request.ValidUntil,
                Metadata = { request.Metadata },
                CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
                UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
            };
            
            var created = await _formulaService.CreateFormulaAsync(formula, context.CancellationToken);
            
            return new CreateRatingFormulaResponse
            {
                FormulaId = created.FormulaId,
                Success = true,
                Formula = created
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating rating formula");
            return new CreateRatingFormulaResponse
            {
                Success = false,
                Errors = { ex.Message }
            };
        }
    }

    public override async Task<GetRatingFormulaResponse> GetRatingFormula(
        GetRatingFormulaRequest request,
        ServerCallContext context)
    {
        try
        {
            RatingFormula? formula = null;
            
            if (!string.IsNullOrEmpty(request.FormulaId))
            {
                formula = await _formulaService.GetFormulaAsync(request.FormulaId, context.CancellationToken);
            }
            else if (!string.IsNullOrEmpty(request.FormulaCode))
            {
                formula = await _formulaService.GetFormulaByCodeAsync(request.FormulaCode, context.CancellationToken);
            }
            
            return new GetRatingFormulaResponse
            {
                Formula = formula,
                Found = formula != null
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting rating formula");
            return new GetRatingFormulaResponse
            {
                Found = false
            };
        }
    }

    public override async Task<UpdateRatingFormulaResponse> UpdateRatingFormula(
        UpdateRatingFormulaRequest request,
        ServerCallContext context)
    {
        try
        {
            var formula = new RatingFormula
            {
                FormulaId = request.FormulaId,
                FormulaName = request.FormulaName,
                Description = request.Description,
                Category = request.Category,
                FormulaExpression = request.FormulaExpression,
                VariablesJson = System.Text.Json.JsonSerializer.Serialize(request.Variables),
                SortOrder = request.SortOrder,
                ValidUntil = request.ValidUntil,
                Metadata = { request.Metadata },
                UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
            };
            
            var updated = await _formulaService.UpdateFormulaAsync(request.FormulaId, formula, context.CancellationToken);
            
            if (updated == null)
            {
                return new UpdateRatingFormulaResponse
                {
                    Success = false,
                    Errors = { "Formula not found" }
                };
            }
            
            return new UpdateRatingFormulaResponse
            {
                Formula = updated,
                Success = true
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating rating formula");
            return new UpdateRatingFormulaResponse
            {
                Success = false,
                Errors = { ex.Message }
            };
        }
    }

    public override async Task<DeleteRatingFormulaResponse> DeleteRatingFormula(
        DeleteRatingFormulaRequest request,
        ServerCallContext context)
    {
        try
        {
            var result = await _formulaService.DeleteFormulaAsync(request.FormulaId, context.CancellationToken);
            
            return new DeleteRatingFormulaResponse
            {
                Success = result,
                Deleted = result
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error deleting rating formula");
            return new DeleteRatingFormulaResponse
            {
                Success = false
            };
        }
    }

    public override async Task<ListRatingFormulasResponse> ListRatingFormulas(
        ListRatingFormulasRequest request,
        ServerCallContext context)
    {
        try
        {
            var formulas = await _formulaService.GetFormulasAsync(
                request.InsuranceType,
                request.Category != FormulaCategory.Unspecified ? request.Category : null,
                request.Status != FormulaStatus.Unspecified ? request.Status : null,
                context.CancellationToken);
            
            var formulaList = formulas.ToList();
            
            // Apply search filter
            if (!string.IsNullOrEmpty(request.SearchQuery))
            {
                formulaList = formulaList
                    .Where(f => f.FormulaName.Contains(request.SearchQuery, StringComparison.OrdinalIgnoreCase) ||
                               f.FormulaCode.Contains(request.SearchQuery, StringComparison.OrdinalIgnoreCase))
                    .ToList();
            }
            
            // Apply pagination
            var pageSize = request.PageSize > 0 ? request.PageSize : 20;
            var skip = 0;
            if (!string.IsNullOrEmpty(request.PageToken))
            {
                int.TryParse(request.PageToken, out skip);
            }
            
            var paginated = formulaList.Skip(skip).Take(pageSize).ToList();
            var nextPageToken = skip + paginated.Count < formulaList.Count 
                ? (skip + paginated.Count).ToString() 
                : string.Empty;
            
            return new ListRatingFormulasResponse
            {
                Formulas = { paginated },
                TotalCount = formulaList.Count,
                NextPageToken = nextPageToken
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error listing rating formulas");
            return new ListRatingFormulasResponse
            {
                TotalCount = 0
            };
        }
    }

    public override async Task<ActivateRatingFormulaResponse> ActivateRatingFormula(
        ActivateRatingFormulaRequest request,
        ServerCallContext context)
    {
        try
        {
            var formula = await _formulaService.ActivateFormulaAsync(request.FormulaId, context.CancellationToken);
            
            if (formula == null)
            {
                return new ActivateRatingFormulaResponse
                {
                    Success = false
                };
            }
            
            return new ActivateRatingFormulaResponse
            {
                Success = true,
                Formula = formula
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error activating rating formula");
            return new ActivateRatingFormulaResponse
            {
                Success = false
            };
        }
    }

    // Calculation history
    public override async Task<GetCalculationResponse> GetCalculation(
        GetCalculationRequest request,
        ServerCallContext context)
    {
        try
        {
            ActuarialCalculation? calculation = null;
            
            if (!string.IsNullOrEmpty(request.CalculationId))
            {
                calculation = await _actuarialService.GetCalculationAsync(request.CalculationId, context.CancellationToken);
            }
            else if (!string.IsNullOrEmpty(request.CalculationReference))
            {
                calculation = await _actuarialService.GetCalculationByReferenceAsync(request.CalculationReference, context.CancellationToken);
            }
            
            return new GetCalculationResponse
            {
                Calculation = calculation,
                Found = calculation != null
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting calculation");
            return new GetCalculationResponse
            {
                Found = false
            };
        }
    }

    public override async Task<ListCalculationsResponse> ListCalculations(
        ListCalculationsRequest request,
        ServerCallContext context)
    {
        try
        {
            var calculations = await _actuarialService.GetCalculationsAsync(
                request.CalculationType != ActuarialCalculationType.Unspecified ? request.CalculationType : null,
                !string.IsNullOrEmpty(request.EntityType) ? request.EntityType : null,
                !string.IsNullOrEmpty(request.EntityId) ? request.EntityId : null,
                request.DateFrom?.ToDateTime(),
                request.DateTo?.ToDateTime(),
                context.CancellationToken);
            
            var calculationList = calculations.ToList();
            
            // Apply pagination
            var pageSize = request.PageSize > 0 ? request.PageSize : 20;
            var skip = 0;
            if (!string.IsNullOrEmpty(request.PageToken))
            {
                int.TryParse(request.PageToken, out skip);
            }
            
            var paginated = calculationList.Skip(skip).Take(pageSize).ToList();
            var nextPageToken = skip + paginated.Count < calculationList.Count 
                ? (skip + paginated.Count).ToString() 
                : string.Empty;
            
            return new ListCalculationsResponse
            {
                Calculations = { paginated },
                TotalCount = calculationList.Count,
                NextPageToken = nextPageToken
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error listing calculations");
            return new ListCalculationsResponse
            {
                TotalCount = 0
            };
        }
    }

    public override Task<RecalculateResponse> Recalculate(
        RecalculateRequest request,
        ServerCallContext context)
    {
        // Simplified implementation - in production, this would retrieve the original calculation,
        // merge updated parameters, and recalculate
        return Task.FromResult(new RecalculateResponse
        {
            NewCalculationId = Guid.NewGuid().ToString(),
            Success = true
        });
    }

    // Actuarial reporting
    public override Task<GenerateActuarialReportResponse> GenerateActuarialReport(
        GenerateActuarialReportRequest request,
        ServerCallContext context)
    {
        // Simplified implementation - in production, this would generate actual reports
        return Task.FromResult(new GenerateActuarialReportResponse
        {
            ReportId = Guid.NewGuid().ToString(),
            ReportUrl = $"/reports/actuarial/{Guid.NewGuid()}.pdf",
            Success = true,
            GeneratedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });
    }
}
