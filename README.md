# 🎵 Music Downloader (Go)

Веб‑сервис для скачивания музыки через **Hitmo** и **yt-dlp** по:

- ссылке на **плейлист или альбом из Яндекс.Музыки**
- **текстовому `.txt` файлу** со списком треков вида `Исполнитель - Трек`

Сервис:

1. Получает список треков (из TXT или по API Яндекс.Музыки)
2. Ищет каждый трек на **Hitmo**
3. При неудаче может **дополнительно пробовать yt-dlp** (YouTube)
4. Скачивает найденные треки
5. Собирает их в единый **ZIP‑архив**

Всё работает как **одна веб‑страница на Go‑сервере** — фронтенд (HTML+CSS+JS) раздаётся из бэкенда.

---

## ✨ Возможности

- � Загрузка **TXT‑файла** со списком треков  
	Формат строки: `Artist - Title`

- 🔗 Загрузка по **ссылке Яндекс.Музыки**:
	- публичный / приватный **плейлист**  
	- **альбом**

- 🎧 Основной источник — **Hitmo** (поиск + скачивание MP3)

- 🎬 Фолбэк через **yt-dlp**  
	Если Hitmo не находит трек, можно включить скачивание через yt‑dlp (YouTube).

- 📦 Автоматическая архивация: все треки упаковываются в `music.zip`  
	Внутри также создаётся `failed_tracks.txt` с треками, которые не удалось скачать.

- ⚡ Параллельная загрузка треков (worker‑пул, настраиваемый в конфиге)

- � Полностью **конфигурируемый** через `config/config.yaml`

---

## 📂 Структура проекта

```bash
music-downloader
│
├── cmd/server            # точка входа (HTTP‑сервер)
├── config                # YAML‑конфиг и структуры
│   ├── config.go
│   └── config.yaml
│
├── internal
│   ├── api               # HTTP API, SSE, статика (HTML/CSS/JS)
│   │   ├── handler.go    # роутер и конструктор
│   │   ├── stream.go     # SSE‑эндпоинты /api/stream/*
│   │   ├── result.go     # хранилище готовых архивов
│   │   ├── middleware.go # логирование, recover
│   │   ├── sse.go        # типы SSE‑сообщений
│   │   ├── static.go     # раздача статики
│   │   ├── utils.go
│   │   └── static/       # index.html, style.css, app.js
│   │
│   ├── domain            # общие типы и утилиты
│   ├── downloader        # сервис, который оркестрирует скачивание
│   ├── sources
│   │   ├── hitmo         # поиск и скачивание с Hitmo
│   │   └── yandex        # парсинг плейлистов/альбомов Яндекс.Музыки
│   │
│   ├── ytdlp             # интеграция с внешним yt-dlp
│   └── zip               # сборка ZIP‑архива
│
└── pkg/logger            # обёртка над slog
```

---

## ⚙️ Конфигурация

Все настройки лежат в `config/config.yaml`.

Пример (упрощённо):

```yaml
app:
	base_url: "https://rus.hitmotop.com"
	search_results_per_page: 48

track:
	items_css: "div.track__info"
	artist_css: "div.track__desc"
	title_css: "div.track__title"
	dl_btn_css: "a.track__download-btn"

client:
	client_timeout: 20s
	mp3_timeout: 30s

ytdlp:
	workers: 3
	rps: 2
	timeout: 2m
	audio_format: "mp3"
	audio_quality: "192K"
	max_audio_size_mb: 30

downloader:
	workers: 5
	hitmo_delay: 1s
	use_ytdlp_fallback: true   # включить/выключить скачивание через yt-dlp

zip:
	compression_level: 6

server:
	addr: ":8080"
	read_timeout: 5m
	write_timeout: 10m

api:
	static_dir: "internal/api/static"
	max_upload_mb: 10
	result_ttl: 30m
	evict_interval: 5m
```

Главные параметры:

- `downloader.workers` — сколько треков качать одновременно
- `downloader.use_ytdlp_fallback` — использовать ли yt-dlp, если Hitmo не нашёл трек
- `zip.compression_level` — уровень сжатия архива
- `server.addr` — адрес/порт HTTP‑сервера
- `api.max_upload_mb` — максимальный размер загружаемого TXT‑файла
- `api.result_ttl` — сколько живут готовые ZIP‑архивы в памяти

---

## 🧰 Зависимости

Требования:

- **Go 1.25+**
- Утилиты в системе:
	- `yt-dlp`
	- `ffmpeg` (нужен yt-dlp для вытаскивания аудио)

