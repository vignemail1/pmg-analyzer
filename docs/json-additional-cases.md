# Cas supplémentaires anonymisés et projections JSON

## Statut et périmètre

Complément de json-format.md. Exemples documentaires, non générés par un parseur et non validés par des golden tests. Les objets ci-dessous sont des projections ciblées des nouveaux champs, pas des dossiers complets conformes au contrat 1.0. Extension 1.1 proposée, à finaliser avec le modèle Go et le schéma JSON. Aucun ancien fichier n'est modifié.

Le premier bloc fourni est réduit à ses deux rejets RCPT et à l'échec STARTTLS, sans prétendre en conserver toutes les sessions. Les autres sessions acceptées de ce bloc ne sont pas intégrées dans ce lot. Les préfixes de recherche fichier/numéro sont retirés : ils ne sont pas du syslog. Les numéros source_lines désignent désormais les lignes des fixtures, à partir de 1.

Nœud, hôtes, adresses, IP publiques, dates, PIDs et identifiants sont remplacés. Les dates sont décalées de manière cohérente avec les fixtures existantes. IP de documentation et domaines example.net/example.org ; loopback et localhost conservés. Scores, règles, tailles, compteurs, durées, ciphers et erreurs techniques sont conservés. Aucun dictionnaire contenant les valeurs originales n'est publié. Les tailles et la chronologie détaillée restent des attributs techniques : il ne s'agit pas d'une anonymisation statistique absolue.

## Extensions du modèle

- sessions : une entrée par connexion, même quand le PID est réutilisé. first_seen n'est pas connected_at ; connected_at=null lorsque connect est absent.
- smtp_decisions : étape, code SMTP, code enrichi, réponse et enveloppe observée ; expéditeur vide représenté par "".
- tls : négociation observée, détails présents uniquement ; failed ne doit pas être traité comme un rejet de mail par règle.
- command_counts : succeeded et attempted, avec attempted=null pour un compteur sans dénominateur explicite.
- queues : tailles et enveloppes par génération de queue, sans écraser les valeurs avant filtrage.
- deliveries : destinataire to distinct de orig_to ; résultat par transport. sent signifie acceptation par le relais suivant, pas preuve de dépôt en boîte.
- remote_queue_id : identifiant attribué par un autre serveur, toujours contextualisé par le relais ; ne pas le confondre avec une queue locale PMG.
- transport_events : liste ordonnée des réponses et erreurs de chaque relais ; une séquence multi-relais ne constitue pas automatiquement plusieurs cycles de retry de la queue.
- final_outcome_known : distinguer le dernier état observé du devenir définitif.

## Deux rejets temporaires et PID réutilisé

Fixture : ../testdata/noqueue-temporary.log. Deux connexions distinctes, un même PID. La première ligne connect est absente, mais TLS et disconnect sont observés. La seconde connexion a ses deux bornes. Aucun queue ID, Message-ID, taille, score ou action PMG ne doit être inventé. Le motif service unavailable ne prouve pas à lui seul quel sous-système est indisponible.

