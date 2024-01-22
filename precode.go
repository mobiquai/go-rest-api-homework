package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

// Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из мапы tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента JSON
	w.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)

	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// Обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	// в заголовок можно НЕ записывать тип контента, т.к. тело ответа не будет возвращаться в ответе

	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	// получаем значение параметра "id" из URL
	id := chi.URLParam(r, "id")

	// Проверяем, есть ли в мапе элемент по ключу id
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// сериализуем данные из найденного элемента мапы
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента JSON
	w.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)

	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// Обработчик удаления задачи по ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	// получаем значение параметра "id" из URL
	id := chi.URLParam(r, "id")

	// Проверяем, есть ли в мапе элемент по ключу id
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// Удаляем из мапы найденный элемент
	delete(tasks, id)

	// в заголовок записываем тип контента "Простой текст"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)

	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte(fmt.Sprint("Задача с ID ", id, " успешно удалена")))
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)

	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `createTask`
	r.Post("/tasks", createTask)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", getTask)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTask`
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
