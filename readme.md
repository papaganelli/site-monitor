# Site Monitor 🚀

Un outil de surveillance de sites web complet avec **dashboard web temps réel**, **export de données**, système d'alertes avancé et interface CLI, écrit en Go.

## ✨ Fonctionnalités

- 🏃 **Surveillance multi-sites** avec goroutines concurrentes
- 💾 **Stockage SQLite** avec historique complet des vérifications
- 📊 **Export de données** en JSON, CSV et HTML pour analyse approfondie
- 🌐 **Dashboard web moderne** avec graphiques temps réel et WebSocket
- 🚨 **Système d'alertes intelligent** (Email, Webhook, Slack, Discord, Teams)
- 🖥️  **CLI avancée** avec 6 commandes puissantes
- ⚡ **Monitoring temps réel** avec mode surveillance
- 📋 **Configuration JSON** flexible et simple
- 🎯 **Validation HTTP** avec codes de statut personnalisables
- 🔍 **Filtrage et pagination** pour l'analyse des données
- 🔔 **Notifications multi-canaux** avec templates personnalisables
- 📈 **Rapports visuels** et analyses statistiques détaillées

## 🚀 Démarrage rapide

### 1. Installation

#### Télécharger le binaire pré-compilé (recommandé)
```bash
# Télécharger la dernière version
wget https://github.com/papaganelli/site-monitor/releases/latest/download/site-monitor

# Rendre exécutable
chmod +x site-monitor

# Installer globalement (optionnel)
sudo mv site-monitor /usr/local/bin/
```

#### Compiler depuis les sources
```bash
git clone https://github.com/papaganelli/site-monitor.git
cd site-monitor
make build
```

### 2. Configuration

Créer un fichier `config.json` :
```json
{
  "sites": [
    {
      "name": "Mon Site Principal",
      "url": "https://monsite.com",
      "interval": "30s",
      "timeout": "10s"
    },
    {
      "name": "API de Production",
      "url": "https://api.monsite.com/health",
      "interval": "60s",
      "timeout": "5s"
    }
  ],
  "alerts": {
    "email": {
      "enabled": true,
      "smtp_server": "smtp.gmail.com",
      "smtp_port": 587,
      "username": "alerts@monsite.com",
      "password": "votre-mot-de-passe-app",
      "from": "Site Monitor <alerts@monsite.com>",
      "recipients": ["admin@monsite.com", "dev@monsite.com"],
      "use_tls": true
    },
    "webhook": {
      "enabled": true,
      "url": "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK",
      "format": "slack",
      "timeout": "30s",
      "retry_count": 3
    },
    "thresholds": {
      "consecutive_failures": 3,
      "response_time_threshold": "5s",
      "uptime_threshold": 95.0,
      "uptime_window": "24h",
      "performance_window": "1h",
      "alert_cooldown": "5m"
    }
  }
}
```

### 3. Utilisation

```bash
# Démarrer la surveillance avec alertes (daemon)
site-monitor run

# 🌐 Dashboard web temps réel
site-monitor dashboard --port 8080

# 📊 NOUVEAU: Export de données
site-monitor export --format json --output data.json

# Voir les statistiques
site-monitor stats

# Consulter l'historique
site-monitor history

# Vérifier le statut actuel
site-monitor status
```

## 📊 Export de Données - NOUVEAU v0.6.0 !

**Site Monitor v0.6.0** introduit un **système d'export complet** pour analyser vos données !

### 🎯 **Fonctionnalités d'Export**

- **📋 3 formats** : JSON (API), CSV (Excel), HTML (Rapports visuels)
- **🔍 Filtrage avancé** : Par site, période, limite de records
- **📈 Statistiques intégrées** : Métriques de performance complètes
- **🖥️ Interface CLI** et **🌐 API REST** complètes
- **📱 Rapports responsives** avec design moderne
- **⚡ Export temps réel** vers fichier ou stdout

### 🚀 **Utilisation CLI Export**

```bash
# Export de base (JSON, 24h, tous sites)
site-monitor export

# Formats spécialisés
site-monitor export --format json --output data.json
site-monitor export --format csv --output report.csv  
site-monitor export --format html --output report.html

# Filtrage avancé
site-monitor export --site "Mon Site" --since 7d --format csv
site-monitor export --since 1h --until "2024-01-15 12:00:00"
site-monitor export --limit 1000 --stats

# Pour pipelines et analyses
site-monitor export --stdout --format json | jq .
site-monitor export --list-formats
site-monitor export --help
```

### 📊 **API REST Export**

