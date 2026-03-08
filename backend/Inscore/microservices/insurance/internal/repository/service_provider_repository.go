package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/lib/pq"

	servicesv1 "github.com/newage-saint/insuretech/gen/go/insuretech/services/entity/v1"
)

type ServiceProviderRepository struct {
	db *gorm.DB
}

func NewServiceProviderRepository(db *gorm.DB) *ServiceProviderRepository {
	return &ServiceProviderRepository{db: db}
}

func (r *ServiceProviderRepository) Create(ctx context.Context, provider *servicesv1.ServiceProvider) (*servicesv1.ServiceProvider, error) {
	if provider.ProviderId == "" {
		return nil, fmt.Errorf("provider_id is required")
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.service_providers
			(provider_id, provider_name, provider_type, address, city, district, phone_number, email,
			 latitude, longitude, services_offered, is_network_provider, supported_product_categories)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		provider.ProviderId,
		provider.ProviderName,
		strings.ToUpper(provider.ProviderType.String()),
		provider.Address,
		provider.City,
		provider.District,
		provider.PhoneNumber,
		provider.Email,
		provider.Latitude,
		provider.Longitude,
		pq.Array(provider.ServicesOffered),
		provider.IsNetworkProvider,
		pq.Array(provider.SupportedProductCategories),
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert service provider: %w", err)
	}

	return r.GetByID(ctx, provider.ProviderId)
}

func (r *ServiceProviderRepository) GetByID(ctx context.Context, providerID string) (*servicesv1.ServiceProvider, error) {
	var (
		prov                        servicesv1.ServiceProvider
		providerTypeStr             sql.NullString
		address                     sql.NullString
		city                        sql.NullString
		district                    sql.NullString
		phoneNumber                 sql.NullString
		email                       sql.NullString
		latitude                    sql.NullFloat64
		longitude                   sql.NullFloat64
		servicesOffered             pq.StringArray
		supportedProductCategories  pq.StringArray
		createdAt                   time.Time
		updatedAt                   time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT provider_id, provider_name, provider_type, address, city, district, phone_number, email,
		       latitude, longitude, COALESCE(services_offered, '{}') as services_offered,
		       is_network_provider, COALESCE(supported_product_categories, '{}') as supported_product_categories,
		       created_at, updated_at
		FROM insurance_schema.service_providers
		WHERE provider_id = $1`,
		providerID,
	).Row().Scan(
		&prov.ProviderId,
		&prov.ProviderName,
		&providerTypeStr,
		&address,
		&city,
		&district,
		&phoneNumber,
		&email,
		&latitude,
		&longitude,
		&servicesOffered,
		&prov.IsNetworkProvider,
		&supportedProductCategories,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get service provider: %w", err)
	}

	if providerTypeStr.Valid {
		k := strings.ToUpper(providerTypeStr.String)
		if v, ok := servicesv1.ServiceProviderType_value[k]; ok {
			prov.ProviderType = servicesv1.ServiceProviderType(v)
		}
	}

	if address.Valid {
		prov.Address = address.String
	}
	if city.Valid {
		prov.City = city.String
	}
	if district.Valid {
		prov.District = district.String
	}
	if phoneNumber.Valid {
		prov.PhoneNumber = phoneNumber.String
	}
	if email.Valid {
		prov.Email = email.String
	}
	if latitude.Valid {
		prov.Latitude = latitude.Float64
	}
	if longitude.Valid {
		prov.Longitude = longitude.Float64
	}

	prov.ServicesOffered = servicesOffered
	prov.SupportedProductCategories = supportedProductCategories

	prov.CreatedAt = timestamppb.New(createdAt)
	prov.UpdatedAt = timestamppb.New(updatedAt)

	return &prov, nil
}

func (r *ServiceProviderRepository) Update(ctx context.Context, provider *servicesv1.ServiceProvider) (*servicesv1.ServiceProvider, error) {
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.service_providers
		SET provider_name = $2,
		    provider_type = $3,
		    address = $4,
		    city = $5,
		    district = $6,
		    phone_number = $7,
		    email = $8,
		    latitude = $9,
		    longitude = $10,
		    services_offered = $11,
		    is_network_provider = $12,
		    supported_product_categories = $13,
		    updated_at = NOW()
		WHERE provider_id = $1`,
		provider.ProviderId,
		provider.ProviderName,
		strings.ToUpper(provider.ProviderType.String()),
		provider.Address,
		provider.City,
		provider.District,
		provider.PhoneNumber,
		provider.Email,
		provider.Latitude,
		provider.Longitude,
		pq.Array(provider.ServicesOffered),
		provider.IsNetworkProvider,
		pq.Array(provider.SupportedProductCategories),
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update service provider: %w", err)
	}

	return r.GetByID(ctx, provider.ProviderId)
}

