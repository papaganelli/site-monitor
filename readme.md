# Site Monitor 🚀

Un outil de surveillance de sites web simple et efficace écrit en Go.

## Fonctionnalités

- ✅ Surveille plusieurs sites web simultanément
- ⚡ Surveillance concurrente avec les goroutines
- 📋 Configuration via fichier JSON
- 🕐 Intervalles de vérification et timeouts personnalisables
- 📊 Rapport de statut en temps réel
- 🎯 Validation des codes de statut HTTP

## Démarrage rapide

1. **Cloner le dépôt**
   ```bash
   git clone https://github.com/ton-username/site-monitor.git
   cd site-monitor
   ```

2. **Configurer vos sites**
   
   Éditer `config.json` :
   ```json
   {
     "sites": [
       {
         "name": "Mon Site Web",
         "url": "https://example.com",
         "interval": "30s",
         "timeout": "10s"
       }
     ]
   }
   ```

3. **Lancer le monitoring**
   ```bash
   go run main.go
   ```

## Configuration

| Champ | Description | Exemple |
|-------|-------------|---------|
| `name` | Nom d'affichage pour le site | `"Mon API"` |
| `url` | URL à surveiller | `"https://api.example.com/health"` |
| `interval` | Fréquence de vérification | `"30s"`, `"5m"`, `"1h"` |
| `timeout` | Timeout des requêtes | `"10s"`, `"30s"` |

### Exemple de configuration complète

```json
{
  "sites": [
    {
      "name": "Site Principal",
      "url": "https://monsite.com",
      "interval": "30s",
      "timeout": "10s"
    },
    {
      "name": "API de Production",
      "url": "https://api.monsite.com/health",
      "interval": "60s",
      "timeout": "5s"
    },
    {
      "name": "Service de Paiement",
      "url": "https://payment.monsite.com/status",
      "interval": "2m",
      "timeout": "15s"
    }
  ]
}
```

## Installation

### Pré-requis
- Go 1.19 ou supérieur

### Depuis les sources
```bash
git clone https://github.com/ton-username/site-monitor.git
cd site-monitor
go build -o site-monitor
```

### Télécharger un binaire pré-compilé
Rendez-vous sur la [page des releases](https://github.com/ton-username/site-monitor/releases) pour télécharger un binaire pour votre plateforme.

## Compilation

```bash
# Compiler pour la plateforme courante
go build -o site-monitor

# Compilation croisée pour différentes plateformes
GOOS=linux GOARCH=amd64 go build -o site-monitor-linux
GOOS=windows GOARCH=amd64 go build -o site-monitor.exe
GOOS=darwin GOARCH=amd64 go build -o site-monitor-macos
```

## Utilisation

### Lancer le monitoring
```bash
./site-monitor
```

### Exemple de sortie
```
🚀 Démarrage du monitoring pour 3 sites
📍 Démarrage de Site Principal (https://monsite.com) - vérification toutes les 30s
📍 Démarrage de API de Production (https://api.monsite.com/health) - vérification toutes les 60s
📍 Démarrage de Service de Paiement (https://payment.monsite.com/status) - vérification toutes les 2m

[14:30:25] ✅ OK (Site Principal) - Status: 200 - Durée: 245ms
[14:30:27] ✅ OK (API de Production) - Status: 200 - Durée: 89ms
[14:30:29] ❌ ERREUR (Service de Paiement) - Status: 500 - Durée: 1.2s
[14:31:25] ✅ OK (Site Principal) - Status: 200 - Durée: 198ms
```

### Codes de statut
- ✅ **OK** : Codes HTTP 200-399 (succès et redirections)
- ❌ **ERREUR** : Codes HTTP 400+ (erreurs client/serveur) ou timeout/erreur réseau

## Tests

```bash
# Lancer tous les tests
go test ./...

# Tests avec couverture
go test -cover ./...

# Tests en mode verbose
go test -v ./...
```

## Architecture du projet

```
site-monitor/
├── main.go              # Point d'entrée
├── config/
│   └── config.go        # Gestion de la configuration JSON
├── monitor/
│   ├── checker.go       # Logique de surveillance
│   ├── checker_test.go  # Tests unitaires
│   └── result.go        # Structure des résultats
├── config.json          # Fichier de configuration
├── .github/
│   └── workflows/
│       └── ci.yml       # Pipeline CI/CD
└── README.md
```

## Contribution

1. Fork le projet
2. Créer une branche pour votre fonctionnalité (`git checkout -b feature/nouvelle-fonctionnalite`)
3. Commit vos changements (`git commit -m 'Ajouter une nouvelle fonctionnalité'`)
4. Push vers la branche (`git push origin feature/nouvelle-fonctionnalite`)
5. Ouvrir une Pull Request

### Guidelines de développement

- Écrire des tests pour les nouvelles fonctionnalités
- Suivre les conventions Go (gofmt, golint)
- Documenter les nouvelles API publiques
- Commenter le code en anglais

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.

## Roadmap

### Version 0.2.0
- [ ] Sauvegarde des données (JSON/SQLite)
- [ ] Alertes par email/webhook
- [ ] Interface en ligne de commande améliorée

### Version 0.3.0
- [ ] Dashboard web
- [ ] Métriques avancées (temps de réponse, disponibilité)
- [ ] Support Docker

### Version 0.4.0
- [ ] Vérifications avancées (contenu, SSL, DNS)
- [ ] Intégrations (Slack, Teams, PagerDuty)
- [ ] Base de données distribuée

## Support

- 🐛 **Bugs** : Ouvrir une [issue](https://github.com/ton-username/site-monitor/issues)
- 💡 **Suggestions** : Ouvrir une [discussion](https://github.com/ton-username/site-monitor/discussions)
- 📖 **Documentation** : Voir le [wiki](https://github.com/ton-username/site-monitor/wiki)

## Changelog

### v0.1.0
- Surveillance multi-sites avec goroutines
- Configuration JSON
- Vérifications HTTP avec timeouts
- Sortie console en temps réel

---

Fait avec ❤️ en Go
```