```bash
# Export JSON avec statistiques
curl "http://localhost:8080/api/export?format=json&since=24h&stats=true"

# Export CSV pour un site spécifique
curl "http://localhost:8080/api/export?format=csv&site=Mon%20Site&since=7d&download=true"

# Export HTML avec rapport visuel
curl "http://localhost:8080/api/export?format=html&since=1h&limit=100"

# Liste des formats disponibles
curl "http://localhost:8080/api/export/formats"
```

### 🎨 **Formats d'Export**

#### 🔸 **JSON** - Intégrations & API
```json
{
  "metadata": {
    "generated_at": "2024-01-15T10:30:00Z",
    "total_records": 150,
    "sites_included": ["Site A", "Site B"],
    "time_range": { "from": "...", "to": "..." }
  },
  "stats": {
    "overall_uptime": 98.7,
    "avg_response_time": 120000000,
    "site_stats": { "Site A": { ... } }
  },
  "history": [
    {
      "timestamp": "2024-01-15T10:25:00Z",
      "site_name": "Site A",
      "success": true,
      "status_code": 200,
      "response_time_ms": 95
    }
  ]
}
```

#### 🔸 **CSV** - Excel & Analyses
```csv
timestamp,site_name,url,success,status_code,response_time_ms,error
2024-01-15T10:25:00Z,Site A,https://site-a.com,true,200,95.00,
2024-01-15T10:24:00Z,Site B,https://site-b.com,false,500,200.00,Internal Server Error
```

#### 🔸 **HTML** - Rapports Visuels
- 🎨 **Design professionnel** avec CSS moderne responsive
- 📊 **Cartes de statistiques** avec métriques colorées  
- 📈 **Tableaux interactifs** avec données détaillées
- 🎯 **Indicateurs visuels** avec icônes et statuts
- 📱 **Mobile-friendly** adaptatif tous écrans

## 🌐 Dashboard Web

Le **dashboard web moderne** avec interface temps réel :

### 🎯 **Fonctionnalités du Dashboard**

- **📊 Vue d'ensemble temps réel** : Métriques globales et par site
- **📈 Graphiques interactifs** : Tendances de temps de réponse et distribution uptime
- **🔴 Statuts visuels** : Indicateurs colorés (Healthy/Degraded/Down/Stale)
- **📋 Activité récente** : Stream en direct des vérifications
- **⚡ WebSocket temps réel** : Mises à jour automatiques sans rechargement
- **📱 Design responsive** : Optimisé mobile et desktop
- **🌙 Mode sombre automatique** : S'adapte aux préférences système
- **📊 Intégration export** : Accès direct aux fonctions d'export

### 🚀 **Démarrer le Dashboard**

```bash
# Port par défaut (8080)
site-monitor dashboard

# Port personnalisé
site-monitor dashboard --port 3000

# Puis ouvrir dans le navigateur
open http://localhost:8080
```

## 🚨 Système d'alertes

Site Monitor intègre un système d'alertes intelligent qui vous notifie automatiquement des problèmes.

### Types d'alertes

- **🔴 Site Down** : Alertes critiques quand un site ne répond plus
- **🟢 Site Up** : Notifications de récupération après une panne  
- **🟡 Slow Response** : Avertissements pour les temps de réponse élevés
- **📉 Low Uptime** : Alertes quand l'uptime passe sous un seuil

### Canaux de notification

#### 📧 **Email (SMTP)**
Emails HTML riches avec détails complets et recommandations d'actions.

#### 🔗 **Webhooks**
Support natif pour Slack, Discord, Microsoft Teams et webhooks génériques.

**Formats supportés :**
- `slack` - Messages Slack avec attachments colorés
- `discord` - Embeds Discord riches  
- `teams` - MessageCards Microsoft Teams
- `generic` - JSON personnalisable

## 🖥️ Interface CLI

Site Monitor propose **6 commandes CLI** pour une gestion complète :

### 🏃 **`run`** - Mode surveillance (par défaut)
Démarre la surveillance continue de tous les sites configurés.

```bash
site-monitor run        # ou simplement: site-monitor
```

### 🌐 **`dashboard`** - Dashboard web
Lance le serveur web avec interface graphique moderne.

```bash
site-monitor dashboard                    # Port 8080 par défaut
site-monitor dashboard --port 3000       # Port personnalisé
```

### 📊 **`export`** - Export de données - NOUVEAU !
Exporte les données de monitoring dans différents formats.

```bash
site-monitor export                           # Export JSON par défaut (24h)
site-monitor export --format csv --since 7d  # CSV dernière semaine
site-monitor export --format html --stats    # Rapport HTML avec stats
site-monitor export --list-formats           # Formats disponibles
site-monitor export --help                   # Aide complète
```

