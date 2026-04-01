namespace InsuranceEngine.SharedKernel.Messaging;

public class InsuranceKafkaOptions
{
    public string BootstrapServers { get; set; } = string.Empty;
    public string GroupId { get; set; } = "insurance-engine-group";
}
