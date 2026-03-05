package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/crypto"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

type AgentRepository struct {
	db            *gorm.DB
	encryptionKey string
}

func NewAgentRepository(db *gorm.DB, encryptionKey string) *AgentRepository {
	return &AgentRepository{db: db, encryptionKey: encryptionKey}
}

// decrypt decrypts an agent's PII fields in-place.
func (r *AgentRepository) decrypt(agent *partnerv1.Agent) {
	if agent == nil {
		return
	}
	if agent.PhoneNumber != "" {
		if plain, err := crypto.DecryptPII(agent.PhoneNumber, r.encryptionKey); err == nil {
			agent.PhoneNumber = plain
		} else {
			logger.Warnf("failed to decrypt agent phone: %v", err)
		}
	}
	if agent.Email != "" {
		if plain, err := crypto.DecryptPII(agent.Email, r.encryptionKey); err == nil {
			agent.Email = plain
		} else {
			logger.Warnf("failed to decrypt agent email: %v", err)
		}
	}
	if agent.NidNumber != "" {
		if plain, err := crypto.DecryptPII(agent.NidNumber, r.encryptionKey); err == nil {
			agent.NidNumber = plain
		} else {
			logger.Warnf("failed to decrypt agent nid: %v", err)
		}
	}
}

// encrypt returns a new PII-encrypted map of fields for saving.
func (r *AgentRepository) encrypt(agent *partnerv1.Agent) (map[string]any, error) {
	encPhone := agent.PhoneNumber
	encEmail := agent.Email
	encNid := agent.NidNumber
	var err error

	if encPhone != "" {
		encPhone, err = crypto.EncryptPII(encPhone, r.encryptionKey)
		if err != nil {
			return nil, err
		}
	}
	if encEmail != "" {
		encEmail, err = crypto.EncryptPII(encEmail, r.encryptionKey)
		if err != nil {
			return nil, err
		}
	}
	if encNid != "" {
		encNid, err = crypto.EncryptPII(encNid, r.encryptionKey)
		if err != nil {
			return nil, err
		}
	}

	return map[string]any{
		"phone_number": encPhone,
		"email":        encEmail,
		"nid_number":   encNid,
	}, nil
}

// Create inserts a new Agent record
func (r *AgentRepository) Create(ctx context.Context, agent *partnerv1.Agent) error {
	if agent.AgentId == "" {
		agent.AgentId = uuid.New().String()
	}
	now := time.Now()
	agent.CreatedAt = timestamppb.New(now)
	agent.UpdatedAt = timestamppb.New(now)
	agent.JoinedAt = timestamppb.New(now)

	if agent.Status == partnerv1.AgentStatus_AGENT_STATUS_UNSPECIFIED {
		agent.Status = partnerv1.AgentStatus_AGENT_STATUS_ACTIVE
	}

	encFields, err := r.encrypt(agent)
	if err != nil {
		return err
	}

	values := map[string]any{
		"agent_id":        agent.AgentId,
		"partner_id":      agent.PartnerId,
		"user_id":         agent.UserId,
		"full_name":       agent.FullName,
		"phone_number":    encFields["phone_number"],
		"email":           encFields["email"],
		"nid_number":      encFields["nid_number"],
		"status":          agent.Status.String(),
		"commission_rate": agent.CommissionRate,
		"joined_at":       now,
		"created_at":      now,
		"updated_at":      now,
	}

	return r.db.WithContext(ctx).Table("partner_schema.agents").Create(values).Error
}

// GetByID finds an agent by their ID
func (r *AgentRepository) GetByID(ctx context.Context, id string) (*partnerv1.Agent, error) {
	var agent partnerv1.Agent
	err := r.db.WithContext(ctx).Table("partner_schema.agents").
		Where("agent_id = ? AND deleted_at IS NULL", id).
		First(&agent).Error
	if err != nil {
		return nil, err
	}
	r.decrypt(&agent)
	return &agent, nil
}

// GetByUserID finds an agent by their related system user ID
func (r *AgentRepository) GetByUserID(ctx context.Context, userID string) (*partnerv1.Agent, error) {
	var agent partnerv1.Agent
	err := r.db.WithContext(ctx).Table("partner_schema.agents").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&agent).Error
	if err != nil {
		return nil, err
	}
	r.decrypt(&agent)
	return &agent, nil
}

// ListByPartner retrieves all agents associated with a specific partner
func (r *AgentRepository) ListByPartner(ctx context.Context, partnerID string, limit, offset int) ([]*partnerv1.Agent, error) {
	var agents []*partnerv1.Agent
	err := r.db.WithContext(ctx).Table("partner_schema.agents").
		Where("partner_id = ? AND deleted_at IS NULL", partnerID).
		Limit(limit).Offset(offset).
		Find(&agents).Error
	if err == nil {
		for _, a := range agents {
			r.decrypt(a)
		}
	}
	return agents, err
}

// UpdateStatus changes an agent's status
func (r *AgentRepository) UpdateStatus(ctx context.Context, agentID string, status partnerv1.AgentStatus) error {
	return r.db.WithContext(ctx).Table("partner_schema.agents").
		Where("agent_id = ?", agentID).
		Updates(map[string]interface{}{
			"status":     status.String(),
			"updated_at": time.Now(),
		}).Error
}
