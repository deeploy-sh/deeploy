# HTTPS/SSL Implementation Plan für Deeploy

## Übersicht

Sichere Verbindungen für Deeploy mit automatischem SSL via Let's Encrypt.

**Kernprinzip:** Zero-Config SSL - Alles automatisch, keine User-Entscheidungen.

**Status:** ✅ Implementiert (2024-12)

---

## Features & Buzzwords

- **Zero-Config SSL** - HTTPS automatisch für alle Domains
- **Auto-Provisioned Certificates** - Let's Encrypt Integration out-of-the-box
- **Secure by Default** - Keine unsicheren HTTP-Verbindungen
- **Wildcard DNS via sslip.io** - Instant Domains ohne DNS-Konfiguration
- **Auto-Renewal** - Zertifikate werden automatisch erneuert
- **Self-Hosted Security** - Volle Kontrolle über eigene Infrastruktur

---

## Implementierte Änderungen

### Phase 1: Traefik SSL-Grundlage ✅

**Dateien:**
- `docker-compose.yml` - HTTPS Entrypoint, Let's Encrypt ACME, Auto-Redirect
- `docker-compose.override.yml` - Dev-Modus ohne SSL (HTTP only)
- `scripts/uninstall.sh` - Löscht auch letsencrypt Volume

**Traefik Konfiguration:**
```yaml
command:
  # ACME / LET'S ENCRYPT
  - "--certificatesresolvers.letsencrypt.acme.httpchallenge=true"
  - "--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web"
  - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"

  # ENTRYPOINTS
  - "--entrypoints.web.address=:80"
  - "--entrypoints.websecure.address=:443"

  # HTTP → HTTPS Redirect
  - "--entrypoints.web.http.redirections.entryPoint.to=websecure"
  - "--entrypoints.web.http.redirections.entryPoint.scheme=https"

volumes:
  - letsencrypt_certs:/letsencrypt
```

### Phase 2: Pod Domains SSL ✅

**Dateien:**
- `internal/server/docker/docker.go` - SSL Labels automatisch für alle Domains in Production
- `internal/server/app/app.go` - isDevelopment an DockerService übergeben
- `internal/server/handler/pod_domain.go` - SSLEnabled=true als Default
- `internal/tui/ui/pages/pod_domains.go` - SSL-Toggle entfernt, zeigt "Automatic (Let's Encrypt)"
- `internal/tui/ui/pages/pod_domains_edit.go` - SSL-Toggle entfernt

**Traefik Labels für Pods:**
```go
// Entrypoint basierend auf Umgebung
entrypoint := "websecure"  // Production: HTTPS
if d.isDevelopment {
    entrypoint = "web"      // Development: HTTP
}

// Labels pro Domain
labels["traefik.http.routers."+routerName+".rule"] = fmt.Sprintf("Host(`%s`)", domain)
labels["traefik.http.routers."+routerName+".entrypoints"] = entrypoint
if !d.isDevelopment {
    labels["traefik.http.routers."+routerName+".tls.certresolver"] = "letsencrypt"
}
```

### Phase 3: Deeploy Admin Domain ✅

**Dateien:**
- `internal/tui/ui/pages/app.go`:
  - `⚠ insecure` Warnung im Header wenn HTTP
  - "Server Connection" Palette-Eintrag
  - `secureConnection` Flag basierend auf `https://` Prefix
- `internal/tui/ui/pages/connect.go`:
  - HTTPS-Hinweis in der UI
  - Placeholder: `https://deeploy.yourdomain.com`

### Phase 4: README ✅

Features mit Buzzwords hinzugefügt:
- Zero-Config SSL
- Auto-Provisioned Certificates
- Secure by Default
- Auto-Renewal

---

## Technische Details

### Wie SSL funktioniert

```
1. User fügt Domain hinzu (auto oder custom)
        ↓
2. Server erstellt Domain mit SSLEnabled=true
        ↓
3. Container wird mit Traefik Labels gestartet:
   - traefik.http.routers.{pod}.rule=Host(`domain`)
   - traefik.http.routers.{pod}.entrypoints=websecure
   - traefik.http.routers.{pod}.tls.certresolver=letsencrypt
        ↓
4. Traefik erkennt neue Domain
        ↓
5. Let's Encrypt HTTP-Challenge:
   - Let's Encrypt besucht http://domain/.well-known/acme-challenge/xxx
   - Traefik antwortet mit Secret Token
   - Let's Encrypt verifiziert und stellt Zertifikat aus
        ↓
6. ✅ HTTPS funktioniert automatisch!
```

### sslip.io für Auto-Generated Domains

```
Format: {subdomain}.{ip}.sslip.io
Beispiel: myapp.192-168-1-1.sslip.io → 192.168.1.1

Vorteile:
- Keine DNS-Konfiguration nötig
- IP ist in Domain eingebettet
- HTTP-Challenge funktioniert (kein DNS-API Token nötig)
- Sofort verfügbar nach Pod-Erstellung
```

### Development vs Production

| Aspekt | Development | Production |
|--------|-------------|------------|
| Entrypoint | `web` (Port 80) | `websecure` (Port 443) |
| SSL | Deaktiviert | Automatisch via Let's Encrypt |
| HTTP→HTTPS Redirect | Nein | Ja |
| Zertifikate | Keine | Auto-provisioned |

---

## Entscheidungen

- **HTTP Handling:** Redirect zu HTTPS (Standard-Verhalten)
- **SSL Toggle:** Entfernt - SSL ist immer aktiv in Production
- **sslip.io:** Für Auto-Generated Domains (kein DNS-Setup nötig)
- **Let's Encrypt:** HTTP-Challenge (kein DNS-API Token nötig)
- **Rate Limits:** Zertifikate in Volume persistiert (max 50/Woche/Domain)

---

## Zum Testen auf VPS

1. **DNS A-Record erstellen:**
   ```
   deeploy.yourdomain.com → VPS-IP
   ```

2. **Server neu starten:**
   ```bash
   sudo docker compose pull
   sudo docker compose up -d
   ```

3. **TUI verbinden:**
   ```
   https://deeploy.yourdomain.com
   ```

4. **Pod deployen:**
   - Auto-Domain: `myapp.192-168-1-1.sslip.io` (sofort HTTPS)
   - Custom Domain: DNS A-Record setzen, dann Domain hinzufügen

---

## Geänderte Dateien (Vollständige Liste)

```
docker-compose.yml
docker-compose.override.yml
scripts/uninstall.sh
internal/server/app/app.go
internal/server/docker/docker.go
internal/server/handler/pod_domain.go
internal/tui/ui/pages/app.go
internal/tui/ui/pages/connect.go
internal/tui/ui/pages/pod_domains.go
internal/tui/ui/pages/pod_domains_edit.go
README.md
```
