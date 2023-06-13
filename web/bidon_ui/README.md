# Nuxt 3 Minimal Starter

Look at the [Nuxt 3 documentation](https://nuxt.com/docs/getting-started/introduction) to learn more.

## Setup

Make sure to install the dependencies:

```bash
yarn install
```

## Development Server

Start the development server on `http://localhost:3000`

```bash
yarn dev
```

## Local api

```bash
go run cmd/bidon-admin/main.go
```

## Static build locally

Frontend can be built statically, and that's how the production version works: frontend is built and then resulting files embedded into the backend, which serves them as-is. Node is not running, only the backend written in Go programming language, which also serves pre-built frontend HTML and JS and CSS files.

Run `yarn generate` inside `./web/bidon_ui/`

Copy static files to embed them into Go binary

```bash
cp -rf web/bidon_ui/.output/public/ .cmd/bidon-admin/web/ui # assume you are in the root dir
```

Run backend

```bash
go run cmd/bidon-admin/main.go
```

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.