**Options complètes :**
- `--format` : json, csv, html
- `--site` : Filtrer par site spécifique
- `--since` : Période (1h, 24h, 7d, 30d)
- `--until` : Date de fin
- `--limit` : Nombre max de records
- `--output` : Fichier de sortie
- `--stats` : Inclure les statistiques
- `--stdout` : Sortie vers stdout

### 📊 **`stats`** - Statistiques détaillées
Affiche les métriques de performance et de disponibilité.

```bash
site-monitor stats                           # Tous les sites (24h par défaut)
site-monitor stats --since 1h               # Dernière heure
site-monitor stats --since 7d               # 7 derniers jours
site-monitor stats --site "Mon Site"        # Site spécifique
```

### 📋 **`history`** - Historique des vérifications
Parcourt l'historique complet avec options de filtrage.

```bash
site-monitor history                         # Historique récent
site-monitor history --limit 50             # Limiter à 50 entrées
site-monitor history --site "Mon Site"      # Site spécifique
site-monitor history --since 2h             # 2 dernières heures
```

### 🔍 **`status`** - Statut temps réel
Affiche l'état actuel de tous les sites surveillés.

```bash
site-monitor status                          # Aperçu unique
site-monitor status --watch                 # Surveillance continue
site-monitor status --watch --interval 10s  # Rafraîchi toutes les 10s
```

## 💡 Cas d'Usage Export

### 📊 **Analyse de Performance**
```bash
# Rapport hebdomadaire complet
site-monitor export --format html --since 7d --stats --output weekly-report.html

# Analyse des tendances (JSON pour scripts)
site-monitor export --format json --since 30d --stdout | jq '.stats.overall_uptime'
```

### 📈 **Reporting Client**
```bash
# Rapport professionnel pour client
site-monitor export \
  --format html \
  --site "Site Client A" \
  --since 24h \
  --stats \
  --output "rapport-client-$(date +%Y%m%d).html"
```

### 🔄 **Intégrations BI/Analytics**
```bash
# Export CSV pour Excel/Google Sheets
site-monitor export --format csv --since 30d --output monthly-data.csv

# Données JSON pour système BI
site-monitor export --format json --stats --stdout | \
  curl -X POST -H "Content-Type: application/json" \
  -d @- https://analytics.monsite.com/api/import
```

### ⚙️ **Automatisation et Pipelines**
```bash
# Export quotidien automatique
#!/bin/bash
DATE=$(date +%Y-%m-%d)
site-monitor export --format csv --since 1d --output "daily-$DATE.csv"

# Alerte si uptime < 95%
UPTIME=$(site-monitor export --stdout --format json | jq '.stats.overall_uptime')
if (( $(echo "$UPTIME < 95" | bc -l) )); then
  echo "⚠️ Uptime critique: $UPTIME%" | mail -s "Site Monitor Alert" admin@monsite.com
fi
```

## ⚙️ Configuration avancée

### API REST Complète

Le dashboard expose maintenant une **API REST étendue** :

#### Endpoints Export
- `GET /api/export` - Export de données avec paramètres
- `GET /api/export/formats` - Liste des formats disponibles

#### Endpoints Monitoring  
- `GET /api/overview` - Vue d'ensemble du système
- `GET /api/stats` - Statistiques détaillées
- `GET /api/history` - Historique des vérifications
- `GET /api/sites` - Liste des sites configurés
- `GET /api/alerts` - Status des alertes

#### WebSocket
- `WS /ws` - Mises à jour temps réel

### Configuration Export
```json
{
  "sites": [ ... ],
  "export": {
    "default_format": "json",
    "max_records": 100000,
    "enable_api": true,
    "rate_limit": "100/hour"
  }
}
```

## 💾 Base de données et stockage

Site Monitor utilise SQLite pour stocker l'historique complet :

- **Fichier** : `site-monitor.db` (créé automatiquement)
- **Schéma optimisé** avec indexes pour les performances
- **Concurrence sûre** pour les accès multiples
- **Mode WAL** pour de meilleures performances
- **Export efficace** avec requêtes optimisées

## 🔧 Développement

### Prérequis
- Go 1.21 ou supérieur
- GCC (pour la compilation SQLite)

### Build et test
```bash
# Cloner le projet
git clone https://github.com/papaganelli/site-monitor.git
cd site-monitor

# Installer les dépendances
go mod tidy

# Compiler
make build

# Lancer les tests (incluant export)
make test

# Tests avec couverture
make test-coverage

# Formater le code
make fmt

# Linter
make lint

# Démo des fonctionnalités CLI
make demo
```

