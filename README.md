# pmg-analyzer

Analyse et corrélation des logs Postfix et Proxmox Mail Gateway.

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

- [Architecture et fonctionnement](docs/architecture.md)
- [Collecte et rotation](docs/collection.md)
- [Modèle et corrélation](docs/model.md)
- [Développement et validation](docs/development.md)
- [Fixtures anonymisées](testdata/README.md)

Les logs bruts contiennent potentiellement des données personnelles. Ne pas publier de logs de production. Le projet ne modifie pas les politiques PMG et ne libère pas de quarantaine.
