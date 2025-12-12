# Use Node.js image
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy dependency files first
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the code
COPY . .

# Build the Next.js app for production
RUN npm run build

# Expose the port
EXPOSE 3000

# Start the app
CMD ["npm", "start"]