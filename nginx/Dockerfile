FROM node:19 AS build
ARG FRONTEND_CODE_PATH="./frontend"
WORKDIR /home/code
COPY $FRONTEND_CODE_PATH .
RUN npm install
RUN npm run build

FROM nginx:1.22.1
RUN apt-get update -y && apt-get install nginx-extras -y
COPY --from=build /home/code/build /var/www/build
CMD ["nginx", "-g", "daemon off;"]
