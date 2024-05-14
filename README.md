# goshell  
git clone https://github.com/u2lentaru/goshell.git  
docker-compose up  

"/", handlers.HandleList).Methods("GET")  
"/commands", handlers.HandleList).Methods("GET")  
"/commands/{id:[0-9]+}", handlers.HandleGetOne).Methods("GET")  
"/commands", handlers.HandlePostExec).Methods("POST")  
"/cmdrun/{id:[0-9]+}", handlers.HandleExecOne).Methods("GET")  
"/cmdrun", handlers.HandleExec).Methods("GET")  
"/results", handlers.HandleResults).Methods("GET")  


/cmdrun&ids=1,2,3 - выполняет команды 1.2,3