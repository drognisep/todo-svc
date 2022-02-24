build-todo:
	docker build -t todo-api:dev -f zarf/todo-api/Dockerfile .

run-todo: build-todo
	docker run --rm --name=test-todo -p 3000:3000 todo-api:dev
