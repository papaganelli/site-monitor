# 🚀 Site Monitor v0.6.0 - Guide d'Utilisation des Nouvelles Fonctionnalités

## ✨ Nouvelles Fonctionnalités

### 🔒 1. Monitoring SSL/TLS avec Alertes d'Expiration

#### Utilisation via CLI
```bash
# Vérifier tous les certificats SSL
site-monitor ssl

# Vérifier un site spécifique 
site-monitor ssl --site "Mon Site"

# Définir le seuil d'alerte (jours avant expiration)
site-monitor ssl --warn-days 14

# Sortie JSON pour intégrations
site-monitor ssl --json
```

#### Configuration
```json
{
  "sites": [
    {
      "name": "Mon Site",
      "url": "https://monsite.com",
      "ssl_check": true,
      "ssl_warn_days": 30
    }
  ],
  "alerts": {
    "thresholds": {
      "ssl_expiry_warning_days": [30, 14, 7, 1]
    }
  }
}
```

#### Exemple de Sortie
```
🔒 SSL Certificate Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ monsite.com
   🏷️  Subject: CN=monsite.com
   🏢 Issuer: Let's Encrypt Authority X3
   📅 Expires: 2024-12-15 23:59 (Valid for 287 days)
   ⚡ Response Time: 142ms
   🔗 Certificate Chain: 2 certificates

⚠️  api.monsite.com
   🏷️  Subject: CN=api.monsite.com
   🏢 Issuer: Let's Encrypt Authority X3  
   📅 Expires: 2024-09-20 14:30 (Expires in 25 days)
   ⚠️  WARNING: Certificate expires in 25 days!
   ⚡ Response Time: 89ms
```

---

### 📈 2. Métriques Avancées (P95, P99, MTTR, MTBF)

#### Utilisation via CLI
```bash
# Métriques avancées pour tous les sites (24h par défaut)
site-monitor metrics

# Site spécifique avec période personnalisée
site-monitor metrics --site "Mon Site" --since 7d

# Métriques détaillées avec tendances
site-monitor metrics --detailed --trends

# Export CSV pour analyse
site-monitor metrics --format csv --output metrics.csv
```

#### Métriques Disponibles

##### **Percentiles de Temps de Réponse**
- **P50 (médiane)** : 50% des requêtes sont plus rapides
- **P95** : 95% des requêtes sont plus rapides (SLA standard)
- **P99** : 99% des requêtes sont plus rapides (SLA premium)
- **P99.9** : 99.9% des requêtes sont plus rapides (ultra-premium)

##### **Métriques de Fiabilité**
- **MTTR** (Mean Time To Recovery) : Temps moyen de récupération après panne
- **MTBF** (Mean Time Between Failures) : Temps moyen entre pannes
- **Availability Nines** : Niveau de disponibilité (99%, 99.9%, 99.99%, etc.)

##### **Analyse des Tendances**
- **Response Time Trend** : Amélioration/dégradation des performances
- **Uptime Trend** : Évolution de la disponibilité
- **Error Pattern Analysis** : Classification automatique des erreurs

#### Exemple de Sortie Détaillée
```
📊 Advanced Metrics for Mon Site API (Last 7 days)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Mon Site API
   📈 Uptime: 99.87% (3 nines) - 1,432/1,435 successful checks
   ⚡ Response Times:
      • P50 (median): 147ms
      • P95: 342ms  
      • P99: 891ms
      • Std Dev: 123ms
   🔧 Reliability:
      • MTTR (Mean Time To Recovery): 4m 23s
      • MTBF (Mean Time Between Failures): 2d 7h
   📊 Trends:
      • Response Time: 📈 improving
      • Uptime: 📊 stable
   🎯 SLA Compliance:
      • 99.9% (8.77h downtime/month): ❌ 99.87%
      • 99.5% (3.65d downtime/month): ✅ 99.87%
   💥 Error Analysis:
      • Timeout: 2 occurrences (66.7%)
      • Network: 1 occurrence (33.3%)
   🕐 Performance Patterns:
      • Best Hour: 03:00 (100.0% uptime)
      • Worst Hour: 14:00 (98.9% uptime)
      • Best Day: Sunday
      • Worst Day: Wednesday
```

---

### 📧 3. Export par Email avec Rapports Périodiques

#### Configuration des Rapports
```json
{
  "reports": {
    "email": {
      "enabled": true,
      "schedules": [
        {
          "name": "Weekly Executive Summary",
          "schedule": "weekly",
          "day_of_week": 1,
          "hour": 9,
          "recipients": ["ceo@monsite.com"],
          "format": "html",
          "sections": ["overview", "sla_compliance", "recommendations"]
        }
      ]
    }
  }
}
```

