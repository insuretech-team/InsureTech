Database Design for the Minimum viable schema (MVS) for the InsureTech platform.

◆	Core Tables
•	Customer
•	Policy
•	Product
•	Payment
•	Claims

✅ 1. User Profile (Master Profile – Common for All Lines) Personal Information
•	Customer ID (System Generated)
•	Full Name
•	National ID / Passport
•	Date of Birth
•	Gender
•	Marital Status
•	Occupation
•	Employer (if applicable)
Contact Details
•	Mobile Number (OTP verified)
•	Email Address
•	Present Address
•	Permanent Address
•	Emergency Contact Name
•	Emergency Contact Number
KYC & Compliance
•	ID Type
•	ID Upload (Front/Back)
•	Photograph / Selfie
•	Proof of Address
•	Consent & Privacy Acceptance

Account
•	Username
•	Password / Biometric Token
•	Preferred Language
•	Notification Preference (SMS / Email / Push)
•	Wallet / Payment Method
 
✅ 2. Enrollment / Policy Purchase Fields (Generic)

Plan Selection
•	Product Type (Life / Health / Motor / Travel / Fire)
•	Plan Name / Code
•	Coverage Amount
•	Policy Term
•	Date of Enrollment
•	Date of Coverage
•	Date of Expiry
•	Premium Amount
•	Payment Frequency

Proposer Details
•	Relationship to Insured
•	Nominee Name
•	Nominee Relation
•	Nominee DOB
•	Nominee Share %
Underwriting (Dynamic)
•	Health Declaration
•	Occupation Risk Class
•	Existing Policies
•	Claims History
•	Lifestyle Questions (Smoking / Alcohol)
Payment
•	Premium Amount
•	VAT / Tax
•	Service Fee
•	Total Payable
•	Payment Gateway Reference
•	Receipt Number
 
✅ 3. Claims Module (Unified Structure)

Claim Registration
•	Policy Number
•	Claim Type (Cashless / Reimbursement)
•	Date of Loss / Event
•	Place of Incident
•	Claim Description
•	Estimated Claim Amount
Documents Upload
•	Claim Form
•	Medical / Repair Bills
•	Discharge Summary / FIR / Survey Report
•	Photos / Videos
•	Bank Details (for payout)
Tracking
•	Claim ID
•	Claim Status
•	Approved Amount
•	Settlement Date
•	Payment Reference
Communication
•	In-App Messages
•	Claim Notes
•	Rejection / Query Reason
•	Appeal Option
 
✅ 4. Services / Customer Experience Layer

Dashboard
•	Active Policies
•	Renewal Alerts
•	Claim Status
•	Wallet Balance
Self-Service
•	Download Policy
•	Update Profile
•	Change Nominee
•	Premium Top-Up
•	Endorsements

Support
•	Chatbot / Live Chat
•	Ticket Creation
•	Call Back Request
•	FAQ
Value Added Services
•	Hospital / Garage Locator
•	Wellness Programs
•	Telemedicine
•	Roadside Assistance
 
✅ 5. Product-Specific Field Extensions

●	Life Insurance
•	Sum Assured
•	Beneficiary Details
•	Medical Examination Result
•	Height / Weight / BMI
•	Income Proof
•	Fund Value
●	Health Insurance
•	Family Members Covered
•	Good Health Declaration
•	Room Rent Limit
•	Network Hospital ID
•	Co-pay % /Deductible Amount
•	In-Patient Coverage
•	OPD Coverage

●	Motor Insurance
•	Vehicle Registration Number
•	Chassis Number
•	Engine Number
•	Make / Model / Year
•	CC / Seating Capacity
•	Driver License
●	Travel Insurance
•	Passport Number
•	Destination Country
•	Travel Start Date
•	Return Date
•	Visa Type
•	Airline / Ticket Number
•	Baggage Coverage
•	Emergency Assistance ID
●	Fire / Property Insurance
•	Property Type
•	Construction Type
•	Location Geo-Tag
•	Floor Area (Sqft)
•	Sum Insured
•	Occupancy Type
•	Asset List
•	Fire Protection Measures
