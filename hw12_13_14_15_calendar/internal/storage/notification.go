package storage

type Notification struct { // TODO
}

// Уведомление - временная сущность, в БД не хранится, складывается в очередь для рассыльщика, содержит поля:
//
// ID события;
// Заголовок события;
// Дата события;
// Пользователь, которому отправлять.