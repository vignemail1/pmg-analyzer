# Rejets permanents, reports résolus et résultats mixtes

## Statut et anonymisation

Complément de json-format.md et json-additional-cases.md. Projections documentaires 1.1-draft : pas sorties générées, pas dossiers complets ni golden tests exécutés. Les source_lines sont les positions dans la fixture indiquée, à partir de 1.

Sélection de cas non redondants dans le résultat fourni. Bandeaux et préfixes de recherche retirés. Dates, nœud, PIDs, queues, adresses, IP, domaines, identifiants distants, URLs et tokens de réponse remplacés. La zone DNSBL incluant un identifiant de compte est remplacée intégralement. Les valeurs SMTP/DSN, étapes, motifs techniques, délais et relations entre adresses restent exploitables. La réponse distante complexe a ses identifiants anonymisés sans supprimer sa structure. Aucun mapping vers les originaux n'est publié. Les attributs techniques conservés ne constituent pas une anonymisation statistique absolue.

Les lignes new mail dont le Message-ID contient /releases/ sont des faux positifs de la recherche de libération. Elles ne sont pas ajoutées comme fixtures de libération. Aucune quarantaine ni libération explicite n'est prouvée par cette sortie.

## Rejet postscreen permanent

Fixture : testdata/noqueue-permanent-postscreen.log. Trois événements indépendants, pas une session unique. Le service est postfix/postscreen, pas smtpd. Conserver le port client, le HELO et la zone DNSBL. Aucun queue ID ni contenu du message observé. PID commun ne justifie aucun regroupement ; bornes de connexion absentes. [0.0.0.0] est conservé comme HELO littéral, pas comme adresse cliente. La relation from=to du troisième événement est conservée.

```json
{
  "schema_version": "1.1-draft",
  "source_file": "testdata/noqueue-permanent-postscreen.log",
  "events": [
    {"timestamp": "2026-01-04T00:19:25.985491+02:00", "node": "pmg-example", "service": "postfix/postscreen", "pid": 600100, "type": "smtp.reject", "stage": "rcpt_to", "decision": "permanent_reject", "queue_id": null, "client": {"ip": "192.0.2.151", "port": 46871, "hostname_logged": null}, "mail_from": "sender-a@example.net", "recipient": "recipient-a@example.org", "protocol": "ESMTP", "helo": "sender-a.example.net", "smtp_code": 550, "enhanced_status_code": "5.7.1", "reason_category": "dnsbl", "dnsbl_zone_logged": "account-placeholder.sbl.dnsbl.example.net", "source_lines": [1]},
    {"timestamp": "2026-01-04T00:19:39.601057+02:00", "node": "pmg-example", "service": "postfix/postscreen", "pid": 600100, "type": "smtp.reject", "stage": "rcpt_to", "decision": "permanent_reject", "queue_id": null, "client": {"ip": "192.0.2.152", "port": 5321, "hostname_logged": null}, "mail_from": "sender-b@example.net", "recipient": "recipient-b@example.org", "protocol": "ESMTP", "helo": "helo-b.example.net", "smtp_code": 550, "enhanced_status_code": "5.7.1", "reason_category": "dnsbl", "dnsbl_zone_logged": "dnsbl-b.example.net", "source_lines": [2]},
    {"timestamp": "2026-01-04T00:24:47.652198+02:00", "node": "pmg-example", "service": "postfix/postscreen", "pid": 600100, "type": "smtp.reject", "stage": "rcpt_to", "decision": "permanent_reject", "queue_id": null, "client": {"ip": "192.0.2.153", "port": 65027, "hostname_logged": null}, "mail_from": "same-address@example.org", "recipient": "same-address@example.org", "protocol": "ESMTP", "helo": "[0.0.0.0]", "smtp_code": 550, "enhanced_status_code": "5.7.1", "reason_category": "dnsbl", "dnsbl_zone_logged": "dnsbl-b.example.net", "source_lines": [3]}
  ],
  "warnings": ["session_boundaries_missing", "message_content_not_observed"]
}
```

## Échecs aval : local au client SMTP ou réponse distante

