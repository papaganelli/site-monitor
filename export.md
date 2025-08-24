# 📊 Site Monitor - Data Export

La fonctionnalité **Export** permet d'extraire les données de monitoring dans différents formats pour une analyse approfondie, des rapports ou des intégrations avec d'autres outils.

## ✨ Fonctionnalités

- **3 formats d'export** : JSON, CSV, HTML
- **Filtrage flexible** : Par site, période, limite de records
- **Statistiques intégrées** : Métriques de performance incluses
- **Interface CLI** et **API REST** complètes
- **Rapports HTML** visuels avec graphiques intégrés
- **Export vers fichier** ou **sortie standard**

## 🚀 Utilisation CLI

### Commande de base

```bash
# Export par défaut (JSON, 24h, tous sites)
site-monitor export

# Export avec format spécifique
site-monitor export --format json --output data.json
site-monitor export --format csv --output report.csv
site-monitor export --format html --output report.html
```

### Options avancées

```bash
# Filtrer par site
site-monitor export --site "Mon Site Principal" --format csv

# Période personnalisée
site-monitor export --since 7d --format json
site-monitor export --since 1h --until "2024-01-15 12:00:00"

# Limiter le nombre de records
site-monitor export --limit 1000 --format csv

# Inclure les statistiques
site-monitor export --stats --format html

# Sortie vers stdout (pour pipelines)
site-monitor export --stdout --format json | jq .

# Combiner plusieurs options
site-monitor export \
  --format html \
  --site "API Production" \
  --since 3d \
  --stats \
  --output "api-report.html"
```

### Aide et formats disponibles

```bash
# Afficher l'aide
site-monitor export --help

# Lister les formats supportés
site-monitor export --list-formats
```

## 🌐 API REST

### Endpoints disponibles

#### GET /api/export
Export des données de monitoring.

**Paramètres de requête :**
- `format` : Format d'export (json, csv, html) [défaut: json]
- `site` : Nom du site à filtrer (optionnel)
- `since` : Période (ex: 24h, 7d) [défaut: 24h]
- `until` : Heure de fin (format RFC3339, optionnel)
- `limit` : Nombre max de records (optionnel)
- `stats` : Inclure les statistiques (true/false)
- `download` : Forcer le téléchargement (true/false)

**Exemples :**
```bash
# Export JSON avec stats
curl "http://localhost:8080/api/export?format=json&since=24h&stats=true"

# Export CSV pour un site spécifique
curl "http://localhost:8080/api/export?format=csv&site=Mon%20Site&since=7d&download=true"

# Export HTML avec limite
curl "http://localhost:8080/api/export?format=html&since=1h&limit=100"
```

#### GET /api/export/formats
Liste des formats d'export disponibles.

**Réponse :**
```json
{
  "formats": [
    {
      "format": "json",
      "description": "JSON - Machine-readable structured data",
      "content_type": "application/json",
      "file_extension": ".json"
    },
    {
      "format": "csv", 
      "description": "CSV - Spreadsheet compatible comma-separated values",
      "content_type": "text/csv",
      "file_extension": ".csv"
    },
    {
      "format": "html",
      "description": "HTML - Human-readable web page report", 
      "content_type": "text/html",
      "file_extension": ".html"
    }
  ],
  "default": "json"
}
```

## 📋 Formats d'export

### 🔸 JSON Format

**Usage :** Intégrations API, analyse programmatique, sauvegarde structurée

**Structure :**
```json
{
  "metadata": {
    "generated_at": "2024-01-15T10:30:00Z",
    "format": "json",
    "total_records": 150,
    "sites_included": ["Site A", "Site B"],
    "time_range": {
      "from": "2024-01-14T10:30:00Z",
      "to": "2024-01-15T10:30:00Z"
    }
  },
  "stats": {
    "total_sites": 2,
    "total_checks": 150,
    "successful_checks": 147,
    "overall_uptime": 98.0,
    "avg_response_time": 120000000,
    "site_stats": {
      "Site A": {
        "total_checks": 75,
        "success_rate": 98.7,
        "avg_response_time": 100000000
      }
    }
  },
  "history": [
    {
      "timestamp": "2024-01-15T10:25:00Z",
      "site_name": "Site A",
      "url": "https://site-a.com",
      "success": true,
      "status_code": 200,
      "response_time_ms": 95,
      "error": ""
    }
  ]
}
```