#### Types de Rapports Disponibles

##### **1. Rapport Exécutif Hebdomadaire**
- Vue d'ensemble des performances
- Conformité SLA
- Certificats SSL
- Recommandations d'amélioration

##### **2. Rapport Opérationnel Quotidien**  
- Métriques détaillées
- Résumé des alertes
- Tendances de performance
- Analyse des erreurs

##### **3. Rapport SLA Mensuel**
- Analyse complète des SLA
- Métriques de conformité
- Export CSV inclus
- Recommandations stratégiques

#### Commandes CLI
```bash
# Envoyer un rapport de test
site-monitor report send-test

# Configurer un rapport hebdomadaire
site-monitor report schedule weekly \
  --name "Rapport Hebdomadaire" \
  --recipients "admin@monsite.com,ops@monsite.com" \
  --sections "overview,metrics,ssl"

# Lister les rapports programmés
site-monitor report list

# Générer un rapport à la demande
site-monitor report generate \
  --period "last-week" \
  --format html \
  --output rapport.html
```

#### Exemple d'Email HTML
```html
🚀 Site Monitor Weekly Report
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Executive Summary (Sep 16-22, 2024)

Overall Performance: ✅ Excellent
• 3 sites monitored
• 99.94% overall uptime  
• 2,847 total checks performed
• 2 sites achieving 99.9%+ SLA

🎯 SLA Compliance Status:
✅ Site Principal: 99.97% (Target: 99.9%)
✅ API Backend: 99.95% (Target: 99.9%)  
⚠️ Service Tiers: 98.23% (Target: 99.5%)

🔒 SSL Certificate Status:
✅ All certificates valid
⚠️ api.monsite.com expires in 23 days

💡 Key Recommendations:
1. Investigate Service Tiers performance issues
2. Renew api.monsite.com SSL certificate
3. Consider adding CDN for improved global performance
```

---

### 🎨 4. Templates d'Alertes Personnalisables

#### Gestion des Templates

##### Lister les Templates Disponibles
```bash
site-monitor template list
```

Output:
```
🎨 Available Alert Templates (8):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📄 Site Down - Email (default)
   ID: default-site-down-email
   Type: site_down → email
   Format: html
   Used: 23 times

📄 Site Down - Slack (default)
   ID: default-site-down-slack  
   Type: site_down → slack
   Format: json
   Used: 23 times

📄 Minimal Site Down - Email
   ID: minimal-site-down-email
   Type: site_down → email
   Format: html
   Used: 0 times
```

##### Tester un Template
```bash
site-monitor template test default-site-down-email
```

##### Créer un Template Personnalisé
```bash
# Exporter un template existant comme base
site-monitor template export default-site-down-email > my-template.json

# Modifier le fichier JSON
# Importer le nouveau template
site-monitor template import my-template.json
```

#### Structure d'un Template

##### Template Email HTML Personnalisé
```json
{
  "id": "corporate-site-down-email",
  "name": "Corporate Site Down Alert",
  "alert_type": "site_down",
  "channel": "email", 
  "format": "html",
  "subject": "🔴 URGENT: {{.SiteName}} Service Disruption",
  "body": "<!DOCTYPE html>...",
  "variables": {
    "SiteName": {
      "name": "SiteName",
      "type": "string", 
      "required": true
    },
    "SiteURL": {
      "name": "SiteURL", 
      "type": "url",
      "required": true
    }
  },
  "style": {
    "theme": "corporate",
    "colors": {
      "primary": "#1f2937",
      "danger": "#dc2626", 
      "success": "#059669"
    },
    "include_logo": true
  }
}
```

##### Variables Disponibles dans les Templates
```
{{.SiteName}}         - Nom du site
{{.SiteURL}}          - URL du site
{{.Message}}          - Message d'alerte
{{.Details}}          - Détails de l'alerte
{{.Timestamp}}        - Horodatage
{{.CurrentStatus}}    - Code de statut HTTP
{{.ConsecutiveFails}} - Nombre d'échecs consécutifs
{{.ResponseTime}}     - Temps de réponse
{{.ErrorMessage}}     - Message d'erreur
{{.UptimePercent}}    - Pourcentage d'uptime
{{.Severity}}         - Niveau de sévérité
```

