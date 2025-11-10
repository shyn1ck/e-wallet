#!/bin/bash

# E-Wallet Management Script
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Database connection
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-e_wallet_db}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}

# Main menu
case "${1}" in
    up)
        print_header "STARTING SERVICES"
        
        print_info "Starting Docker containers..."
        docker-compose up -d
        
        print_info "Waiting for PostgreSQL..."
        sleep 5
        
        print_info "Waiting for Redis..."
        sleep 2
        
        print_success "Services started!"
        docker-compose ps
        ;;
    
    down)
        print_header "STOPPING SERVICES"
        docker-compose down
        print_success "Services stopped!"
        ;;
    
    restart)
        print_header "RESTARTING SERVICES"
        docker-compose restart
        print_success "Services restarted!"
        docker-compose ps
        ;;
    
    logs)
        SERVICE=${2:-}
        if [ -z "$SERVICE" ]; then
            print_info "Showing logs for all services..."
            docker-compose logs -f --tail=100
        else
            print_info "Showing logs for $SERVICE..."
            docker-compose logs -f --tail=100 "$SERVICE"
        fi
        ;;
    
    status)
        print_header "SERVICE STATUS"
        docker-compose ps
        ;;
    
    init)
        print_header "INITIALIZING DATABASE"
        
        print_info "Creating database..."
        PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -f scripts/init.sql
        
        print_success "Database initialized!"
        ;;
    
    seed)
        print_header "SEEDING DATABASE"
        
        print_info "Truncating existing data..."
        PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
            -c "TRUNCATE api_clients, wallets, transactions RESTART IDENTITY CASCADE;" 2>/dev/null || true
        
        print_info "Inserting seed data..."
        PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f scripts/seed.sql
        
        print_success "Database seeded successfully!"
        ;;
    
    build)
        print_header "BUILDING APPLICATION"
        
        print_info "Building binary..."
        go build -o bin/e-wallet cmd/server/main.go
        
        print_success "Build completed: bin/e-wallet"
        ;;
    
    run)
        print_header "RUNNING APPLICATION"
        
        if [ ! -f bin/e-wallet ]; then
            print_info "Binary not found, building..."
            go build -o bin/e-wallet cmd/server/main.go
        fi
        
        print_info "Starting application..."
        ./bin/e-wallet
        ;;
    
    dev)
        print_header "DEVELOPMENT MODE"
        
        print_info "Starting Docker containers (PostgreSQL + Redis)..."
        docker-compose up -d postgres redis
        
        print_info "Waiting for PostgreSQL..."
        sleep 5
        
        print_info "Initializing database..."
        PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -f scripts/init.sql 2>/dev/null || true
        
        print_info "Seeding database..."
        PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
            -c "TRUNCATE api_clients, wallets, transactions RESTART IDENTITY CASCADE;" 2>/dev/null || true
        PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f scripts/seed.sql
        
        print_success "Environment ready!"
        print_info "Starting application..."
        go run cmd/server/main.go
        ;;
    
    prod)
        print_header "PRODUCTION MODE"
        
        print_info "Building and starting all services..."
        docker-compose up -d --build
        
        print_info "Waiting for services to be healthy..."
        sleep 10
        
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            print_success "Production deployment successful!"
            print_info "API: http://localhost:8080"
            docker-compose ps
        else
            print_error "Health check failed!"
            docker-compose logs app
            exit 1
        fi
        ;;
    
    test)
        print_header "RUNNING API TESTS"
        
        # Check if app is running
        if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
            print_error "Application is not running!"
            print_info "Start the app first: ./scripts/manage.sh dev"
            exit 1
        fi
        
        print_info "Executing test suite..."
        bash scripts/test-api.sh
        ;;
    
    swagger)
        print_header "GENERATING SWAGGER DOCS"
        
        if ! command -v swag &> /dev/null; then
            print_error "swag not found!"
            print_info "Install: go install github.com/swaggo/swag/cmd/swag@latest"
            exit 1
        fi
        
        print_info "Generating documentation..."
        swag init -g cmd/server/main.go -o docs
        
        print_success "Swagger docs generated!"
        print_info "View at: http://localhost:8080/swagger/index.html"
        ;;
    
    clean)
        print_header "CLEANUP"
        
        print_info "Stopping containers..."
        docker-compose down -v
        
        print_info "Removing binary..."
        rm -f bin/e-wallet
        
        print_info "Removing logs..."
        rm -rf logs/*
        
        print_success "Cleanup completed!"
        ;;
    
    backup)
        print_header "DATABASE BACKUP"
        
        BACKUP_DIR="./backups"
        mkdir -p $BACKUP_DIR
        BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"
        
        print_info "Creating backup..."
        PGPASSWORD=$POSTGRES_PASSWORD pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"
        
        if [ -f "$BACKUP_FILE" ]; then
            gzip "$BACKUP_FILE"
            print_success "Backup created: ${BACKUP_FILE}.gz"
            
            # Keep only last 7 backups
            # shellcheck disable=SC2012
            ls -t $BACKUP_DIR/backup_*.sql.gz 2>/dev/null | tail -n +8 | xargs -r rm
            print_info "Old backups cleaned up"
        else
            print_error "Backup failed!"
            exit 1
        fi
        ;;
    
    restore)
        if [ -z "$2" ]; then
            print_error "Please specify backup file"
            echo "Usage: $0 restore <backup_file>"
            exit 1
        fi
        
        print_header "DATABASE RESTORE"
        print_info "Restoring from: $2"
        
        if [[ $2 == *.gz ]]; then
            gunzip -c "$2" | PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME"
        else
            PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" < "$2"
        fi
        
        print_success "Restore completed!"
        ;;
    
    hmac)
        print_header "HMAC GENERATOR"
        go run tools/hmac-gen/main.go
        ;;
    
    *)
        echo "E-Wallet Management Script"
        echo ""
        echo "Usage: $0 {command}"
        echo ""
        echo "Docker Commands:"
        echo "  up        - Start Docker containers (PostgreSQL + Redis)"
        echo "  down      - Stop Docker containers"
        echo "  restart   - Restart Docker containers"
        echo "  logs      - View logs (optional: service name)"
        echo "  status    - Show service status"
        echo ""
        echo "Application Commands:"
        echo "  dev       - Development mode (docker DBs + local app)"
        echo "  prod      - Production mode (docker all services)"
        echo "  build     - Build application binary"
        echo "  run       - Run application from binary"
        echo "  test      - Run API tests"
        echo "  hmac      - Run HMAC generator tool"
        echo ""
        echo "Database Commands:"
        echo "  init      - Initialize database (create DB + extensions)"
        echo "  seed      - Seed database with test data"
        echo "  backup    - Create database backup"
        echo "  restore   - Restore database from backup"
        echo ""
        echo "Utility Commands:"
        echo "  swagger   - Generate Swagger documentation"
        echo "  clean     - Clean up containers, volumes, and binaries"
        echo ""
        echo "Quick Start:"
        echo "  $0 dev                    # Start everything (recommended)"
        echo ""
        echo "Examples:"
        echo "  $0 up && $0 init && $0 seed  # Manual setup"
        echo "  $0 test                      # Run API tests"
        echo "  $0 logs postgres             # View PostgreSQL logs"
        echo "  $0 restore backup.sql.gz     # Restore database"
        exit 1
        ;;
esac