Установка (macOS, через Homebrew):

```bash
brew install yt-dlp
brew install ffmpeg
```

---

## 🚀 Установка и запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/THENEAL24/Music-Downloader.git
cd Music-Downloader/music-downloader
```

### 2. Установить Go‑зависимости

```bash
go mod download
```

### 3. Проверить / поправить `config/config.yaml`

- Убедитесь, что `server.addr` свободен (по умолчанию `:8080`)
- При необходимости настройте `downloader.workers`, `use_ytdlp_fallback` и т.д.

### 4. Запустить сервер

```bash
go run ./cmd/server
```

Сервер поднимается на адресе из `config.yaml` (по умолчанию `http://localhost:8080`).

---

## 📥 Использование (через веб‑интерфейс)

Откройте в браузере:

```text
http://localhost:8080
```

На странице есть две вкладки:

### 1. TXT‑файл

1. Выберите или перетащите `.txt` файл.
2. Убедитесь, что строки в формате:

	 ```text
	 Linkin Park - Numb
	 Rammstein - Sonne
	 Daft Punk - Harder, Better, Faster, Stronger
	 ```

3. При необходимости задайте:
	 - количество потоков (`Потоков`)
	 - лимит треков (`Лимит`, 0 = все)
4. Нажмите **«Начать скачивание»**.
5. Внизу появится лог и прогресс‑бар.
6. По завершении станет доступна кнопка **«Скачать архив»** — это `music.zip`.

### 2. Яндекс.Музыка

Поддерживаются:

- публичные / приватные **плейлисты**:
	```text
	https://music.yandex.ru/users/<login>/playlists/<id>
	```
- **альбомы**:
	```text
	https://music.yandex.ru/album/<id>
	```

Шаги:

1. Вставьте ссылку на плейлист или альбом.
2. Если плейлист приватный — укажите OAuth‑токен Яндекса (для публичных можно оставить пустым).
3. Нажмите **«Начать скачивание»**.
4. Сервис:
	 - получит список треков через API Яндекс.Музыки,
	 - попробует найти и скачать их с Hitmo (и при необходимости через yt-dlp),
	 - соберёт архив.
5. По готовности появится кнопка **«Скачать архив»**.

---

## 📦 Структура результата

Вы скачиваете один ZIP‑файл (имя по умолчанию `music.zip`), внутри:

```text
Artist 1 - Track 1.mp3
Artist 2 - Track 2.mp3
...
failed_tracks.txt    # список треков, которые не удалось найти/скачать
```

`failed_tracks.txt` содержит строки вида:

```text
Original Query  ←  причина ошибки
```

---

## 🔧 API (если хочется дергать без фронта)

### `POST /api/stream/txt`

- Тип: `multipart/form-data`
- Поля:
	- `file` — TXT‑файл
	- (опционально) `workers`, `limit` — переопределяют настройки из конфига
- Ответ: **SSE‑поток** (`text/event-stream`)  
	Клиенту приходят события `start`, `track`, `done`, `error`.

### `POST /api/stream/yandex`

- Тип: `application/json`
- Тело:

	```json
	{
		"url": "https://music.yandex.ru/users/login/playlists/123",
		"token": "optional-oauth-token"
	}
	```

- Ответ: SSE‑поток с теми же типами событий.

### `GET /api/result/{id}`

- Возвращает готовый ZIP‑архив по `job_id`, который пришёл в событии `done`.

---

## ⚠️ Отказ от ответственности

Проект создан **исключительно в образовательных и исследовательских целях**.

Автор не поощряет нарушение авторских прав.  
Используйте инструмент только в рамках законодательства вашей страны и условий сервисов (Yandex, YouTube, др.).

---

## 🛠 Технологии

- **Go** (HTTP‑сервер, бизнес‑логика)
- **Server‑Sent Events (SSE)** для стриминга прогресса
- **Hitmo** (поиск и скачивание MP3)
- **yt-dlp + ffmpeg** (фолбэк‑загружчик)
- **Yandex Music API parsing** (плейлисты и альбомы)
- Чистая фронтенд‑страница на **HTML + CSS + vanilla JS**

---

## 💡 Идея проекта

Многие музыкальные сервисы ограничивают доступ к музыке из-за:

- региональных ограничений
- цензуры
- удаления контента

Этот проект создан как **технический инструмент для получения доступа к уже существующей музыке**, используя альтернативные источники.

---

## 📜 Лицензия

MIT License

---