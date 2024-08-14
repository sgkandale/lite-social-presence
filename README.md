## socialite

Backend for a simple online gaming party

- users can register and login
- they can send friend requests to each other
- act on the received friend requests (accept or reject)
- view their friend list
- remove any user as a friend

--

- users can create a party
- invite other users to join a party
- join the party they have been invited to
- leave the party
- remove users from the party

--

- subscribe to a websocket to ping their online status periodically
- receive a list of their friends who are currently online
- subscribe to another websocket with a party_id to see which party members, who are also their friends, are currently online

### API
- postman collection is there to interact with the backend
- all the APIs supported by this service and created in the collection, which is ready to test

### Deployment
- this service can be deployed as a single binary on any host server along with a config file, a sample of which is present in the codebase
- this service can also be deployed via docker compose
- Dockerfile and docker-compose file are located in build directory
- to deploy on kubernetes, deployment and service are also located in the same directory