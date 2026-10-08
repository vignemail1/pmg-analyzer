# Architecture

## Périmètre et état

Les composants ci-dessous sont une spécification d'implémentation. Seule la CLI `version` est implémentée dans ce premier commit.

## Flux

Sources → découverte → lecture → stockage brut/checkpoint → parsing → événements → corrélation → dossiers → recherche/export/API.

## Composants

### CLI et configuration

CLI d'import, suivi, recherche et export à construire. Configuration fichier, surcharges par variables d'environnement et options explicites. Validation au démarrage ; une configuration invalide ne doit pas lancer une collecte partielle silencieuse.

### Collecteur

Découvre les fichiers actifs et archives. Suit des générations de fichiers, pas uniquement leurs chemins. Enregistre chaque ligne et sa provenance avant d'avancer son checkpoint. Voir collection.md.

### Stockage

SQLite envisagé pour le MVP. Tables prévues : sources, générations, checkpoints, lignes brutes, événements, sessions, transactions, messages, traitements de filtre, queues, destinataires et liens. Lignes et checkpoint dans la même transaction. Migrations versionnées ; contraintes d'unicité pour les relectures d'une même génération. Interface de stockage remplaçable.

### Parseurs

Parseur d'enveloppe syslog RFC3339 puis parseurs Postfix et PMG. Chaque parseur produit des événements typés sans décider de leur appartenance à un dossier. Conserver texte brut, champs observés, version du parseur et erreurs. Une ligne inconnue reste stockée. Prévoir fuzzing et limites de taille configurables ; une ligne dépassant la limite déclenche une erreur visible, jamais une troncature silencieuse.

### Corrélateur

Privilégie les liens explicites entre queue et identifiant PMG. Corrélation heuristique séparée et signalée. État persistant, indépendant de la rotation ; dossiers enrichissables par des événements tardifs. Voir model.md.

### Recherche et export

Filtres structurés par période, nœud, IP/CIDR, enveloppe, identifiants, règle, étape et résultat. JSON consolidé par dossier et JSONL par événement. Pagination et bornes de requêtes à prévoir. Ne pas confondre expéditeur d'enveloppe et en-tête From.

### API et interface web

Phase ultérieure : consultation de la chronologie et des preuves. Lecture seule, authentification, autorisation et audit des exports avant exposition réseau. Aucune action de remise en queue ou modification PMG.

### Supervision

Retard de collecte, dernière découverte/transaction réussie, erreurs de lecture/compression, archives en attente, lignes inconnues, événements orphelins et espace disque. Les échecs de stockage interdisent l'avancement du checkpoint.

## Confidentialité

Droits minimaux pour lire les sources ; stockage non accessible publiquement. Rétention distincte pour brut et données normalisées avec impact documenté sur les preuves et le retraitement. Ne pas journaliser des contenus sensibles inutilement. Ne pas exposer la base directement.

## Phases

1. Collecte durable et reprise testée.
2. Parseurs et golden tests.
3. Corrélation des queues et actions par destinataire.
4. Recherche CLI et exports.
5. API/web, statistiques et alertes.
