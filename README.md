# FlexLiving Reviews Backend

A Go-based REST API for managing property reviews from multiple sources (Hostaway, Google Reviews).

## Tech Stack

- **Language**: Go 1.21+
- **Router**: Gorilla Mux
- **CORS**: rs/cors
- **Configuration**: godotenv
- **Architecture**: Service-oriented with clear separation of concerns

## Project Structure

```
backend/
├── cmd/server/main.go              # Application entry point
├── internal/
│   ├── config/config.go            # Configuration management
│   ├── models/
│   │   ├── hostaway.go             # Hostaway API models
│   │   └── normalized.go           # Normalized review models
│   ├── services/
│   │   ├── hostaway_service.go     # Hostaway API integration
│   │   ├── normalization_service.go # Review normalization
│   │   └── approval_service.go     # Review approval workflow
│   └── handlers/
│       └── reviews_handler.go      # HTTP request handlers
├── mockdata/
│   └── hostaway_mock.json          # Mock review data
├── go.mod
├── go.sum
└── .env                            # Environment variables
```

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. **Clone the repository** (if not already done):
```bash
git clone <repository-url>
cd flexliving-reviews/backend
```

2. **Install dependencies**:
```bash
go mod download
```

3. **Configure environment variables**:
Create a `.env` file in the `backend/` directory:
```bash
HOSTAWAY_API_KEY=f94377ebbbb479490bb3ec364649168dc443dda2e4830facaf5de2e74ccc9152
HOSTAWAY_ACCOUNT_ID=61148
HOSTAWAY_BASE_URL=https://api.hostaway.com/v1
SERVER_PORT=8080
USE_MOCK_DATA=false
```

4. **Run the server**:
```bash
# From backend directory
go run cmd/server/main.go
```

Or build and run:
```bash
go build -o server cmd/server/main.go
./server
```

The API will start on `http://localhost:8080`

### Using Mock Data

To test without hitting the real Hostaway API:
```bash
# Set in .env
USE_MOCK_DATA=true
```

Or set environment variable:
```bash
USE_MOCK_DATA=true go run cmd/server/main.go
```

## API Endpoints

### Core Endpoints

#### 1. Get Hostaway Reviews (Raw)
```http
GET /api/reviews/hostaway
```

Returns raw reviews from Hostaway API in their original format.

**Response Example**:
```json
{
  "status": "success",
  "count": 5,
  "data": [...]
}
```

#### 2. Get Normalized Reviews
```http
GET /api/reviews/normalized
GET /api/reviews/normalized?status=pending
```

Returns reviews in standardized format. Supports filtering by status.

**Query Parameters**:
- `status` (optional): Filter by status (`pending`, `approved`, `rejected`)

**Response Example**:
```json
{
  "status": "success",
  "count": 5,
  "data": [
    {
      "id": "hostaway-7453",
      "source": "hostaway",
      "propertyId": "prop-12",
      "propertyName": "2B N1 A - 29 Shoreditch Heights",
      "guestName": "Shane Finkelstein",
      "rating": 10,
      "reviewText": "Shane and family are wonderful!",
      "categories": {
        "cleanliness": 10,
        "communication": 10,
        "respect_house_rules": 10
      },
      "submittedAt": "2020-08-21T22:45:14Z",
      "status": "approved",
      "approvalStatus": {
        "isApproved": true,
        "isRejected": false
      }
    }
  ]
}
```

#### 3. Approve Review
```http
POST /api/reviews/{id}/approve
Content-Type: application/json

{
  "approvedBy": "admin@flexliving.com"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Review approved successfully",
  "data": {...}
}
```

#### 4. Reject Review
```http
POST /api/reviews/{id}/reject
Content-Type: application/json

{
  "reason": "Contains inappropriate content"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Review rejected successfully",
  "data": {...}
}
```

#### 5. Get Review Statistics
```http
GET /api/reviews/stats
```

Returns aggregated statistics across all reviews.

