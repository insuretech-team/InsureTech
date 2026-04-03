using InsuranceEngine.Grpc.Clients;
using InsuranceEngine.Infrastructure.Notifications;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Infrastructure.Payments;

public enum PaymentMethod
{
    Unspecified = 0,
    Bkash = 1,
    Nagad = 2,
    Rocket = 3,
    Card = 4,
    BankTransfer = 5,
    SslCommerz = 6,
    Manual = 7
}

public class PaymentRequest
{
    public string UserId { get; set; } = string.Empty;
    public string PolicyId { get; set; } = string.Empty;
    public decimal Amount { get; set; }
    public string Currency { get; set; } = "BDT";
    public string PaymentMethod { get; set; } = string.Empty;
    public string CallbackUrl { get; set; } = string.Empty;
    public Dictionary<string, string>? Metadata { get; set; }
    public string? OrderId { get; set; }
    public string? InvoiceId { get; set; }
    public CustomerDetails? Customer { get; set; }
}

public class CustomerDetails
{
    public string Name { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string Phone { get; set; } = string.Empty;
    public string Address { get; set; } = string.Empty;
    public string City { get; set; } = string.Empty;
    public string Postcode { get; set; } = string.Empty;
    public string Country { get; set; } = "Bangladesh";
}

public class PaymentResult
{
    public string PaymentId { get; set; } = string.Empty;
    public string TransactionId { get; set; } = string.Empty;
    public string PaymentUrl { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public DateTime? ExpiresAt { get; set; }
    public string Provider { get; set; } = string.Empty;
    public string? GatewayPageUrl { get; set; }
    public string? TranId { get; set; }
    public string? SessionKey { get; set; }
}

public class PaymentVerificationResult
{
    public string PaymentId { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public bool Verified { get; set; }
    public decimal Amount { get; set; }
    public string Provider { get; set; } = string.Empty;
}

public interface IPaymentService
{
    Task<PaymentResult> InitiatePaymentAsync(PaymentRequest request, CancellationToken ct = default);
    Task<PaymentVerificationResult> VerifyPaymentAsync(string paymentId, string transactionId, string paymentMethod, CancellationToken ct = default);
    Task<PaymentVerificationResult> GetPaymentAsync(string paymentId, CancellationToken ct = default);
    Task<List<PaymentVerificationResult>> ListPaymentsAsync(string? userId = null, string? policyId = null, string? status = null, CancellationToken ct = default);
    Task HandleGatewayWebhookAsync(string provider, byte[] rawPayload, CancellationToken ct = default);
    Task<PaymentVerificationResult> SubmitManualPaymentProofAsync(string paymentId, string fileId, string submittedBy, string? notes = null, CancellationToken ct = default);
    Task<PaymentVerificationResult> ReviewManualPaymentAsync(string paymentId, bool approved, string reviewedBy, string? reviewNotes = null, string? rejectionReason = null, CancellationToken ct = default);
    Task<string> GenerateReceiptAsync(string paymentId, CancellationToken ct = default);
    Task NotifyPaymentStatusAsync(string userId, string paymentId, string status, decimal amount, CancellationToken ct = default);
}

public class PaymentService : IPaymentService
{
    private readonly InsuranceServiceClient _client;
    private readonly INotificationService _notificationService;
    private readonly ILogger<PaymentService> _logger;

    public PaymentService(
        InsuranceServiceClient client,
        INotificationService notificationService,
        ILogger<PaymentService> logger)
    {
        _client = client;
        _notificationService = notificationService;
        _logger = logger;
    }