```json
{
  "schema_version": "1.1-draft",
  "node": "pmg-example",
  "sessions": [
    {
      "pid": 300101,
      "client": {"ip": "192.0.2.101", "hostname_logged": "mx-a.example.net", "helo": "mx-a.example.net"},
      "first_seen": "2026-01-06T11:59:33.212040+02:00",
      "connected_at": null,
      "disconnected_at": "2026-01-06T11:59:33.907590+02:00",
      "tls": {"status": "established", "protocol": "TLSv1.3", "cipher": "TLS_AES_128_GCM_SHA256", "key_exchange": "x25519", "server_signature": "RSA-PSS", "server_signature_bits": 2048, "server_digest": "SHA256", "source_lines": [1]},
      "smtp_decisions": [{"stage": "rcpt_to", "decision": "temporary_reject", "smtp_code": 450, "enhanced_status_code": "4.7.1", "mail_from": "", "recipient": "list-owner@example.org", "response_text": "<list-owner@example.org>: Recipient address rejected: Service is unavailable (try later)", "source_lines": [2]}],
      "command_counts": {"ehlo": {"succeeded": 2, "attempted": null}, "starttls": {"succeeded": 1, "attempted": null}, "mail": {"succeeded": 1, "attempted": null}, "rcpt": {"succeeded": 0, "attempted": 1}, "quit": {"succeeded": 1, "attempted": null}, "commands": {"succeeded": 5, "attempted": 6}},
      "queue_ids": [],
      "warnings": ["connect_line_missing"]
    },
    {
      "pid": 300101,
      "client": {"ip": "192.0.2.102", "hostname_logged": "mx-b.example.net", "helo": "helo-b.example.net"},
      "first_seen": "2026-01-06T11:59:36.648587+02:00",
      "connected_at": "2026-01-06T11:59:36.648587+02:00",
      "disconnected_at": "2026-01-06T11:59:36.727055+02:00",
      "tls": {"status": "established", "protocol": "TLSv1.2", "cipher": "ECDHE-RSA-AES256-GCM-SHA384", "key_exchange": null, "source_lines": [5]},
      "smtp_decisions": [{"stage": "rcpt_to", "decision": "temporary_reject", "smtp_code": 450, "enhanced_status_code": "4.7.1", "mail_from": "", "recipient": "list-owner@example.org", "response_text": "<list-owner@example.org>: Recipient address rejected: Service is unavailable (try later)", "source_lines": [6]}],
      "command_counts": {"ehlo": {"succeeded": 2, "attempted": null}, "starttls": {"succeeded": 1, "attempted": null}, "mail": {"succeeded": 1, "attempted": null}, "rcpt": {"succeeded": 0, "attempted": 1}, "data": {"succeeded": 0, "attempted": 1}, "rset": {"succeeded": 1, "attempted": null}, "quit": {"succeeded": 1, "attempted": null}, "commands": {"succeeded": 6, "attempted": 8}},
      "queue_ids": [],
      "warnings": []
    }
  ]
}
```

## Échec STARTTLS sans message observé

Fixture : ../testdata/starttls-failure.log. L'alerte bad certificate est une preuve d'échec TLS, pas une preuve suffisante de sa cause exacte ou de l'identité du certificat en cause. Pas de MAIL FROM ni RCPT observés. Ne pas fabriquer une transaction de message ou un code SMTP de rejet.

```json
{
  "schema_version": "1.1-draft",
  "node": "pmg-example",
  "session": {
    "pid": 300101,
    "client": {"ip": "192.0.2.103", "hostname_logged": "mx-tls.example.net", "helo": null},
    "connected_at": "2026-01-06T12:03:51.151201+02:00",
    "disconnected_at": "2026-01-06T12:03:51.678884+02:00",
    "tls": {"status": "failed", "protocol": null, "cipher": null, "error_text": "ssl/tls alert bad certificate", "alert_number": 42, "source_lines": [2, 3]},
    "connection_outcome": "lost_after_starttls",
    "smtp_decisions": [],
    "queue_ids": [],
    "mail_from": null,
    "recipients": [],
    "command_counts": {"ehlo": {"succeeded": 1, "attempted": null}, "starttls": {"succeeded": 0, "attempted": 1}, "commands": {"succeeded": 1, "attempted": 2}},
    "proof_lines": [1, 2, 3, 4, 5]
  }
}
```

## Acceptation aval, expéditeur vide, deux modifications d'en-tête

Fixture : ../testdata/accepted-downstream.log. DMARC_QUAR et KAM_DMARC_QUARANTINE sont des hits SA, pas des actions de quarantaine. L'action terminale observée est accept/default-accept. Le score 4 doit rester la valeur journalisée. Les tailles des deux queues sont différentes ; elles ne doivent pas s'écraser.

