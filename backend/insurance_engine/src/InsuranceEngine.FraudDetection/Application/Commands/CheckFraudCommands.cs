using Insuretech.Fraud.Services.V1;
using MediatR;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.FraudDetection.Application.Commands;

public sealed record CheckFraudCommand(
    string EntityType,
    string EntityId,
    Struct Data) : IRequest<CheckFraudResponse>;
