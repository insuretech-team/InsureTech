package service

import (
	"context"
	"strings"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	vehiclev1 "github.com/newage-saint/insuretech/gen/go/insuretech/vehicle/entity/v1"
	vehicleservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/vehicle/services/v1"
	"gorm.io/gorm"
)

const (
	vehicleTable            = "vehicle_schema.vehicles"
	vehicleRegistrationTable = "vehicle_schema.vehicle_registrations"
	vehicleCalcTable        = "vehicle_schema.vehicle_premium_calculations"
)

type VehicleService struct {
	vehicleservicev1.UnimplementedVehicleServiceServer
	db *gorm.DB
}

func NewVehicleService(db *gorm.DB) *VehicleService {
	return &VehicleService{db: db}
}

func (s *VehicleService) GetVehicle(ctx context.Context, req *vehicleservicev1.GetVehicleRequest) (*vehicleservicev1.GetVehicleResponse, error) {
	var vehicle vehiclev1.Vehicle
	err := s.db.WithContext(ctx).
		Table(vehicleTable).
		Where("vehicle_id = ? AND deleted_at IS NULL", req.GetVehicleId()).
		First(&vehicle).Error
	if err != nil {
		return &vehicleservicev1.GetVehicleResponse{
			Error: errorResponse("VEHICLE_NOT_FOUND", "vehicle not found", 404),
		}, nil
	}

	return &vehicleservicev1.GetVehicleResponse{Vehicle: &vehicle}, nil
}

func (s *VehicleService) GetVehicleByModel(ctx context.Context, req *vehicleservicev1.GetVehicleByModelRequest) (*vehicleservicev1.GetVehicleResponse, error) {
	var vehicle vehiclev1.Vehicle
	err := s.db.WithContext(ctx).
		Table(vehicleTable).
		Where("LOWER(model) = LOWER(?) AND deleted_at IS NULL", req.GetModel()).
		First(&vehicle).Error
	if err != nil {
		return &vehicleservicev1.GetVehicleResponse{
			Error: errorResponse("VEHICLE_NOT_FOUND", "vehicle not found", 404),
		}, nil
	}

	return &vehicleservicev1.GetVehicleResponse{Vehicle: &vehicle}, nil
}

func (s *VehicleService) ListVehicles(ctx context.Context, req *vehicleservicev1.ListVehiclesRequest) (*vehicleservicev1.ListVehiclesResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(vehicleTable).Where("deleted_at IS NULL")
	if req.GetOnlyActive() {
		query = query.Where("is_active = ?", true)
	}
	if req.GetVehicleType() != vehiclev1.VehicleType_VEHICLE_TYPE_UNSPECIFIED {
		query = query.Where("type = ?", req.GetVehicleType())
	}
	if manufacturer := strings.TrimSpace(req.GetManufacturer()); manufacturer != "" {
		query = query.Where("LOWER(manufacturer) = LOWER(?)", manufacturer)
	}
	if req.GetYear() > 0 {
		query = query.Where("\"year\" = ?", req.GetYear())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &vehicleservicev1.ListVehiclesResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	var vehicles []*vehiclev1.Vehicle
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&vehicles).Error; err != nil {
		return &vehicleservicev1.ListVehiclesResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &vehicleservicev1.ListVehiclesResponse{
		Vehicles:       vehicles,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(vehicles), int(total)),
	}, nil
}