func (r *ServiceProviderRepository) Delete(ctx context.Context, providerID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.service_providers
		WHERE provider_id = $1`,
		providerID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete service provider: %w", err)
	}

	return nil
}

func (r *ServiceProviderRepository) List(ctx context.Context, providerType, city string, page, pageSize int) ([]*servicesv1.ServiceProvider, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := make([]interface{}, 0)
	argIndex := 1

	if providerType != "" {
		whereClause += fmt.Sprintf(" AND provider_type = $%d", argIndex)
		args = append(args, strings.ToUpper(providerType))
		argIndex++
	}

	if city != "" {
		whereClause += fmt.Sprintf(" AND city = $%d", argIndex)
		args = append(args, city)
		argIndex++
	}

	// Get total count
	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM insurance_schema.service_providers %s`, whereClause)
	err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count service providers: %w", err)
	}

	// Get providers
	query := fmt.Sprintf(`
		SELECT provider_id, provider_name, provider_type, address, city, district, phone_number, email,
		       latitude, longitude, COALESCE(services_offered, '{}') as services_offered,
		       is_network_provider, COALESCE(supported_product_categories, '{}') as supported_product_categories,
		       created_at, updated_at
		FROM insurance_schema.service_providers
		%s
		ORDER BY provider_name ASC LIMIT %d OFFSET %d`, whereClause, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list service providers: %w", err)
	}
	defer rows.Close()

	providers := make([]*servicesv1.ServiceProvider, 0)
	for rows.Next() {
		var (
			prov                       servicesv1.ServiceProvider
			providerTypeStr            sql.NullString
			address                    sql.NullString
			city                       sql.NullString
			district                   sql.NullString
			phoneNumber                sql.NullString
			email                      sql.NullString
			latitude                   sql.NullFloat64
			longitude                  sql.NullFloat64
			servicesOffered            pq.StringArray
			supportedProductCategories pq.StringArray
			createdAt                  time.Time
			updatedAt                  time.Time
		)

		err := rows.Scan(
			&prov.ProviderId,
			&prov.ProviderName,
			&providerTypeStr,
			&address,
			&city,
			&district,
			&phoneNumber,
			&email,
			&latitude,
			&longitude,
			&servicesOffered,
			&prov.IsNetworkProvider,
			&supportedProductCategories,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan service provider: %w", err)
		}

		if providerTypeStr.Valid {
			k := strings.ToUpper(providerTypeStr.String)
			if v, ok := servicesv1.ServiceProviderType_value[k]; ok {
				prov.ProviderType = servicesv1.ServiceProviderType(v)
			}
		}

		if address.Valid {
			prov.Address = address.String
		}
		if city.Valid {
			prov.City = city.String
		}
		if district.Valid {
			prov.District = district.String
		}
		if phoneNumber.Valid {
			prov.PhoneNumber = phoneNumber.String
		}
		if email.Valid {
			prov.Email = email.String
		}
		if latitude.Valid {
			prov.Latitude = latitude.Float64
		}
		if longitude.Valid {
			prov.Longitude = longitude.Float64
		}

		prov.ServicesOffered = servicesOffered
		prov.SupportedProductCategories = supportedProductCategories

		prov.CreatedAt = timestamppb.New(createdAt)
		prov.UpdatedAt = timestamppb.New(updatedAt)

		providers = append(providers, &prov)
	}

	return providers, total, nil
}
