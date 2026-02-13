# PDF Render Payload Documentation

This document describes the JSON payload structure used for rendering PDF documents for both **Personal** and **Company** application types.

---

## Table of Contents

1. [Overview](#overview)
2. [Personal Application Payload](#personal-application-payload)
3. [Company Application Payload](#company-application-payload)
4. [Shared Objects Reference](#shared-objects-reference)

---

## Overview

The PDF Renderer API accepts JSON payloads to generate application documents. There are two main types:

| Type | Template | Description |
|------|----------|-------------|
| Personal | `data-application-personal-v2.html` | For individual/personal loan applications |
| Company | `data-application-company-v18.html` | For company/corporate loan applications |

---

## Personal Application Payload

### Root Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `lead_id` | string | Yes | Unique lead identifier (e.g., `LD20260204820130000543`) |
| `npwp_number` | string | Yes | NPWP number |
| `npwp_type_name` | string | Yes | NPWP type name (e.g., `SP-NPWP`) |
| `created_at` | string | Yes | Application creation date (e.g., `10 Februari 2026`) |
| `status` | string | Yes | Application status (e.g., `Menunggu Persetujuan`) |
| `personal` | object | Yes | Personal information object |
| `assets` | array | Yes | Array of asset objects |
| `owners` | array | Yes | Array of owner objects |
| `showrooms` | array | Yes | Array of showroom objects |
| `banks` | array | Yes | Array of bank objects |
| `bank_summary` | array | Yes | Array of bank summary objects |
| `grand_total_average` | object | Yes | Grand total average object |
| `plafond` | object | Yes | Plafond/credit limit information |
| `background_showroom` | object | Yes | Showroom background information |
| `operational_showroom` | object | Yes | Showroom operational information |
| `emergency_contacts` | array | Yes | Array of emergency contact objects |
| `billing_address` | object | Yes | Billing address information |
| `legal_data` | object | Yes | Legal data and documents |
| `financial_documents` | array | Yes | Array of financial document objects |
| `sales_documents` | array | Yes | Array of sales document objects |
| `additional_documents` | array | Yes | Array of additional document objects |

### Personal Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Full name |
| `nik` | string | Yes | NIK (National ID number) |
| `birth_date` | string | Yes | Birth date (format: `YYYY-MM-DD`) |
| `birth_place` | string | Yes | Birth place |
| `gender` | string | Yes | Gender (`M` or `F`) |
| `phone_number` | string | Yes | Primary phone number |
| `additional_phone_number` | string | No | Secondary phone number |
| `email` | string | Yes | Email address |
| `marital_status_name` | string | Yes | Marital status (e.g., `Belum Kawin`, `Kawin`) |
| `occupation_name` | string | Yes | Occupation name |
| `economy_sector_name` | string | Yes | Economy sector |
| `industry_type_name` | string | Yes | Industry type |
| `ktp_document_url` | string | Yes | URL to KTP document image |
| `legal_address` | string | Yes | Legal/KTP address |
| `domicile_address` | string | Yes | Current domicile address |
| `domicile_address_detail` | string | No | Additional domicile address details |
| `company_name` | string | No | Company name (if employed) |
| `company_address` | string | No | Company address |
| `company_address_detail` | string | No | Additional company address details |
| `religion_name` | string | Yes | Religion |
| `education_name` | string | Yes | Education level |
| `mother_name` | string | Yes | Mother's maiden name |
| `home_status_name` | string | Yes | Home ownership status |
| `home_location_code` | string | Yes | Home location code |
| `home_location_name` | string | Yes | Home location name |
| `stay_since_year` | number | Yes | Year started living at current address |
| `home_price` | string | Yes | Home price (formatted currency) |
| `spouse` | object | No | Spouse information (null if not married) |
| `negative_list` | object | Yes | Negative list check result |
| `pefindo` | object | Yes | PEFINDO check result |

### Personal Spouse Object (if applicable)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Spouse full name |
| `id_number` | string | Yes | Spouse NIK |
| `birth_place` | string | Yes | Spouse birth place |
| `birth_date` | string | Yes | Spouse birth date |
| `phone_number` | string | Yes | Spouse phone number |
| `negative_list` | object | Yes | Negative list check result |
| `pefindo` | object | Yes | PEFINDO check result |

### Personal Legal Data Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `capital_source` | string | Yes | Capital source information |
| `family_card_document_url` | string | Yes | URL to family card document |
| `property_ownership_document_url` | string | Yes | URL to property ownership document |
| `legal_business_license_document_url` | string | Yes | URL to business license document |

---

## Company Application Payload

### Root Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `lead_id` | string | Yes | Unique lead identifier |
| `company` | object | Yes | Company information object |
| `npwp_number` | string | Yes | Company NPWP number |
| `npwp_type_name` | string | Yes | NPWP type name |
| `created_at` | string | Yes | Application creation date |
| `status` | string | Yes | Application status |
| `assets` | array | Yes | Array of asset objects |
| `owners` | array | Yes | Array of owner objects |
| `showrooms` | array | Yes | Array of showroom objects |
| `banks` | array | Yes | Array of bank objects |
| `bank_summary` | array | Yes | Array of bank summary objects |
| `grand_total_average` | object | Yes | Grand total average object |
| `plafond` | object | Yes | Plafond/credit limit information |
| `background_showroom` | object | Yes | Showroom background information |
| `operational_showroom` | object | Yes | Showroom operational information |
| `emergency_contacts` | array | Yes | Array of emergency contact objects |
| `billing_address` | object | Yes | Billing address information |
| `legal_data` | object | Yes | Legal data and documents |
| `financial_documents` | array | Yes | Array of financial document objects |
| `sales_documents` | array | Yes | Array of sales document objects |
| `additional_documents` | array | Yes | Array of additional document objects |

### Company Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `phone_number` | string | Yes | Company phone number |
| `company_name` | string | Yes | Company name |
| `legal_address` | string | Yes | Company legal address |
| `address_detail` | string | No | Additional address details |
| `npwp_document_url` | string | Yes | URL to company NPWP document |
| `industry_group_name` | string | Yes | Industry group (e.g., `Perseroan Terbatas`) |
| `economy_sector_name` | string | Yes | Economy sector |
| `industry_name` | string | Yes | Industry name |
| `pic` | object | Yes | Person In Charge (PIC) information |

### Company PIC Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | PIC full name |
| `email` | string | Yes | PIC email |
| `phone_number` | string | Yes | PIC phone number |
| `position_name` | string | Yes | PIC position |
| `id_number` | string | Yes | PIC NIK |
| `birth_place` | string | Yes | PIC birth place |
| `birth_date` | string | Yes | PIC birth date (formatted: `1 Februari 1987`) |
| `gender` | string | Yes | PIC gender |
| `marital_status` | string | Yes | PIC marital status |
| `negative_list` | object | Yes | Negative list check result |
| `pefindo` | object | Yes | PEFINDO check result |
| `spouse` | object | No | Spouse information (if married) |

### Company PIC Spouse Object (if applicable)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Spouse full name |
| `id_number` | string | Yes | Spouse NIK |
| `birth_place` | string | Yes | Spouse birth place |
| `birth_date` | string | Yes | Spouse birth date |
| `phone_number` | string | Yes | Spouse phone number |
| `negative_list` | object | Yes | Negative list check result |
| `pefindo` | object | Yes | PEFINDO check result |

### Company Legal Data Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `capital_source` | string | Yes | Capital source information |
| `establishment_deed_document_url` | string | Yes | URL to establishment deed document |
| `rups_approval_document_url` | string | Yes | URL to RUPS approval document |
| `minister_approval_document_url` | string | Yes | URL to minister approval document |

---

## Shared Objects Reference

The following objects are shared between Personal and Company payloads with identical structure.

### Asset Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Asset index number |
| `brand_name` | string | Yes | Vehicle brand (e.g., `DAIHATSU`) |
| `model_name` | string | Yes | Vehicle model (e.g., `AYLA`) |
| `variant_code` | string | Yes | Variant code |
| `variant_name` | string | Yes | Variant name |
| `manufacturing_year` | number | Yes | Manufacturing year |
| `showroom_otr` | string | Yes | Showroom OTR price (formatted) |
| `pl_bfi` | string | Yes | BFI price limit (formatted) |
| `chassis_number` | string | Yes | Chassis number |
| `engine_number` | string | Yes | Engine number |
| `price_comparison_average` | string | No | Average price comparison |
| `price_comparisons` | array | No | Array of price comparison objects |
| `rapindo` | object | Yes | RAPINDO check result |
| `stock_unit_photo_1` | string | Yes | URL to stock unit photo 1 |
| `stock_unit_photo_2` | string | Yes | URL to stock unit photo 2 |
| `stock_unit_photo_3` | string | Yes | URL to stock unit photo 3 |
| `stock_unit_photo_4` | string | Yes | URL to stock unit photo 4 |
| `payment_receipt_document` | string | No | URL to payment receipt document |

### Price Comparison Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `notes` | string | Yes | Notes |
| `source` | string | Yes | Price source |
| `pic_check` | string | Yes | PIC who checked |
| `price` | string | Yes | Price (formatted) |

### RAPINDO Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `result` | string | Yes | Check result |
| `checked_at` | string | Yes | Check date |

### Owner Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Owner index number |
| `name` | string | Yes | Owner full name |
| `nik` | string | Yes | Owner NIK |
| `position_name` | string | Yes | Position (e.g., `Owner`) |
| `phone_number` | string | Yes | Phone number |
| `relationship_name` | string | Yes | Relationship to applicant |
| `share_percentage` | string | Yes | Share percentage |
| `gender` | string | Yes | Gender |
| `birth_date` | string | Yes | Birth date |
| `birth_place` | string | Yes | Birth place |
| `npwp_number` | string | Yes | NPWP number |
| `marital_status_name` | string | Yes | Marital status |
| `legal_address` | string | Yes | Legal address |
| `rt` | string | Yes | RT number |
| `rw` | string | Yes | RW number |
| `postal_code` | string | Yes | Postal code |
| `address_detail` | string | No | Additional address details |
| `ktp_photo_url` | string | Yes | URL to KTP photo |
| `negative_list` | object | Yes | Negative list check result |
| `pefindo` | object | Yes | PEFINDO check result |

### Showroom Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Showroom index number |
| `name` | string | Yes | Showroom name |
| `is_main` | boolean | Yes | Is main showroom |
| `showroom_type` | string | Yes | Showroom type (e.g., `Perorangan`) |
| `is_pkd` | string | Yes | PKD status (`Yes`/`No`) |
| `employee_count` | number | Yes | Number of employees |
| `persona` | string | Yes | Persona classification |
| `business_duration` | string | Yes | Business duration |
| `category` | string | Yes | Category (e.g., `Silver`) |
| `criteria` | string | Yes | Criteria |
| `legal_address` | string | Yes | Legal address |
| `rt` | string | Yes | RT number |
| `rw` | string | Yes | RW number |
| `postal_code` | string | Yes | Postal code |
| `address_detail` | string | No | Additional address details |
| `selfie_photo_url` | string | Yes | URL to selfie photo |
| `workshop_photo_url` | string | Yes | URL to workshop photo |
| `showroom_photo_url` | string | Yes | URL to showroom photo |

### Bank Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Bank index number |
| `bank_name` | string | Yes | Bank name |
| `account_number` | string | Yes | Account number |
| `account_holder_name` | string | Yes | Account holder name |
| `account_status` | string | Yes | Account validation status |
| `judol_status` | string | No | Online gambling check status |
| `statements` | array | No | Array of bank statement objects |
| `average` | object | No | Average calculations |

### Bank Statement Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Statement index |
| `opening_balance` | string | Yes | Opening balance (formatted) |
| `year` | string | Yes | Year |
| `month` | string | Yes | Month |
| `debit_mutation` | string | Yes | Debit mutation (formatted) |
| `credit_mutation` | string | Yes | Credit mutation (formatted) |
| `closing_balance` | string | Yes | Closing balance (formatted) |
| `statement_photo_url` | string | Yes | URL to statement photo |

### Bank Average Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `debit_mutation` | string | Yes | Average debit mutation |
| `credit_mutation` | string | Yes | Average credit mutation |
| `closing_balance` | string | Yes | Average closing balance |

### Bank Summary Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `bank_name` | string | Yes | Bank name |
| `account_number` | string | Yes | Account number |
| `account_holder_name` | string | Yes | Account holder name |
| `average_debit_mutation` | string | Yes | Average debit mutation |
| `average_credit_mutation` | string | Yes | Average credit mutation |

### Grand Total Average Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `debit_mutation` | string | Yes | Total average debit mutation |
| `credit_mutation` | string | Yes | Total average credit mutation |
| `closing_balance` | string | Yes | Total average closing balance |

### Plafond Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `request_amount` | string | Yes | Requested amount (formatted) |

### Background Showroom Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `before_showroom_business` | string | Yes | Business background |
| `showroom_start` | string | Yes | How showroom started |
| `showroom_development` | string | Yes | Showroom development history |

### Operational Showroom Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `stock_purchase_source` | string | Yes | Stock purchase source |
| `quality_control` | string | Yes | Quality control method |
| `unit_maintenance` | string | Yes | Unit maintenance schedule |
| `unit_marketing` | string | Yes | Marketing method |
| `management_planning_stock` | string | Yes | Stock management plan |

### Emergency Contact Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Contact index |
| `relationship` | string | Yes | Relationship to applicant |
| `name` | string | Yes | Contact name |
| `nik` | string | Yes | Contact NIK |
| `phone_number` | string | Yes | Contact phone number |
| `legal_address` | string | Yes | Contact address |
| `rt` | string | Yes | RT number |
| `rw` | string | Yes | RW number |
| `postal_code` | string | Yes | Postal code |
| `address_detail` | string | No | Additional address details |

### Billing Address Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `legal_address` | string | Yes | Billing address |
| `rt` | string | Yes | RT number |
| `rw` | string | Yes | RW number |
| `postal_code` | string | Yes | Postal code |
| `address_detail` | string | No | Additional address details |

### Financial Document Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Document index |
| `month` | string | Yes | Month |
| `document_url` | string | Yes | URL to document |

### Sales Document Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Document index |
| `month` | string | Yes | Month |
| `document_url` | string | Yes | URL to document |

### Additional Document Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | number | Yes | Document index |
| `type` | string | Yes | Document type (e.g., `KTP`, `NPWP`) |
| `document_url` | string | Yes | URL to document |

### Negative List Check Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `checked_at` | string | Yes | Check date |
| `result` | string | Yes | Check result (`Clean` / `Tidak Clean`) |

### PEFINDO Check Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `checked_at` | string | Yes | Check date |
| `result` | string | Yes | Check result (`Clean` / `Tidak Clean`) |

---

## Key Differences Between Personal and Company Payloads

| Aspect | Personal | Company |
|--------|----------|---------|
| Main Entity | `personal` object | `company` object |
| Identity | Individual's `nik`, `name` | Company's `company_name`, `pic` |
| Legal Documents | Family card, property ownership, business license | Establishment deed, RUPS approval, minister approval |
| NPWP Document | Stored in `personal.ktp_document_url` | Stored in `company.npwp_document_url` |
| Spouse Info | Under `personal.spouse` | Under `company.pic.spouse` |

---

## Sample Payload Structure

### Personal Application (Minimal)

```json
{
  "lead_id": "LD20260204820130000543",
  "personal": {
    "name": "John Doe",
    "nik": "3173051205722306",
    "birth_date": "2000-01-01",
    "birth_place": "Jakarta",
    "gender": "M",
    "phone_number": "6281284907886",
    "email": "john@example.com",
    "marital_status_name": "Belum Kawin",
    "occupation_name": "Wiraswasta",
    "economy_sector_name": "INDUSTRI/MANUFACTURING",
    "industry_type_name": "KENDARAAN BERMOTOR",
    "ktp_document_url": "https://example.com/ktp.jpg",
    "legal_address": "Jl. Contoh No. 1",
    "domicile_address": "Jl. Contoh No. 1",
    "religion_name": "Islam",
    "education_name": "Strata 1",
    "mother_name": "Jane Doe",
    "home_status_name": "Milik Sendiri",
    "home_location_code": "K",
    "home_location_name": "Perkampungan",
    "stay_since_year": 2020,
    "home_price": "Rp 500.000.000",
    "negative_list": { "checked_at": "10 Februari 2026", "result": "Clean" },
    "pefindo": { "checked_at": "10 Februari 2026", "result": "Clean" }
  },
  "npwp_number": "029904939344856",
  "npwp_type_name": "SP-NPWP",
  "created_at": "10 Februari 2026",
  "status": "Menunggu Persetujuan",
  "assets": [],
  "owners": [],
  "showrooms": [],
  "banks": [],
  "bank_summary": [],
  "grand_total_average": { "debit_mutation": "Rp 0", "credit_mutation": "Rp 0", "closing_balance": "Rp 0" },
  "plafond": { "request_amount": "Rp 100.000.000" },
  "background_showroom": { "before_showroom_business": "", "showroom_start": "", "showroom_development": "" },
  "operational_showroom": { "stock_purchase_source": "", "quality_control": "", "unit_maintenance": "", "unit_marketing": "", "management_planning_stock": "" },
  "emergency_contacts": [],
  "billing_address": { "legal_address": "", "rt": "", "rw": "", "postal_code": "" },
  "legal_data": { "capital_source": "", "family_card_document_url": "", "property_ownership_document_url": "", "legal_business_license_document_url": "" },
  "financial_documents": [],
  "sales_documents": [],
  "additional_documents": []
}
```

### Company Application (Minimal)

```json
{
  "lead_id": "LD20260204820130000543",
  "company": {
    "phone_number": "628123123123",
    "company_name": "PT Example",
    "legal_address": "Jl. Contoh No. 1",
    "npwp_document_url": "https://example.com/npwp.jpg",
    "industry_group_name": "Perseroan Terbatas",
    "economy_sector_name": "INDUSTRI/MANUFACTURING",
    "industry_name": "KENDARAAN BERMOTOR",
    "pic": {
      "name": "John Doe",
      "email": "john@example.com",
      "phone_number": "6281284907886",
      "position_name": "Manager",
      "id_number": "3173051205722306",
      "birth_place": "Jakarta",
      "birth_date": "1 Februari 1987",
      "gender": "Laki-laki",
      "marital_status": "Belum Kawin",
      "negative_list": { "checked_at": "10 Februari 2026", "result": "Clean" },
      "pefindo": { "checked_at": "10 Februari 2026", "result": "Clean" }
    }
  },
  "npwp_number": "029904939344856",
  "npwp_type_name": "SP-NPWP",
  "created_at": "10 Februari 2026",
  "status": "Menunggu Persetujuan",
  "assets": [],
  "owners": [],
  "showrooms": [],
  "banks": [],
  "bank_summary": [],
  "grand_total_average": { "debit_mutation": "Rp 0", "credit_mutation": "Rp 0", "closing_balance": "Rp 0" },
  "plafond": { "request_amount": "Rp 100.000.000" },
  "background_showroom": { "before_showroom_business": "", "showroom_start": "", "showroom_development": "" },
  "operational_showroom": { "stock_purchase_source": "", "quality_control": "", "unit_maintenance": "", "unit_marketing": "", "management_planning_stock": "" },
  "emergency_contacts": [],
  "billing_address": { "legal_address": "", "rt": "", "rw": "", "postal_code": "" },
  "legal_data": { "capital_source": "", "establishment_deed_document_url": "", "rups_approval_document_url": "", "minister_approval_document_url": "" },
  "financial_documents": [],
  "sales_documents": [],
  "additional_documents": []
}
```

---

## Notes

1. **Currency Format**: All monetary values should be formatted as Indonesian Rupiah (e.g., `Rp 100.000.000`)
2. **Date Format**: Dates displayed in UI are formatted in Indonesian (e.g., `10 Februari 2026`), while internal dates use ISO format (`YYYY-MM-DD`)
3. **URLs**: All document URLs should be valid and accessible signed URLs
4. **Arrays**: Empty arrays `[]` should be provided if no data is available
5. **Conditional Rendering**: Spouse information sections are only rendered if `spouse` object is not null
