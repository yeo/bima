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

  pipeline :authenticate do
    #plug BasicAuth, username: "bima", password: "bima"
    plug BimaWeb.Plugs.Auth
  end

  pipeline :admin_only do
    plug BimaWeb.Plugs.Auth, admin: true
  end

	scope "/", BimaWeb do
		pipe_through :browser

		get "/", PageController, :index
    get "/buy", SubscriptionController, :index

    delete "/logout", SessionController, :delete
  end

  scope "/", BimaWeb do
		pipe_through [:browser, BimaWeb.Plugs.Guest]

    resources "/register", UserController, only: [:create, :new, :show]
    get "/login", SessionController, :new
    post "/login", SessionController, :create
  end

	# Other scopes may use custom stacks.
	scope "/api", BimaWeb do
		pipe_through :api

		post "/sync", Api.SyncController, :sync
		get "/me", Api.MeController, :index
    #post "/app", Api.AppController, :create

		post "/blob", Api.BlobController, :create
		get "/blob/:code", Api.BlobController, :show
	end

  use Kaffy.Routes, scope: "/admin", pipe_through: [:browser, :admin_only]
end
