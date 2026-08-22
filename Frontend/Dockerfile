# Multi-stage Dockerfile pour le Frontend Vue 3 + Nginx

# Étape 1 : Compilation des assets avec Vite
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Étape 2 : Serveur Nginx de production
FROM nginx:alpine

# Copie du build compilé
COPY --from=builder /app/dist /usr/share/nginx/html

# Copie de la configuration Nginx Docker
COPY nginx-docker.conf /etc/nginx/conf.d/default.conf

# Copie des pages d'erreur personnalisées (404, 500)
COPY errors /usr/share/nginx/errors

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
