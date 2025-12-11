U can use:

[https://localhost:8080/posts]:
---> POST(/posts): add the post. u must add information about post. 
example: JSON: {
    "title": "First Blo1g123",
    "content": "its my first blog",
    "category": "default",
    "tags": ["first", "default", "blog"]
}
---> PUT(/posts/<id>): update u post. u must give information about update(JSON)
example:
JSON: {
    "title": "First Blo1g123",
    "content": "its my first blog",
    "category": "default",
    "tags": ["first", "default", "blog"]
}
--> DELETE(/posts/<id>): delete u post.
-->  (in develop)

https://roadmap.sh/projects/blogging-platform-api
