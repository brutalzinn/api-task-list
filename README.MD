# github.com/brutalzinn/api-task-list
# STILL IN DEVELOPMENT

### A api to manage tasks by repositories and sync between a CLI and mobile application. Utils for people that likes to wrotten many annotations.


[Swagger documentation](https://brutalzinn.github.io/api-task-list/)
# This api its a backend for the bellow projects

- https://github.com/brutalzinn/web-task-list
- https://github.com/brutalzinn/cli-task-list


# Implementations

- [x] Authentication
- [x] Tasks route
- [x] Users route
- [x] Authentication route
- [x] Api key route ( revoke and generate )
- [x] Table users
- [x] Table tasks
- [x] Table api keys
- [x] Hypermedia links base
- [ ] Swagger
- [ ] Firebase integration

# Purpose of existence

## This project is designed to my estudies case of hypermedia REST apis.

#### This example shows apikey list route. This routes list all api keys of current authenticated user. The frontend application can uses the hypermedia references to perform actions with the api keys.
<img src="docs/imgs/postman_example_apikeys_list.png" alt="postman example of api keys list" width="500"/>

#### This example shows a repository paging route. These routes list all repositories created by the user. The cool part here is that the API client can see all available actions to perform. This means that the organic flow of route consumption is in the hands of the API.
<img src="docs/imgs/postman_example_repo_pagination.png" alt="postman example of api keys list" width="500"/>

### References

- [api design with hypermedia rest](https://www.mscharhag.com/api-design/hypermedia-rest)