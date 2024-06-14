# Dockerfile for the frontend

FROM node:20-alpine as builder

WORKDIR /frontend

COPY package.json ./
COPY package-lock.json ./

RUN npm install

COPY . ./

RUN npm run build

FROM nginx:1.21-alpine as final

COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /frontend/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
