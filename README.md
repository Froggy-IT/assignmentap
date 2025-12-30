API endpoints:
POST   /data
GET    /data
GET    /data/{key}
DELETE /data/{key}
GET    /stats

Features:
- Thread-safe in-memory storage
- Background worker printing stats every 5 seconds
- Graceful shutdown on SIGINT/SIGTERM
- Generic Store implementation
