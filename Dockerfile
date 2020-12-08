FROM node:14-alpine

WORKDIR /home/node/api/src

COPY package*.json ./
RUN npm install

# Bundle app source
COPY . .

EXPOSE 31337
CMD [ "npm", "start" ]
