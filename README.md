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
CREATE TABLE games(id INT PRIMARY KEY AUTO_INCREMENT, gamemode VARCHAR(255), winner INT);
INSERT INTO games(gamemode, winner) VALUES('MultiPlayer', 0);
INSERT INTO games(gamemode, winner) VALUES('IA', 0);
INSERT INTO games(gamemode, winner) VALUES('IA', 2);
INSERT INTO games(gamemode, winner) VALUES('IA', 0);
INSERT INTO games(gamemode, winner) VALUES('IA', 2);
INSERT INTO games(gamemode, winner) VALUES('IA', 1);
INSERT INTO games(gamemode, winner) VALUES('IA', 1);
INSERT INTO games(gamemode, winner) VALUES('MultiPlayer', 0);
INSERT INTO games(gamemode, winner) VALUES('MultiPlayer', 0);
INSERT INTO games(gamemode, winner) VALUES('MultiPlayer', 1);
INSERT INTO games(gamemode, winner) VALUES('MultiPlayer', 0);
INSERT INTO games(gamemode, winner) VALUES('MultiPlayer', 0);
INSERT INTO games(gamemode, winner) VALUES('MultiPlayer', 2);
```

winner values :
- 0 : tie
- 1 : player1
- 2 : player2 (or IA)