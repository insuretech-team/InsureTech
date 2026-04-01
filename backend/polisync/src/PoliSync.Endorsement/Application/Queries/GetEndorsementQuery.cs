using PoliSync.SharedKernel.CQRS;
using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;

namespace PoliSync.Endorsement.Application.Queries;

public sealed record GetEndorsementQuery(string EndorsementId) : IQuery<EndorsementEntity>;
