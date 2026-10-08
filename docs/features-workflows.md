# Fonctionnalités et workflows prévus

## Statut et décision de périmètre

Spécification, pas liste de fonctionnalités déjà opérationnelles. Seule la CLI version est implémentée à ce stade. Go 1.25 est la version minimale du projet.

Source primaire : fichiers mail.log et rotations, y compris GZIP. Le cœur fonctionne hors ligne à partir de fichiers locaux ou copiés, sans API PMG ni Tracking Center. Une copie doit être identifiée et rapprochée des sources existantes pour éviter les doublons d'import. L'intégration PMG API est facultative, ultérieure et hors MVP.

Ce document précise le périmètre des autres documents. La mention API/web dans architecture.md désigne l'API propre à pmg-analyzer ; elle ne constitue pas une dépendance à l'API PMG. D'autres sources de logs pourront être ajoutées explicitement, sans collecte parallèle implicite de syslog qui dupliquerait mail.log.

## Fonctionnalités prévues

| Fonctionnalité | Phase | Résultat attendu |
|---|---|---|
| Import historique de fichiers et archives | MVP | Lignes brutes persistées avec provenance et reprise |
| Suivi continu des fichiers actifs | MVP | Capture incrémentale indépendante du parsing |
| Rotation, compression et checkpoints | MVP | Reprise testée ; pertes connues et ambiguïtés signalées |
| Parsing Postfix et PMG | MVP | Événements typés, champs observés, lignes inconnues conservées |
| Reconstruction des sessions et transactions | MVP | Séparation malgré PID réutilisé ; dossiers partiels autorisés |
| Corrélation queues/traitements/destinataires | MVP | Liens explicites documentés ; heuristiques signalées |
| Recherche CLI structurée et consultation | MVP | Chronologie, résultats et preuves brutes |
| Export JSON et JSONL | MVP | Dossiers et événements exportés dans des flux distincts |
| Retraitement des logs déjà collectés | MVP | Nouveau parsing sans dépendance aux fichiers source encore présents |
| Supervision de collecte et couverture | MVP | Retard, erreurs, lignes inconnues et événements orphelins visibles |
| API locale de consultation et interface web | Après MVP | Consultation du même stockage, en lecture seule |
| Statistiques et alertes | Après MVP | Comptages distincts par session, transaction, message et destinataire |
| Agent multi-nœuds et transmission centralisée | Évolution | Livraison durable, authentifiée, avec provenance du nœud |
| Enrichissement via API PMG | Option ultérieure | Métadonnées supplémentaires, séparées des observations de logs |

MVP désigne une fonctionnalité à implémenter et valider, pas une fonctionnalité déjà disponible. Le support de syntaxe doit être attesté par les fixtures et les tests ; une action inconnue n'est pas silencieusement convertie en acceptation.

## Workflow 1 : import historique hors ligne

1. Choisir les fichiers et déclarer le nœud/source logique ainsi que les paramètres temporels nécessaires.
2. Valider configuration, permissions et capacité du stockage avant import.
3. Découvrir les sources et leurs générations ; détecter les fichiers déjà importés et les copies ambiguës.
4. Lire les lignes et enregistrer brut/provenance/checkpoint dans la même transaction.
5. Parser les données persistées ; conserver les lignes non reconnues et les erreurs.
6. Corréler les événements par preuves, indépendamment de l'ordre d'import des fichiers.
7. Publier les dossiers consultables et un bilan d'import : lectures, erreurs, ambiguïtés et couverture.

La fin d'un fichier ou d'un import n'est pas la preuve de la fin de vie d'un message. Les événements tardifs peuvent enrichir un dossier. Les fichiers source ne sont jamais modifiés ni supprimés par l'outil.

## Workflow 2 : suivi continu et rotation

1. Reprendre les checkpoints persistants ; découvrir l'historique restant et le fichier actif.
2. Suivre simultanément le fichier actif et les générations tournées encore pertinentes.
3. Lors d'un renommage, finir l'ancien lecteur et ouvrir le nouveau sans identifier un fichier par son chemin seul.
4. Utiliser les notifications comme optimisation et les scans périodiques comme réconciliation.
5. Persister les lignes avant d'avancer le checkpoint ; ne pas valider une ligne partiellement écrite.
6. Rapprocher les archives compressées des générations déjà lues lorsque la continuité est démontrable.
7. En cas de crash, relire les lots non validés ; en cas d'échec de stockage, arrêter l'avancement et alerter.

