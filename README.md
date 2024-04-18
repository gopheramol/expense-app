# expense-app
Expense App is a simple web application built with Go (Golang) and MongoDB, designed to help users keep track of their expenses. Users can add, view, and edit their expenses, as well as see the total amount spent.
## Features:
- Add Expenses: Users can add new expenses by providing a name and amount.
- View Expenses: All expenses are displayed in a table format, showing the name, amount, and creation time.
- Edit Expenses: Users can edit existing expenses to update the name or amount.
- Total Amount: The total amount spent is calculated and displayed at the bottom of the expense list.
  
* MongoDB Integration: Expenses are stored in a MongoDB database, allowing for easy storage and retrieval.

## Technologies Used:
- Go (Golang): The backend of the application is written in Go, providing fast and efficient handling of HTTP requests.
- MongoDB: Expenses are stored in a MongoDB database, providing a flexible and scalable data storage solution.
- Gin: The Gin web framework is used to handle routing and middleware for the web server.
- HTML/CSS: Simple HTML templates are used for rendering the user interface, styled with CSS for a clean and user-friendly experience.

## How to Run:
* Clone the repository: git clone https://github.com/gopheramol/expense-app.git
* Install dependencies: ```go mod tidy```
* Start the server: ```go run main.go```
* Access the application in your browser at `http://localhost:8080`