**Response Example**:
```json
{
  "status": "success",
  "data": {
    "totalReviews": 5,
    "averageRating": 9.2,
    "ratingsBySource": {
      "hostaway": 9.2
    },
    "categoryAverages": {
      "cleanliness": 8.8,
      "communication": 9.6,
      "respect_house_rules": 8.6
    },
    "statusBreakdown": {
      "approved": 3,
      "pending": 1,
      "rejected": 1
    },
    "recentReviews": 2
  }
}
```

#### 6. Health Check
```http
GET /health
```

Returns API health status.

## Key Design Decisions

### 1. Service Layer Architecture
- **Separation of Concerns**: Clear separation between HTTP handlers, business logic, and data access
- **Testability**: Services can be tested independently
- **Maintainability**: Easy to modify or extend individual services

### 2. Review Normalization
- Converts different review formats (Hostaway, Google, etc.) into a unified schema
- Calculates overall ratings from category ratings when not provided
- Standardizes date formats and status values
- Makes it easy to add new review sources

### 3. In-Memory Storage
- Uses in-memory maps with mutex locks for thread-safe operations
- Suitable for MVP and demonstration purposes
- Easy to replace with database in production (PostgreSQL, MongoDB, etc.)
- Fast read/write operations

### 4. API Response Format
- Consistent response structure across all endpoints
- Always includes `status` field (`success` or `error`)
- Includes metadata like `count` for list responses
- Detailed error messages for debugging

### 5. CORS Configuration
- Allows frontend development on different ports
- Configured for common development ports (5173 for Vite, 3000 for create-react-app)
- Easy to modify for production domains

## API Behaviors

### Review Status Flow
1. **Pending**: Initial state for new reviews requiring approval
2. **Approved**: Reviews marked as approved by admin
3. **Rejected**: Reviews rejected with reason

### Rating Calculation
- If Hostaway provides overall rating, use it
- Otherwise, calculate average from category ratings
- Normalized to 0-10 scale

### Error Handling
- Returns appropriate HTTP status codes
- Includes descriptive error messages
- Logs errors for debugging

### Concurrency
- Thread-safe operations using sync.RWMutex
- Multiple requests can read simultaneously
- Write operations are serialized

## Testing the API

### Using curl

```bash
# Get all reviews
curl http://localhost:8080/api/reviews/normalized

# Get pending reviews
curl "http://localhost:8080/api/reviews/normalized?status=pending"

# Approve a review
curl -X POST http://localhost:8080/api/reviews/hostaway-7453/approve \
  -H "Content-Type: application/json" \
  -d '{"approvedBy": "admin"}'

# Get statistics
curl http://localhost:8080/api/reviews/stats

# Health check
curl http://localhost:8080/health
```

### Using Postman/Insomnia
Import the following collection for quick testing (create a new collection with these endpoints).

## Production Considerations

For production deployment, consider:

1. **Database**: Replace in-memory storage with PostgreSQL/MongoDB
2. **Authentication**: Add JWT or OAuth for API security
3. **Rate Limiting**: Implement rate limiting for API endpoints
4. **Logging**: Add structured logging (e.g., zerolog, logrus)
5. **Monitoring**: Add metrics and health checks
6. **Caching**: Implement Redis for frequently accessed data
7. **Pagination**: Add pagination for large result sets
8. **Validation**: Add input validation middleware
9. **Environment**: Use proper secrets management (AWS Secrets Manager, Vault)
10. **Docker**: Containerize for consistent deployments

## Development

### Adding a New Review Source

1. Create model in `internal/models/`
2. Create service in `internal/services/`
3. Add normalization logic to `normalization_service.go`
4. Update handlers to expose new endpoints

### Running Tests

```bash
go test ./...
```

## Troubleshooting

### Port Already in Use
```bash
# Change port in .env
SERVER_PORT=8081
```

### Mock Data Not Loading
```bash
# Ensure mockdata directory exists and contains hostaway_mock.json
ls -la mockdata/
```

### CORS Issues
Update allowed origins in `cmd/server/main.go` if frontend runs on different port.

## License

MIT