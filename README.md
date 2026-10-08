# pmg-analyzer

Analyse et corrélation des logs Postfix et Proxmox Mail Gateway, à partir des fichiers mail.log et de leurs rotations.

## Source primaire et périmètre

Les fichiers mail.log, mail.log.1 et les archives compressées sont les données primaires. Import historique, suivi, parsing, corrélation, recherche et export doivent fonctionner sans accès réseau à PMG, sans identifiants API et sans dépendance au Tracking Center.

L'accès à l'API PMG est une option d'enrichissement pour une phase ultérieure, désactivée par défaut et hors MVP. Il ne remplace pas les preuves issues des logs. L'API de consultation propre à pmg-analyzer, envisagée pour l'interface web, est un composant distinct.

## État du projet

Initialisation : CLI de version et fixtures disponibles. Collecte, parsing, stockage, corrélation, recherche et API sont spécifiés mais ne sont pas encore implémentés. Aucune garantie de collecte en production à ce stade.

## Prérequis et commandes

Go 1.25 ou ultérieur.

```sh
go test -race ./...
go vet ./...
go run ./cmd/pmg-analyzer version
go build -o pmg-analyzer ./cmd/pmg-analyzer
```

## Documentation

- [Fonctionnalités, périmètre et workflows prévus](docs/features-workflows.md)
- [Architecture et fonctionnement](docs/architecture.md)
- [Collecte et rotation](docs/collection.md)
- [Modèle et corrélation](docs/model.md)
- [Format JSON prévu](docs/json-format.md)
- [Cas JSON supplémentaires](docs/json-additional-cases.md)
- [Développement et validation](docs/development.md)
- [Fixtures anonymisées](testdata/README.md)

Les logs bruts contiennent potentiellement des données personnelles. Ne pas publier de logs de production. Le projet ne modifie pas les politiques PMG et ne libère pas de quarantaine.
