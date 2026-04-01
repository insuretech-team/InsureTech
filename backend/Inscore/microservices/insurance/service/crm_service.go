package service

import (
	"context"
	"sort"
	"strings"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	crmv1 "github.com/newage-saint/insuretech/gen/go/insuretech/crm/entity/v1"
	crmservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/crm/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	leadTable    = "crm_schema.leads"
	contactTable = "crm_schema.contacts"
)

type CrmService struct {
	crmservicev1.UnimplementedCrmServiceServer
	db *gorm.DB
}

func NewCrmService(db *gorm.DB) *CrmService {
	return &CrmService{db: db}
}

func (s *CrmService) CreateLead(ctx context.Context, req *crmservicev1.CreateLeadRequest) (*crmservicev1.CreateLeadResponse, error) {
	if strings.TrimSpace(req.GetFirstName()) == "" || strings.TrimSpace(req.GetLastName()) == "" {
		return &crmservicev1.CreateLeadResponse{
			Error: errorResponse("INVALID_ARGUMENT", "first_name and last_name are required", 400),
		}, nil
	}

	now := nowTS()
	lead := &crmv1.Lead{
		LeadId:                newID(),
		Title:                 req.GetTitle(),
		FirstName:             req.GetFirstName(),
		LastName:              req.GetLastName(),
		DateOfBirth:           req.GetDateOfBirth(),
		Gender:                req.GetGender(),
		PhoneNumber:           req.GetPhoneNumber(),
		EmailAddress:          req.GetEmailAddress(),
		Address:               req.GetAddress(),
		LeadSource:            req.GetLeadSource(),
		LeadStatus:            crmv1.LeadStatus_LEAD_STATUS_NEW,
		LeadPriority:          req.GetLeadPriority(),
		AssignedAgentId:       req.GetAssignedAgentId(),
		DesiredInsuranceType:  req.GetDesiredInsuranceType(),
		DesiredCoverageAmount: req.GetDesiredCoverageAmount(),
		SpecificRequirements:  req.GetSpecificRequirements(),
		QualificationStatus:   crmv1.QualificationStatus_QUALIFICATION_STATUS_IN_PROGRESS,
		Metadata:              ensureMetadata(req.GetMetadata()),
		CreatedBy:             req.GetCreatedBy(),
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := s.db.WithContext(ctx).Table(leadTable).Create(lead).Error; err != nil {
		return &crmservicev1.CreateLeadResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.CreateLeadResponse{Lead: lead}, nil
}

func (s *CrmService) GetLead(ctx context.Context, req *crmservicev1.GetLeadRequest) (*crmservicev1.GetLeadResponse, error) {
	var lead crmv1.Lead
	err := s.db.WithContext(ctx).
		Table(leadTable).
		Where("lead_id = ? AND deleted_at IS NULL", req.GetLeadId()).
		First(&lead).Error
	if err != nil {
		return &crmservicev1.GetLeadResponse{
			Error: errorResponse("LEAD_NOT_FOUND", "lead not found", 404),
		}, nil
	}

	return &crmservicev1.GetLeadResponse{Lead: &lead}, nil
}

func (s *CrmService) ListLeads(ctx context.Context, req *crmservicev1.ListLeadsRequest) (*crmservicev1.ListLeadsResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(leadTable).Where("deleted_at IS NULL")
	if req.GetStatus() != crmv1.LeadStatus_LEAD_STATUS_UNSPECIFIED {
		query = query.Where("lead_status = ?", req.GetStatus())
	}
	if req.GetSource() != crmv1.LeadSource_LEAD_SOURCE_UNSPECIFIED {
		query = query.Where("lead_source = ?", req.GetSource())
	}
	if agentID := strings.TrimSpace(req.GetAssignedAgentId()); agentID != "" {
		query = query.Where("assigned_agent_id = ?", agentID)
	}
	if search := strings.TrimSpace(req.GetSearchQuery()); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email_address) LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &crmservicev1.ListLeadsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	var leads []*crmv1.Lead
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&leads).Error; err != nil {
		return &crmservicev1.ListLeadsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.ListLeadsResponse{
		Leads:          leads,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(leads), int(total)),
	}, nil
}

