# aggregator

## Purpose
The purpose of this project is to provide a simple configuration management tool for Go applications. It allows users to read and write configuration settings from a JSON file located in the user's home directory.

## Setup
1. Clone the repository:
   ```
   git clone https://github.com/mrjkey/aggregator.git
   ```
2. Navigate to the project directory:
   ```
   cd aggregator
   ```
3. Install the dependencies:
   ```
   go mod tidy
   ```

## Usage
1. Create a `.gatorconfig.json` file in your home directory with the following content:
   ```json
   {
     "db_url": "your_database_url",
     "current_user_name": "your_username"
   }
   ```
2. Run the application:
   ```
   go run main.go
   ```

## Dependencies
- Go 1.24.0 or higher

## Contact
For any questions or suggestions, please contact the project maintainer at [your_email@example.com].

## Contribution Guidelines
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push the branch to your fork.
4. Create a pull request with a detailed description of your changes.
