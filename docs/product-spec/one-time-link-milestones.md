# One-Time Link Implementation Milestones

This roadmap prioritizes working software over architectural complexity, with clear learning objectives for each phase.

## Milestone 1: Foundation and Local Development ✅

### Goal
Establish a working development environment with clear service boundaries but simple deployment.

### Completed
- [x] Monorepo structure with frontend/backend separation
- [x] React + TypeScript frontend shell with Vite
- [x] Go backend with internal service boundaries
- [x] Docker Compose for local Redis
- [x] Health endpoint and basic HTTP routing
- [x] Structured logging and CORS middleware
- [x] API contract documentation

### Key Learning Outcomes
- Go project structure and HTTP server basics
- React + TypeScript setup with Vite
- Docker Compose for local development
- API contract design and documentation

## Milestone 2: Client-Side Encryption and Secret Creation

### Goal
Implement the complete secret creation flow with proper client-side encryption.

### Tasks
- [ ] **Frontend crypto implementation**
  - Generate AES-GCM keys using Web Crypto API
  - Encrypt plaintext with random nonces
  - Handle Base64 encoding for transmission
  - Build create secret form with TTL selection

- [ ] **Backend secret storage**
  - Implement `POST /api/secrets` endpoint
  - Validate ciphertext, nonce, and TTL parameters
  - Store encrypted data in Redis with TTL
  - Return secret ID for URL generation

- [ ] **URL generation and fragment handling**
  - Generate shareable URLs with fragment keys
  - Ensure fragment keys never reach server
  - Test URL copying and sharing functionality

### Success Criteria
- User can create a secret and receive a shareable link
- Ciphertext is stored in Redis with proper TTL
- Fragment key remains client-side only
- No plaintext ever reaches the server

### Key Learning Outcomes
- Web Crypto API usage and best practices
- Redis TTL and key management
- URL fragment handling in browsers
- Client-server security boundaries

## Milestone 3: Reveal Gate and Status Checking

### Goal
Build the reveal page with preview bot protection and clear status communication.

### Tasks
- [ ] **Status endpoint implementation**
  - Implement `GET /api/secrets/{id}/status`
  - Return pending/not_found states without revealing content
  - Handle expired secrets (Redis TTL cleanup)

- [ ] **Reveal gate UI**
  - Build reveal page that loads without consuming secret
  - Extract fragment key from URL safely
  - Display clear "Click to reveal" interaction gate
  - Handle all error states with user-friendly messages

- [ ] **Client-side decryption**
  - Implement AES-GCM decryption with Web Crypto API
  - Handle decryption failures gracefully
  - Display decrypted secret securely (no logging)

### Success Criteria
- Opening a link shows reveal gate, doesn't consume secret
- Status endpoint correctly identifies pending/expired/not_found
- Preview bots cannot accidentally consume secrets
- Invalid fragment keys show appropriate error messages

### Key Learning Outcomes
- React routing and page state management
- URL parsing and fragment extraction
- Web Crypto API decryption
- User experience design for security features

## Milestone 4: Atomic Consumption and Race Condition Prevention

### Goal
Implement exactly-once reveal semantics with proper concurrency handling.

### Tasks
- [ ] **Atomic consumption endpoint**
  - Implement `POST /api/secrets/{id}/consume`
  - Use Redis GETDEL for atomic read-and-delete
  - Handle concurrent requests correctly
  - Return appropriate error codes for different scenarios

- [ ] **Rate limiting implementation**
  - Add IP-based rate limiting for create and consume endpoints
  - Use Redis counters with TTL for rate limit tracking
  - Return proper HTTP 429 responses with retry-after headers

- [ ] **Comprehensive error handling**
  - Implement all error states: already_used, expired, not_found
  - Add request logging without sensitive data
  - Test concurrent consumption scenarios

### Success Criteria
- Only first reveal attempt succeeds, subsequent attempts fail safely
- Rate limiting prevents abuse without blocking legitimate users
- All error states return clear, actionable messages
- System handles high concurrency correctly

### Key Learning Outcomes
- Redis atomic operations and concurrency
- HTTP status codes and error handling
- Rate limiting strategies and implementation
- Concurrent programming concepts

## Milestone 5: Production Deployment and Security Hardening

### Goal
Deploy the application to production with proper security measures and monitoring.

### Tasks
- [ ] **Production environment setup**
  - Configure Vietnamese VPS with Go binary, Redis, and Caddy
  - Set up Oracle Cloud standby VPS
  - Configure DNS records for `secret.quorix.io.vn` and `api.secret.quorix.io.vn`

