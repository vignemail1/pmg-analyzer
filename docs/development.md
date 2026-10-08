# Développement

## Outils

Go 1.25 ; bibliothèque standard pour l'initialisation. Pas encore de dépendance SQLite ni de framework web. Les commandes du README et la CI vérifient format, vet, tests race et build. Leur présence ne signifie pas qu'une exécution a déjà réussi.

## Structure cible

cmd/pmg-analyzer : interface CLI.
internal/collector : découverte et lecteurs.
internal/store : transactions et migrations.
internal/parser : événements normalisés.
internal/correlator : liens et dossiers.
internal/search : requêtes et export.
internal/api : consultation future.

Créer ces packages avec leur implémentation, pas sous forme de composants factices.

## Tests à construire

Golden tests des fixtures avec résultats attendus ; tests de reprise et idempotence ; scénarios multi-destinataires, plusieurs transactions/session, PIDs réutilisés, queues liées, livraison différée puis réussie, ordre d'arrivée inversé, journaux tronqués et inconnus. Fuzzing des parseurs, benchmarks et tests de concurrence.

## Livraison

Migrations versionnées et tests de compatibilité avant modification du schéma. Pas de promesse production avant couverture de collecte/rotation et validation sur un corpus représentatif. Ajouter une licence uniquement après choix explicite du propriétaire.
