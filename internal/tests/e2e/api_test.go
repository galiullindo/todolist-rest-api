package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/galiullindo/todolist-rest-api/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

func createTask(t *testing.T, db *sql.DB, task model.Task) {
	t.Helper()

	result, err := db.Exec("INSERT INTO tasks (title, status) VALUES (@title, @status)",
		sql.Named("title", task.Title),
		sql.Named("status", task.Status),
	)

	if err != nil {
		t.Error(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		t.Error(err)
	}
	if rows != 1 {
		t.Errorf("expected single row affected, got %d rows affected", rows)
	}
}

func TestAPI(t *testing.T) {
	t.Run("GET /tasks возвращает пустой список задач", func(t *testing.T) {
		want := []model.Task{}
		wantStatusCode := http.StatusOK
		wantContentType := "application/json"

		resp, err := testClient.Get(testServer.URL + "/tasks")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}

		var got []model.Task
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Error(err)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v\n", got, want)
		}
	})

	t.Run("GET /tasks возвращает список задач", func(t *testing.T) {
		defer rollbackToSavepoint(testDB)

		want := []model.Task{
			{ID: 1, Title: "Новая задача 1", Status: "Статус задачи 1"},
			{ID: 2, Title: "Новая задача 2", Status: "Статус задачи 2"},
			{ID: 3, Title: "Новая задача 3", Status: "Статус задачи 3"},
		}
		wantStatusCode := http.StatusOK
		wantContentType := "application/json"

		for _, task := range want {
			createTask(t, testDB, task)
		}

		resp, err := testClient.Get(testServer.URL + "/tasks")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}

		var got []model.Task
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Error(err)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v\n", got, want)
		}
	})

	t.Run("HEAD /tasks возвращает метаданные списка задач", func(t *testing.T) {
		wantStatusCode := http.StatusOK
		wantContentType := "application/json"

		resp, err := testClient.Head(testServer.URL + "/tasks")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}
	})

	t.Run("POST /tasks создает задачу", func(t *testing.T) {
		defer rollbackToSavepoint(testDB)

		want := model.Task{ID: 1, Title: "Новая задача", Status: "Статус задачи"}
		wantStatusCode := http.StatusCreated
		wantContentType := "application/json"

		body := bytes.NewBuffer([]byte(`{"title":"Новая задача","status":"Статус задачи"}`))
		resp, err := testClient.Post(testServer.URL+"/tasks", "application/json", body)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}

		var got model.Task
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Error(err)
		}

		if got != want {
			t.Errorf("got %v, want %v\n", got, want)
		}
	})

	t.Run("POST /tasks возвращает 400 для невалидного JSON", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		// wantContentType := "application/json"

		body := bytes.NewBuffer([]byte(`{invalid json}`))
		resp, err := testClient.Post(testServer.URL+"/tasks", "application/json", body)
		if err != nil {
			t.Errorf("client post: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("GET /tasks/{id} возвращает задачу по id", func(t *testing.T) {
		defer rollbackToSavepoint(testDB)

		want := model.Task{ID: 1, Title: "Новая задача", Status: "Статус задачи"}
		wantStatusCode := http.StatusOK
		wantContentType := "application/json"

		createTask(t, testDB, want)

		resp, err := testClient.Get(testServer.URL + "/tasks/1")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}

		var got model.Task
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Error(err)
		}

		if got != want {
			t.Errorf("got %v, want %v\n", got, want)
		}
	})

	t.Run("GET /tasks/{id} возвращает 400 для невалидного id", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		// wantContentType := "application/json"

		resp, err := testClient.Get(testServer.URL + "/tasks/abc")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("GET /tasks/{id} возвращает 404 для несуществующей задчи", func(t *testing.T) {
		wantStatusCode := http.StatusNotFound
		// wantContentType := "application/json"

		resp, err := testClient.Get(testServer.URL + "/tasks/1")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("HEAD /tasks/{id} возвращает 200 для существующей задачи", func(t *testing.T) {
		defer rollbackToSavepoint(testDB)

		wantStatusCode := http.StatusOK
		wantContentType := "application/json"

		task := model.Task{ID: 1, Title: "Новая задача", Status: "Статус задачи"}
		createTask(t, testDB, task)

		resp, err := testClient.Head(testServer.URL + "/tasks/1")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}
	})

	t.Run("HEAD /tasks/{id} возвращает 400 для невалидного id", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		// wantContentType := "application/json"

		resp, err := testClient.Head(testServer.URL + "/tasks/abc")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("HEAD /tasks/{id} возвращает 404 для несуществующей задачи", func(t *testing.T) {
		wantStatusCode := http.StatusNotFound
		// wantContentType := "application/json"

		resp, err := testClient.Head(testServer.URL + "/tasks/1")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("PUT /tasks/{id} заменяет все поля задачи по id", func(t *testing.T) {
		defer rollbackToSavepoint(testDB)

		want := model.Task{ID: 1, Title: "Измененная задача", Status: "Измененный статус"}
		wantStatusCode := http.StatusOK
		wantContentType := "application/json"

		task := model.Task{ID: 1, Title: "Новая задача", Status: "Статус задачи"}
		createTask(t, testDB, task)

		body := bytes.NewBuffer([]byte(`{"title":"Измененная задача","status":"Измененный статус"}`))
		req, err := http.NewRequest(http.MethodPut, testServer.URL+"/tasks/1", body)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}

		var got model.Task
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Error(err)
		}

		if got != want {
			t.Errorf("got %v, want %v\n", got, want)
		}
	})

	t.Run("PUT /tasks/{id} возвращает 400 для невалидного id", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		// wantContentType := "application/json"

		body := bytes.NewBuffer([]byte(`{"title":"Измененная задача","status":"Измененный статус"}`))
		req, err := http.NewRequest(http.MethodPut, testServer.URL+"/tasks/abc", body)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("PUT /tasks/{id} возвращает 400 для невалидного JSON", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		// wantContentType := "application/json"

		body := bytes.NewBuffer([]byte(`{invalid json}`))
		req, err := http.NewRequest(http.MethodPut, testServer.URL+"/tasks/1", body)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("PUT /tasks/{id} возвращает 404 для несуществующего id", func(t *testing.T) {
		wantStatusCode := http.StatusNotFound
		// wantContentType := "application/json"

		body := bytes.NewBuffer([]byte(`{"title":"Измененная задача","status":"Измененный статус"}`))
		req, err := http.NewRequest(http.MethodPut, testServer.URL+"/tasks/1", body)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("DELETE /tasks/{id} удаляет задачу по id", func(t *testing.T) {
		defer rollbackToSavepoint(testDB)

		wantStatusCode := http.StatusNoContent
		wantContentType := ""

		task := model.Task{ID: 1, Title: "Новая задача", Status: "Статус задачи"}
		createTask(t, testDB, task)

		req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/tasks/1", nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != wantContentType {
			t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		}
	})

	t.Run("DELETE /tasks/{id} возвращает 400 для невалидного id", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		// wantContentType := "application/json"

		req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/tasks/abc", nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})

	t.Run("DELETE /tasks/{id} возвращает 404 для несуществующего id", func(t *testing.T) {
		wantStatusCode := http.StatusNotFound
		// wantContentType := "application/json"

		req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/tasks/1", nil)
		if err != nil {
			t.Error(err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != wantStatusCode {
			t.Errorf("status code got %d, want %d\n", resp.StatusCode, wantStatusCode)
		}

		// contentType := resp.Header.Get("Content-Type")
		// if contentType != wantContentType {
		// 	t.Errorf("content type got %s, want %s\n", contentType, wantContentType)
		// }
	})
}
