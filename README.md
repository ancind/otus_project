# Image Previewer

## Описание сервиса
Сервис предназначен для изготовления preview (создания из изображения уменьшенной копии)

## Команды для работы с проектом
1. make build-server - билдим бинарник приложения
2. make test-single - запуск одного прохода тестов с таймаутом 1м
3. make test-race - запускает 100 проходов теста с таймаутом 7м
4. make test-coverage - посмотреть % покрытия тестами проекта 
5. make test-integration - запуск интеграционных тестов
6. make lint - запуск линтера в проекте

### Запуск докера

Докер запускается с флагом `--remove-orphans`, чтобы не захламлять старыми контейнерами

1. make start - запускает проект 
2. make stop - остановка контейнеров

## Описание работы сервиса
Отправляем на url resize параметры для изменения размера изображения и ссылку для получения исходного изображения которое хотим изменить

## Как проверить работоспособность сервиса.

### Параметры url
`host/fill/{width}/{height}/{imageUrl}`

1. width - ширина картинки, которую хотим получить
2. height - высота картинки, которую хотим получить
3. imageUrl - ссылка на исходное изображение для уменьшения по заданным параметрам

### Пример для получения картинки после запуска докера
`http://127.0.0.1/fill/200/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_1024x252.jpg`