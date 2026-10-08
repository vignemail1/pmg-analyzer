# Format JSON prévu

## Statut

Contrat proposé version 1.0, non encore implémenté par la CLI. Les exemples sont des projections documentaires des fixtures testdata/accepted.log et testdata/blocked.log, pas des sorties générées ni des golden tests exécutés. Ce document complète model.md.

## Formats de fichiers

- Dossier consolidé : un objet JSON par fichier, avec extension .json.
- Export de plusieurs dossiers : JSONL, un objet complet par ligne.
- Export d'événements : JSONL distinct, un événement par ligne ; ne pas mélanger événements et dossiers dans un même flux.
- Encodage UTF-8. Dates RFC3339 avec précision et décalage disponibles dans la source ; comparaisons chronologiques sur les instants normalisés.

## Champs du dossier

| Champ | Type | Sens |
|---|---|---|
| schema_version | string | Version du contrat, ici 1.0 |
| node | string | Nœud journalisant |
| first_seen / last_seen | string | Première et dernière observation du dossier, pas bornes de connexion |
| client | object | ip, hostname_logged et helo : string ou null |
| identifiers | object | initial_queue_id et pmg_filter_id : string ou null ; post_filter_queue_ids : tableau de strings ; message_id : string ou null |
| envelope | object | mail_from : string ou null ; recipient_count_logged : entier ou null ; recipients : tableau |
| size_bytes_logged | integer ou null | Taille journalisée, en octets |
| spam | object ou null | Analyse SA observée |
| filter | object ou null | Temps de traitement et actions ordonnées |
| filter_handoffs | array | Transferts internes vers le filtre, par queue et destinataire |
| outcome | object | Synthèse du filtre et de la livraison ; mixed autorisé pour plusieurs destinataires |
| correlation | object | Méthode, preuves, confiance et complétude par dimension |

Tous les champs de cette projection sont présents. null signifie donnée non observée ou non applicable ; un tableau vide signifie aucune entrée observée, pas preuve de l'absence universelle d'événements. unknown est un état métier indéterminé. L'expéditeur vide SMTP doit être représenté par une chaîne vide, et non null. Aucun champ absent ne doit être inventé.

### Destinataires

Chaque entrée contient address, filter_action, filter_rule, post_filter_queue_id et delivery_status. Les actions et livraisons restent individuelles. États de livraison : unknown, pending, deferred, sent, bounced ou not_attempted_due_to_filter_block. Ce dernier décrit le chemin bloqué observé, pas l'impossibilité qu'une autre tentative distincte ait existé.

### Analyse antispam

score_logged et threshold_logged sont des nombres tels que journalisés. analysis_time_seconds est un nombre. bayes_logged et autolearn_logged conservent leurs chaînes originales. hits est un tableau de {name, score}. Ne pas remplacer le score global journalisé par la somme des contributions ni déduire les conditions complètes d'une règle.

### Actions

Une action contient type, recipient, rule et source_lines. Pour modify_header : header et value, avec value=null si non journalisée. Pour accept : post_filter_queue_id. Les actions sont dans l'ordre observé ; modify_header ne remplace pas accept ou block.

### Transferts

Chaque entrée contient queue_id, recipient, protocol, relay_logged, postfix_status, smtp_code, enhanced_status_code, response_text, delay_seconds_logged et initial_queue_removed. status=sent vers le filtre ne constitue pas une preuve de livraison finale. Une réponse 250 BLOCKED n'est pas un rejet SMTP permanent 5xx.

### Corrélation et provenance

confidence : high, medium ou low ; method : explicit ou heuristic. proofs référence le fichier et les numéros de lignes, à partir de 1. completeness sépare session, filtering et delivery : partial, complete ou not_applicable. complete est relatif à la phase et au destinataire observés, pas une preuve d'exhaustivité de tous les logs.

