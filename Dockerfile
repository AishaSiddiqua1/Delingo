# Use official Node.js image
FROM node:18

# Set working directory
WORKDIR /app

# Copy package.json and yarn.lock first to install dependencies
COPY package.json yarn.lock ./

# Install dependencies using Yarn
RUN yarn install

# Copy the rest of the application code
COPY . .

# Build the React app for production
RUN yarn build

# Expose the React app port
EXPOSE 3000

# Start the React app
CMD ["yarn", "start"]
