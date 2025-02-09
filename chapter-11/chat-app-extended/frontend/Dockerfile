# Use the official Node.js image to build the app
FROM node:18-alpine as build-stage

# Set the working directory for the build process
WORKDIR /app

# Install Angular CLI
RUN npm install -g @angular/cli

# Copy the application source code into the container
COPY ./ /app

# Install the dependencies
RUN npm install

# Build the Angular application for production
RUN ng build --prod

# Stage 2: Serve the app using Nginx
FROM nginx:alpine as production-stage

# Copy the built Angular app from the build-stage
COPY --from=build-stage /app/dist/chat-frontend /usr/share/nginx/html

# Expose the port Nginx is running on
EXPOSE 80

# Start the Nginx server to serve the Angular app
CMD ["nginx", "-g", "daemon off;"]