La projection finale pourra ajouter des identifiants internes stables et des événements normalisés. Chaque événement devra conserver type, timestamp, node, service, pid, identifiants observés, attributs, version du parseur et référence à la ligne brute. Les exemples ci-dessous n'inventent pas ces identifiants non encore générés. Les lignes citées restent consultables dans la fixture.

## Cas accepté : accepted.log

Le lien initial queue → PMG vient de la ligne 9 ; le lien PMG → nouvelle queue vient de la ligne 7. La livraison aval n'est pas présente.

```json
{
  "schema_version": "1.0",
  "node": "pmg-example",
  "first_seen": "2026-01-03T21:06:25.559770+02:00",
  "last_seen": "2026-01-03T21:06:27.230683+02:00",
  "client": {"ip": "2001:db8::12", "hostname_logged": "mx.example.net", "helo": null},
  "identifiers": {
    "initial_queue_id": "A1B2C3D4E5",
    "pmg_filter_id": "F100000000000001",
    "post_filter_queue_ids": ["B2C3D4E5F6"],
    "message_id": "<accepted-001@sender.example.net>"
  },
  "envelope": {
    "mail_from": "sender@example.net",
    "recipient_count_logged": 1,
    "recipients": [{"address": "recipient@example.org", "filter_action": "accept", "filter_rule": "Whitelist", "post_filter_queue_id": "B2C3D4E5F6", "delivery_status": "unknown"}]
  },
  "size_bytes_logged": 36751,
  "spam": {
    "score_logged": 1,
    "threshold_logged": 5,
    "analysis_time_seconds": 1.523,
    "bayes_logged": "undefined",
    "autolearn_logged": "disabled",
    "hits": [
      {"name": "DMARC_NONE", "score": 0.1},
      {"name": "HTML_EXTRA_CLOSE", "score": 0.001},
      {"name": "HTML_MESSAGE", "score": 0.001},
      {"name": "KAM_DMARC_NONE", "score": 0.25},
      {"name": "KAM_DMARC_STATUS", "score": 0.01},
      {"name": "KAM_LAZY_DOMAIN_SECURITY", "score": 1},
      {"name": "SPF_HELO_PASS", "score": -0.001},
      {"name": "SPF_NONE", "score": 0.001}
    ]
  },
  "filter": {
    "processing_time_seconds": 1.619,
    "actions": [
      {"type": "modify_header", "header": "X-SPAM-LEVEL", "value": null, "recipient": "recipient@example.org", "rule": "Modify Header", "source_lines": [6]},
      {"type": "accept", "recipient": "recipient@example.org", "rule": "Whitelist", "post_filter_queue_id": "B2C3D4E5F6", "source_lines": [7]}
    ]
  },
  "filter_handoffs": [{
    "queue_id": "A1B2C3D4E5",
    "recipient": "recipient@example.org",
    "protocol": "lmtp",
    "relay_logged": "127.0.0.1[127.0.0.1]:10024",
    "postfix_status": "sent",
    "smtp_code": 250,
    "enhanced_status_code": "2.5.0",
    "response_text": "OK (F100000000000001)",
    "delay_seconds_logged": 1.7,
    "initial_queue_removed": true
  }],
  "outcome": {"filter_status": "accepted", "delivery_status": "unknown"},
  "correlation": {
    "confidence": "high",
    "method": "explicit",
    "proofs": [{"file": "testdata/accepted.log", "lines": [7, 9]}],
    "completeness": {"session": "partial", "filtering": "complete", "delivery": "partial"},
    "warnings": ["smtp_session_boundaries_missing", "downstream_delivery_missing"]
  }
}
```

## Cas bloqué : blocked.log

Le lien initial queue → PMG vient de la ligne 9 ; la décision block est explicite ligne 7. Le transfert interne réussit mais le filtre bloque le destinataire.

