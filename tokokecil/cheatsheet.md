# 1 Lakukan deploy db mysql railways
# 2 Adjust Port dinamis mengikuti port railways => main.go
# 3 REST API ke heroku via cli
# 4 Set ENV = heroku-cli


## cheatsheet

1. heroku login
2. heroku create project_name
3. git init 
4. git add .
5. git commit -m "message_commit"
6. heroku git:remote -a tokokecil-app 
7. git remote -v
8. git push heroku main


## cheatsheet - GO Framework ( Echo/Gin/Fiber ) 

1. heroku login 
2. (Opsional: create app di web) - tidak perlu, karena udah dibuat di step 1. db postgre
3. Inisialisasi & commit git ( git init, git add . , git commit)
4. Set remote Heroku (heroku git:remote -a toko-mini-app)
5. Cek git remote (git remote -v) 
6. Set environment variable di dashboard () - Database URL dibiarkan saja karena itu set dari POSTGRE
7. Push ke Heroku ( git push heroku main)
8. Cek log (heroku logs --tail)
9. Akses dan test API

# kondisi : ada penambahan fitur / troubleshoot 

1. git add .
2. git commit -m "message_commit"
3. heroku git:remote -a tokokecil-app 
4. git remote -v
5. git push heroku main

# restart heroku 

heroku config:set DB_USER= DB_PASS= DB_HOST= DB_PORT DB_NAME -a project_name
heroku restart -a tokokecil-app