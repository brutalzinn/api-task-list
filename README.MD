# github.com/brutalzinn/api-task-list
# STILL IN DEVELOPMENT

### A api to manage tasks by repositories and sync between a CLI and mobile application. Utils for people that likes to wrotten many annotations.


[Swagger documentation](https://brutalzinn.github.io/api-task-list/)
# This api its a backend for the bellow projects

- https://github.com/brutalzinn/web-task-list
- https://github.com/brutalzinn/cli-task-list


# Goals

- [x] Tasks routes ( put, patch, delete, get, post)
- [x] Users routes ( put, delete, get )
- [x] Login routes ( post )
- [x] Register routes ( post )
- [x] Api key routes ( delete, post, get )
- [x] OAuth2 routes ( delete, post, get )
- [x] Table users
- [x] Table tasks
- [x] Table repos
- [x] Table api_keys
- [x] Table users_oauth_client
- [x] Table oauth_client_application
- [x] Hypermedia links base
- [x] Hypermedia links middleware
- [ ] Unit tests
- [ ] Database tests using postgree temporary container
- [ ] Makefile to execute tests and easy test database queries
- [ ] CI/CD to digital ocean with Dokku
- [ ] Swagger auto documentation
- [ ] Notification channel using web hook at IFTTT
- [ ] Error log with MongoDB
- [ ] Generic error middleware handler
- [x] Configure Dockerfile to hot reload
- [ ] Error middleware



# Table relations

#### Version 4 by dbeaver ( still in development )
<img src="docs/imgs/dbeaver.png" alt="postman example of api keys list" height="50%"/>

The orange tables is created by go package oauth2. So i dont have this managed resource. To do my own oauth2 implementation of my business rule, i made a many to many helper table. 
# Purpose of existence

This project is designed to my estudies case of hypermedia REST apis(HATEOAS). This project is a overthink project. Why? Because i can.

# HATEOAS Example with repo paginate route

    {
        "error": false,
        "message": "",
        "data": {
            "repos": [
                {
                    "id": 1,
                    "title": "My first repo",
                    "description": "This is a example repo",
                    "create_at": "2023-03-25T01:26:29.854407Z",
                    "update_at": null,
                    "links": [
                        {
                            "rel": "delete",
                            "href": "/repo/1",
                            "type": "DELETE"
                        },
                        {
                            "rel": "update_one",
                            "href": "/repo/1",
                            "type": "PATCH"
                        },
                        {
                            "rel": "_self",
                            "href": "/repo/1",
                            "type": "GET"
                        }
                    ]
                }
            ],
            "totalItems": 1,
            "totalPages": 0,
            "currentPage": 1
        }
    }

# HATEOAS Example with api keys route

    {
        "error": false,
        "message": "",
        "data": [
            {
                "id": "574b32d5-f788-4b2c-aefb-979ef17a6602",
                "name": "Flutter app",
                "expire_at": "2023-04-24T01:32:07.641258Z",
                "scopes": "task_manager,repo_manager",
                "create_at": "2023-03-25T01:32:07.656142Z",
                "update_at": null,
                "links": [
                    {
                        "rel": "delete",
                        "href": "/apikey/delete/574b32d5-f788-4b2c-aefb-979ef17a6602",
                        "type": "DELETE"
                    },
                    {
                        "rel": "regenerate",
                        "href": "/apikey/regenerate/574b32d5-f788-4b2c-aefb-979ef17a6602",
                        "type": "POST"
                    }
                ]
            }
        ]
    }


# Web Example under construction

[Web App](https://github.com/brutalzinn/web-task-list)

Web app is a react frontend application to perform login and manage user api keys and oauth2.
# Dart CLI under construction

[Dart CLI](https://github.com/brutalzinn/cli-task-list)

Dart cli is a project designed to use this API.
 This is the principal frontend application that will interact with this api. The ideia is use oauth2 auth flow to CLi authenticates and perform actions. The user can create tasks and repos using this. The user can choose a preference default editor like noteppad to edit the tasks and send to server.
### References

- [api design with hypermedia rest](https://www.mscharhag.com/api-design/hypermedia-rest)


# Libs

- [go-oauth2/go-oauth2](https://github.com/go-oauth2/oauth2)
- [vgarvardt/go-pg-adapter](https://github.com/vgarvardt/go-pg-adapter)
- [brutalzinn/go-auth2-pg my fork of vgarvardt/go-oauth2-pg](https://github.com/brutalzinn/go-oauth2-pg)
- [go-chi/chi](github.com/go-chi/chi)