    public async Task<PaymentResult> InitiatePaymentAsync(PaymentRequest request, CancellationToken ct = default)
    {
        try
        {
            var protoRequest = new Insuretech.Payment.Services.V1.InitiatePaymentRequest
            {
                UserId = request.UserId,
                PolicyId = request.PolicyId,
                Amount = new Insuretech.Common.V1.Money { Amount = (long)(request.Amount * 100), Currency = request.Currency },
                Currency = request.Currency,
                PaymentMethod = request.PaymentMethod.ToUpper(),
                CallbackUrl = request.CallbackUrl,
                IdempotencyKey = Guid.NewGuid().ToString()
            };

            if (!string.IsNullOrEmpty(request.OrderId))
                protoRequest.OrderId = request.OrderId;
            if (!string.IsNullOrEmpty(request.InvoiceId))
                protoRequest.InvoiceId = request.InvoiceId;

            if (request.Metadata != null)
            {
                foreach (var kvp in request.Metadata)
                    protoRequest.Metadata[kvp.Key] = kvp.Value;
            }

            if (request.Customer != null)
            {
                protoRequest.CustomerName = request.Customer.Name;
                protoRequest.CustomerEmail = request.Customer.Email;
                protoRequest.CustomerPhone = request.Customer.Phone;
                protoRequest.CustomerAddressLine1 = request.Customer.Address;
                protoRequest.CustomerCity = request.Customer.City;
                protoRequest.CustomerPostcode = request.Customer.Postcode;
                protoRequest.CustomerCountry = request.Customer.Country;
            }

            var response = await _client.Payments.InitiatePaymentAsync(protoRequest, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to initiate payment: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Payment initiation failed: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Payment initiated: {PaymentId} for {Amount} {Currency} via {Method}",
                response.PaymentId, request.Amount, request.Currency, request.PaymentMethod);

            return new PaymentResult
            {
                PaymentId = response.PaymentId,
                TransactionId = response.TransactionId,
                PaymentUrl = response.PaymentUrl,
                Status = response.Status,
                ExpiresAt = response.ExpiresAt?.ToDateTime(),
                Provider = response.Provider,
                GatewayPageUrl = response.GatewayPageUrl,
                TranId = response.TranId,
                SessionKey = response.SessionKey
            };
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error initiating payment for policy {PolicyId}", request.PolicyId);
            throw;
        }
    }

    public async Task<PaymentVerificationResult> VerifyPaymentAsync(
        string paymentId,
        string transactionId,
        string paymentMethod,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Payment.Services.V1.VerifyPaymentRequest
            {
                PaymentId = paymentId,
                TransactionId = transactionId,
                PaymentMethod = paymentMethod.ToUpper(),
                IdempotencyKey = Guid.NewGuid().ToString()
            };

            var response = await _client.Payments.VerifyPaymentAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to verify payment: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Payment verification failed: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Payment verified: {PaymentId}, Verified={Verified}, Status={Status}",
                paymentId, response.Verified, response.Status);

            return new PaymentVerificationResult
            {
                PaymentId = response.PaymentId,
                Status = response.Status,
                Verified = response.Verified,
                Amount = (decimal)(response.Payment?.Amount?.Amount ?? 0),
                Provider = response.Payment?.Provider ?? string.Empty
            };
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error verifying payment {PaymentId}", paymentId);
            throw;
        }
    }

    public async Task<PaymentVerificationResult> GetPaymentAsync(string paymentId, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Payment.Services.V1.GetPaymentRequest
            {
                PaymentId = paymentId
            };

            var response = await _client.Payments.GetPaymentAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to get payment: {response.Error.Message}");
            }