Fixture : testdata/bounced-downstream.log. Trois queues distinctes. Premier cas : absence de capacité SMTPUTF8, aucun code de réponse SMTP de rejet journalisé ; smtp_code=null, et non 550 inventé. Deuxième : réponse distante 550. Troisième : DSN Postfix 5.0.0 distinct du code enrichi 5.7.1 cité dans le texte distant. Conserver les deux sans les confondre.

```json
{
  "schema_version": "1.1-draft",
  "source_file": "testdata/bounced-downstream.log",
  "deliveries": [
    {"queue_id": "C100000001", "protocol": "smtp", "to": "target-support@example.org", "orig_to": "support@example.org", "relay": {"hostname_logged": "mailbox.example.org", "ip": "2001:db8:1::10", "port": 25}, "postfix_status": "bounced", "postfix_dsn": "5.6.7", "smtp_code": null, "remote_enhanced_status_code": null, "smtp_stage": null, "reason_origin": "smtp_client", "reason_category": "smtputf8_unavailable", "source_lines": [1]},
    {"queue_id": "C100000002", "protocol": "smtp", "to": "unknown-user@example.org", "orig_to": null, "relay": {"hostname_logged": "mx-remote.example.net", "ip": "192.0.2.154", "port": 25}, "postfix_status": "bounced", "postfix_dsn": "5.1.1", "smtp_code": 550, "remote_enhanced_status_code": "5.1.1", "smtp_stage": "rcpt_to", "reason_origin": "remote_reply", "reason_category": "recipient_unknown", "source_lines": [2]},
    {"queue_id": "C100000003", "protocol": "smtp", "to": "notification@example.org", "orig_to": null, "relay": {"hostname_logged": "relay-policy.example.net", "ip": "192.0.2.155", "port": 25}, "postfix_status": "bounced", "postfix_dsn": "5.0.0", "smtp_code": 553, "remote_enhanced_status_code": null, "remote_status_codes_mentioned": ["5.7.1"], "smtp_stage": "rcpt_to", "reason_origin": "remote_reply", "reason_category": "relay_policy", "source_lines": [3]}
  ],
  "warnings": ["initial_queue_metadata_missing", "filter_actions_missing", "bounce_notification_not_observed"]
}
```

Les catégories normalisées sont des propositions ; la réponse brute et sa provenance restent nécessaires. status=bounced ne prouve pas qu'une notification de retour a été livrée. La taille et l'expéditeur ne sont pas disponibles dans ces lignes.

## Échec du transfert vers le filtre local

Fixture : testdata/bounced-local-filter.log. Ne pas compter ce cas comme un refus du serveur aval ni comme une règle block PMG. La réponse contient deux codes enrichis ; aucun ne doit écraser le DSN de transport.

```json
{
  "schema_version": "1.1-draft",
  "source_file": "testdata/bounced-local-filter.log",
  "filter_handoffs": [{"queue_id": "C100000004", "recipient": "contact-list@example.org", "protocol": "lmtp", "relay": {"ip": "127.0.0.1", "port": 10024}, "postfix_status": "bounced", "postfix_dsn": "5.0.0", "smtp_code": 554, "remote_enhanced_status_code": "5.0.0", "remote_status_codes_mentioned": ["5.4.0"], "smtp_stage": "end_of_message", "reason_category": "too_many_hops", "source_lines": [1]}],
  "outcome": {"filter_handoff_status": "failed", "filter_status": "unknown", "delivery_status": "unknown"}
}
```

## Report suivi de succès

Fixture : testdata/deferred-then-sent.log. Deux tentatives observées, même queue et même destinataire. Ne pas prétendre connaître le nombre total de tentatives ni les métadonnées absentes. delay est le délai journalisé de la queue, pas simplement l'écart entre ces deux timestamps.

