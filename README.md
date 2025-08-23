# inventory

1. Create db manually
```
mysql -u root -p
create database inventory;
use inventory;
create table products(ID int NOT NULL PRIMARY KEY AUTO_INCREMENT, Name varchar(255), Quantity int, Price float(10,7));
insert into products values(1,"apple",2,10.99);
insert into products values(2,"banana",2,2.99);
desc table products;
```
funfact: When you enter 10.99, MySQL stores the closest representable value, which might be 10.9899998, not exactly 10.99

2. `go mod init github.com/anushasgorawar/inventory`

3. app.go will have data/method related to routes
A struct to store mux router and db details 
Once created, run `go get`

4. For that struct pointer, create 2 methods. Initialse and Run which we call in main.go

5. Install mysql driver:
Run `go get "github.com/go-sql-driver/mysql"` 
Add `_ "github.com/go-sql-driver/mysql"` in imports of main and do go mod 

5. app.Router.HandleFunc("/products", app.allProducts).Methods("GET") 
Methods function registers a new route for this allProducts method
App.allProducts will be called only for APP struct
it'll call allProducts which is present in (7)
in app.allProducts-> we are calling the allProducts method 
if error, return status code 500 with err.Error() 
err.Error() -> return the error as string
else, call sendResponse function explained in 6

6. Create a method that responds with status code if error, else with a response with the payload i.e. products and write it to w

7. models.go will have all the db related methods
allProducts will have th function to get all products

8. Use postman instead of browser.
Create workspace -> Create collection -> create a request 

9. app.go has the routing functions and models.go has all the db related methods

10. 