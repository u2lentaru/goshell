Установка и запуск приложения  
 
1. git clone https://github.com/u2lentaru/goshell.git  
2. cd goshell 
3. docker-compose up  
  
Должны быть установлены Docker/Docker Desktop и Docker-compose  
  
GET  
http://localhost:8080/ - список всех команд  
http://localhost:8080/commands - список всех команд  
http://localhost:8080/commands/{id} - выводит команду с id={id}  
http://localhost:8080/cmdrun?ids=1,2,3 - выполняет команды с id=1,2,3  
http://localhost:8080/cmdrun/{id} - выполняет команду с id={id}  
http://localhost:8080/results - список всех результатов  
http://localhost:8080/health-check - health-check  
  
POST  
http://localhost:8080/commands - добавляет в базу скрипт из тела запроса и выполняет его