package app

import (
	"errors"
	"log"
	"math"

	"github.com/google/uuid"
)

type TicTacToeService struct{}

const (
	Empty = 0
	Cross = 1 // X
	Naught = 2 // O
)

func (s *TicTacToeService) NewGame(Computer bool, Uuid string) (*CurrentGame){
	newID := uuid.New()
	px, err := uuid.Parse(Uuid)
	if err != nil{
		log.Fatal("Bad uuid for player X")
	}
    newGame := &CurrentGame{
        UUID:    newID,
        Field: [][]int{
            {0, 0, 0},
            {0, 0, 0},
            {0, 0, 0},
        },
		Status: Wait,
		Computer: Computer,
		PlayerX: px,
		PlayerO: uuid.Nil,
    }
	return newGame
}

func (s *TicTacToeService) Connect(game *CurrentGame, Uuidgame string, Uuidplayero string) (*CurrentGame){
	tmp, err := uuid.Parse(Uuidplayero)
	if err != nil {
		log.Fatalf("Bad PlayerO uuid")
		return game
	}
	game.PlayerO = tmp
	return game
}

func (s *TicTacToeService) FieldValidation(game *CurrentGame) (bool, error){
	if game == nil || len(game.Field) != 3{
		return false, errors.New("invalid field size")
	}

	CrossCount, NaughtCount := 0,0

	for _, row := range game.Field{
		if len(row) != 3 {
			return false, errors.New("invalid rows size")
		}
		for _, elem := range row{
			switch (elem){
			case Cross:
				CrossCount++
			case Naught:
				NaughtCount++
			case Empty:
			default:
				return false, errors.New("invalid value")
			}
		}
	}

	if NaughtCount > CrossCount || CrossCount - NaughtCount > 1 {
		return false, errors.New("more cell than can be")
	}

	return true, nil
}

func (s *TicTacToeService) NextMove(game *CurrentGame) (*CurrentGame, error){
	
	if !game.Computer{
		if game.Status == Wait {
			game.Status = MoveX
		}
		if game.Status == MoveX{
			game.Status = MoveO
		}
		if game.Status == MoveO {
			game.Status = MoveX
		}
		return game, nil
	}

	bestScore := math.Inf(-1) // Ищем максимальный score, изначально минус бесконечность
	var moveX, moveY int      // Лучшая найденная позиция

	// Перебираем все клетки поля
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Если клетка пустая — можно попробовать походить
			if game.Field[i][j] == Empty {
				// Пробуем поставить нолик
				game.Field[i][j] = Naught

				// Рекурсивно оцениваем этот ход
				score := minimax(game.Field, 0, false, game)

				// Откатываем ход обратно
				game.Field[i][j] = Empty

				// Если ход лучше предыдущих — сохраняем координаты
				if score > bestScore {
					bestScore = score
					moveX = i
					moveY = j
				}
			}
		}
	}

	// Делаем лучший найденный ход
	game.Field[moveX][moveY] = Naught
	game.Status = MoveX
	return game, nil
}

// minimax — классический алгоритм поиска оптимального хода в игре нолики-крестики.
// Он перебирает все возможные ходы, оценивает их и выбирает лучший.
// isMaximizing — true, если сейчас ходит ИИ (нолик).
func minimax(field Field, depth int, isMaximizing bool, game *CurrentGame) float64 {
	winner := checkwinner(game)

	// Если победил ИИ — хорошо
	if winner == Naught {
		return 1
	}
	// Если победил игрок — плохо
	if winner == Cross {
		return -1
	}

	// Если ничья
	isFull := true
	for _, row := range field {
		for _, cell := range row {
			if cell == Empty {
				isFull = false
			}
		}
	}
	if isFull {
		return 0
	}

	// Ход ИИ — максимизируем
	if isMaximizing {
		best := math.Inf(-1)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if field[i][j] == Empty {
					field[i][j] = Naught
					best = math.Max(best, minimax(field, depth+1, false, game))
					field[i][j] = Empty
				}
			}
		}
		return best
	}

	// Ход игрока — минимизируем
	best := math.Inf(1)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if field[i][j] == Empty {
				field[i][j] = Cross
				best = math.Min(best, minimax(field, depth+1, true, game))
				field[i][j] = Empty
			}
		}
	}
	return best
}

func (s *TicTacToeService) GameIsOver(game *CurrentGame) bool{
	winner := checkwinner(game)

	if winner == Cross{
		game.Status = WinX
	} else if winner == Naught{
		game.Status = WinO
	} else if winner != 0 {
		game.Status = Draw
	}

	if winner != 0{
		return true
	}

	for _, row := range game.Field{
		for _, cell := range row {
			if cell == Empty{
				return false
			}
		}
	}

	return true
}

func checkwinner(game *CurrentGame) int{

	for _, row := range game.Field{
		counter_cross := 0
		counter_naught := 0
		for _, cell := range row{
			if cell == Cross{
				counter_cross++
			}
			if cell == Naught{
				counter_naught++
			}
		}
		if counter_cross == 3{
			return Cross
		}
		if counter_naught == 3{
			return Naught
		}
	}

	for i:=0; i<3; i++{
		counter_cross := 0
		counter_naught := 0
		for j:=0;j<3;j++{
			if game.Field[j][i] == Cross{
				counter_cross++
			}
			if game.Field[j][i] == Naught{
				counter_naught++
			}
		}
		if counter_cross == 3{
			return Cross
		}
		if counter_naught == 3{
			return Naught
		}
	}

	if game.Field[0][0] == game.Field[1][1] && 
	game.Field[1][1] == game.Field[2][2] && 
	game.Field[0][0] != Empty{
		if game.Field[0][0] == Cross{
			return Cross
		}
		if game.Field[0][0] == Naught{
			return Naught
		}
	} 

	if game.Field[0][2] == game.Field[1][1] &&
	game.Field[1][1] == game.Field[2][0] &&
	game.Field[0][2] != Empty {

	if game.Field[0][2] == Cross {
		return Cross
	}
	if game.Field[0][2] == Naught {
		return Naught
	}
	}


	return 0
}

