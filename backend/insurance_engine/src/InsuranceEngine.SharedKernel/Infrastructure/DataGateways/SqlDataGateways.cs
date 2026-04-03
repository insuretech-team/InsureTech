namespace InsuranceEngine.SharedKernel.Infrastructure.DataGateways;

using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using System;
using System.Collections.Generic;
using System.Linq.Expressions;
using System.Threading;
using System.Threading.Tasks;

// NOTE: SQL-based gateways for Policy and Claim are deprecated in favor of the 'GoDataGateway' (gRPC Proxy) pattern.
// If local persistence is required in the future, these implementations should be restored to match the updated I*DataGateway interfaces.
