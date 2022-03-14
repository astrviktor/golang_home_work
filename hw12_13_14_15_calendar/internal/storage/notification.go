package storage

import "time"

type Notification struct {
	ID        string    `json:"id"`        // ID - уникальный идентификатор события (UUID)
	Title     string    `json:"title"`     // Заголовок - короткий
	DateStart time.Time `json:"dateStart"` // Дата и время начала события
	UserID    int       `json:"usedId"`    // ID пользователя, владельца события
}

// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика, содержит поля:
//
// ID события;
// Заголовок события;
// Дата события;
// Пользователь, которому отправлять.