func (s *VehicleService) CreateVehicle(ctx context.Context, req *vehicleservicev1.CreateVehicleRequest) (*vehicleservicev1.CreateVehicleResponse, error) {
	if strings.TrimSpace(req.GetModel()) == "" {
		return &vehicleservicev1.CreateVehicleResponse{
			Error: errorResponse("INVALID_ARGUMENT", "model is required", 400),
		}, nil
	}

	now := time.Now().UTC()
	vehicle := &vehiclev1.Vehicle{
		VehicleId:    newID(),
		Model:        req.GetModel(),
		Type:         req.GetType(),
		Price:        req.GetPrice(),
		Manufacturer: req.GetManufacturer(),
		Year:         req.GetYear(),
		ImageUri:     req.GetImageUri(),
		IsActive:     true,
		Metadata:     ensureMetadata(req.GetMetadata()),
		CreatedAt:    nowTS(),
		UpdatedAt:    nowTS(),
	}
	if err := s.db.WithContext(ctx).Table(vehicleTable).Create(vehicle).Error; err != nil {
		return &vehicleservicev1.CreateVehicleResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	_ = now
	return &vehicleservicev1.CreateVehicleResponse{Vehicle: vehicle}, nil
}

func (s *VehicleService) UpdateVehicle(ctx context.Context, req *vehicleservicev1.UpdateVehicleRequest) (*vehicleservicev1.UpdateVehicleResponse, error) {
	var vehicle vehiclev1.Vehicle
	err := s.db.WithContext(ctx).
		Table(vehicleTable).
		Where("vehicle_id = ? AND deleted_at IS NULL", req.GetVehicleId()).
		First(&vehicle).Error
	if err != nil {
		return &vehicleservicev1.UpdateVehicleResponse{
			Error: errorResponse("VEHICLE_NOT_FOUND", "vehicle not found", 404),
		}, nil
	}

	vehicle.Model = req.GetModel()
	vehicle.Type = req.GetType()
	vehicle.Price = req.GetPrice()
	vehicle.Manufacturer = req.GetManufacturer()
	vehicle.Year = req.GetYear()
	vehicle.ImageUri = req.GetImageUri()
	vehicle.IsActive = req.GetIsActive()
	if req.GetMetadata() != nil {
		vehicle.Metadata = ensureMetadata(req.GetMetadata())
	}
	vehicle.UpdatedAt = nowTS()

	if err := s.db.WithContext(ctx).Table(vehicleTable).
		Where("vehicle_id = ?", vehicle.GetVehicleId()).
		Save(&vehicle).Error; err != nil {
		return &vehicleservicev1.UpdateVehicleResponse{
			Error: errorResponse("UPDATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &vehicleservicev1.UpdateVehicleResponse{Vehicle: &vehicle}, nil
}

func (s *VehicleService) DeleteVehicle(ctx context.Context, req *vehicleservicev1.DeleteVehicleRequest) (*vehicleservicev1.DeleteVehicleResponse, error) {
	query := s.db.WithContext(ctx).Table(vehicleTable).Where("vehicle_id = ?", req.GetVehicleId())
	var result *gorm.DB
	if req.GetPermanent() {
		result = query.Delete(&vehiclev1.Vehicle{})
	} else {
		result = query.Update("deleted_at", nowTS())
	}
	if result.Error != nil {
		return &vehicleservicev1.DeleteVehicleResponse{
			Error: errorResponse("DELETE_FAILED", result.Error.Error(), 500),
		}, nil
	}

	return &vehicleservicev1.DeleteVehicleResponse{Success: result.RowsAffected > 0}, nil
}

func (s *VehicleService) GetVehicleModels(ctx context.Context, req *vehicleservicev1.GetVehicleModelsRequest) (*vehicleservicev1.GetVehicleModelsResponse, error) {
	query := s.db.WithContext(ctx).Table(vehicleTable).Where("deleted_at IS NULL")
	if req.GetVehicleType() != vehiclev1.VehicleType_VEHICLE_TYPE_UNSPECIFIED {
		query = query.Where("type = ?", req.GetVehicleType())
	}

	var models []string
	if err := query.Distinct("model").Order("model").Pluck("model", &models).Error; err != nil {
		return &vehicleservicev1.GetVehicleModelsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &vehicleservicev1.GetVehicleModelsResponse{Models: models}, nil
}

func (s *VehicleService) RegisterVehicle(ctx context.Context, req *vehicleservicev1.RegisterVehicleRequest) (*vehicleservicev1.RegisterVehicleResponse, error) {
	vehicle, errResp := s.getVehicle(ctx, req.GetVehicleId())
	if errResp != nil {
		return &vehicleservicev1.RegisterVehicleResponse{Error: errResp}, nil
	}

	year := req.GetRegistrationYear()
	currentYear := int32(time.Now().UTC().Year())
	age := int32(0)
	if year > 0 && year <= currentYear {
		age = currentYear - year
	}

	registration := &vehiclev1.VehicleRegistration{
		RegistrationId:     newID(),
		VehicleId:          req.GetVehicleId(),
		OwnerId:            req.GetOwnerId(),
		RegistrationNumber: req.GetRegistrationNumber(),
		RegistrationYear:   year,
		RegistrationState:  extractRegistrationState(req.GetRegistrationNumber()),
		Status:             vehiclev1.RegistrationStatus_REGISTRATION_STATUS_ACTIVE,
		VehicleAge:         age,
		CurrentValue:       depreciatedVehicleValue(vehicle.GetPrice(), age),
		AdditionalInfo:     ensureMetadata(req.GetAdditionalInfo()),
		CreatedAt:          nowTS(),
		UpdatedAt:          nowTS(),
	}

	if err := s.db.WithContext(ctx).Table(vehicleRegistrationTable).Create(registration).Error; err != nil {
		return &vehicleservicev1.RegisterVehicleResponse{
			Error: errorResponse("CREATE_FAILED", err.Error(), 500),
		}, nil
	}

	return &vehicleservicev1.RegisterVehicleResponse{Registration: registration}, nil
}

func (s *VehicleService) GetVehicleRegistration(ctx context.Context, req *vehicleservicev1.GetVehicleRegistrationRequest) (*vehicleservicev1.GetVehicleRegistrationResponse, error) {
	var registration vehiclev1.VehicleRegistration
	err := s.db.WithContext(ctx).
		Table(vehicleRegistrationTable).
		Where("registration_id = ?", req.GetRegistrationId()).
		First(&registration).Error
	if err != nil {
		return &vehicleservicev1.GetVehicleRegistrationResponse{
			Error: errorResponse("REGISTRATION_NOT_FOUND", "vehicle registration not found", 404),
		}, nil
	}

	vehicle, _ := s.getVehicle(ctx, registration.GetVehicleId())
	return &vehicleservicev1.GetVehicleRegistrationResponse{
		Registration: &registration,
		Vehicle:      vehicle,
	}, nil
}

func (s *VehicleService) CalculatePremium(ctx context.Context, req *vehicleservicev1.CalculatePremiumRequest) (*vehicleservicev1.CalculatePremiumResponse, error) {
	registration, errResp := s.getRegistration(ctx, req.GetRegistrationId(), req.GetRegistrationNumber())
	if errResp != nil {
		return &vehicleservicev1.CalculatePremiumResponse{Error: errResp}, nil
	}

	vehicle, vehicleErr := s.getVehicle(ctx, registration.GetVehicleId())
	if vehicleErr != nil {
		return &vehicleservicev1.CalculatePremiumResponse{Error: vehicleErr}, nil
	}

	calculation := calculateVehiclePremium(vehicle, registration, req.GetAccidentalCover())
	if err := s.db.WithContext(ctx).Table(vehicleCalcTable).Create(calculation).Error; err != nil {
		return &vehicleservicev1.CalculatePremiumResponse{
			Error: errorResponse("CALCULATION_FAILED", err.Error(), 500),
		}, nil
	}

	return &vehicleservicev1.CalculatePremiumResponse{Calculation: calculation}, nil
}

func (s *VehicleService) GetPremiumCalculations(ctx context.Context, req *vehicleservicev1.GetPremiumCalculationsRequest) (*vehicleservicev1.GetPremiumCalculationsResponse, error) {
	offset := pageOffset(req.GetPageToken())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 50
	}

	query := s.db.WithContext(ctx).Table(vehicleCalcTable).Where("registration_id = ?", req.GetRegistrationId())

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &vehicleservicev1.GetPremiumCalculationsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	var calculations []*vehiclev1.VehiclePremiumCalculation
	if err := query.Order("calculated_at DESC").Offset(offset).Limit(pageSize).Find(&calculations).Error; err != nil {
		return &vehicleservicev1.GetPremiumCalculationsResponse{
			Error: errorResponse("LIST_FAILED", err.Error(), 500),
		}, nil
	}

	return &vehicleservicev1.GetPremiumCalculationsResponse{
		Calculations:   calculations,
		TotalCount:     int32(total),
		NextPageToken:  nextToken(offset, len(calculations), int(total)),
	}, nil
}

func (s *VehicleService) GetVehicleImages(_ context.Context, _ *vehicleservicev1.GetVehicleImagesRequest) (*vehicleservicev1.GetVehicleImagesResponse, error) {
	return &vehicleservicev1.GetVehicleImagesResponse{
		Images: map[string]string{
			"bike":       "/assets/images/bike.jpg",
			"car":        "/assets/images/car.jpg",
			"commercial": "/assets/images/commercial.jpg",
		},
	}, nil
}

func (s *VehicleService) getVehicle(ctx context.Context, vehicleID string) (*vehiclev1.Vehicle, *commonv1.Error) {
	var vehicle vehiclev1.Vehicle
	err := s.db.WithContext(ctx).
		Table(vehicleTable).
		Where("vehicle_id = ? AND deleted_at IS NULL", vehicleID).
		First(&vehicle).Error
	if err != nil {
		return nil, errorResponse("VEHICLE_NOT_FOUND", "vehicle not found", 404)
	}
	return &vehicle, nil
}

func (s *VehicleService) getRegistration(ctx context.Context, registrationID, registrationNumber string) (*vehiclev1.VehicleRegistration, *commonv1.Error) {
	query := s.db.WithContext(ctx).Table(vehicleRegistrationTable)
	switch {
	case strings.TrimSpace(registrationID) != "":
		query = query.Where("registration_id = ?", registrationID)
	case strings.TrimSpace(registrationNumber) != "":
		query = query.Where("registration_number = ?", registrationNumber)
	default:
		return nil, errorResponse("INVALID_ARGUMENT", "registration_id or registration_number is required", 400)
	}

	var registration vehiclev1.VehicleRegistration
	if err := query.First(&registration).Error; err != nil {
		return nil, errorResponse("REGISTRATION_NOT_FOUND", "vehicle registration not found", 404)
	}
	return &registration, nil
}

func calculateVehiclePremium(vehicle *vehiclev1.Vehicle, registration *vehiclev1.VehicleRegistration, accidentalCover bool) *vehiclev1.VehiclePremiumCalculation {
	typeMultiplier := float32(1.0)
	switch vehicle.GetType() {
	case vehiclev1.VehicleType_VEHICLE_TYPE_CAR:
		typeMultiplier = 1.2
	case vehiclev1.VehicleType_VEHICLE_TYPE_COMMERCIAL, vehiclev1.VehicleType_VEHICLE_TYPE_TRUCK, vehiclev1.VehicleType_VEHICLE_TYPE_BUS:
		typeMultiplier = 1.35
	}

	ageMultiplier := float32(1.5)
	switch age := registration.GetVehicleAge(); {
	case age <= 2:
		ageMultiplier = 1.0
	case age <= 5:
		ageMultiplier = 1.2
	}

	currentValue := depreciatedVehicleValue(vehicle.GetPrice(), registration.GetVehicleAge())
	valueMultiplier := float32(1.2)
	if currentValue >= vehicle.GetPrice() {
		valueMultiplier = 1.0
	}

	location := strings.ToUpper(registration.GetRegistrationState())
	locationMultiplier := float32(1.2)
	switch location {
	case "TN":
		locationMultiplier = 1.0
	case "AP", "TL", "KL", "KA":
		locationMultiplier = 1.1
	}

	basePremiumAmount := int64(float64(vehicle.GetPrice()) * 0.10)
	totalPremium := int64(float64(basePremiumAmount) * float64(typeMultiplier) * float64(ageMultiplier) * float64(valueMultiplier) * float64(locationMultiplier))

	tp1 := totalPremium
	tp2 := totalPremium * 2
	tp3 := totalPremium * 3
	compFactor := int64(12)
	compDivisor := int64(10)
	comp1 := tp1 * compFactor / compDivisor
	comp2 := tp2 * compFactor / compDivisor
	comp3 := tp3 * compFactor / compDivisor
	if accidentalCover {
		comp1 += 500
		comp2 += 1000
		comp3 += 1500
	}

	return &vehiclev1.VehiclePremiumCalculation{
		CalculationId:          newID(),
		RegistrationId:         registration.GetRegistrationId(),
		BasePremium:            money(basePremiumAmount, "BDT"),
		TypeMultiplier:         typeMultiplier,
		AgeMultiplier:          ageMultiplier,
		ValueMultiplier:        valueMultiplier,
		LocationMultiplier:     locationMultiplier,
		TpPremium_1Year:        money(tp1, "BDT"),
		TpPremium_2Year:        money(tp2, "BDT"),
		TpPremium_3Year:        money(tp3, "BDT"),
		CompPremium_1Year:      money(comp1, "BDT"),
		CompPremium_2Year:      money(comp2, "BDT"),
		CompPremium_3Year:      money(comp3, "BDT"),
		AccidentalCover:        accidentalCover,
		CalculatedAt:           nowTS(),
		CalculationDurationMs:  1,
	}
}

func depreciatedVehicleValue(price int64, age int32) int64 {
	if price <= 0 {
		return 0
	}

	current := price
	switch {
	case age <= 2:
		current = price
	case age <= 5:
		current = int64(float64(price) * 0.8)
	default:
		current = int64(float64(price) * 0.5)
	}

	minValue := int64(float64(price) * 0.5)
	if current < minValue {
		return minValue
	}
	return current
}

func extractRegistrationState(registrationNumber string) string {
	number := strings.TrimSpace(registrationNumber)
	if number == "" {
		return "NA"
	}

	parts := strings.Split(number, "-")
	if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
		return strings.ToUpper(strings.TrimSpace(parts[0]))
	}

	if len(number) >= 2 {
		return strings.ToUpper(number[:2])
	}

	return strings.ToUpper(number)
}