##### Fonctions Helper Disponibles
```
{{.Timestamp | formatTime}}           - Format: 2024-09-22 14:30:15
{{.ResponseTime | formatDuration}}    - Format: 1.23s ou 456ms
{{.Timestamp | unixTime}}             - Timestamp Unix
{{.SiteName | upper}}                 - MAJUSCULES
{{.ErrorMessage | title}}             - Title Case
```

#### Templates pour Différents Canaux

##### Template Slack Avancé
```json
{
  "body": "{
    \"text\": \"🚨 *SITE DOWN ALERT*\",
    \"attachments\": [
      {
        \"color\": \"danger\",
        \"title\": \"{{.SiteName}} is experiencing issues\",
        \"title_link\": \"{{.SiteURL}}\",
        \"text\": \"{{.Details}}\",
        \"fields\": [
          {
            \"title\": \"Status Code\",
            \"value\": \"{{.CurrentStatus}}\",
            \"short\": true
          },
          {
            \"title\": \"Consecutive Failures\", 
            \"value\": \"{{.ConsecutiveFails}}\",
            \"short\": true
          }
        ],
        \"actions\": [
          {
            \"type\": \"button\",
            \"text\": \"View Dashboard\",
            \"url\": \"https://monitor.monsite.com/dashboard\"
          },
          {
            \"type\": \"button\", 
            \"text\": \"Check Site\",
            \"url\": \"{{.SiteURL}}\"
          }
        ],
        \"footer\": \"Site Monitor\",
        \"ts\": {{.Timestamp | unixTime}}
      }
    ]
  }"
}
```

##### Template Discord avec Embed
```json
{
  "body": "{
    \"content\": \"🚨 **CRITICAL ALERT**\",
    \"embeds\": [
      {
        \"title\": \"{{.SiteName}} Status Alert\",
        \"description\": \"{{.Message}}\",
        \"color\": 15158332,
        \"fields\": [
          {
            \"name\": \"🌐 Site\",
            \"value\": \"[{{.SiteName}}]({{.SiteURL}})\",
            \"inline\": true
          },
          {
            \"name\": \"📊 Status\", 
            \"value\": \"{{.CurrentStatus}}\",
            \"inline\": true
          },
          {
            \"name\": \"⏱️ Time\",
            \"value\": \"{{.Timestamp | formatTime}}\",
            \"inline\": true
          }
        ],
        \"footer\": {
          \"text\": \"Site Monitor Alert\"
        },
        \"timestamp\": \"{{.Timestamp | date \\\"2006-01-02T15:04:05Z07:00\\\"}}\"
      }
    ]
  }"
}
```

---

## 🛠️ Configuration Avancée

### Configuration Complète Enhanced
```json
{
  "sites": [
    {
      "name": "Production API",
      "url": "https://api.monsite.com/health", 
      "interval": "30s",
      "timeout": "10s",
      "ssl_check": true,
      "ssl_warn_days": 14,
      "expected_status_codes": [200],
      "headers": {
        "Authorization": "Bearer monitor-token",
        "User-Agent": "SiteMonitor/0.6.0"
      }
    }
  ],

  "alerts": {
    "email": {
      "enabled": true,
      "smtp_server": "smtp.gmail.com:587",
      "username": "monitoring@monsite.com", 
      "password": "app-password",
      "from": "Site Monitor <monitoring@monsite.com>",
      "recipients": ["ops@monsite.com"]
    },
    
    "templates": {
      "site_down_email": "corporate-site-down-email",
      "ssl_expiry_email": "corporate-ssl-expiry-email"
    },
    
    "thresholds": {
      "ssl_expiry_warning_days": [30, 14, 7, 1]
    }
  },

  "reports": {
    "email": {
      "enabled": true,
      "schedules": [
        {
          "name": "Weekly Executive Report",
          "schedule": "weekly", 
          "day_of_week": 1,
          "hour": 9,
          "recipients": ["ceo@monsite.com"],
          "sections": ["overview", "sla_compliance", "recommendations"]
        }
      ]
    }
  },

  "metrics": {
    "advanced": {
      "enabled": true,
      "percentiles": [50, 95, 99],
      "sla_targets": {
        "99.9": "Premium SLA",
        "99.5": "Standard SLA"
      }
    }
  }
}
```

---

## 📋 Nouveaux Commandes CLI

### Commandes SSL
```bash
site-monitor ssl                    # Status SSL tous sites
site-monitor ssl --site "Mon Site" # Site spécifique
site-monitor ssl --json            # Output JSON
site-monitor ssl --warn-days 7     # Seuil d'alerte personnalisé
```