func (s *CrmService) UpdateLead(ctx context.Context, req *crmservicev1.UpdateLeadRequest) (*crmservicev1.UpdateLeadResponse, error) {
	var lead crmv1.Lead
	err := s.db.WithContext(ctx).
		Table(leadTable).
		Where("lead_id = ? AND deleted_at IS NULL", req.GetLeadId()).
		First(&lead).Error
	if err != nil {
		return &crmservicev1.UpdateLeadResponse{
			Error: errorResponse("LEAD_NOT_FOUND", "lead not found", 404),
		}, nil
	}

	lead.Title = req.GetTitle()
	lead.FirstName = req.GetFirstName()
	lead.LastName = req.GetLastName()
	lead.PhoneNumber = req.GetPhoneNumber()
	lead.EmailAddress = req.GetEmailAddress()
	lead.Address = req.GetAddress()
	lead.LeadStatus = req.GetLeadStatus()
	lead.LeadPriority = req.GetLeadPriority()
	lead.LeadScore = req.GetLeadScore()
	lead.AssignedAgentId = req.GetAssignedAgentId()
	lead.DesiredInsuranceType = req.GetDesiredInsuranceType()
	lead.DesiredCoverageAmount = req.GetDesiredCoverageAmount()
	lead.QualificationStatus = req.GetQualificationStatus()
	if req.GetMetadata() != nil {
		lead.Metadata = ensureMetadata(req.GetMetadata())
	}
	lead.UpdatedAt = nowTS()

	if err := s.db.WithContext(ctx).Table(leadTable).
		Where("lead_id = ?", lead.GetLeadId()).
		Save(&lead).Error; err != nil {
		return &crmservicev1.UpdateLeadResponse{
			Error: errorResponse("UPDATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.UpdateLeadResponse{Lead: &lead}, nil
}

func (s *CrmService) DeleteLead(ctx context.Context, req *crmservicev1.DeleteLeadRequest) (*crmservicev1.DeleteLeadResponse, error) {
	query := s.db.WithContext(ctx).Table(leadTable).Where("lead_id = ?", req.GetLeadId())
	var result *gorm.DB
	if req.GetPermanent() {
		result = query.Delete(&crmv1.Lead{})
	} else {
		result = query.Update("deleted_at", nowTS())
	}
	if result.Error != nil {
		return &crmservicev1.DeleteLeadResponse{
			Error: errorResponse("DELETE_FAILED", result.Error.Error(), 500),
		}, nil
	}

	return &crmservicev1.DeleteLeadResponse{Success: result.RowsAffected > 0}, nil
}

func (s *CrmService) ConvertLeadToContact(ctx context.Context, req *crmservicev1.ConvertLeadToContactRequest) (*crmservicev1.ConvertLeadToContactResponse, error) {
	lead, errResp := s.getLead(ctx, req.GetLeadId())
	if errResp != nil {
		return &crmservicev1.ConvertLeadToContactResponse{Error: errResp}, nil
	}

	now := nowTS()
	contact := &crmv1.Contact{
		ContactId:              newID(),
		FirstName:              lead.GetFirstName(),
		LastName:               lead.GetLastName(),
		DateOfBirth:            lead.GetDateOfBirth(),
		Gender:                 lead.GetGender(),
		PhoneNumber:            lead.GetPhoneNumber(),
		EmailAddress:           lead.GetEmailAddress(),
		Address:                lead.GetAddress(),
		AlternatePhoneNumber:   lead.GetAlternatePhoneNumber(),
		AdditionalEmailAddress: lead.GetAdditionalEmailAddress(),
		AssignedAgentId:        lead.GetAssignedAgentId(),
		ContactType:            crmv1.ContactType_CONTACT_TYPE_INDIVIDUAL,
		ContactStatus:          crmv1.ContactStatus_CONTACT_STATUS_ACTIVE,
		ReferralSource:         lead.GetLeadSource().String(),
		PreferredLanguage:      "English",
		Metadata:               ensureMetadata(lead.GetMetadata()),
		CreatedBy:              lead.GetCreatedBy(),
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if err := s.db.WithContext(ctx).Table(contactTable).Create(contact).Error; err != nil {
		return &crmservicev1.ConvertLeadToContactResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	lead.LeadStatus = crmv1.LeadStatus_LEAD_STATUS_CONVERTED
	lead.ConvertedContactId = contact.GetContactId()
	lead.ConvertedAt = now
	lead.ConversionReason = req.GetConversionReason()
	lead.UpdatedAt = now
	if err := s.db.WithContext(ctx).Table(leadTable).Where("lead_id = ?", lead.GetLeadId()).Save(lead).Error; err != nil {
		return &crmservicev1.ConvertLeadToContactResponse{
			Error: errorResponse("UPDATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.ConvertLeadToContactResponse{
		Contact: contact,
		Lead:    lead,
	}, nil
}

func (s *CrmService) AssignLead(ctx context.Context, req *crmservicev1.AssignLeadRequest) (*crmservicev1.AssignLeadResponse, error) {
	lead, errResp := s.getLead(ctx, req.GetLeadId())
	if errResp != nil {
		return &crmservicev1.AssignLeadResponse{Error: errResp}, nil
	}

	lead.AssignedAgentId = req.GetAgentId()
	lead.UpdatedAt = nowTS()
	if err := s.db.WithContext(ctx).Table(leadTable).Where("lead_id = ?", lead.GetLeadId()).Save(lead).Error; err != nil {
		return &crmservicev1.AssignLeadResponse{
			Error: errorResponse("UPDATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.AssignLeadResponse{Lead: lead}, nil
}

func (s *CrmService) CreateContact(ctx context.Context, req *crmservicev1.CreateContactRequest) (*crmservicev1.CreateContactResponse, error) {
	if strings.TrimSpace(req.GetFirstName()) == "" || strings.TrimSpace(req.GetLastName()) == "" {
		return &crmservicev1.CreateContactResponse{
			Error: errorResponse("INVALID_ARGUMENT", "first_name and last_name are required", 400),
		}, nil
	}

	now := nowTS()
	contact := &crmv1.Contact{
		ContactId:              newID(),
		FirstName:              req.GetFirstName(),
		LastName:               req.GetLastName(),
		DateOfBirth:            req.GetDateOfBirth(),
		Gender:                 req.GetGender(),
		PhoneNumber:            req.GetPhoneNumber(),
		EmailAddress:           req.GetEmailAddress(),
		Address:                req.GetAddress(),
		PreferredContactMethod: req.GetPreferredContactMethod(),
		ReferralSource:         req.GetReferralSource(),
		PreferredLanguage:      req.GetPreferredLanguage(),
		MarketingConsent:       req.GetMarketingConsent(),
		AssignedAgentId:        req.GetAssignedAgentId(),
		ContactType:            req.GetContactType(),
		ContactStatus:          crmv1.ContactStatus_CONTACT_STATUS_ACTIVE,
		Metadata:               ensureMetadata(req.GetMetadata()),
		CreatedBy:              req.GetCreatedBy(),
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.db.WithContext(ctx).Table(contactTable).Create(contact).Error; err != nil {
		return &crmservicev1.CreateContactResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.CreateContactResponse{Contact: contact}, nil
}

func (s *CrmService) GetContact(ctx context.Context, req *crmservicev1.GetContactRequest) (*crmservicev1.GetContactResponse, error) {
	var contact crmv1.Contact
	err := s.db.WithContext(ctx).
		Table(contactTable).
		Where("contact_id = ? AND deleted_at IS NULL", req.GetContactId()).
		First(&contact).Error
	if err != nil {
		return &crmservicev1.GetContactResponse{
			Error: errorResponse("CONTACT_NOT_FOUND", "contact not found", 404),
		}, nil
	}

	return &crmservicev1.GetContactResponse{Contact: &contact}, nil
}

func (s *CrmService) ListContacts(ctx context.Context, req *crmservicev1.ListContactsRequest) (*crmservicev1.ListContactsResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(contactTable).Where("deleted_at IS NULL")
	if req.GetStatus() != crmv1.ContactStatus_CONTACT_STATUS_UNSPECIFIED {
		query = query.Where("contact_status = ?", req.GetStatus())
	}
	if req.GetContactType() != crmv1.ContactType_CONTACT_TYPE_UNSPECIFIED {
		query = query.Where("contact_type = ?", req.GetContactType())
	}
	if agentID := strings.TrimSpace(req.GetAssignedAgentId()); agentID != "" {
		query = query.Where("assigned_agent_id = ?", agentID)
	}
	if search := strings.TrimSpace(req.GetSearchQuery()); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email_address) LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &crmservicev1.ListContactsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	var contacts []*crmv1.Contact
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&contacts).Error; err != nil {
		return &crmservicev1.ListContactsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.ListContactsResponse{
		Contacts:       contacts,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(contacts), int(total)),
	}, nil
}

func (s *CrmService) UpdateContact(ctx context.Context, req *crmservicev1.UpdateContactRequest) (*crmservicev1.UpdateContactResponse, error) {
	var contact crmv1.Contact
	err := s.db.WithContext(ctx).
		Table(contactTable).
		Where("contact_id = ? AND deleted_at IS NULL", req.GetContactId()).
		First(&contact).Error
	if err != nil {
		return &crmservicev1.UpdateContactResponse{
			Error: errorResponse("CONTACT_NOT_FOUND", "contact not found", 404),
		}, nil
	}

	contact.FirstName = req.GetFirstName()
	contact.LastName = req.GetLastName()
	contact.PhoneNumber = req.GetPhoneNumber()
	contact.EmailAddress = req.GetEmailAddress()
	contact.Address = req.GetAddress()
	contact.PreferredContactMethod = req.GetPreferredContactMethod()
	contact.MarketingConsent = req.GetMarketingConsent()
	contact.AssignedAgentId = req.GetAssignedAgentId()
	contact.ContactStatus = req.GetContactStatus()
	if req.GetMetadata() != nil {
		contact.Metadata = ensureMetadata(req.GetMetadata())
	}
	contact.UpdatedAt = nowTS()

	if err := s.db.WithContext(ctx).Table(contactTable).Where("contact_id = ?", contact.GetContactId()).Save(&contact).Error; err != nil {
		return &crmservicev1.UpdateContactResponse{
			Error: errorResponse("UPDATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &crmservicev1.UpdateContactResponse{Contact: &contact}, nil
}

func (s *CrmService) DeleteContact(ctx context.Context, req *crmservicev1.DeleteContactRequest) (*crmservicev1.DeleteContactResponse, error) {
	query := s.db.WithContext(ctx).Table(contactTable).Where("contact_id = ?", req.GetContactId())
	var result *gorm.DB
	if req.GetPermanent() {
		result = query.Delete(&crmv1.Contact{})
	} else {
		result = query.Update("deleted_at", nowTS())
	}
	if result.Error != nil {
		return &crmservicev1.DeleteContactResponse{
			Error: errorResponse("DELETE_FAILED", result.Error.Error(), 500),
		}, nil
	}

	return &crmservicev1.DeleteContactResponse{Success: result.RowsAffected > 0}, nil
}

func (s *CrmService) GetCrmStatistics(ctx context.Context, req *crmservicev1.GetCrmStatisticsRequest) (*crmservicev1.GetCrmStatisticsResponse, error) {
	leads, err := s.listLeadsForStats(ctx, req.GetAgentId(), req.GetStartDate(), req.GetEndDate())
	if err != nil {
		return &crmservicev1.GetCrmStatisticsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	contacts, err := s.listContactsForStats(ctx, req.GetAgentId())
	if err != nil {
		return &crmservicev1.GetCrmStatisticsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	resp := &crmservicev1.GetCrmStatisticsResponse{
		TotalLeads:     int32(len(leads)),
		TotalContacts:  int32(len(contacts)),
		LeadsBySource:  map[string]int32{},
	}
	for _, lead := range leads {
		switch lead.GetLeadStatus() {
		case crmv1.LeadStatus_LEAD_STATUS_NEW:
			resp.NewLeads++
		case crmv1.LeadStatus_LEAD_STATUS_CONTACTED:
			resp.ContactedLeads++
		case crmv1.LeadStatus_LEAD_STATUS_QUALIFIED:
			resp.QualifiedLeads++
		case crmv1.LeadStatus_LEAD_STATUS_CONVERTED:
			resp.ConvertedLeads++
		case crmv1.LeadStatus_LEAD_STATUS_LOST:
			resp.LostLeads++
		}
		resp.LeadsBySource[lead.GetLeadSource().String()]++
	}
	for _, contact := range contacts {
		if contact.GetContactStatus() == crmv1.ContactStatus_CONTACT_STATUS_ACTIVE {
			resp.ActiveContacts++
		}
	}
	if resp.GetTotalLeads() > 0 {
		resp.ConversionRate = float64(resp.GetConvertedLeads()) / float64(resp.GetTotalLeads()) * 100
	}

	return resp, nil
}

func (s *CrmService) GetAgentDashboard(ctx context.Context, req *crmservicev1.GetAgentDashboardRequest) (*crmservicev1.GetAgentDashboardResponse, error) {
	leads, err := s.listLeadsForStats(ctx, req.GetAgentId(), nil, nil)
	if err != nil {
		return &crmservicev1.GetAgentDashboardResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}
	contacts, err := s.listContactsForStats(ctx, req.GetAgentId())
	if err != nil {
		return &crmservicev1.GetAgentDashboardResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	startOfToday := time.Now().UTC().Truncate(24 * time.Hour)
	startOfMonth := time.Date(startOfToday.Year(), startOfToday.Month(), 1, 0, 0, 0, 0, time.UTC)
	resp := &crmservicev1.GetAgentDashboardResponse{
		MyLeadsCount:    int32(len(leads)),
		MyContactsCount: int32(len(contacts)),
	}

	for _, lead := range leads {
		if lead.GetCreatedAt() != nil && !lead.GetCreatedAt().AsTime().Before(startOfToday) {
			resp.NewLeadsToday++
		}
		if lead.GetLeadStatus() == crmv1.LeadStatus_LEAD_STATUS_CONVERTED && lead.GetConvertedAt() != nil && !lead.GetConvertedAt().AsTime().Before(startOfMonth) {
			resp.LeadsConvertedThisMonth++
		}
	}

	sortLeadPointersByCreatedDesc(leads)
	sortContactPointersByCreatedDesc(contacts)
	resp.RecentLeads = takeLeads(leads, 5)
	resp.RecentContacts = takeContacts(contacts, 5)

	return resp, nil
}

func (s *CrmService) getLead(ctx context.Context, leadID string) (*crmv1.Lead, *commonv1.Error) {
	var lead crmv1.Lead
	err := s.db.WithContext(ctx).
		Table(leadTable).
		Where("lead_id = ? AND deleted_at IS NULL", leadID).
		First(&lead).Error
	if err != nil {
		return nil, errorResponse("LEAD_NOT_FOUND", "lead not found", 404)
	}
	return &lead, nil
}

func (s *CrmService) listLeadsForStats(ctx context.Context, agentID string, start, end *timestamppb.Timestamp) ([]*crmv1.Lead, error) {
	query := s.db.WithContext(ctx).Table(leadTable).Where("deleted_at IS NULL")
	if strings.TrimSpace(agentID) != "" {
		query = query.Where("assigned_agent_id = ?", agentID)
	}
	if start != nil {
		query = query.Where("created_at >= ?", start.AsTime())
	}
	if end != nil {
		query = query.Where("created_at <= ?", end.AsTime())
	}
	var leads []*crmv1.Lead
	return leads, query.Find(&leads).Error
}

func (s *CrmService) listContactsForStats(ctx context.Context, agentID string) ([]*crmv1.Contact, error) {
	query := s.db.WithContext(ctx).Table(contactTable).Where("deleted_at IS NULL")
	if strings.TrimSpace(agentID) != "" {
		query = query.Where("assigned_agent_id = ?", agentID)
	}
	var contacts []*crmv1.Contact
	return contacts, query.Find(&contacts).Error
}

func sortLeadPointersByCreatedDesc(leads []*crmv1.Lead) {
	sort.Slice(leads, func(i, j int) bool {
		return leads[i].GetCreatedAt().AsTime().After(leads[j].GetCreatedAt().AsTime())
	})
}

func sortContactPointersByCreatedDesc(contacts []*crmv1.Contact) {
	sort.Slice(contacts, func(i, j int) bool {
		return contacts[i].GetCreatedAt().AsTime().After(contacts[j].GetCreatedAt().AsTime())
	})
}

func takeLeads(leads []*crmv1.Lead, count int) []*crmv1.Lead {
	if len(leads) <= count {
		return leads
	}
	return leads[:count]
}

func takeContacts(contacts []*crmv1.Contact, count int) []*crmv1.Contact {
	if len(contacts) <= count {
		return contacts
	}
	return contacts[:count]
}
