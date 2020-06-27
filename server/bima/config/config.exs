# This file is responsible for configuring your application
# and its dependencies with the aid of the Mix.Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
use Mix.Config

config :bima,
  ecto_repos: [Bima.Repo]

# Configures the endpoint
config :bima, BimaWeb.Endpoint,
  url: [host: "localhost"],
  secret_key_base: "O5r24zIyg6Dmc4lQg5wYBtjCntu6Yp5e86rRYA9WSwp3Z67yfyvlHxnXCzQdehat",
  render_errors: [view: BimaWeb.ErrorView, accepts: ~w(html json)],
  pubsub: [name: Bima.PubSub, adapter: Phoenix.PubSub.PG2]

# Configures Elixir's Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

config :bima_versions,
  latest: [1,1,2],
  url: %{mac: 'macurl', linux: 'linuxurl', window: 'windowurl'}

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{Mix.env()}.exs"
