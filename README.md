# Library Management System (MVC)

A Library Management System built with a Go backend and a React frontend.

##  Setup Instructions

Clone the repo. From the root directory of the cloned repo:
```bash
go mod tidy
cp sample.env .env 
```
*(Make sure to update the `.env` file with your actual database credentials and a secure `JWT_SECRET`)*

### MYSQL Setup
1. `mysql -u root -p` : and enter password
2. Create a new database 'library_management': `CREATE DATABASE library_management;`
3. Connect to the database: `USE library_management;`

### Database Migrations
This project uses `golang-migrate` to manage database schemas.
1. Ensure that you have [golang-migrate](https://github.com/golang-migrate/migrate) installed.
2. Change the username and password in the connection string and run the migrations:
```bash
migrate -path database/migration/ -database "mysql://username:password@tcp(localhost:3306)/library_management" -verbose up
```

### First Admin Setup
Because this application uses manual salted hashing, you cannot insert an admin directly using plaintext SQL. 
1. Run the frontend and backend (instructions below).
2. Register a new user via the frontend UI (`/register`) to securely hash your password. 
3. Open your MySQL terminal and manually promote your account:
```sql
UPDATE users SET role = 'admin', request_status = 'accepted' WHERE email = 'your_email@example.com';
```

---

## Running the Application

**1. Running the server:**
Open a terminal in the root directory:
```bash
go build -o mvc ./cmd/main.go
./mvc
```
*(The backend runs on `http://localhost:8080`)*

**2. Running the frontend:**
Open a separate terminal and navigate to the `frontend` folder:
```bash
cd frontend
npm install
npm run dev
```
*(The UI runs on `http://localhost:5173`)*

---

## Hosting

Install or Update apache on your system:
```bash
sudo apt update
sudo apt install apache2
```

Enable required modules in apache:
```bash
sudo a2enmod proxy
sudo a2enmod proxy_http
sudo systemctl restart apache2
```

Create Virtual host configuration file:
```bash
sudo nano /etc/apache2/sites-available/library.conf
```

Add the proxy configuration mapping to your Go server port (8080):
```apache
<VirtualHost *:80>
    ServerName library.local
    ProxyPreserveHost On
    ProxyPass /api/ http://localhost:8080/
    ProxyPassReverse /api/ http://localhost:8080/
    ErrorLog /var/log/library-error.log
    CustomLog /var/log/library-access.log combined
</VirtualHost>
```

Enable the virtual host:
```bash
sudo a2ensite library.conf
sudo systemctl restart apache2
```

Run the go binary using `./mvc` and access the hosted site on `library.local`.