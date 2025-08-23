# Site Monitor 🚀

Un outil de surveillance de sites web complet avec **dashboard web temps réel**, système d'alertes avancé et interface CLI, écrit en Go.

## ✨ Fonctionnalités

- 🏃 **Surveillance multi-sites** avec goroutines concurrentes
- 💾 **Stockage SQLite** avec historique complet des vérifications
- 📊 **Statistiques détaillées** (uptime, temps de réponse, SLA)
- 🌐 **Dashboard web moderne** avec graphiques temps réel et WebSocket
- 🚨 **Système d'alertes intelligent** (Email, Webhook, Slack, Discord, Teams)
- 🖥️  **CLI avancée** avec 5 commandes puissantes
- ⚡ **Monitoring temps réel** avec mode surveillance
- 📋 **Configuration JSON** flexible et simple
- 🎯 **Validation HTTP** avec codes de statut personnalisables
- 🔍 **Filtrage et pagination** pour l'analyse des données
- 🔔 **Notifications multi-canaux** avec templates personnalisables

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

# 🌐 NOUVEAU: Lancer le dashboard web
site-monitor dashboard --port 8080

# Voir les statistiques
site-monitor stats

# Consulter l'historique
site-monitor history

# Vérifier le statut actuel
site-monitor status
```

## 🌐 Dashboard Web - NOUVEAU !

**Site Monitor v0.5.0** introduit un **dashboard web moderne** avec interface temps réel !

### 🎯 **Fonctionnalités du Dashboard**

- **📊 Vue d'ensemble temps réel** : Métriques globales et par site
- **📈 Graphiques interactifs** : Tendances de temps de réponse et distribution uptime
- **🔴 Statuts visuels** : Indicateurs colorés (Healthy/Degraded/Down/Stale)
- **📋 Activité récente** : Stream en direct des vérifications
- **⚡ WebSocket temps réel** : Mises à jour automatiques sans rechargement
- **📱 Design responsive** : Optimisé mobile et desktop
- **🌙 Mode sombre automatique** : S'adapte aux préférences système

### 🚀 **Démarrer le Dashboard**

```bash
# Port par défaut (8080)
site-monitor dashboard

# Port personnalisé
site-monitor dashboard --port 3000

