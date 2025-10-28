FROM python:3.11-slim

WORKDIR /tests

# Copy only requirements first for better layer caching
COPY requirements.txt .

# Install dependencies (cached as separate layer)
RUN pip install -r requirements.txt

# Copy test files
COPY . .

# Set environment variable for base URL
ENV BASE_URL=http://server:8080

# Run tests by default with 4 parallel workers for optimal performance
CMD ["pytest", "-v", "--tb=short", "-n", "4"]
