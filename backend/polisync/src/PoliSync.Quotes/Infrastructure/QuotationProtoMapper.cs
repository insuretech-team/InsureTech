using System.Text.Json;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using ProtoQuotationStatus = Insuretech.Policy.Entity.V1.QuotationStatus;
using DomainQuotation = PoliSync.Quotes.Domain.Quotation;
using DomainQuotationStatus = PoliSync.Quotes.Domain.QuotationStatus;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;

namespace PoliSync.Quotes.Infrastructure;

internal static class QuotationProtoMapper
{
    public static QuotationEntity ToProto(DomainQuotation quotation)
    {
        var metadata = new PersistedQuotationMetadata(
            quotation.Status.ToString(),
            quotation.ProductId,
            quotation.PlanId,
            quotation.CustomerId,
            quotation.BasePremium,
            quotation.RiderPremium,
            quotation.LoadingAmount,
            quotation.DiscountAmount,
            quotation.VatTax,
            quotation.ServiceFee,
            quotation.TotalPayable);

        return new QuotationEntity
        {
            QuotationId = quotation.Id == Guid.Empty ? Guid.NewGuid().ToString() : quotation.Id.ToString(),
            BusinessId = quotation.TenantId == Guid.Empty ? string.Empty : quotation.TenantId.ToString(),
            DepartmentId = quotation.ProductId == Guid.Empty ? string.Empty : quotation.ProductId.ToString(),
            PlanId = quotation.PlanId == Guid.Empty ? string.Empty : quotation.PlanId.ToString(),
            CreatedByUserId = quotation.CustomerId == Guid.Empty ? string.Empty : quotation.CustomerId.ToString(),
            QuotationNumber = string.IsNullOrWhiteSpace(quotation.QuotationNumber)
                ? $"QUO-{DateTime.UtcNow:yyyyMMdd}-{quotation.Id.ToString("N")[..8].ToUpper()}"
                : quotation.QuotationNumber,
            // B2C quotations default to LIFE insurance type (personal accident / health);
            // The DB stores InsuranceType enum string e.g. "INSURANCE_TYPE_LIFE".
            InsuranceCategory = Insuretech.Common.V1.InsuranceType.Life,
            Status = ToProtoStatus(quotation.Status),
            ValidUntil = Timestamp.FromDateTime(quotation.ExpiryDate.ToUniversalTime()),
            CreatedAt = Timestamp.FromDateTime(quotation.CreatedAt.ToUniversalTime()),
            UpdatedAt = Timestamp.FromDateTime((quotation.UpdatedAt ?? quotation.CreatedAt).ToUniversalTime()),
            RejectionReason = quotation.RejectionReason ?? string.Empty,
            EstimatedPremium = NewMoney(quotation.BasePremium),
            QuotedAmount = NewMoney(quotation.TotalPayable),
            InsurerName = JsonSerializer.Serialize(metadata)
        };
    }

    public static DomainQuotation ToDomain(QuotationEntity quotation)
    {
        var metadata = ParseMetadata(quotation);

        return DomainQuotation.Rehydrate(
            id: ParseGuid(quotation.QuotationId),
            tenantId: ParseGuid(quotation.BusinessId),
            quotationNumber: quotation.QuotationNumber,
            productId: metadata.ProductId != Guid.Empty ? metadata.ProductId : ParseGuid(quotation.DepartmentId),
            planId: metadata.PlanId != Guid.Empty ? metadata.PlanId : ParseGuid(quotation.PlanId),
            customerId: metadata.CustomerId != Guid.Empty ? metadata.CustomerId : ParseGuid(quotation.CreatedByUserId),
            status: ToDomainStatus(quotation.Status, metadata.DomainStatus),
            expiryDate: quotation.ValidUntil?.ToDateTime() ?? DateTime.UtcNow,
            basePremium: metadata.BasePremium != 0 ? metadata.BasePremium : quotation.EstimatedPremium?.Amount ?? 0,
            riderPremium: metadata.RiderPremium,
            loadingAmount: metadata.LoadingAmount,
            discountAmount: metadata.DiscountAmount,
            vatTax: metadata.VatTax,
            serviceFee: metadata.ServiceFee,
            totalPayable: metadata.TotalPayable != 0 ? metadata.TotalPayable : quotation.QuotedAmount?.Amount ?? 0,
            rejectionReason: string.IsNullOrWhiteSpace(quotation.RejectionReason) ? null : quotation.RejectionReason,
            createdAt: quotation.CreatedAt?.ToDateTime() ?? DateTime.UtcNow,
            updatedAt: quotation.UpdatedAt?.ToDateTime());
    }

    private static PersistedQuotationMetadata ParseMetadata(QuotationEntity quotation)
    {
        if (string.IsNullOrWhiteSpace(quotation.InsurerName))
        {
            return PersistedQuotationMetadata.Empty;
        }

        try
        {
            return JsonSerializer.Deserialize<PersistedQuotationMetadata>(quotation.InsurerName) ?? PersistedQuotationMetadata.Empty;
        }
        catch (JsonException)
        {
            return PersistedQuotationMetadata.Empty;
        }
    }

    private static Guid ParseGuid(string? value)
        => Guid.TryParse(value, out var guid) ? guid : Guid.Empty;

    private static Money NewMoney(long amount)
        => new() { Amount = amount, Currency = "BDT" };

    private static ProtoQuotationStatus ToProtoStatus(DomainQuotationStatus status)
        => status switch
        {
            DomainQuotationStatus.Draft => ProtoQuotationStatus.Draft,
            DomainQuotationStatus.Submitted => ProtoQuotationStatus.Submitted,
            DomainQuotationStatus.Received => ProtoQuotationStatus.Received,
            DomainQuotationStatus.Approved => ProtoQuotationStatus.Approved,
            DomainQuotationStatus.Rejected => ProtoQuotationStatus.Rejected,
            DomainQuotationStatus.Expired => ProtoQuotationStatus.Rejected,
            _ => ProtoQuotationStatus.Unspecified
        };

    private static DomainQuotationStatus ToDomainStatus(ProtoQuotationStatus status, string? persistedStatus)
    {
        if (System.Enum.TryParse<DomainQuotationStatus>(persistedStatus, ignoreCase: true, out var storedStatus))
        {
            return storedStatus;
        }

        return status switch
        {
            ProtoQuotationStatus.Submitted => DomainQuotationStatus.Submitted,
            ProtoQuotationStatus.Received => DomainQuotationStatus.Received,
            ProtoQuotationStatus.Approved => DomainQuotationStatus.Approved,
            ProtoQuotationStatus.Rejected => DomainQuotationStatus.Rejected,
            _ => DomainQuotationStatus.Draft
        };
    }

    private sealed record PersistedQuotationMetadata(
        string DomainStatus,
        Guid ProductId,
        Guid PlanId,
        Guid CustomerId,
        long BasePremium,
        long RiderPremium,
        long LoadingAmount,
        long DiscountAmount,
        long VatTax,
        long ServiceFee,
        long TotalPayable)
    {
        public static PersistedQuotationMetadata Empty => new(
            string.Empty,
            Guid.Empty,
            Guid.Empty,
            Guid.Empty,
            0,
            0,
            0,
            0,
            0,
            0,
            0);
    }
}
