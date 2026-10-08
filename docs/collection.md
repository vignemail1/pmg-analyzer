# Collecte, rotation et reprise

## Modes

Import historique et suivi continu. Au démarrage, découvrir l'historique et ouvrir le fichier actif sans attendre la fin d'un gros backfill. Une source canonique par flux ; ne pas importer mail.log et syslog comme indépendants lorsqu'ils recopient les mêmes événements sans politique explicite.

## Identité

Identifiant interne de génération, nœud, source logique, device/inode, empreinte d'un préfixe stable, taille et chemins observés. Inode et empreinte sont des indices : l'inode peut être réutilisé et deux préfixes peuvent être identiques. Ne pas fusionner des sources ambiguës silencieusement.

## Renommage/création

Lors de mail.log → mail.log.1 : maintenir le lecteur de l'ancien fichier et ouvrir le nouveau. EOF n'est pas une preuve de clôture ; attendre une période d'inactivité et permettre une réouverture. Notifications système comme optimisation ; scan périodique et réconciliation au démarrage comme garde-fous. La compression ou suppression avant capture complète peut rendre des données inaccessibles.

## Checkpoints

Pour chaque lot, stocker les lignes brutes puis avancer l'offset après la dernière ligne complète, dans une même transaction. Une contrainte (generation_id, offset_start) rend la relecture idempotente. Reprise après crash : relire le lot non validé. Les identités des copies doivent être réconciliées séparément.

Conserver les octets d'une ligne incomplète ; ne pas avancer le checkpoint au-delà de la dernière ligne complète. Une source close avec ligne finale sans newline nécessite une politique explicite et un indicateur.

## Compression

Lire les .gz en streaming, vérifier les erreurs jusqu'à la fin. Offset logique décompressé ; reprise MVP par redécompression depuis le début jusqu'au checkpoint. Rapprocher une archive avec sa génération non compressée par continuité de contenu, pas par nom. Ne jamais dédupliquer globalement des lignes identiques.

## copytruncate

Détecter diminution de taille et rupture de continuité ; nouvelle génération après troncature. Taille < offset seule est insuffisante si le fichier a repoussé entre scans. Rapprocher avec la copie tournée lorsque démontrable. copytruncate peut perdre des écritures entre copie et troncature : aucune garantie zéro perte. Préférer renommage et réouverture si compatibles avec la configuration réelle.

## Rétention et capacité

La conservation source doit dépasser l'indisponibilité maximale envisagée. Si stockage indisponible ou disque plein : ne pas avancer les checkpoints, alerter. Ne jamais supprimer les logs source. Rendre visibles les trous connus ; ne pas prétendre détecter tous les événements perdus.

## Validation avant production

Inspecter la configuration logrotate réelle, les permissions, le producteur des logs et sa procédure de réouverture. Tester rotation concurrente, écritures tardives, compression, archives supprimées, réutilisation d'inode, troncature, arrêt brutal, redémarrage, disque plein et import répété.

## Références

- https://man.archlinux.org/man/logrotate.8.fr
- https://www.elastic.co/docs/reference/beats/filebeat/file-log-rotation
- https://www.elastic.co/docs/reference/beats/filebeat/inode-reuse-issue
