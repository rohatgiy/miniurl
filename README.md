# miniurl
Tired of limitations on how many URLs you can shorten at [bit.ly](https://bitly.com/pages/products/url-shortener) or [tinyurl.com](https://tinyurl.com)? miniurl is a self-hosted, no limit URL shortener that runs on your local machine!
miniurl supports shortening up to 62^10 or about 830 quadrillion non-expiring URLs (that's a lot).

![image](https://github.com/user-attachments/assets/449b03b5-99cb-42ef-bf0d-3522634cbf66)

### Recommended Specs
- **CPU**: >= 1 core
- **Memory**: >= 4GB RAM
- **Storage**: >= 10 GB
- **Network**: Reliable network connection

## How to use

### Prerequisites
You must have [Docker and Docker Compose]((https://docs.docker.com/engine/install/)) installed on your machine.

### Procedure:
1. Clone source code
   ```
   git clone https://github.com/rohatgiy/miniurl.git
   ```
2. Run the following command to spin up the Docker containers
   ```
   make dc_up
   ```
3. Visit [localhost](http://localhost/) and start shortening!

## Technologies used
- Go: REST API written in Go, using Gin server framework
- PostgreSQL: Underlying database for storing shortened URLs
- Redis: Cache for frequently accessed URLs (default 2GB max cache size)
- Nginx: Load balancer to split traffic between server instances
- HTMX: Lightweight FE to shorten URLs

