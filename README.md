# TicTacToe project - Group 9
---
## Basic informations
- HEIA-FR student mini-project
- Goal: discover Go language by developing a game
- Initial repo: https://github.com/LempekPL/GoTicTacToe

## Students
- Margueron William
- Roch-Neirey Martin

## Objectifs
- Mettre à jour le programme
- Identifier des faiblesses de code
- Contrôler la nomenclature et les constantes
- Génération d'image
- Refactoring
- Analyse de performance

## MySQL infos

Création de la table :
```mysql
CREATE TABLE games3(
                       id int auto_increment primary key,
                       date datetime,
                       properties json
);
```
Requête pour récupérer les 5 dernières games :
```mysql
SELECT * FROM
    (
        SELECT * FROM games3 ORDER BY id DESC LIMIT 5
    ) AS sub
ORDER BY id ASC;
```


winner values :
- 0 : tie
- 1 : player1
- 2 : player2 (or IA)