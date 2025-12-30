# Customization Guide

This template is designed so downstream apps can add features without rewriting core auth, RBAC, and profile systems. Keep your changes inside the custom hooks whenever possible to reduce merge conflicts when pulling template updates.

## Frontend Customization

### Branding, Navigation, and Home Content
Edit `frontend/src/custom/appConfig.ts` to override template defaults:

- Brand name and short name (logo mark)
- Header navigation items
- Sidebar navigation items
- Footer links and social icons
- Home page title, tagline, CTAs, and feature list

Example:

```ts
export const customAppConfig = {
  brand: {
    name: 'My Blog',
    shortName: 'MB',
  },
  navigation: {
    header: [
      { path: '/', label: 'Home' },
      { path: '/blog', label: 'Blog' },
    ],
  },
  home: {
    title: 'My Blog',
    subtitle: 'Notes and essays',
    ctas: [{ label: 'Read Posts', to: '/blog', variant: 'primary' }],
  },
}
```

### Add Custom Pages and Routes
Add routes in `frontend/src/custom/routes.ts`:

```ts
import type { RouteRecordRaw } from 'vue-router'

export const customRoutes: RouteRecordRaw[] = [
  {
    path: '/blog',
    name: 'blog',
    component: () => import('@/views/Blog.vue'),
  },
]
```

## Backend Customization

### Add API Routes
Register API routes in `backend/internal/custom/routes.go`. This runs after core routes are created and gives you access to the database, config, JWT service, and RBAC middleware:

```go
func RegisterRoutes(v1 *gin.RouterGroup, deps Dependencies) {
  blog := v1.Group("/blog")
  blog.GET("", func(c *gin.Context) {
    c.JSON(200, gin.H{"posts": []string{}})
  })
}
```

Use `deps.JWT` and `deps.RBAC` if you want to protect routes with the same auth and permission systems as the core API.

## Keeping Template Updates Easy

- Prefer adding new functionality in `frontend/src/custom/*` and `backend/internal/custom/*`.
- Avoid editing core files unless necessary; template updates are less likely to conflict when your changes are isolated.
- `node_modules/` is ignored via `.gitignore` and should not be committed.