### Makefile - Commandes utiles
```bash
make build          # Compiler le binaire
make run            # Lancer en mode surveillance
make stats          # Voir les statistiques
make history        # Voir l'historique  
make status         # Voir le statut
make dashboard      # Lancer le dashboard web
make export         # Tester l'export
make clean          # Nettoyer les artifacts
make install        # Installer globalement
make demo           # Démonstration CLI complète
```

## 📦 Installation système

### Installation globale
```bash
# Après compilation ou téléchargement
make install
# ou manuellement :
sudo cp site-monitor /usr/local/bin/

# Vérification
site-monitor --version  # v0.6.0
```

### Service systemd (Linux)
Créer `/etc/systemd/system/site-monitor.service` :
```ini
[Unit]
Description=Site Monitor
After=network.target

[Service]
Type=simple
User=monitoring
WorkingDirectory=/opt/site-monitor
ExecStart=/usr/local/bin/site-monitor run
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Puis :
```bash
sudo systemctl enable site-monitor
sudo systemctl start site-monitor
sudo systemctl status site-monitor
```

## 🎯 Codes de statut

- ✅ **OK (200-399)** : Succès et redirections HTTP
- ❌ **ERREUR (400+)** : Erreurs client/serveur
- 🔄 **TIMEOUT** : Dépassement du timeout configuré
- 🌐 **RÉSEAU** : Erreurs de connectivité réseau

## 🚀 Roadmap

### ✅ Version 0.6.0 (Actuelle)
- ✅ **Export complet** JSON, CSV, HTML avec API REST
- ✅ **Statistiques avancées** dans les exports
- ✅ **Filtrage et pagination** des données
- ✅ **Rapports HTML** visuels et responsives
- ✅ **Interface CLI étendue** avec 6 commandes
- ✅ **Dashboard web moderne** avec WebSocket temps réel
- ✅ Système d'alertes complet (Email, Webhook)
- ✅ Support Slack, Discord, Microsoft Teams
- ✅ CLI avancée et stockage SQLite complet

### 🔮 Version 0.7.0 (Prochaine)
- [ ] 🔒 **Monitoring SSL/TLS** avec alertes d'expiration certificats
- [ ] 📈 **Métriques avancées** (P95, P99, MTTR, MTBF)
- [ ] 📊 **Export Excel** (.xlsx) natif avec graphiques
- [ ] 🔄 **Export programmé** (cron-like) automatique
- [ ] 📧 **Export par email** avec rapports périodiques
- [ ] 🎨 **Templates d'alertes** personnalisables

### 🔮 Version 0.8.0
- [ ] 🐳 **Support Docker et Kubernetes** complet
- [ ] ☁️  **Déploiement cloud** (AWS, GCP, Azure)  
- [ ] 🔗 **Intégrations** (Grafana, Prometheus, DataDog)
- [ ] 🌍 **Monitoring multi-régions** et géo-distribué
- [ ] 📱 **Application mobile** companion
- [ ] 🗜️ **Compression** automatique des exports

## 🤝 Contribution

1. **Fork** le projet
2. **Créer** une branche (`git checkout -b feature/export-excel`)
3. **Committer** (`git commit -m 'feat: ajouter export Excel avec graphiques'`)
4. **Push** (`git push origin feature/export-excel`)
5. **Ouvrir** une Pull Request

### Guidelines
- ✅ Tests unitaires pour les nouvelles fonctionnalités
- ✅ Code formaté avec `gofmt`
- ✅ Documentation des APIs publiques
- ✅ Commentaires en anglais dans le code
- ✅ Messages de commit descriptifs

## 📊 Architecture

```
site-monitor/
├── main.go                    # Point d'entrée et CLI parsing
├── cmd/                       # Commandes CLI
│   ├── app.go                # Structure principale CLI
│   ├── stats.go              # Commande statistiques
│   ├── history.go            # Commande historique
│   ├── status.go             # Commande statut
│   ├── dashboard.go          # Commande dashboard web
│   └── export.go             # Commande export - NOUVEAU !
├── config/
│   └── config.go             # Configuration (sites + alertes)
├── monitor/
│   ├── checker.go            # Logique de surveillance
│   └── result.go             # Structure des résultats
├── storage/
│   ├── storage.go            # Interface générique
│   └── sqlite.go             # Implémentation SQLite
├── export/                   # Système d'export - NOUVEAU !
│   ├── types.go              # Types et structures d'export
│   ├── exporter.go           # Logique principale d'export
│   ├── formatters.go         # Formatters JSON/CSV/HTML
│   ├── exporter_test.go      # Tests unitaires
│   └── formatters_test.go    # Tests des formatters
├── alerts/                   # Système d'alertes
│   ├── types.go              # Types et interfaces d'alertes
│   ├── manager.go            # Gestionnaire central d'alertes
│   ├── email.go              # Canal d'alerte email (SMTP)
│   └── webhook.go            # Canal webhook (Slack/Discord/Teams)
├── web/                      # Dashboard web
│   ├── server.go             # Serveur HTTP et API REST (+ export API)
│   ├── dashboard.go          # Templates HTML/CSS/JS
│   └── types.go              # Types pour API REST
├── config.json               # Configuration des sites et alertes
└── site-monitor.db           # Base SQLite (auto-créée)
```

## 📄 Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.

## 💬 Support

- 🐛 **Bugs** : [Issues GitHub](https://github.com/papaganelli/site-monitor/issues)
- 💡 **Fonctionnalités** : [Discussions](https://github.com/papaganelli/site-monitor/discussions)
- 📖 **Documentation** : [Wiki](https://github.com/papaganelli/site-monitor/wiki)
- 📊 **Guide Export** : [EXPORT.md](EXPORT.md)

## 📈 Changelog

### v0.6.0 - Export de Données Complet 📊
- 📊 **Système d'export complet** avec 3 formats (JSON, CSV, HTML)
- 🖥️  **Commande CLI export** avec options avancées et aide intégrée
- 🌐 **API REST export** (`/api/export`, `/api/export/formats`)
- 📈 **Statistiques étendues** dans tous les exports
- 🎨 **Rapports HTML** professionnels avec design responsive
- 🔍 **Filtrage avancé** (site, période, limite, until)
- ⚡ **Support stdout** pour pipelines et intégrations
- 📝 **Documentation complète** avec exemples d'usage
- 🧪 **Tests unitaires** complets pour tous les composants

### v0.5.0 - Dashboard Web Moderne 🌐
- 🌐 **Dashboard web complet** avec interface graphique moderne et responsive
- ⚡ **WebSocket temps réel** pour mises à jour automatiques sans rechargement
- 📊 **Graphiques interactifs** (Chart.js) - tendances temps de réponse et distribution uptime
- 📱 **Design responsive** optimisé pour mobile et desktop avec mode sombre
- 🔗 **API REST complète** (/api/overview, /api/stats, /api/history, /api/sites, /api/alerts)
- 🎨 **Interface utilisateur moderne** avec animations, toasts et indicateurs visuels
- 🚀 **Commande dashboard** ajoutée : `site-monitor dashboard --port 8080`
- 🛡️  **Gestion d'erreurs améliorée** dans toutes les couches (HTTP, WebSocket, JSON)

### v0.4.0 - Système d'Alertes Intelligent
- 🚨 **Système d'alertes complet** avec 4 types d'alertes automatiques
- 📧 **Canal Email** avec templates HTML riches et SMTP configurable
- 🔗 **Canal Webhook** avec support Slack, Discord, Microsoft Teams
- ⚙️  **Seuils configurables** pour tous les types d'alertes
- 🛡️  **Anti-spam** avec cooldown et logique d'état intelligente
- 🔄 **Retry automatique** pour les webhooks avec backoff exponentiel
- 📊 **Intégration complète** avec le système de monitoring existant

### v0.3.0 - CLI Avancée
- ✨ Ajout CLI complète avec 4 commandes puissantes
- 📊 Statistiques détaillées avec métriques de performance
- 📋 Navigation dans l'historique avec filtres et pagination
- 🔍 Monitoring temps réel avec mode surveillance
- 🎨 Interface colorée avec emojis et formatage professionnel

### v0.2.0 - Stockage SQLite
- 💾 Ajout stockage SQLite pour l'historique complet
- 📈 Calcul automatique des statistiques de disponibilité
- 🔧 Interface Storage générique pour extensibilité future
- ⚡ Optimisations avec indexes et mode WAL

### v0.1.1 - Corrections CI/CD
- 🔧 Correction permissions GitHub Actions
- 🌐 Commentaires code en anglais
- 📝 Améliorations notes de release

### v0.1.0 - Version initiale
- 🚀 Surveillance multi-sites concurrent
- ⚙️ Configuration JSON flexible
- 📊 Affichage temps réel console
- 🔍 Validation codes de statut HTTP

---

**Fait avec ❤️ en Go** • [Site Monitor v0.6.0](https://github.com/papaganelli/site-monitor)