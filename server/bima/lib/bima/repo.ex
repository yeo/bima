defmodule Bima.Repo do
  use Ecto.Repo,
    otp_app: :bima,
    adapter: Ecto.Adapters.Postgres
end
