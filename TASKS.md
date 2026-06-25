# TASKS.md

# STOPS Development Tasks

## Milestone 0: Project Initialization

- [x] Read attached SRS
- [x] Read attached SDS
- [x] Read attached AGENTS rulebook
- [x] Initialize `AGENTS.md`
- [x] Initialize `PROJECT_CONTEXT.md`
- [x] Initialize `TASKS.md`
- [x] Initialize `CHANGELOG_AI.md`
- [x] Initialize `DECISIONS.md`
- [x] Initialize `SESSION_HANDOFF.md`
- [x] Initialize `README.md`
- [x] Developer approval to proceed beyond planning
- [x] Resolve stack conflict: SDS Next.js/Go vs AGENTS Vite/Express/Prisma

## Milestone 1: Repository Scaffold

- [x] Create approved frontend app structure
- [x] Create approved backend app structure
- [x] Add strict TypeScript configuration where applicable
- [x] Add Tailwind CSS setup
- [x] Add environment variable templates
- [x] Add local development instructions
- [x] Add linting and formatting scripts

## Authentication

- [x] Register with email or phone number
- [x] Login with JWT
- [ ] Password reset
- [x] bcrypt password hashing
- [x] Role-based access control: Guest, User, Admin
- [x] Protected private endpoints
- [x] Admin role middleware
- [x] Authentication tests
- [x] Frontend login/register screen

## Database

- [x] Select Go migration tool
- [x] Enable PostgreSQL/PostGIS
- [x] Create users table/model
- [x] Create transit stops table/model
- [x] Create routes table/model
- [x] Create fares table/model
- [x] Create trips/trip history table/model
- [x] Create traffic data table/model
- [x] Create predictions table/model
- [x] Create crowdsource reports table/model
- [x] Create admin audit log table/model
- [x] Add spatial indexes
- [x] Add foreign keys
- [x] Seed initial Addis Ababa data
- [x] Execute migrations against local database

## Interactive Map

- [x] Render OpenStreetMap map using Leaflet
- [x] Support zoom and pan
- [x] Support GPS positioning
- [x] Render transit stop markers
- [x] Render taxi stand markers
- [x] Display empty-region prompt
- [x] Responsive map layout

## Transit Stop Discovery

- [x] Use GPS coordinates
- [x] Support manually entered location
- [x] Query nearby bus stops and taxi stands
- [x] Show name, distance, walking time
- [x] Show walking directions to selected stop
- [x] Nearby search tests

## Route Planning

- [x] Source and destination input
- [ ] OSRM route integration
- [ ] Route options with walking segments and transfers
- [x] Best-route recommendation by time, congestion, and availability
- [ ] Alternative routes when congestion is detected
- [x] Route planning tests

## Fare Estimation

- [x] Estimate bus fares
- [x] Estimate taxi fares
- [ ] Account for route distance where applicable
- [x] Show unavailable fare message when fare data is missing
- [x] Fare estimation tests

## Travel Time Estimation

- [x] Estimate route travel time
- [ ] Include historical travel data
- [x] Include traffic level
- [ ] Include walking distance
- [ ] Update estimates from congestion signals
- [x] Travel time tests

## AI Prediction Engine

- [x] Bus arrival ETA heuristic
- [x] Taxi availability heuristic
- [x] Congestion level heuristic
- [x] Structured JSON prediction output
- [x] Degrade gracefully when historical data is sparse
- [x] Prediction tests

## Gemini AI Assistant

- [x] Send structured prediction JSON to Gemini
- [x] Explain route options
- [x] Compare alternatives
- [x] Answer transit questions
- [x] Avoid direct numerical prediction in Gemini
- [x] Handle Gemini failures gracefully
- [x] AI assistant tests

## Crowdsourced Transit Mapping

- [x] Authenticated stop submission
- [x] Pending submission state
- [x] Upvote/downvote reports
- [x] Auto-verify at +5 net votes
- [ ] Admin verification override
- [x] Rate limit: 5 pending submissions per 24 hours
- [x] Crowdsourcing tests

## User Profile and Personalization

- [ ] Save favorite locations: home, work, custom
- [ ] Store trip history
- [ ] Set preferred transport type: bus or taxi
- [ ] Recommend routes from preferences and history
- [ ] Profile tests

## Admin Dashboard

- [ ] Manage bus stops
- [ ] Manage taxi stands
- [ ] Manage routes
- [ ] Manage fares
- [ ] Verify or reject submissions
- [ ] View usage analytics
- [ ] View contribution logs
- [ ] Log all admin actions
- [ ] Admin tests

## Testing and Quality

- [x] Unit tests
- [x] Integration tests
- [ ] System tests
- [ ] Performance tests for latency targets
- [x] Security tests for JWT, input validation, rate limiting
- [x] Manual responsive UI verification

## Deployment

- [x] Docker/Docker Compose local setup if approved
- [x] Frontend deployment setup
- [x] Backend deployment setup
- [ ] Database deployment setup
- [x] Environment variable documentation
- [ ] Production HTTPS configuration
- [ ] Actual provider deployment
- [ ] Production deployment smoke test