```json
{
  "schema_version": "1.1-draft",
  "node": "pmg-example",
  "identifiers": {"initial_queue_id": "D4E5F6A7B8", "pmg_filter_id": "F300000000000003", "post_filter_queue_ids": ["E5F6A7B8C9"], "message_id": "<downstream-001@sender.example.net>"},
  "queues": [
    {"queue_id": "D4E5F6A7B8", "mail_from": "", "size_bytes_logged": 17129, "removed": true, "source_lines": [3, 14]},
    {"queue_id": "E5F6A7B8C9", "mail_from": "", "size_bytes_logged": 18070, "removed": true, "source_lines": [10, 16]}
  ],
  "filter": {
    "actions": [
      {"type": "modify_header", "header": "X-SPAM-LEVEL", "value": null, "recipient": "list-owner@example.org", "rule": "Modify Header", "source_lines": [6]},
      {"type": "modify_header", "header": "X-Spam-Status", "value": null, "recipient": "list-owner@example.org", "rule": "SPAM Tag Header (Level 3)", "source_lines": [7]},
      {"type": "accept", "recipient": "list-owner@example.org", "rule": "default-accept", "post_filter_queue_id": "E5F6A7B8C9", "source_lines": [11]}
    ]
  },
  "deliveries": [{
    "queue_id": "E5F6A7B8C9", "protocol": "smtp", "to": "list-owner@example.org", "orig_to": null,
    "relay": {"hostname_logged": "relay-list.example.org", "ip": "192.0.2.105", "port": 25},
    "postfix_status": "sent", "smtp_code": 250, "enhanced_status_code": "2.0.0",
    "response_text": "Ok: queued as A100000001", "remote_queue_id": "A100000001",
    "delay_seconds_logged": 0.15, "delays_seconds_logged": [0.04, 0, 0.07, 0.03],
    "acceptance_scope": "next_hop", "mailbox_delivery_confirmed": null, "source_lines": [15]
  }],
  "outcome": {"filter_status": "accepted", "delivery_status": "sent", "acceptance_scope": "next_hop"},
  "correlation": {"method": "explicit", "proof_lines": [11, 13, 15], "session_completeness": "partial"}
}
```

Le client de réinjection localhost ne remplace pas le client original mx-c.example.net. orig_client est une observation distincte à conserver. Le JSON ci-dessus n'affiche qu'une projection des champs nouveaux, pas l'analyse SA complète présente dans le fichier.

## Complément de la fixture accepted.log : destinataire réécrit

Fixture : ../testdata/accepted-downstream-rewritten.log à traiter avec ../testdata/accepted.log. Le mapping préexistant est conservé : B2C3D4E5F6 et accepted-001. Les deux fichiers se chevauchent dans le temps ; l'ordre des événements doit être reconstruit, pas déduit de l'ordre d'import. orig_to et to différents ne prouvent pas plusieurs destinataires ni l'origine exacte de la réécriture.

```json
{
  "schema_version": "1.1-draft",
  "node": "pmg-example",
  "identifiers": {"initial_queue_id": "A1B2C3D4E5", "pmg_filter_id": "F100000000000001", "post_filter_queue_ids": ["B2C3D4E5F6"], "message_id": "<accepted-001@sender.example.net>"},
  "queues": [
    {"queue_id": "A1B2C3D4E5", "size_bytes_logged": 36751, "source": {"file": "testdata/accepted.log", "lines": [3]}},
    {"queue_id": "B2C3D4E5F6", "size_bytes_logged": 37633, "source": {"file": "testdata/accepted-downstream-rewritten.log", "lines": [3]}}
  ],
  "deliveries": [{
    "queue_id": "B2C3D4E5F6", "protocol": "smtp", "to": "target@example.org", "orig_to": "recipient@example.org",
    "relay": {"hostname_logged": "mailbox.example.org", "ip": "2001:db8:1::10", "port": 25},
    "postfix_status": "sent", "smtp_code": 250, "enhanced_status_code": "2.0.0",
    "response_text": "Ok: queued as A200000002", "remote_queue_id": "A200000002",
    "delay_seconds_logged": 0.18, "delays_seconds_logged": [0.05, 0, 0.07, 0.06],
    "acceptance_scope": "next_hop", "mailbox_delivery_confirmed": null,
    "source": {"file": "testdata/accepted-downstream-rewritten.log", "lines": [4]}
  }],
  "outcome": {"filter_status": "accepted", "delivery_status": "sent", "acceptance_scope": "next_hop"}
}
```

