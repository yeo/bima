defmodule BimaWeb.Router do
  use BimaWeb, :router

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_flash
    plug :protect_from_forgery
    plug :put_secure_browser_headers
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/", BimaWeb do
    pipe_through :browser

    get "/", PageController, :index
    get "/buy", SubscriptionController, :index
  end

  # Other scopes may use custom stacks.
  scope "/api", BimaWeb do
    pipe_through :api

    post "/sync", Api.SyncController, :sync
    get "/me", Api.MeController, :index
    post "/app", Api.AppController, :create

    post "/blob", Api.BlobController, :create
    get "/blob/:code", Api.BlobController, :show
  end
end
