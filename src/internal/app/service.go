package app

import (
	"errors"
	"math"
	"github.com/google/uuid"
)

type TicTacToeService struct{}

const (
	Empty = 0
	Cross = 1 // X
	Naught = 2 // O
)

func (s *TicTacToeService) NewGame(Computer bool) (*CurrentGame){
	newID := uuid.New()
    newGame := &CurrentGame{
        UUID:    newID,
        Field: [][]int{
            {0, 0, 0},
            {0, 0, 0},
            {0, 0, 0},
        },
		Status: Wait,
		Computer: Computer,
    }
	return newGame
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