L'exemple 1.0 de accepted.log demeure juste pour ce fichier seul : livraison inconnue. L'import du complément permet l'enrichissement ; il ne faut pas modifier rétroactivement la fixture initiale.

## Report aval et événements sur deux relais

Fixture : ../testdata/deferred-downstream.log. Queue initiale connue par LMTP, mais lignes initiales de réception et suppression absentes. Expéditeur et taille après filtrage connus ; taille avant filtrage inconnue. Un seul destinataire original et sa cible réécrite : pas un résultat multi-destinataires.

```json
{
  "schema_version": "1.1-draft",
  "node": "pmg-example",
  "identifiers": {"initial_queue_id": "B7C8D9E0F1", "pmg_filter_id": "F400000000000004", "post_filter_queue_ids": ["A6B7C8D9E0"], "message_id": "<deferred-001@sender.example.net>"},
  "client": {"ip": "192.0.2.106", "hostname_logged": "mx-d.example.net", "observed_via": "orig_client", "source_lines": [4]},
  "queues": [
    {"queue_id": "B7C8D9E0F1", "mail_from": null, "size_bytes_logged": null, "removed": null},
    {"queue_id": "A6B7C8D9E0", "mail_from": "sender@example.net", "size_bytes_logged": 18603, "removed": null, "source_lines": [6]}
  ],
  "filter": {"terminal_action": "accept", "recipient": "contact@example.org", "rule": "default-accept", "source_lines": [7]},
  "transport_events": [
    {"type": "remote_reply", "relay": {"hostname_logged": "mx2.relay.example.org", "ip": "192.0.2.107", "port": null}, "smtp_stage": "rcpt_to", "smtp_code": 421, "enhanced_status_code": "4.7.1", "response_text": "<target-deferred@example.org>: Recipient address rejected: blacklisted domain", "source_lines": [10]},
    {"type": "connection_lost", "relay": {"hostname_logged": "mx2.relay.example.org", "ip": "192.0.2.107", "port": null}, "smtp_stage": "data", "source_lines": [11]}
  ],
  "deliveries": [{
    "queue_id": "A6B7C8D9E0", "protocol": "smtp", "to": "target-deferred@example.org", "orig_to": "contact@example.org",
    "relay": {"hostname_logged": "mx1.relay.example.org", "ip": "192.0.2.108", "port": 25},
    "postfix_status": "deferred", "smtp_code": 421, "enhanced_status_code": "4.7.1", "smtp_stage": "rcpt_to",
    "response_text": "<target-deferred@example.org>: Recipient address rejected: blacklisted domain",
    "delay_seconds_logged": 0.44, "delays_seconds_logged": [0.05, 0, 0.37, 0.02], "source_lines": [12]
  }],
  "outcome": {"filter_status": "accepted", "delivery_status": "deferred", "final_outcome_known": false},
  "correlation": {"method": "explicit", "proof_lines": [7, 9, 12], "warnings": ["initial_queue_metadata_missing", "subsequent_delivery_outcome_missing"]}
}
```

Le texte blacklisted domain provient du relais distant, pas d'une règle PMG. Conserver le code 421 : ne pas transformer ce motif en échec permanent 5xx. Ne pas inventer une réussite ultérieure ni un nombre de retries.

## Couverture restante

Toujours absents : rejet permanent NOQUEUE, action de quarantaine explicite, libération de quarantaine, report suivi de succès, échec permanent aval et résultats mixtes par destinataire. Aucun de ces cas ne doit être déduit des seuls noms de tests SA, des adresses réécrites ou de plusieurs relais.
