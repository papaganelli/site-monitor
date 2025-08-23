# Site Monitor 🚀

Un outil de surveillance de sites web complet avec interface CLI avancée, écrit en Go.

## ✨ Fonctionnalités

- 🏃 **Surveillance multi-sites** avec goroutines concurrentes
- 💾 **Stockage SQLite** avec historique complet des vérifications
- 📊 **Statistiques détaillées** (uptime, temps de réponse, SLA)
- 🖥️  **CLI avancée** avec 4 commandes puissantes
- ⚡ **Monitoring temps réel** avec mode surveillance
- 📋 **Configuration JSON** flexible et simple
- 🎯 **Validation HTTP** avec codes de statut personnalisables
- 🔍 **Filtrage et pagination** pour l'analyse des données

## 🚀 Démarrage rapide

### 1. Installation

#### Télécharger le binaire pré-compilé (recommandé)
```bash
# Télécharger la dernière version
wget https://github.com/ton-username/site-monitor/releases/latest/download/site-monitor

# Rendre exécutable
chmod +x site-monitor

# Installer globalement (optionnel)
sudo mv site-monitor /usr/local/bin/
```

#### Compiler depuis les sources
```bash
git clone https://github.com/ton-username/site-monitor.git
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
  ]
}
```

### 3. Utilisation

```bash
# Démarrer la surveillance (daemon)
site-monitor run

# Voir les statistiques
site-monitor stats

# Consulter l'historique
site-monitor history

# Vérifier le statut actuel
site-monitor status
```

## 📋 Commandes CLI

### 🏃 **`run`** - Mode surveillance (par défaut)
Démarre la surveillance continue de tous les sites configurés.

```bash
site-monitor run        # ou simplement: site-monitor
```

**Sortie :**
```
🚀 Starting monitoring for 2 sites
💾 Database initialized: site-monitor.db
📍 Starting Mon Site Principal (https://monsite.com) - checking every 30s
📍 Starting API de Production (https://api.monsite.com/health) - checking every 60s

[14:30:25] ✅ OK (Mon Site Principal) - Status: 200 - Duration: 245ms
[14:30:27] ✅ OK (API de Production) - Status: 200 - Duration: 89ms
[14:31:25] ✅ OK (Mon Site Principal) - Status: 200 - Duration: 198ms
```

### 📊 **`stats`** - Statistiques détaillées
Affiche les métriques de performance et de disponibilité.

```bash
site-monitor stats                           # Tous les sites (24h par défaut)
site-monitor stats --since 1h               # Dernière heure
site-monitor stats --since 7d               # 7 derniers jours
site-monitor stats --site "Mon Site"        # Site spécifique
```

**Exemple de sortie :**
```
📊 Monitoring Statistics (Last 24 hours)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Mon Site Principal
   📈 Uptime: 99.2% (1,432/1,444 checks)
   ⚡ Response: 245ms avg (min: 89ms, max: 1.2s)
   🕐 Last Check: 2 minutes ago
   📅 Monitoring Duration: 24h

⚠️ API de Production  
   📈 Uptime: 97.8% (1,411/1,444 checks)
   ⚡ Response: 156ms avg (min: 45ms, max: 2.1s)
   🕐 Last Check: 1 minute ago
   💥 Failed Checks: 33

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Summary: 2 sites monitored
🎯 Overall Uptime: 98.5% (2,843/2,888 checks)
💚 Healthy Sites: 1/2 (≥99% uptime)
```

### 📋 **`history`** - Historique des vérifications
Parcourt l'historique complet avec options de filtrage.

```bash
site-monitor history                         # Historique récent
site-monitor history --limit 50             # Limiter à 50 entrées
site-monitor history --site "Mon Site"      # Site spécifique
site-monitor history --since 2h             # 2 dernières heures
```

**Exemple de sortie :**
```
📋 Monitoring History (Last 24 hours) - Limited to 20 entries
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🌐 Mon Site Principal (15 entries)
──────────────────────────────────────────────────────
   [14:35:25] ✅ OK - 200 - 198ms
   [14:34:55] ✅ OK - 200 - 201ms
   [14:34:25] ❌ FAIL - 0 - 10s - context deadline exceeded
   [14:33:55] ✅ OK - 200 - 187ms

🌐 API de Production (5 entries)
──────────────────────────────────────────────────────
   [14:35:15] ✅ OK - 200 - 89ms
   [14:34:15] ✅ OK - 200 - 92ms

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Summary: 20 entries from 2 sites
✅ Success Rate: 95.0% (19/20)
⚡ Response Times: 168ms avg (min: 89ms, max: 10s)
⏱️  Time Span: 1h15m
```

### 🔍 **`status`** - Statut temps réel
Affiche l'état actuel de tous les sites surveillés.

```bash
site-monitor status                          # Aperçu unique
site-monitor status --watch                 # Surveillance continue
site-monitor status --watch --interval 10s  # Rafraîchi toutes les 10s
```

**Exemple de sortie :**
```
🚀 Site Monitor Status - 2025-08-23 14:35:42
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ HEALTHY      Mon Site Principal
   📈 Recent Success: 99.2% (143/144 checks)
   ⚡ Response Time: 201ms avg
   🕐 Last Check: 17 seconds ago

⚠️ DEGRADED     API de Production
   📈 Recent Success: 96.5% (138/143 checks)
   ⚡ Response Time: 89ms avg
   🕐 Last Check: 45 seconds ago
   💥 Recent Failures: 5
   🚨 Issues: Some failures

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Overall: 1/2 sites healthy ⚠️ Some issues detected
```

