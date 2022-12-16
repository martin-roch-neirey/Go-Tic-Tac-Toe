# TicTacToe project - Group 9
---
## Students:
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

```mysql
CREATE TABLE games3(
                       id int auto_increment primary key,
                       date datetime,
                       properties json
);
```

winner values :
- 0 : tie
- 1 : player1
- 2 : player2 (or IA)