Pas de promesse zéro perte : copytruncate, suppression précoce des sources, archives illisibles et indisponibilité prolongée peuvent rendre des données irrécupérables. Voir collection.md pour les mécanismes et tests.

## Workflow 3 : parsing et reconstruction

1. Extraire timestamp, nœud, service, PID et message de la ligne.
2. Produire un événement normalisé avec sa version de parseur et sa provenance.
3. Délimiter les sessions SMTP par bornes observées ; ne jamais fusionner sur le PID seul.
4. Relier la queue initiale au traitement PMG via les identifiants explicites des réponses LMTP.
5. Relier les actions accept aux queues après filtrage et à leurs destinataires.
6. Ajouter actions ordonnées, analyses, transferts et résultats de transport.
7. Conserver indépendamment confiance, complétude et avertissements.

NOQUEUE n'est pas une queue. Un échec STARTTLS peut ne contenir aucun message. Message-ID est un attribut de recherche et de cohérence, pas une clé primaire universelle.

## Workflow 4 : investigation et recherche

1. Rechercher par période/nœud, IP ou CIDR, enveloppe, identifiants, règle, code SMTP, étape ou résultat.
2. Ouvrir le dossier et lire sa chronologie.
3. Examiner chaque destinataire, ses éventuelles réécritures et ses résultats.
4. Consulter les lignes brutes justifiant les décisions et liens.
5. Lire les informations absentes ou ambiguës avant de conclure.
6. Exporter le dossier ou les événements sélectionnés, avec provenance et avertissements.

Ne pas assimiler status=sent vers le filtre à une livraison finale. Une réponse SMTP aval positive prouve une acceptation par le prochain relais, pas le dépôt en boîte. Une action de quarantaine doit être observée explicitement ; un hit SA nommé QUARANTINE ne suffit pas.

## Workflow 5 : retraitement et évolution du modèle

1. Conserver brut et versions de parsing suffisamment longtemps selon la rétention définie.
2. Exécuter le nouveau parseur sur les données persistées, avec un identifiant de traitement/version.
3. Reconstruire les événements et dossiers dérivés de manière idempotente, sans doubler les lignes brutes.
4. Comparer couverture, dossiers et invariants avec les fixtures de référence.
5. Publier les vues mises à jour seulement après succès ; une interruption ne doit pas exposer une reconstruction partielle comme complète.

Le mécanisme de publication/versionnement sera fixé avec le stockage. Une rétention ayant supprimé le brut limite le retraitement ; cette limitation doit être visible.

## Workflow 6 : enrichissement optionnel PMG API

Phase ultérieure uniquement. Activation explicite, secrets protégés, délais et erreurs bornés. Le cœur de collecte/recherche doit continuer en cas d'indisponibilité ou d'absence d'API.

Toute donnée enrichie conservera origine, date de récupération et méthode de rattachement. La configuration actuelle d'une règle ne prouve pas sa configuration lors d'un événement historique. L'enrichissement ne doit pas écraser les preuves brutes ni inventer des liens. Les ressources exactes de l'API seront définies après vérification des versions PMG prises en charge.

Aucune modification de politique, libération de quarantaine ou remise en queue n'est prévue dans ce périmètre.

## Configuration et exploitation

Configuration par fichier avec surcharges d'environnement et options explicites ; noms d'options à fixer avec l'implémentation. Paramètres prévus : sources, nœud, stockage, scan, lots, limites de lignes, fuseau/année si absents des logs, rétention et niveau de journalisation. Aucune URL ou credential PMG requis pour le MVP.

Droits minimaux, stockage non public, contrôles d'accès avant exposition web, audit des exports, pagination et limites de requêtes. La CLI minimale existante ne fournit pas encore import, follow, search ou reprocess ; ces noms sont des intentions, pas des commandes utilisables.

## Corpus et critères de validation

Fixtures disponibles : acceptation/blocage par filtre, rejets temporaires NOQUEUE, échec STARTTLS, acceptation aval, réécriture de destinataire et report aval multi-relais. Voir json-format.md et json-additional-cases.md.

À compléter : rejet permanent, quarantaine explicite, libération, report puis succès, échec permanent aval et résultats mixtes. Une liste de fonctionnalités ne remplace pas ces cas de validation.

Avant déclaration d'un MVP exploitable : tests de rotation/reprise/idempotence, golden tests du parsing et de la corrélation, imports dans des ordres différents, tests de concurrence et erreurs de stockage, migrations versionnées, invariants de résultat et documentation des limites. Les performances seront mesurées sur un corpus représentatif ; aucun débit ou dimensionnement n'est garanti à ce stade.