### 🔸 CSV Format

**Usage :** Excel, Google Sheets, analyses statistiques, rapports

**Structure :**
```csv
timestamp,site_name,url,success,status_code,response_time_ms,error
2024-01-15T10:25:00Z,Site A,https://site-a.com,true,200,95.00,
2024-01-15T10:24:00Z,Site B,https://site-b.com,false,500,200.00,Internal Server Error
```

**Colonnes :**
- `timestamp` : Date/heure de la vérification (RFC3339)
- `site_name` : Nom du site surveillé
- `url` : URL vérifiée
- `success` : Succès de la vérification (true/false)
- `status_code` : Code de statut HTTP
- `response_time_ms` : Temps de réponse en millisecondes
- `error` : Message d'erreur (si applicable)

### 🔸 HTML Format

**Usage :** Rapports visuels, partage avec équipes, archivage

**Fonctionnalités :**
- 🎨 **Design professionnel** avec CSS moderne
- 📊 **Statistiques visuelles** sous forme de cartes
- 📈 **Tableaux interactifs** avec données détaillées
- 📱 **Responsive design** adaptatif mobile/desktop
- 🎯 **Indicateurs de statut** colorés et iconiques
- 🔍 **Métriques par site** avec détails complets

## 📊 Statistiques incluses

Quand l'option `--stats` est activée, l'export inclut :

### Métriques globales
- Nombre total de sites
- Total des vérifications
- Vérifications réussies/échouées
- Uptime global (pourcentage)
- Temps de réponse (moyen/min/max)

### Métriques par site
- Uptime individuel par site
- Nombre de vérifications par site
- Temps de réponse moyen par site
- Dernière vérification
- Durée de monitoring

### Distributions temporelles
- Répartition par heure
- Répartition par jour
- Tendances de performance

## 🛠️ Cas d'usage

### Analyse de performance
```bash
# Rapport hebdomadaire pour analyse
site-monitor export \
  --format html \
  --since 7d \
  --stats \
  --output "weekly-report.html"
```

### Export de données pour BI
```bash
# Données pour tableau de bord BI
site-monitor export \
  --format json \
  --since 30d \
  --stats \
  --output "bi-data.json"
```

### Rapport client
```bash
# Rapport visuel pour client
site-monitor export \
  --format html \
  --site "Site Client A" \
  --since 1d \
  --stats \
  --output "rapport-client.html"
```

### Analyse en pipeline
```bash
# Analyse avec jq
site-monitor export --stdout --format json \
  | jq '.stats.overall_uptime'

# Import dans Python/pandas
site-monitor export --format csv --stdout > data.csv
python -c "import pandas as pd; df = pd.read_csv('data.csv'); print(df.describe())"
```

### Intégrations API
```bash
# Webhook vers autre système
curl -X POST "https://monitoring-system.com/api/import" \
  -H "Content-Type: application/json" \
  -d @<(site-monitor export --format json --stdout)
```

## ⚡ Performance et limites

### Optimisations
- **Index SQLite** optimisés pour les requêtes d'export
- **Streaming** pour les gros volumes de données
- **Pagination automatique** pour éviter les timeouts
- **Compression** optionnelle pour les gros fichiers

### Limites recommandées
- **JSON/CSV** : Jusqu'à 100k records par export
- **HTML** : Jusqu'à 10k records (pour performance navigateur)
- **API** : Timeout de 30s par défaut

### Conseils performance
```bash
# Pour de gros volumes, utiliser des filtres
site-monitor export --site "Site spécifique" --limit 10000

# Exports périodiques plutôt qu'un gros export
site-monitor export --since 1d --format csv >> daily-exports.csv

# Utiliser stdout pour éviter l'écriture fichier
site-monitor export --stdout | gzip > compressed-export.json.gz
```

## 🔧 Intégration dans le code

### Utilisation programmatique

