using Grpc.Core;
using Insuretech.Policy.Services.V1;

namespace Insuretech.Policy.Services.V1;

/// <summary>
/// Partial class to add missing methods to the generated PolicyServiceClient.
/// Must be in the same assembly as the generated code (InsuranceEngine.Proto).
/// </summary>
public static partial class PolicyService
{
    public partial class PolicyServiceClient
    {
        public virtual AsyncUnaryCall<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CallOptions options)
        {
            return CallInvoker.AsyncUnaryCall(__Method_ApproveCancellation, null, options, request);
        }

        public virtual ApproveCancellationResponse ApproveCancellation(ApproveCancellationRequest request, CallOptions options)
        {
            return CallInvoker.BlockingUnaryCall(__Method_ApproveCancellation, null, options, request);
        }
    }
}
