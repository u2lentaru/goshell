Установка приложения  
git clone https://github.com/u2lentaru/goshell.git  
docker-compose up  
  
Должны бытьустановлены Docker/Docker Desktop и Docker-compose  
  
GET  
/ - список всех команд  
/commands - список всех команд  
/commands/{id} - выводит команду с id={id}  
/cmdrun&ids=1,2,3 - выполняет команды с id=1,2,3  
/cmdrun/{id} - выполняет команду с id={id}  
/results - список всех результатов  
  
POST  
/commands - добавляет в базу скрипт из тела запроса и выполняет его  