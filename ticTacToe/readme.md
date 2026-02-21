# Tic-Tac-Toe

---

## ✅ Core Requirements (assumptions)

* 2 players
* NxN board (can be 3x3 by default)
* Turn-based
* Check winner / draw
* Pluggable strategies (optional for bot later)

---

## 🧩 Main Classes

### 1️⃣ Game (Orchestrator)

Responsible for:

* Game flow
* Turns
* Starting & ending game

**Attributes**

* `board: Board`
* `players: List<Player>`
* `currentPlayerIndex: int`
* `status: GameStatus`
* `winner: Player`

**Methods**

* `start()`
* `makeMove(row, col)`
* `switchTurn()`
* `checkGameOver()`
* `getWinner()`

---

### 2️⃣ Board

Responsible for:

* Maintaining grid
* Validating moves
* Checking win/draw

**Attributes**

* `size: int`
* `grid: Cell[][]`

**Methods**

* `isValidMove(row, col): boolean`
* `applyMove(row, col, symbol)`
* `checkWinner(): Symbol`
* `isFull(): boolean`
* `reset()`

---

### 3️⃣ Cell

Represents each square

**Attributes**

* `row: int`
* `col: int`
* `symbol: Symbol` (X / O / EMPTY)

**Methods**

* `isEmpty(): boolean`
* `setSymbol(Symbol)`

---

### 4️⃣ Player (Abstract / Interface)

**Attributes**

* `name: String`
* `symbol: Symbol`

**Methods**

* `makeMove(board): Move`

---

### 5️⃣ HumanPlayer extends Player

**Methods**

* `makeMove(board)`

---

### 6️⃣ BotPlayer extends Player (Optional for extensibility)

**Attributes**

* `strategy: MoveStrategy`

**Methods**

* `makeMove(board)`

---

### 7️⃣ Move

**Attributes**

* `row: int`
* `col: int`
* `symbol: Symbol`

---

### 8️⃣ Enums

```text
Symbol
- X
- O
- EMPTY

GameStatus
- IN_PROGRESS
- DRAW
- WON
```

---

## 📊 UML Class Diagram (Textual Representation)

```
+------------------+
|      Game        |
+------------------+
| - board: Board   |
| - players: List<Player> |
| - currentPlayerIndex: int |
| - status: GameStatus |
| - winner: Player |
+------------------+
| + start()        |
| + makeMove(r,c) |
| + switchTurn()   |
| + checkGameOver()|
| + getWinner()    |
+------------------+
        |
        | 1
        |
        | has
        |
        v
+------------------+
|      Board       |
+------------------+
| - size: int      |
| - grid: Cell[][] |
+------------------+
| + isValidMove()  |
| + applyMove()    |
| + checkWinner()  |
| + isFull()       |
+------------------+

+------------------+
|      Cell        |
+------------------+
| - row: int       |
| - col: int       |
| - symbol: Symbol |
+------------------+
| + isEmpty()      |
| + setSymbol()    |
+------------------+

+------------------+
|     Player       |<---------+
+------------------+          |
| - name: String   |          |
| - symbol: Symbol |          |
+------------------+          |
| + makeMove()     |          |
+------------------+          |
     ▲                        ▲
     |                        |
+-----------+        +----------------+
| Human     |        | BotPlayer      |
| Player    |        +----------------+
+-----------+        | - strategy     |
                     +----------------+

+------------------+
|      Move        |
+------------------+
| - row: int       |
| - col: int       |
| - symbol: Symbol |
+------------------+
```

---

## 🧠 Design Decisions (Interview Talking Points)

✔ **Single Responsibility**

* `Game` → flow control
* `Board` → game rules
* `Player` → move provider

✔ **Open-Closed Principle**

* Can add `BotPlayer`, `RemotePlayer`, `AIPlayer` without changing Game logic.

✔ **Extensible**

* Board size NxN
* Can add:

  * WinningStrategy interface
  * Replay
  * Observer for UI updates

---

## ✨ Optional Advanced UML (Strategy Pattern for Winner Check)

```
WinningStrategy (interface)
  |
  |-- RowWinningStrategy
  |-- ColumnWinningStrategy
  |-- DiagonalWinningStrategy
```

Board uses:

```
List<WinningStrategy>
```

