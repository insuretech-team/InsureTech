namespace InsuranceEngine.Infrastructure.Messaging;

public static class KafkaEventTopics
{
    public const string PolicyPrefix = "insuretech.policy.v1.";
    
    public static class Policy
    {
        public const string Created = PolicyPrefix + "created";
        public const string Issued = PolicyPrefix + "issued";
        public const string Renewed = PolicyPrefix + "renewed";
        public const string Cancelled = PolicyPrefix + "cancelled";
        public const string Lapsed = PolicyPrefix + "lapsed";
        public const string Endorsed = PolicyPrefix + "endorsed";
    }

    public static class Claims
    {
        private const string ClaimsPrefix = "insuretech.claims.v1.";
        
        public const string Submitted = ClaimsPrefix + "submitted";
        public const string Approved = ClaimsPrefix + "approved";
        public const string Rejected = ClaimsPrefix + "rejected";
        public const string Settled = ClaimsPrefix + "settled";
        public const string DocumentsRequested = ClaimsPrefix + "documents_requested";
        public const string Disputed = ClaimsPrefix + "disputed";
    }

    public static class Payment
    {
        private const string PaymentPrefix = "insuretech.payment.v1.";
        
        public const string Initiated = PaymentPrefix + "initiated";
        public const string Confirmed = PaymentPrefix + "confirmed";
        public const string Failed = PaymentPrefix + "failed";
        public const string RefundInitiated = PaymentPrefix + "refund_initiated";
        public const string RefundCompleted = PaymentPrefix + "refund_completed";
    }

    public static class Notification
    {
        private const string NotificationPrefix = "insuretech.notification.v1.";
        
        public const string OtpSent = NotificationPrefix + "otp_sent";
        public const string SmsSent = NotificationPrefix + "sms_sent";
        public const string EmailSent = NotificationPrefix + "email_sent";
        public const string PushSent = NotificationPrefix + "push_sent";
    }

    public static class Fraud
    {
        private const string FraudPrefix = "insuretech.fraud.v1.";
        
        public const string Alert = FraudPrefix + "alert";
        public const string Confirmed = FraudPrefix + "confirmed";
    }

    public static class Commission
    {
        private const string CommissionPrefix = "insuretech.commission.v1.";
        
        public const string Calculated = CommissionPrefix + "calculated";
        public const string PayoutCreated = CommissionPrefix + "payout_created";
        public const string PayoutProcessed = CommissionPrefix + "payout_processed";
    }

    public static class Underwriting
    {
        private const string UnderwritingPrefix = "insuretech.underwriting.v1.";
        
        public const string QuoteRequested = UnderwritingPrefix + "quote_requested";
        public const string DecisionMade = UnderwritingPrefix + "decision_made";
        public const string ProposalSubmitted = UnderwritingPrefix + "proposal_submitted";
    }
}