```go
import (
    "site-monitor/export"
    "site-monitor/storage"
)

// Créer un exporter
exporter := export.NewExporter(storage)

// Options d'export
opts := export.ExportOptions{
    Format:       export.FormatJSON,
    SiteName:     "Mon Site",
    Since:        24 * time.Hour,
    IncludeStats: true,
}

// Exporter les données
data, err := exporter.Export(opts)
if err != nil {
    return err
}

// Formater et sauvegarder
formatter, _ := export.GetFormatter(export.FormatJSON)
file, _ := os.Create("export.json")
defer file.Close()
formatter.Format(data, file)
```

### Extension avec nouveaux formats

```go
// Implémenter l'interface Formatter
type XMLFormatter struct{}

func (f *XMLFormatter) Format(data *ExportData, writer io.Writer) error {
    // Logique de formatage XML
    return nil
}

func (f *XMLFormatter) ContentType() string {
    return "application/xml"
}

func (f *XMLFormatter) FileExtension() string {
    return ".xml"
}
```

## 🧪 Tests

### Lancer les tests
```bash
# Tests unitaires
go test ./export/... -v

# Tests avec couverture
go test ./export/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Tests d'intégration CLI
```bash
# Tester les différents formats
site-monitor export --format json --stdout > /dev/null
site-monitor export --format csv --stdout > /dev/null  
site-monitor export --format html --stdout > /dev/null

# Tester les options
site-monitor export --list-formats
site-monitor export --help
```

## 🚨 Gestion d'erreurs

### Erreurs communes et solutions

#### "Database not found"
```bash
# Solution : Lancer d'abord le monitoring
site-monitor run
# Puis dans un autre terminal
site-monitor export
```

#### "No data found"
```bash
# Vérifier la période
site-monitor export --since 168h  # 7 jours

# Vérifier les sites
site-monitor stats
```

#### "Invalid format"
```bash
# Lister les formats valides
site-monitor export --list-formats
```

#### "Permission denied" (écriture fichier)
```bash
# Utiliser stdout
site-monitor export --stdout > /tmp/export.json

# Ou spécifier un répertoire accessible
site-monitor export --output /tmp/export.json
```

## 📈 Roadmap

### v0.6.0 (Actuel)
- ✅ Export JSON, CSV, HTML
- ✅ CLI et API REST
- ✅ Statistiques intégrées
- ✅ Filtrage avancé

### v0.7.0 (Futur)
- 📊 **Format Excel** (.xlsx) natif
- 🔄 **Export programmé** (cron-like)
- 📧 **Export par email** automatique
- 🗜️ **Compression** automatique des gros exports
- 📊 **Formats spécialisés** (Prometheus, InfluxDB)

### v0.8.0 (Futur)
- 🌐 **Export cloud** (S3, GCS, Azure)
- 📋 **Templates personnalisés** pour HTML
- 🔗 **Webhooks** sur export terminé
- 📊 **Métriques P95/P99** étendues

## 🤝 Contribution

### Ajouter un nouveau format

1. **Implémenter** l'interface `Formatter`
2. **Ajouter** le format dans `GetFormatter()`
3. **Créer** les tests correspondants
4. **Documenter** le nouveau format

### Améliorer les performances

1. **Profiler** avec `go tool pprof`
2. **Optimiser** les requêtes SQL
3. **Tester** avec de gros volumes
4. **Benchmarker** les changements

## 🔗 Liens utiles

- [Configuration](../config.json) - Exemple de configuration
- [API Documentation](../web/server.go) - Endpoints détaillés
- [Storage Layer](../storage/storage.go) - Interface de stockage
- [Tests](./export_test.go) - Tests unitaires complets

## 📞 Support

Pour questions et support :

- 🐛 **Bugs** : [GitHub Issues](https://github.com/papaganelli/site-monitor/issues)
- 💡 **Demandes** : [GitHub Discussions](https://github.com/papaganelli/site-monitor/discussions)
- 📖 **Documentation** : [Wiki](https://github.com/papaganelli/site-monitor/wiki)

---

**Site Monitor Export v0.6.0** - Extraction et analyse de données de monitoring simplifiées 🚀