```json
{
  "schema_version": "1.1-draft",
  "source_file": "testdata/deferred-then-sent.log",
  "queue_id": "C200000001",
  "recipient": "retry-recipient@example.org",
  "attempts": [
    {"timestamp": "2026-01-04T03:24:16.004811+02:00", "relay": {"hostname_logged": "mx-remote.example.net", "ip": "192.0.2.154", "port": 25}, "postfix_status": "deferred", "postfix_dsn": "4.7.1", "smtp_code": 450, "remote_enhanced_status_code": "4.7.1", "smtp_stage": "end_of_message", "delay_seconds_logged": 6.6, "delays_seconds_logged": [0.05, 0, 6.4, 0.15], "source_lines": [1]},
    {"timestamp": "2026-01-04T03:31:48.050567+02:00", "relay": {"hostname_logged": "mx-remote.example.net", "ip": "192.0.2.154", "port": 25}, "postfix_status": "sent", "postfix_dsn": "2.0.0", "smtp_code": 250, "remote_enhanced_status_code": "2.0.0", "delay_seconds_logged": 459, "delays_seconds_logged": [452, 0.02, 6.3, 0.23], "remote_queue_id": "remoteRetry001", "source_lines": [2]}
  ],
  "outcome": {"delivery_status": "sent", "acceptance_scope": "next_hop", "mailbox_delivery_confirmed": null},
  "warnings": ["initial_queue_metadata_missing", "queue_removal_not_observed"]
}
```

## Trois cibles, deux succès et un échec

Fixture : testdata/mixed-delivery-results.log. Les lignes sont triées par timestamp, contrairement à l'affichage du script. Un même orig_to est conservé pour trois cibles distinctes. L'extrait prouve les résultats individuels sous le même queue ID ; confirmer la continuité de la queue avec qmgr/cleanup avant d'en faire une preuve complète de message unique. Ne pas inférer nrcpt=3 : le nombre journalisé n'est pas fourni.

```json
{
  "schema_version": "1.1-draft",
  "source_file": "testdata/mixed-delivery-results.log",
  "queue_id": "C300000001",
  "recipient_count_logged": null,
  "observed_delivery_targets": 3,
  "deliveries": [
    {"to": "target-success-a@example.org", "orig_to": "contact-mixed@example.org", "relay": {"hostname_logged": "mx-success-a.example.net", "ip": "2001:db8:2::10", "port": 25}, "postfix_status": "sent", "postfix_dsn": "2.0.0", "smtp_code": 250, "remote_enhanced_status_code": "2.0.0", "remote_queue_id": null, "source_lines": [1]},
    {"to": "target-failed@example.org", "orig_to": "contact-mixed@example.org", "relay": {"hostname_logged": "mx-failed.example.net", "ip": "192.0.2.156", "port": 25}, "postfix_status": "bounced", "postfix_dsn": "5.1.0", "smtp_code": 501, "remote_enhanced_status_code": "5.1.0", "smtp_stage": "mail_from", "reason_origin": "remote_reply", "reason_category": "sender_policy", "source_lines": [2]},
    {"to": "target-success-b@example.org", "orig_to": "contact-mixed@example.org", "relay": {"hostname_logged": "mx-success-b.example.net", "ip": "192.0.2.157", "port": 25}, "postfix_status": "sent", "postfix_dsn": "2.6.0", "smtp_code": 250, "remote_enhanced_status_code": "2.6.0", "remote_status_codes_mentioned": ["2.1.5"], "remote_queue_id": null, "remote_internal_id": "10000000000001", "response_message_id": "<mixed-001@sender.example.net>", "source_lines": [3]}
  ],
  "outcome": {"delivery_status": "mixed", "counts_observed": {"sent": 2, "bounced": 1}, "acceptance_scope": "next_hop"},
  "warnings": ["queue_lifecycle_metadata_missing", "filter_actions_missing", "mailbox_delivery_not_confirmed"]
}
```

Ne pas transformer un token de réponse, InternalId ou Message-ID distant en remote_queue_id sans libellé explicite queued as. Un échec MAIL FROM est attribuable à l'émetteur/politique distante, même si la ligne Postfix contient to. Le texte distant ne prouve pas à lui seul la validité de son diagnostic SPF.

Le second candidat mixte fourni, séparé de plusieurs dizaines de minutes et sans métadonnées de queue, est volontairement exclu jusqu'à confirmation de l'absence de réutilisation du queue ID.

## Couverture mise à jour

Ajoutés : rejet permanent postscreen/DNSBL, échecs aval avec ou sans réponse SMTP distante, échec LMTP local, report puis succès, résultats de transport mixtes. Toujours manquants : action de quarantaine explicite, libération et actions PMG mixtes. Les mentions de cas manquants dans les documents antérieurs reflètent leur corpus au moment de rédaction ; ce document actualise cette couverture.