### 📖 **`--help`** - Aide complète
```bash
site-monitor --help         # Aide générale
site-monitor --version      # Version du logiciel
```

## ⚙️ Configuration avancée

### Format du fichier config.json

```json
{
  "sites": [
    {
      "name": "Nom d'affichage du site",
      "url": "https://example.com/health",
      "interval": "30s",
      "timeout": "10s"
    }
  ]
}
```

| Champ | Description | Exemples | Obligatoire |
|-------|-------------|----------|-------------|
| `name` | Nom affiché dans les rapports | `"API Production"` | ✅ |
| `url` | URL à surveiller | `"https://api.com/health"` | ✅ |
| `interval` | Fréquence des vérifications | `"30s"`, `"5m"`, `"1h"` | ✅ |
| `timeout` | Timeout des requêtes HTTP | `"10s"`, `"30s"` | ✅ |

### Exemples de configurations

#### Configuration basique
```json
{
  "sites": [
    {
      "name": "Site Web",
      "url": "https://monsite.com",
      "interval": "60s",
      "timeout": "15s"
    }
  ]
}
```

#### Configuration multi-environnements
```json
{
  "sites": [
    {
      "name": "Production - Site Principal",
      "url": "https://monsite.com",
      "interval": "30s",
      "timeout": "10s"
    },
    {
      "name": "Production - API",
      "url": "https://api.monsite.com/health",
      "interval": "60s",
      "timeout": "5s"
    },
    {
      "name": "Staging - API",
      "url": "https://staging-api.monsite.com/health",
      "interval": "5m",
      "timeout": "15s"
    },
    {
      "name": "Service de Paiement",
      "url": "https://payments.monsite.com/status",
      "interval": "2m",
      "timeout": "20s"
    }
  ]
}
```

## 💾 Base de données et stockage

Site Monitor utilise SQLite pour stocker l'historique complet :

- **Fichier** : `site-monitor.db` (créé automatiquement)
- **Schéma optimisé** avec indexes pour les performances
- **Concurrence sûre** pour les accès multiples
- **Mode WAL** pour de meilleures performances

### Requêtes manuelles (optionnel)
```bash
# Examiner la base de données
sqlite3 site-monitor.db

# Quelques requêtes utiles
.schema                                    # Structure des tables
SELECT COUNT(*) FROM results;             # Nombre total de vérifications
SELECT * FROM results ORDER BY timestamp DESC LIMIT 10;  # 10 dernières vérifications
```

## 🔧 Développement

### Prérequis
- Go 1.21 ou supérieur
- GCC (pour la compilation SQLite)

### Build et test
```bash
# Cloner le projet
git clone https://github.com/ton-username/site-monitor.git
cd site-monitor

# Installer les dépendances
go mod tidy

# Compiler
make build

# Lancer les tests
make test

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
make clean          # Nettoyer les artifacts
make install        # Installer globalement
make demo           # Démonstration CLI
```

## 📦 Installation système

### Installation globale
```bash
# Après compilation ou téléchargement
make install
# ou manuellement :
sudo cp site-monitor /usr/local/bin/

# Vérification
site-monitor --version
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

### ✅ Version 0.3.0 (Actuelle)
- ✅ CLI avancée avec 4 commandes
- ✅ Stockage SQLite complet
- ✅ Statistiques détaillées
- ✅ Historique avec filtrage
- ✅ Monitoring temps réel

### 🔮 Version 0.4.0 (Prochaine)
- [ ] 📧 Système d'alertes (email, webhook, Slack)
- [ ] 🌐 Dashboard web avec graphiques
- [ ] 📊 Export des données (JSON, CSV)
- [ ] 🔔 Notifications push
- [ ] 📈 Métriques avancées (p95, p99)

### 🔮 Version 0.5.0
- [ ] 🐳 Support Docker complet
- [ ] ☁️  Déploiement cloud (AWS, GCP)
- [ ] 🔗 Intégrations (Grafana, Prometheus)
- [ ] 🛡️  Vérifications SSL/TLS
- [ ] 🌍 Monitoring multi-régions

## 🤝 Contribution

1. **Fork** le projet
2. **Créer** une branche (`git checkout -b feature/nouvelle-fonctionnalite`)
3. **Committer** (`git commit -m 'feat: ajouter nouvelle fonctionnalité'`)
4. **Push** (`git push origin feature/nouvelle-fonctionnalite`)
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
│   └── status.go             # Commande statut
├── config/
│   └── config.go             # Gestion configuration JSON
├── monitor/
│   ├── checker.go            # Logique de surveillance
│   └── result.go             # Structure des résultats
├── storage/
│   ├── storage.go            # Interface générique
│   └── sqlite.go             # Implémentation SQLite
├── config.json               # Configuration des sites
└── site-monitor.db           # Base SQLite (auto-créée)
```

## 📄 Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.

## 💬 Support

- 🐛 **Bugs** : [Issues GitHub](https://github.com/ton-username/site-monitor/issues)
- 💡 **Fonctionnalités** : [Discussions](https://github.com/ton-username/site-monitor/discussions)
- 📖 **Documentation** : [Wiki](https://github.com/ton-username/site-monitor/wiki)

## 📈 Changelog

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

**Fait avec ❤️ en Go** • [Site Monitor](https://github.com/ton-username/site-monitor)