# Puis ouvrir dans le navigateur
open http://localhost:8080
```

### 📸 **Aperçu du Dashboard**

Le dashboard affiche :
- **Cartes de résumé** : Sites totaux, sites sains, uptime global, vérifications totales
- **Grille des sites** : Statut, uptime et temps de réponse par site
- **Graphiques temps réel** : Tendances des temps de réponse sur 24h
- **Activité live** : Stream des dernières vérifications avec statuts
- **Indicateur de connexion** : WebSocket connecté/déconnecté

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

```json
{
  "email": {
    "enabled": true,
    "smtp_server": "smtp.gmail.com",
    "smtp_port": 587,
    "username": "alerts@monsite.com",
    "password": "votre-mot-de-passe-app",
    "from": "Site Monitor <alerts@monsite.com>",
    "recipients": ["admin@monsite.com", "ops@monsite.com"],
    "use_tls": true
  }
}
```

#### 🔗 **Webhooks**
Support natif pour Slack, Discord, Microsoft Teams et webhooks génériques.

```json
{
  "webhook": {
    "enabled": true,
    "url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
    "format": "slack",
    "timeout": "30s",
    "retry_count": 3,
    "headers": {
      "Authorization": "Bearer your-token"
    }
  }
}
```

**Formats supportés :**
- `slack` - Messages Slack avec attachments colorés
- `discord` - Embeds Discord riches  
- `teams` - MessageCards Microsoft Teams
- `generic` - JSON personnalisable

## 🖥️ Interface CLI

Site Monitor propose 5 commandes CLI pour une gestion complète :

### 🏃 **`run`** - Mode surveillance (par défaut)
Démarre la surveillance continue de tous les sites configurés.

```bash
site-monitor run        # ou simplement: site-monitor
```

### 🌐 **`dashboard`** - Dashboard web - NOUVEAU !
Lance le serveur web avec interface graphique moderne.

```bash
site-monitor dashboard                    # Port 8080 par défaut
site-monitor dashboard --port 3000       # Port personnalisé
```

**Fonctionnalités du dashboard :**
- Interface web moderne et responsive
- Graphiques temps réel avec Chart.js
- WebSocket pour mises à jour automatiques
- Vue d'ensemble et détails par site
- Stream d'activité en direct

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

### 📖 **`--help`** - Aide complète
```bash
site-monitor --help         # Aide générale
site-monitor --version      # Version du logiciel
```

## ⚙️ Configuration avancée

### Fichier config.json complet

```json
{
  "sites": [
    {
      "name": "Site Principal Production",
      "url": "https://monsite.com",
      "interval": "30s",
      "timeout": "10s"
    },
    {
      "name": "API REST v1",
      "url": "https://api.monsite.com/v1/health",
      "interval": "60s",
      "timeout": "5s"
    },
    {
      "name": "Service de Paiement",
      "url": "https://payments.monsite.com/status",
      "interval": "2m",
      "timeout": "15s"
    }
  ],
  "alerts": {
    "email": {
      "enabled": true,
      "smtp_server": "smtp.gmail.com",
      "smtp_port": 587,
      "username": "monitoring@monsite.com",
      "password": "app-specific-password",
      "from": "Site Monitor <noreply@monsite.com>",
      "recipients": [
        "admin@monsite.com",
        "devops@monsite.com",
        "on-call@monsite.com"
      ],
      "use_tls": true
    },
    "webhook": {
      "enabled": true,
      "url": "https://hooks.slack.com/services/T123/B456/xyz789",
      "format": "slack",
      "timeout": "30s",
      "retry_count": 3,
      "headers": {
        "User-Agent": "SiteMonitor/0.5.0"
      }
    },
    "thresholds": {
      "consecutive_failures": 3,
      "response_time_threshold": "5s",
      "uptime_threshold": 95.0,
      "uptime_window": "24h",
      "performance_window": "1h",
      "alert_cooldown": "10m"
    }
  }
}
```

## 💾 Base de données et stockage

Site Monitor utilise SQLite pour stocker l'historique complet :

- **Fichier** : `site-monitor.db` (créé automatiquement)
- **Schéma optimisé** avec indexes pour les performances
- **Concurrence sûre** pour les accès multiples
- **Mode WAL** pour de meilleures performances

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
make dashboard      # Lancer le dashboard web
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

### ✅ Version 0.5.0 (Actuelle)
- ✅ **Dashboard web moderne** avec interface graphique complète
- ✅ **WebSocket temps réel** pour mises à jour automatiques
- ✅ **Graphiques interactifs** (Chart.js) - temps de réponse et uptime
- ✅ **Design responsive** optimisé mobile et desktop
- ✅ **API REST complète** pour intégrations tierces
- ✅ Système d'alertes complet (Email, Webhook)
- ✅ Support Slack, Discord, Microsoft Teams
- ✅ Seuils configurables et logique intelligente
- ✅ Templates d'emails HTML riches
- ✅ CLI avancée avec 5 commandes
- ✅ Stockage SQLite complet

### 🔮 Version 0.6.0 (Prochaine)
- [ ] 📊 **Export des données** (JSON, CSV, API REST étendue)
- [ ] 🔔 **Notifications push** et intégrations mobiles
- [ ] 📈 **Métriques avancées** (p95, p99, MTTR, MTBF)
- [ ] 🎨 **Templates d'alertes** personnalisables
- [ ] 🛡️  **Vérifications SSL/TLS** et monitoring certificats

### 🔮 Version 0.7.0
- [ ] 🐳 **Support Docker et Kubernetes** complet
- [ ] ☁️  **Déploiement cloud** (AWS, GCP, Azure)
- [ ] 🔗 **Intégrations** (Grafana, Prometheus, DataDog)
- [ ] 🌍 **Monitoring multi-régions** et géo-distribué
- [ ] 📱 **Application mobile** companion

## 🤝 Contribution

1. **Fork** le projet
2. **Créer** une branche (`git checkout -b feature/nouvelle-fonctionnalite`)
3. **Committer** (`git commit -m 'feat: ajouter dashboard web moderne'`)
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
│   ├── status.go             # Commande statut
│   └── dashboard.go          # Commande dashboard web
├── config/
│   └── config.go             # Configuration (sites + alertes)
├── monitor/
│   ├── checker.go            # Logique de surveillance
│   └── result.go             # Structure des résultats
├── storage/
│   ├── storage.go            # Interface générique
│   └── sqlite.go             # Implémentation SQLite
├── alerts/                   # Système d'alertes
│   ├── types.go              # Types et interfaces d'alertes
│   ├── manager.go            # Gestionnaire central d'alertes
│   ├── email.go              # Canal d'alerte email (SMTP)
│   └── webhook.go            # Canal webhook (Slack/Discord/Teams)
├── web/                      # Dashboard web - NOUVEAU !
│   ├── server.go             # Serveur HTTP et API REST
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

## 📈 Changelog

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

**Fait avec ❤️ en Go** • [Site Monitor](https://github.com/papaganelli/site-monitor)