using PoliSync.SharedKernel.CQRS;
using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;

namespace PoliSync.Endorsement.Application.Queries;

public sealed record ListEndorsementsQuery(string PolicyId, int PageNumber, int PageSize)
    : IQuery<EndorsementListResult>;

public sealed record EndorsementListResult(List<EndorsementEntity> Endorsements, int TotalCount);