### Commandes Métriques Avancées  
```bash
site-monitor metrics                           # Métriques avancées
site-monitor metrics --site "API" --since 7d  # Site + période
site-monitor metrics --percentiles             # Détail percentiles
site-monitor metrics --trends                  # Analyse tendances
```

### Commandes Rapports
```bash
site-monitor report send-test                  # Test email
site-monitor report schedule weekly            # Programme rapport
site-monitor report list                       # Liste rapports
site-monitor report generate --period last-week # Génère à la demande
```

### Commandes Templates
```bash
site-monitor template list                     # Liste templates  
site-monitor template test <template-id>       # Test rendu
site-monitor template export <template-id>     # Export JSON
site-monitor template import template.json     # Import JSON
```

---

## 🎯 Cas d'Usage Pratiques

### 1. **Monitoring E-commerce Critique**
```bash
# Configuration pour site e-commerce avec SLA 99.95%
site-monitor metrics --site "Boutique" --since 30d
# → Vérifier P95 < 200ms, P99 < 500ms, Uptime > 99.95%

site-monitor ssl --site "Boutique" --warn-days 14  
# → Certificats e-commerce doivent être renouvelés plus fréquemment

site-monitor report schedule daily \
  --name "E-commerce Daily" \
  --recipients "ops@boutique.com" \
  --sections "overview,metrics,ssl,alerts"
```

### 2. **Monitoring API avec SLA Clients**
```bash
# Métriques pour justifier SLA clients
site-monitor metrics --format csv --output sla-report.csv

# Template personnalisé pour alertes API
cat > api-alert-template.json << EOF
{
  "name": "API Production Alert",
  "subject": "🚨 API {{.SiteName}} - SLA Impact", 
  "body": "Critical API issue affecting customer SLA..."
}
EOF

site-monitor template import api-alert-template.json
```

### 3. **Rapports Exécutifs Automatisés**
```bash
# Rapport hebdomadaire direction
site-monitor report schedule weekly \
  --name "Executive Weekly" \
  --recipients "ceo@monsite.com,cto@monsite.com" \
  --day monday --hour 9 \
  --sections "overview,sla_compliance,recommendations" \
  --format html
```

---

## 🔧 Dépannage et FAQ

### SSL Issues
**Q: "Certificate verification failed"**
```bash
# Vérifier la chaîne de certificats
site-monitor ssl --site "MonSite" --verify-chain

# Ignorer temporairement les erreurs SSL
site-monitor ssl --insecure
```

### Métriques Incomplètes
**Q: "Not enough data for percentiles"**
```bash
# Vérifier la période de données disponibles
site-monitor history --site "MonSite" --limit 10

# Réduire la période d'analyse
site-monitor metrics --since 6h  # Au lieu de 24h
```

### Rapports Email  
**Q: "Failed to send email report"**
```bash
# Tester la configuration email
site-monitor report send-test

# Vérifier les logs détaillés
site-monitor --debug report send-test
```

---

## 📊 Migration depuis v0.5.0

### Étapes de Migration

1. **Sauvegarder la configuration**
```bash
cp config.json config.json.backup
cp site-monitor.db site-monitor.db.backup
```

2. **Mettre à jour la configuration**
```bash
# Ajouter les nouvelles sections
cat >> config.json << EOF
  "ssl": {"enabled": true, "warning_thresholds": [30, 14, 7]},
  "reports": {"email": {"enabled": false}},
  "metrics": {"advanced": {"enabled": true}}
EOF
```

3. **Tester les nouvelles fonctionnalités**
```bash
site-monitor ssl                    # Test SSL monitoring
site-monitor metrics --since 1h    # Test advanced metrics  
site-monitor template list          # Vérifier templates
```

4. **Activer graduellement**
```bash
# Phase 1: SSL monitoring seulement
# Phase 2: Métriques avancées  
# Phase 3: Rapports automatiques
# Phase 4: Templates personnalisés
```

### Compatibilité
- ✅ **Base de données** : Compatible, migration automatique
- ✅ **Configuration** : Rétro-compatible, nouvelles options optionnelles
- ✅ **CLI** : Toutes les anciennes commandes fonctionnent
- ✅ **Dashboard** : Nouvelles métriques ajoutées automatiquement

---

Cette mise à jour v0.6.0 transforme Site Monitor en une solution de monitoring professionnel avec des capacités d'analyse avancées, de reporting automatique et de personnalisation poussée ! 🚀