            var payment = response.Payment;
            return new PaymentVerificationResult
            {
                PaymentId = payment?.PaymentId ?? paymentId,
                Status = payment?.Status.ToString() ?? "UNKNOWN",
                Verified = payment?.Status == Insuretech.Payment.Entity.V1.PaymentStatus.Success,
                Amount = (decimal)(payment?.Amount?.Amount ?? 0),
                Provider = payment?.Provider ?? string.Empty
            };
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error getting payment {PaymentId}", paymentId);
            throw;
        }
    }

    public async Task<List<PaymentVerificationResult>> ListPaymentsAsync(
        string? userId = null,
        string? policyId = null,
        string? status = null,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Payment.Services.V1.ListPaymentsRequest();

            if (!string.IsNullOrEmpty(userId))
                request.UserId = userId;
            if (!string.IsNullOrEmpty(policyId))
                request.PolicyId = policyId;
            if (!string.IsNullOrEmpty(status))
                request.Status = status;

            var response = await _client.Payments.ListPaymentsAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to list payments: {response.Error.Message}");
            }

            return response.Payments.Select(p => new PaymentVerificationResult
            {
                PaymentId = p.PaymentId,
                Status = p.Status.ToString(),
                Verified = p.Status == Insuretech.Payment.Entity.V1.PaymentStatus.Success,
                Amount = (decimal)(p.Amount?.Amount ?? 0),
                Provider = p.Provider ?? string.Empty
            }).ToList();
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error listing payments");
            throw;
        }
    }

    public async Task HandleGatewayWebhookAsync(string provider, byte[] rawPayload, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Payment.Services.V1.HandleGatewayWebhookRequest
            {
                Provider = provider.ToLower(),
                RawPayload = Google.Protobuf.ByteString.CopyFrom(rawPayload),
                ReceivedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
            };

            var response = await _client.Payments.HandleGatewayWebhookAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogWarning("Gateway webhook processed with error: {Error}", response.Error.Message);
            }

            _logger.LogInformation(
                "Gateway webhook handled: Provider={Provider}, PaymentId={PaymentId}, Accepted={Accepted}",
                provider, response.PaymentId, response.Accepted);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error handling gateway webhook for {Provider}", provider);
            throw;
        }
    }

    public async Task<PaymentVerificationResult> SubmitManualPaymentProofAsync(
        string paymentId,
        string fileId,
        string submittedBy,
        string? notes = null,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Payment.Services.V1.SubmitManualPaymentProofRequest
            {
                PaymentId = paymentId,
                ManualProofFileId = fileId,
                SubmittedBy = submittedBy,
                Notes = notes ?? string.Empty
            };

            var response = await _client.Payments.SubmitManualPaymentProofAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to submit proof: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Manual payment proof submitted: PaymentId={PaymentId}, Status={Status}",
                paymentId, response.Status);

            return new PaymentVerificationResult
            {
                PaymentId = paymentId,
                Status = response.Status,
                Verified = false,
                Amount = 0,
                Provider = "MANUAL"
            };
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error submitting manual proof for {PaymentId}", paymentId);
            throw;
        }
    }

    public async Task<PaymentVerificationResult> ReviewManualPaymentAsync(
        string paymentId,
        bool approved,
        string reviewedBy,
        string? reviewNotes = null,
        string? rejectionReason = null,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Payment.Services.V1.ReviewManualPaymentRequest
            {
                PaymentId = paymentId,
                Approved = approved,
                ReviewedBy = reviewedBy,
                ReviewNotes = reviewNotes ?? string.Empty,
                RejectionReason = rejectionReason ?? string.Empty
            };

            var response = await _client.Payments.ReviewManualPaymentAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to review payment: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Manual payment reviewed: PaymentId={PaymentId}, Approved={Approved}, By={ReviewedBy}",
                paymentId, approved, reviewedBy);

            return new PaymentVerificationResult
            {
                PaymentId = paymentId,
                Status = response.Status,
                Verified = approved,
                Amount = 0,
                Provider = "MANUAL"
            };
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error reviewing manual payment {PaymentId}", paymentId);
            throw;
        }
    }

    public async Task<string> GenerateReceiptAsync(string paymentId, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Payment.Services.V1.GenerateReceiptRequest
            {
                PaymentId = paymentId
            };

            var response = await _client.Payments.GenerateReceiptAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to generate receipt: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Payment receipt generated: PaymentId={PaymentId}, ReceiptNumber={ReceiptNumber}",
                paymentId, response.ReceiptNumber);

            return response.ReceiptNumber;
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error generating receipt for {PaymentId}", paymentId);
            throw;
        }
    }

    public async Task NotifyPaymentStatusAsync(
        string userId,
        string paymentId,
        string status,
        decimal amount,
        CancellationToken ct = default)
    {
        var statusMessage = status.ToUpper() switch
        {
            "INITIATED" => $"Your payment (ID: {paymentId}) of BDT {amount:N2} has been initiated. Please complete the payment.",
            "VERIFIED" or "SUCCESS" or "COMPLETED" => $"Your payment (ID: {paymentId}) of BDT {amount:N2} has been verified successfully. Thank you!",
            "FAILED" => $"Your payment (ID: {paymentId}) has failed. Please try again or contact support.",
            "REFUNDED" => $"Your refund of BDT {amount:N2} for payment (ID: {paymentId}) has been processed.",
            _ => $"Payment update for {paymentId}: Status changed to {status}."
        };

        var subject = status.ToUpper() switch
        {
            "VERIFIED" or "SUCCESS" or "COMPLETED" => $"Payment Confirmed - {paymentId}",
            "FAILED" => $"Payment Failed - {paymentId}",
            "REFUNDED" => $"Refund Processed - {paymentId}",
            _ => $"Payment Update - {paymentId}"
        };

        await _notificationService.SendEmailAsync(
            userId,
            subject,
            statusMessage,
            new Dictionary<string, string>
            {
                ["payment_id"] = paymentId,
                ["status"] = status,
                ["amount"] = amount.ToString("N2")
            },
            ct);

        _logger.LogInformation("Payment status notification sent to user {UserId}: {PaymentId} -> {Status}",
            userId, paymentId, status);
    }
}

public class MockPaymentService : IPaymentService
{
    private readonly INotificationService _notificationService;
    private readonly ILogger<MockPaymentService> _logger;