```json
{
  "schema_version": "1.0",
  "node": "pmg-example",
  "first_seen": "2026-01-01T10:00:01.274254+02:00",
  "last_seen": "2026-01-01T10:00:04.919414+02:00",
  "client": {"ip": "192.0.2.42", "hostname_logged": "mx-out.example.net", "helo": null},
  "identifiers": {
    "initial_queue_id": "C3D4E5F6A7",
    "pmg_filter_id": "F200000000000002",
    "post_filter_queue_ids": [],
    "message_id": "<blocked-001@sender.example.net>"
  },
  "envelope": {
    "mail_from": "sender@example.net",
    "recipient_count_logged": 1,
    "recipients": [{"address": "list@example.org", "filter_action": "block", "filter_rule": "SPAM Drop (Level 5)", "post_filter_queue_id": null, "delivery_status": "not_attempted_due_to_filter_block"}]
  },
  "size_bytes_logged": 444601,
  "spam": {
    "score_logged": 8,
    "threshold_logged": 5,
    "analysis_time_seconds": 3.151,
    "bayes_logged": "undefined",
    "autolearn_logged": "disabled",
    "hits": [
      {"name": "DKIM_INVALID", "score": 0.1},
      {"name": "DKIM_SIGNED", "score": 0.1},
      {"name": "DMARC_REJECT", "score": 0.1},
      {"name": "GB_GEN_REDIR_URL", "score": 0.5},
      {"name": "HTML_IMAGE_RATIO_08", "score": 0.001},
      {"name": "HTML_MESSAGE", "score": 0.001},
      {"name": "HTTPS_HTTP_MISMATCH", "score": 0.1},
      {"name": "KAM_DMARC_REJECT", "score": 7},
      {"name": "KAM_DMARC_STATUS", "score": 0.01},
      {"name": "KAM_SHORT", "score": 0.001},
      {"name": "RCVD_IN_DNSWL_NONE", "score": -0.0001},
      {"name": "SPF_HELO_PASS", "score": -0.001},
      {"name": "SPF_SOFTFAIL", "score": 0.972}
    ]
  },
  "filter": {
    "processing_time_seconds": 3.523,
    "actions": [
      {"type": "modify_header", "header": "X-SPAM-LEVEL", "value": null, "recipient": "list@example.org", "rule": "Modify Header", "source_lines": [6]},
      {"type": "block", "recipient": "list@example.org", "rule": "SPAM Drop (Level 5)", "source_lines": [7]}
    ]
  },
  "filter_handoffs": [{
    "queue_id": "C3D4E5F6A7",
    "recipient": "list@example.org",
    "protocol": "lmtp",
    "relay_logged": "127.0.0.1[127.0.0.1]:10024",
    "postfix_status": "sent",
    "smtp_code": 250,
    "enhanced_status_code": "2.7.0",
    "response_text": "BLOCKED (F200000000000002)",
    "delay_seconds_logged": 3.7,
    "initial_queue_removed": true
  }],
  "outcome": {"filter_status": "blocked", "delivery_status": "not_attempted_due_to_filter_block"},
  "correlation": {
    "confidence": "high",
    "method": "explicit",
    "proofs": [{"file": "testdata/blocked.log", "lines": [7, 9]}],
    "completeness": {"session": "partial", "filtering": "complete", "delivery": "not_applicable"},
    "warnings": ["smtp_session_boundaries_missing", "original_client_smtp_response_missing"]
  }
}
```

## Évolution et cas non couverts

Les deux fixtures ne couvrent ni rejet NOQUEUE, ni quarantaine, ni livraison aval, ni destinataires à résultats mixtes. Ne pas fabriquer de sorties de référence pour ces cas sans traces représentatives. Les traces tardives enrichissent les dossiers ; toute livraison doit être rattachée à la bonne queue et au bon destinataire. Changements incompatibles : nouvelle version majeure du contrat. Ajouts compatibles : champs supplémentaires ; consommateurs tolérants aux champs inconnus. Un schéma JSON formel et des tests de conformité seront ajoutés avec le modèle Go.
