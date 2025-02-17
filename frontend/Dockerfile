FROM node:20-alpine

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of the application
COPY . .

# Expose port 3000
EXPOSE 3000

# Run the development server with the host flag for Docker
CMD ["npm", "run", "dev", "--", "-H", "0.0.0.0"] 