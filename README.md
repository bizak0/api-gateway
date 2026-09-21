# API Gateway — Go Microservices Architecture

Architecture de microservices en Go composée de trois briques indépendantes — un **reverse proxy TLS**, un **API Gateway**, et un **Load Balancer** — reliées par une chaîne de requêtes complète, avec un pipeline **DevSecOps** intégré (SAST, secret scanning, SCA, DAST).

```
Client → Reverse Proxy (:8443, TLS) → API Gateway (:8081) → Load Balancer (:8082) → Microservices
```

## Fonctionnalités

### Reverse Proxy
- Terminaison TLS (certificat auto-signé fourni pour le dev)
- Cache public des réponses

### API Gateway
- Authentification par JWT
- Autorisation par rôle (RBAC)
- Rate limiting (token bucket, par IP)
- Versioning d'API (`/v1/...`, `/v2/...`)
- Transformation de requêtes/réponses
- Adaptateur de protocole interne
- Cache privé (par utilisateur)

### Load Balancer
- Répartition de charge round-robin
- Health checks périodiques des services enregistrés
- Circuit breaker
- Retry automatique
- Registre de services (discovery) et routage par préfixe de chemin

### CI/CD — Pipeline de sécurité (GitHub Actions)
- **SAST** : analyse statique du code Go avec CodeQL
- **Secret Scanning** : détection de secrets commités avec Gitleaks
- **SCA** : scan des dépendances (`go.mod`/`go.sum`) avec Trivy
- **DAST** : scan dynamique de la chaîne complète (Proxy → Gateway → LB → mocks) avec OWASP ZAP

## Stack technique

- **Go 1.25**
- [`golang-jwt/jwt`](https://github.com/golang-jwt/jwt) — authentification JWT
- [`golang.org/x/time`](https://pkg.go.dev/golang.org/x/time) — rate limiting
- [`google/uuid`](https://github.com/google/uuid)
- Docker & Docker Compose
- Nginx (microservices simulés pour le développement local)
- GitHub Actions (CodeQL, Gitleaks, Trivy, OWASP ZAP)

## Structure du projet

```
cmd/
  gateway/          # point d'entrée de l'API Gateway
  loadbalancer/      # point d'entrée du Load Balancer
  proxy/              # point d'entrée du Reverse Proxy
internal/
  gateway/
    adaptor/          # adaptateur de protocole interne
    cache/            # cache privé
    middleware/       # auth, RBAC, rate limit, versioning, transform
  loadbalancer/
    balancer/         # round-robin
    health/            # health checks
    registry/          # service discovery
    resilience/        # circuit breaker, retry
    router/             # routage par préfixe
  proxy/
    cache/              # cache public
    tls.go               # terminaison TLS
deployments/
  Dockerfile.gateway
  Dockerfile.loadbalancer
  Dockerfile.proxy
  docker-compose.yml
  nginx/                # configs des microservices simulés
.github/workflows/
  go-security-pipeline.yml   # pipeline SAST/Secret/SCA/DAST
```

## Lancer le projet

### Avec Docker Compose (recommandé)

```bash
cd deployments
docker compose up --build
```

Cela démarre les 3 microservices simulés (nginx), le Load Balancer, le Gateway, puis le Reverse Proxy, dans cet ordre (via `depends_on` + health checks).

Le point d'entrée public est ensuite : `https://localhost:8443`

### En local, sans Docker

Trois microservices doivent tourner sur `:9091`, `:9092` et `:9093` (ou adapter les variables d'environnement ci-dessous), puis :

```bash
go run ./cmd/loadbalancer   # démarre sur :8082
go run ./cmd/gateway         # démarre sur :8081
go run ./cmd/proxy            # démarre sur :8443 (TLS)
```

Variables d'environnement disponibles :

| Variable | Défaut | Composant |
|---|---|---|
| `LOADBALANCER_URL` | `http://localhost:8082` | Gateway |
| `GATEWAY_URL` | `http://localhost:8081` | Proxy |
| `USERS_SERVICE_URL` | `http://localhost:9091` | Load Balancer |
| `ORDERS_SERVICE_URL` | `http://localhost:9092` | Load Balancer |
| `ADMIN_SERVICE_URL` | `http://localhost:9093` | Load Balancer |

## Contributeurs

- [**bizak0**](https://github.com/bizak0) — architecture initiale, Gateway (auth JWT, RBAC, rate limiting, versioning, cache), Load Balancer (round-robin, health check, circuit breaker, retry), Reverse Proxy (TLS, cache public)
- [**midoumkt02**](https://github.com/midoumkt02) — pipeline DevSecOps CI/CD complet (CodeQL, Gitleaks, Trivy, OWASP ZAP), chaînage Docker Compose des trois services et correctifs réseau associés
