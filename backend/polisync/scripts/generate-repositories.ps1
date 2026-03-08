# Generate CRUD repositories for all insurance schema entities
# This script creates repository interfaces and implementations for proto-generated entities

$entities = @(
    @{ Name = "Quote"; IdField = "QuoteId"; NumberField = "QuoteNumber"; Namespace = "Insuretech.Products.Entity.V1"; Status = "QuoteStatus" },
    @{ Name = "Quotation"; IdField = "QuotationId"; NumberField = "QuotationNumber"; Namespace = "Insuretech.Products.Entity.V1"; Status = "QuotationStatus" },
    @{ Name = "Endorsement"; IdField = "EndorsementId"; NumberField = "EndorsementNumber"; Namespace = "Insuretech.Policy.Entity.V1"; Status = "EndorsementStatus" },
    @{ Name = "Refund"; IdField = "RefundId"; NumberField = "RefundNumber"; Namespace = "Insuretech.Policy.Entity.V1"; Status = "RefundStatus" },
    @{ Name = "RenewalSchedule"; IdField = "ScheduleId"; NumberField = $null; Namespace = "Insuretech.Policy.Entity.V1"; Status = "RenewalStatus" },
    @{ Name = "RenewalReminder"; IdField = "ReminderId"; NumberField = $null; Namespace = "Insuretech.Policy.Entity.V1"; Status = "ReminderStatus" },
    @{ Name = "GracePeriod"; IdField = "GracePeriodId"; NumberField = $null; Namespace = "Insuretech.Policy.Entity.V1"; Status = "GracePeriodStatus" },
    @{ Name = "UnderwritingDecision"; IdField = "DecisionId"; NumberField = $null; Namespace = "Insuretech.Policy.Entity.V1"; Status = "DecisionStatus" },
    @{ Name = "HealthDeclaration"; IdField = "DeclarationId"; NumberField = $null; Namespace = "Insuretech.Policy.Entity.V1"; Status = $null },
    @{ Name = "Beneficiary"; IdField = "BeneficiaryId"; NumberField = $null; Namespace = "Insuretech.Policy.Entity.V1"; Status = "BeneficiaryStatus" },
    @{ Name = "FraudAlert"; IdField = "AlertId"; NumberField = $null; Namespace = "Insuretech.Claims.Entity.V1"; Status = "AlertStatus" },
    @{ Name = "FraudCase"; IdField = "CaseId"; NumberField = "CaseNumber"; Namespace = "Insuretech.Claims.Entity.V1"; Status = "CaseStatus" },
    @{ Name = "FraudRule"; IdField = "RuleId"; NumberField = $null; Namespace = "Insuretech.Claims.Entity.V1"; Status = "RuleStatus" },
    @{ Name = "Insurer"; IdField = "InsurerId"; NumberField = $null; Namespace = "Insuretech.Products.Entity.V1"; Status = "InsurerStatus" },
    @{ Name = "InsurerConfig"; IdField = "ConfigId"; NumberField = $null; Namespace = "Insuretech.Products.Entity.V1"; Status = $null },
    @{ Name = "InsurerProduct"; IdField = "InsurerProductId"; NumberField = $null; Namespace = "Insuretech.Products.Entity.V1"; Status = $null },
    @{ Name = "PricingConfig"; IdField = "ConfigId"; NumberField = $null; Namespace = "Insuretech.Products.Entity.V1"; Status = $null },
    @{ Name = "ProductPlan"; IdField = "PlanId"; NumberField = $null; Namespace = "Insuretech.Products.Entity.V1"; Status = "PlanStatus" },
    @{ Name = "Rider"; IdField = "RiderId"; NumberField = $null; Namespace = "Insuretech.Products.Entity.V1"; Status = "RiderStatus" }
)

$outputDir = "backend/polisync/src/PoliSync.Infrastructure/Repositories"

foreach ($entity in $entities) {
    $entityName = $entity.Name
    $idField = $entity.IdField
    $numberField = $entity.NumberField
    $namespace = $entity.Namespace
    $statusEnum = $entity.Status
    
    $fileName = "$outputDir/${entityName}Repository.cs"
    
    $hasNumberField = $numberField -ne $null
    $hasStatus = $statusEnum -ne $null
    
    $numberMethod = if ($hasNumberField) {
        @"
    Task<$entityName?> GetByNumberAsync(string $($numberField.ToLower()), CancellationToken cancellationToken = default);
"@
    } else { "" }
    
    $statusMethod = if ($hasStatus) {
        @"
    Task<List<$entityName>> GetByStatusAsync($statusEnum status, CancellationToken cancellationToken = default);
"@
    } else { "" }
    
    $numberImpl = if ($hasNumberField) {
        @"

    public async Task<$entityName?> GetByNumberAsync(string $($numberField.ToLower()), CancellationToken cancellationToken = default)
    {
        return await _context.${entityName}s
            .FirstOrDefaultAsync(e => e.$numberField == $($numberField.ToLower()), cancellationToken);
    }
"@
    } else { "" }
    
    $statusImpl = if ($hasStatus) {
        @"

    public async Task<List<$entityName>> GetByStatusAsync($statusEnum status, CancellationToken cancellationToken = default)
    {
        return await _context.${entityName}s
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
"@
    } else { "" }

    $content = @"
using Google.Protobuf.WellKnownTypes;
using $namespace;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface I${entityName}Repository
{
    Task<$entityName> CreateAsync($entityName entity, CancellationToken cancellationToken = default);
    Task<$entityName?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
$numberMethod$statusMethod    Task<List<$entityName>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<$entityName> UpdateAsync($entityName entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class ${entityName}Repository : I${entityName}Repository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<${entityName}Repository> _logger;

    public ${entityName}Repository(PoliSyncDbContext context, ILogger<${entityName}Repository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<$entityName> CreateAsync($entityName entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.$idField))
        {
            entity.$idField = Guid.NewGuid().ToString();
        }

        _context.${entityName}s.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created $entityName {Id}", entity.$idField);
        return entity;
    }

    public async Task<$entityName?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.${entityName}s
            .FirstOrDefaultAsync(e => e.$idField == id, cancellationToken);
    }
$numberImpl$statusImpl
    public async Task<List<$entityName>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.${entityName}s
            .OrderByDescending(e => e.$idField)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<$entityName> UpdateAsync($entityName entity, CancellationToken cancellationToken = default)
    {
        _context.${entityName}s.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated $entityName {Id}", entity.$idField);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.${entityName}s.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted $entityName {Id}", id);
        }
    }
}
"@

    Write-Host "Generating $fileName..."
    $content | Out-File -FilePath $fileName -Encoding UTF8
}

Write-Host "`nGenerated $($entities.Count) repository files successfully!" -ForegroundColor Green