- [ ] **Security implementation**
  - Add HTTPS-only configuration with security headers
  - Implement proper CORS for production domains
  - Add input validation and request size limits
  - Configure Redis authentication and localhost binding

- [ ] **Frontend production deployment**
  - Deploy React app to Vercel with production API URL
  - Configure custom domain and SSL
  - Test full production flow

- [ ] **Monitoring and logging**
  - Add structured logging with request IDs
  - Implement health checks and basic monitoring
  - Set up log rotation and basic alerting

### Success Criteria
- Application accessible at `https://secret.quorix.io.vn`
- All security headers present and HTTPS enforced
- Create → reveal → already_used flow works in production
- Monitoring shows healthy system status

### Key Learning Outcomes
- Production deployment and server management
- HTTPS configuration and security headers
- DNS management and domain configuration
- Basic monitoring and operational practices

## Milestone 6: Operational Excellence and Failover

### Goal
Implement failover capabilities and operational procedures for reliability.

### Tasks
- [ ] **Failover implementation**
  - Document manual failover procedures
  - Test DNS switching between primary and standby
  - Create activation scripts for standby VPS
  - Verify failover completes within acceptable timeframe

- [ ] **Operational procedures**
  - Create deployment runbooks and checklists
  - Implement backup procedures for configuration
  - Set up basic monitoring and alerting
  - Document troubleshooting procedures

- [ ] **Performance optimization**
  - Add response caching where appropriate
  - Optimize bundle size and loading performance
  - Monitor and tune Redis performance
  - Test system under realistic load

### Success Criteria
- Failover from primary to standby completes within 15 minutes
- All operational procedures documented and tested
- System performs well under expected load
- Recovery procedures verified through testing

### Key Learning Outcomes
- Disaster recovery and failover planning
- Operational documentation and procedures
- Performance monitoring and optimization
- System reliability engineering basics

## Milestone 7: Future Enhancements (Optional)

### Goal
Add advanced features only after core functionality is stable and proven.

### Potential Enhancements
- [ ] **Advanced security features**
  - Custom passphrase protection
  - CAPTCHA integration for suspicious activity
  - Advanced bot detection heuristics

- [ ] **User experience improvements**
  - Email delivery integration
  - Custom expiration times
  - Secret preview without consumption

- [ ] **Operational enhancements**
  - Admin dashboard for monitoring
  - Advanced analytics and metrics
  - Automated backup and recovery

- [ ] **Scalability features**
  - Multi-region deployment
  - Database replication
  - CDN integration for global performance

### Key Learning Outcomes
- Feature prioritization and product management
- Advanced security implementation
- Scalability planning and implementation
- Product analytics and user behavior analysis

## Implementation Strategy

### Development Approach
1. **Start simple**: Begin with monolithic deployment for faster iteration
2. **Maintain boundaries**: Keep service logic separated even in single binary
3. **Test thoroughly**: Each milestone should have comprehensive tests
4. **Document everything**: Maintain clear documentation for future reference

### Learning Priorities
1. **Security first**: Understand crypto, HTTPS, and secure coding practices
2. **Reliability**: Learn about race conditions, atomic operations, and error handling
3. **Operations**: Gain experience with deployment, monitoring, and troubleshooting
4. **Performance**: Understand caching, optimization, and scalability

### Success Metrics
- **Functionality**: All user journeys work correctly
- **Security**: No plaintext leakage, proper encryption, secure deployment
- **Reliability**: System handles failures gracefully, failover works
- **Performance**: Fast response times, handles expected load
- **Maintainability**: Code is clean, documented, and testable

## Recommended Learning Path

### Phase 1: Core Development (Milestones 1-4)
Focus on building working software with proper security foundations.

**Key Skills:**
- Go HTTP servers and middleware
- React + TypeScript development
- Web Crypto API usage
- Redis operations and atomic commands
- Testing strategies for concurrent systems

### Phase 2: Production Readiness (Milestones 5-6)
Learn operational skills and production deployment.

**Key Skills:**
- Linux server administration
- HTTPS and security configuration
- DNS management and failover
- Monitoring and logging
- Incident response procedures

### Phase 3: Enhancement and Scale (Milestone 7)
Add advanced features and prepare for growth.

**Key Skills:**
- Feature prioritization
- Performance optimization
- Advanced security measures
- Scalability planning
- Product analytics

This milestone structure balances learning objectives with practical implementation, ensuring each phase builds valuable skills while delivering working software.