    public MockPaymentService(
        INotificationService notificationService,
        ILogger<MockPaymentService> logger)
    {
        _notificationService = notificationService;
        _logger = logger;
    }

    public Task<PaymentResult> InitiatePaymentAsync(PaymentRequest request, CancellationToken ct = default)
    {
        var paymentId = Guid.NewGuid().ToString();
        var transactionId = $"TXN-{DateTime.UtcNow:yyyyMMddHHmmss}-{paymentId[..8]}";

        _logger.LogInformation(
            "[MOCK] Payment initiated: {PaymentId}, Amount={Amount} {Currency}, Method={Method}",
            paymentId, request.Amount, request.Currency, request.PaymentMethod);

        return Task.FromResult(new PaymentResult
        {
            PaymentId = paymentId,
            TransactionId = transactionId,
            PaymentUrl = $"https://mock-gateway.labaid.com/pay/{paymentId}",
            Status = "INITIATED",
            ExpiresAt = DateTime.UtcNow.AddMinutes(30),
            Provider = request.PaymentMethod.ToUpper(),
            GatewayPageUrl = $"https://mock-gateway.labaid.com/checkout/{transactionId}"
        });
    }

    public Task<PaymentVerificationResult> VerifyPaymentAsync(
        string paymentId,
        string transactionId,
        string paymentMethod,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Payment verified: {PaymentId}, TxnId={TransactionId}",
            paymentId, transactionId);

        return Task.FromResult(new PaymentVerificationResult
        {
            PaymentId = paymentId,
            Status = "VERIFIED",
            Verified = true,
            Amount = 10000,
            Provider = paymentMethod.ToUpper()
        });
    }

    public Task<PaymentVerificationResult> GetPaymentAsync(string paymentId, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Get payment: {PaymentId}", paymentId);

        return Task.FromResult(new PaymentVerificationResult
        {
            PaymentId = paymentId,
            Status = "VERIFIED",
            Verified = true,
            Amount = 10000,
            Provider = "BKASH"
        });
    }

    public Task<List<PaymentVerificationResult>> ListPaymentsAsync(
        string? userId = null,
        string? policyId = null,
        string? status = null,
        CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] List payments: User={UserId}, Policy={PolicyId}", userId, policyId);

        return Task.FromResult(new List<PaymentVerificationResult>
        {
            new()
            {
                PaymentId = Guid.NewGuid().ToString(),
                Status = "VERIFIED",
                Verified = true,
                Amount = 10000,
                Provider = "BKASH"
            }
        });
    }

    public Task HandleGatewayWebhookAsync(string provider, byte[] rawPayload, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Gateway webhook handled: Provider={Provider}", provider);
        return Task.CompletedTask;
    }

    public Task<PaymentVerificationResult> SubmitManualPaymentProofAsync(
        string paymentId,
        string fileId,
        string submittedBy,
        string? notes = null,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Manual proof submitted: PaymentId={PaymentId}, FileId={FileId}",
            paymentId, fileId);

        return Task.FromResult(new PaymentVerificationResult
        {
            PaymentId = paymentId,
            Status = "PROCESSING",
            Verified = false,
            Amount = 0,
            Provider = "MANUAL"
        });
    }

    public Task<PaymentVerificationResult> ReviewManualPaymentAsync(
        string paymentId,
        bool approved,
        string reviewedBy,
        string? reviewNotes = null,
        string? rejectionReason = null,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Manual payment reviewed: PaymentId={PaymentId}, Approved={Approved}, By={ReviewedBy}",
            paymentId, approved, reviewedBy);

        return Task.FromResult(new PaymentVerificationResult
        {
            PaymentId = paymentId,
            Status = approved ? "VERIFIED" : "FAILED",
            Verified = approved,
            Amount = 10000,
            Provider = "MANUAL"
        });
    }

    public Task<string> GenerateReceiptAsync(string paymentId, CancellationToken ct = default)
    {
        var receiptNumber = $"RCP-{DateTime.UtcNow:yyyyMMdd}-{paymentId[..8]}";
        _logger.LogInformation("[MOCK] Receipt generated: {ReceiptNumber}", receiptNumber);
        return Task.FromResult(receiptNumber);
    }

    public Task NotifyPaymentStatusAsync(
        string userId,
        string paymentId,
        string status,
        decimal amount,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Payment notification: User={UserId}, Payment={PaymentId}, Status={Status}",
            userId, paymentId, status);
        return Task.CompletedTask;
    }
}
