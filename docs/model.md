# Modèle et corrélation

## Entités

Session SMTP : nœud/service/PID/intervalle, client et bornes si observées. Transaction : tentative SMTP, enveloppe et décisions. Message : traitement et identifiants. Destinataire : actions et livraison individuelles. Queue : identité contextualisée par nœud et génération temporelle ; ne pas traiter un queue ID comme éternellement unique.

Ne pas utiliser PID seul ni Message-ID seul comme clé. Conserver événements orphelins et dossiers partiels. Valeur absente : null/unknown, jamais une valeur inventée.

## Liens explicites dans les fixtures

accepted.log : A1B2C3D4E5 → réponse LMTP avec F100000000000001 → accept avec nouvelle queue B2C3D4E5F6.

blocked.log : C3D4E5F6A7 → réponse LMTP avec F200000000000002 → action block explicite.

Le Message-ID est un contrôle de cohérence, pas la clé primaire.

## Résultats séparés

Décision SMTP, action du filtre et livraison sont des dimensions indépendantes. `status=sent` vers 127.0.0.1:10024 décrit le transfert au filtre, pas la livraison finale. `250 ... BLOCKED` est un succès du transfert interne accompagné d'un blocage métier, pas un rejet SMTP 5xx.

accepted.log : accepté par Whitelist, nouvelle queue connue, livraison aval inconnue.
blocked.log : blocage par SPAM Drop (Level 5), aucune livraison aval observée. Réponse au client SMTP d'origine absente des deux fixtures.

## JSON prévu

schema_version, identifiant interne, nœud, première/dernière observation, client, enveloppe, queues, pmg_filter_id, message_id, taille journalisée, analyse spam, actions ordonnées, transferts, résultats par destinataire, événements et provenance.

Conserver score_logged, threshold_logged et hits séparément. Ne pas remplacer le score affiché par la somme des contributions. Conserver bayes/autolearn tels que journalisés. Ne pas inventer la valeur d'un en-tête modifié ni les conditions d'une règle.

## Preuves et complétude

Chaque lien doit indiquer événement justificatif et méthode : explicit ou heuristic. Confiance et complétude distinctes. Prévoir completeness pour session, filtrage et livraison plutôt qu'un seul statut global. Rotation et timeout ne rendent pas un dossier définitivement clos ; enrichissement tardif autorisé.

## Invariants

- BLOCKED ne devient jamais delivered du fait de status=sent.
- Une action modify_header ne remplace pas l'action terminale.
- Une queue removed ne prouve pas à elle seule une livraison.
- Absence de trace antivirus ne signifie pas antivirus clean.
- Les tests SA observés ne reconstruisent pas toute la configuration PMG.
