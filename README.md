# 📝 Blogging Platform API

You can use the following endpoints to work with posts.

---

## 📌 Base URL
https://localhost:8080/posts

pgsql
Копіювати код

---

## ➕ POST `/posts`
Add a new post.  
You must provide JSON with post information.

### Example JSON
```json
{
  "title": "First Blo1g123",
  "content": "its my first blog",
  "category": "default",
  "tags": ["first", "default", "blog"]
}
```
✏️ PUT /posts/<id>
Update an existing post.
You must provide updated JSON.

Example JSON
```json
{
  "title": "First Blo1g123",
  "content": "its my first blog",
  "category": "default",
  "tags": ["first", "default", "blog"]
}
```
❌ DELETE /posts/<id>
Delete a post by ID.

🚧 In Development
More functionality is being worked on.

📎 Project Source
https://roadmap.sh/projects/blogging-platform-api
